package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/twoBoots/antimatter/pkg/policy"
	"github.com/twoBoots/antimatter/pkg/version"
	"github.com/twoBoots/bender/pkg/mcp"
)

var toolHandlerRegistry = make(map[string]mcp.ToolHandler)

// GetToolHandler retrieves a registered MCP tool handler for testing or direct invocation
func GetToolHandler(srv *mcp.Server, name string) (mcp.ToolHandler, bool) {
	h, ok := toolHandlerRegistry[name]
	return h, ok
}

// BuildMCPServer constructs an MCP server with all Antimatter tools registered
func BuildMCPServer(policyPath string) *mcp.Server {
	cwd, _ := os.Getwd()
	srv := mcp.NewServer("antimatter", version.Version, cwd)

	// Tool 1: antimatter_get_whitelist
	getWhitelistTool := mcp.Tool{
		Name:        "antimatter_get_whitelist",
		Description: "Returns all active whitelisted commands and tool permissions in the Antimatter policy.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
		},
	}
	getWhitelistHandler := func(ctx context.Context, args map[string]interface{}) (mcp.CallToolResult, error) {
		p, err := loadOrCreatePolicy(policyPath)
		if err != nil {
			return mcp.NewErrorResult(fmt.Sprintf("Failed to load policy: %v", err)), nil
		}
		data, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			return mcp.NewErrorResult(fmt.Sprintf("Failed to marshal policy: %v", err)), nil
		}
		return mcp.NewTextResult(string(data), false), nil
	}
	srv.RegisterTool(getWhitelistTool, getWhitelistHandler)
	toolHandlerRegistry["antimatter_get_whitelist"] = getWhitelistHandler

	// Tool 2: antimatter_check_command
	checkCommandTool := mcp.Tool{
		Name:        "antimatter_check_command",
		Description: "Evaluates whether a candidate shell command or tool call is whitelisted under active policy.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"command": map[string]interface{}{
					"type":        "string",
					"description": "The shell command line to evaluate (e.g. 'git status')",
				},
				"tool": map[string]interface{}{
					"type":        "string",
					"description": "Optional tool name to evaluate instead of command",
				},
				"agent": map[string]interface{}{
					"type":        "string",
					"description": "Agent evaluating the rule (e.g. 'antigravity', 'claude', 'all')",
				},
			},
		},
	}
	checkCommandHandler := func(ctx context.Context, args map[string]interface{}) (mcp.CallToolResult, error) {
		p, err := loadOrCreatePolicy(policyPath)
		if err != nil {
			return mcp.NewErrorResult(fmt.Sprintf("Failed to load policy: %v", err)), nil
		}

		agent, _ := args["agent"].(string)
		if agent == "" {
			agent = "all"
		}

		if toolName, ok := args["tool"].(string); ok && toolName != "" {
			res := policy.EvaluateTool(p, agent, toolName)
			out := fmt.Sprintf("DECISION: %s\nReason: %s", strings.ToUpper(string(res.Decision)), res.Reason)
			return mcp.NewTextResult(out, false), nil
		}

		command, _ := args["command"].(string)
		if command == "" {
			return mcp.NewErrorResult("Either 'command' or 'tool' argument is required"), nil
		}

		res := policy.EvaluateCommand(p, agent, command)
		out := fmt.Sprintf("DECISION: %s\nReason: %s", strings.ToUpper(string(res.Decision)), res.Reason)
		if res.MatchedRule != nil {
			out += fmt.Sprintf("\nMatched Rule ID: %s (%s: %s)", res.MatchedRule.ID, res.MatchedRule.MatchMode, res.MatchedRule.Pattern)
		}
		return mcp.NewTextResult(out, false), nil
	}
	srv.RegisterTool(checkCommandTool, checkCommandHandler)
	toolHandlerRegistry["antimatter_check_command"] = checkCommandHandler

	// Tool 3: antimatter_add_whitelist_rule
	addRuleTool := mcp.Tool{
		Name:        "antimatter_add_whitelist_rule",
		Description: "Adds a new whitelisted command pattern or tool call into the version-controlled policy.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Optional unique rule ID",
				},
				"command": map[string]interface{}{
					"type":        "string",
					"description": "Command line pattern to whitelist",
				},
				"tool": map[string]interface{}{
					"type":        "string",
					"description": "Tool name pattern to whitelist",
				},
				"mode": map[string]interface{}{
					"type":        "string",
					"description": "Match mode: 'exact', 'prefix', 'glob', or 'regex'",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Rationale / note for the rule",
				},
			},
		},
	}
	addRuleHandler := func(ctx context.Context, args map[string]interface{}) (mcp.CallToolResult, error) {
		p, err := loadOrCreatePolicy(policyPath)
		if err != nil {
			return mcp.NewErrorResult(fmt.Sprintf("Failed to load policy: %v", err)), nil
		}

		commandPat, _ := args["command"].(string)
		toolPat, _ := args["tool"].(string)
		if commandPat == "" && toolPat == "" {
			return mcp.NewErrorResult("Must specify 'command' or 'tool' pattern"), nil
		}

		ruleType := policy.RuleTypeCommand
		pattern := commandPat
		if toolPat != "" {
			ruleType = policy.RuleTypeTool
			pattern = toolPat
		}

		id, _ := args["id"].(string)
		if id == "" {
			id = fmt.Sprintf("%s-%d", ruleType, time.Now().Unix())
		}

		modeStr, _ := args["mode"].(string)
		if modeStr == "" {
			modeStr = "exact"
		}
		desc, _ := args["description"].(string)

		rule := policy.PolicyRule{
			ID:          id,
			Type:        ruleType,
			Pattern:     pattern,
			MatchMode:   policy.MatchMode(modeStr),
			Decision:    policy.DecisionAllow,
			Description: desc,
			Agents:      []string{"all"},
			CreatedAt:   time.Now().UTC(),
		}

		p.Rules = append(p.Rules, rule)
		if err := policy.SaveToFile(p, policyPath); err != nil {
			return mcp.NewErrorResult(fmt.Sprintf("Failed to save policy: %v", err)), nil
		}

		return mcp.NewTextResult(fmt.Sprintf("Successfully added rule '%s' (%s: %s)", id, ruleType, pattern), false), nil
	}
	srv.RegisterTool(addRuleTool, addRuleHandler)
	toolHandlerRegistry["antimatter_add_whitelist_rule"] = addRuleHandler

	return srv
}

// NewMCPCmd returns the MCP server and installer command group
func NewMCPCmd() *cobra.Command {
	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Launch Model Context Protocol (MCP) server or configure AI clients",
		Long:  "Runs a stdio JSON-RPC 2.0 MCP server for Antigravity, Claude Code, and other AI assistants.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			srv := BuildMCPServer(path)
			return srv.Serve(cmd.Context(), os.Stdin, os.Stdout)
		},
	}

	mcpCmd.AddCommand(newMCPInstallCmd())

	return mcpCmd
}

func newMCPInstallCmd() *cobra.Command {
	var (
		listClients bool
		clients     []string
		allClients  bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Antimatter MCP server configuration into local AI assistants",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, _ := os.Getwd()
			home, _ := os.UserHomeDir()
			supported := mcp.GetSupportedClients(cwd, home)

			if listClients {
				cmd.Println("Supported AI Clients for MCP auto-configuration:")
				for _, c := range supported {
					detectedStatus := "not found"
					if c.Detected {
						detectedStatus = "detected"
					}
					cmd.Printf("• %-16s (%s) [%s] - %s\n", c.DisplayName, c.ID, c.ConfigPath, detectedStatus)
				}
				return nil
			}

			clientIDs := clients
			if allClients {
				clientIDs = nil // nil in Bender InstallerOptions installs all supported clients
			}

			results, err := mcp.InstallClients(mcp.InstallerOptions{
				ServerName: "antimatter",
				Command:    "antimatter",
				Args:       []string{"mcp"},
				ClientIDs:  clientIDs,
				Cwd:        cwd,
				HomeDir:    home,
			})
			if err != nil {
				return fmt.Errorf("failed to install MCP configurations: %w", err)
			}

			for _, r := range results {
				if r.Error == nil {
					cmd.Printf("✔ Successfully configured %s (%s)\n", r.DisplayName, r.ConfigPath)
				} else {
					cmd.Printf("✖ Failed to configure %s: %v\n", r.DisplayName, r.Error)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&listClients, "list", false, "list all supported AI clients")
	cmd.Flags().StringSliceVar(&clients, "client", []string{}, "comma-separated list of clients to configure (antigravity, claude, cursor)")
	cmd.Flags().BoolVar(&allClients, "all", false, "configure all discovered AI clients")

	return cmd
}
