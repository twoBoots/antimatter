package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/twoBoots/antimatter/pkg/policy"
)

// NewWhitelistCmd returns the whitelist command group
func NewWhitelistCmd() *cobra.Command {
	whitelistCmd := &cobra.Command{
		Use:     "whitelist",
		Aliases: []string{"wl", "policy"},
		Short:   "Manage, inspect, and evaluate whitelist policies",
		Long:    "Inspect, add, remove, and evaluate command and tool call permissions against active policies.",
	}

	whitelistCmd.AddCommand(newListCmd())
	whitelistCmd.AddCommand(newAddCmd())
	whitelistCmd.AddCommand(newRemoveCmd())
	whitelistCmd.AddCommand(newCheckCmd())

	return whitelistCmd
}

func getPolicyPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "policy.yaml"
	}
	return filepath.Join(home, ".antimatter", "policy.yaml")
}

func loadOrCreatePolicy(path string) (*policy.Policy, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		p := policy.NewPolicy("default", "Default Antimatter Policy")
		if err := policy.SaveToFile(p, path); err != nil {
			return nil, err
		}
		return p, nil
	}
	return policy.LoadFromFile(path)
}

func newListCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all active whitelist rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			p, err := loadOrCreatePolicy(path)
			if err != nil {
				return err
			}

			if jsonOutput {
				data, err := json.MarshalIndent(p.Rules, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(data))
				return nil
			}

			cmd.Printf("Policy: %s (%s) [%d rules]\n", p.Name, path, len(p.Rules))
			cmd.Println(strings.Repeat("-", 60))
			for _, r := range p.Rules {
				cmd.Printf("• [%s] %-16s (%s:%s) -> %s\n", r.Type, r.ID, r.MatchMode, r.Pattern, strings.ToUpper(string(r.Decision)))
				if r.Description != "" {
					cmd.Printf("  Desc: %s\n", r.Description)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	return cmd
}

func newAddCmd() *cobra.Command {
	var (
		id          string
		commandPat  string
		toolPat     string
		matchMode   string
		decision    string
		description string
		agents      []string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new whitelist rule to policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if commandPat == "" && toolPat == "" {
				return fmt.Errorf("either --command or --tool pattern must be specified")
			}

			ruleType := policy.RuleTypeCommand
			pattern := commandPat
			if toolPat != "" {
				ruleType = policy.RuleTypeTool
				pattern = toolPat
			}

			if id == "" {
				// generate slug ID
				id = fmt.Sprintf("%s-%d", ruleType, time.Now().Unix())
			}

			if matchMode == "" {
				matchMode = string(policy.MatchExact)
			}
			if decision == "" {
				decision = string(policy.DecisionAllow)
			}

			rule := policy.PolicyRule{
				ID:          id,
				Type:        ruleType,
				Pattern:     pattern,
				MatchMode:   policy.MatchMode(matchMode),
				Decision:    policy.Decision(decision),
				Description: description,
				Agents:      agents,
				CreatedAt:   time.Now().UTC(),
			}

			path := getPolicyPath()
			p, err := loadOrCreatePolicy(path)
			if err != nil {
				return err
			}

			p.Rules = append(p.Rules, rule)
			if err := policy.SaveToFile(p, path); err != nil {
				return fmt.Errorf("failed to save policy: %w", err)
			}

			cmd.Printf("Successfully added rule '%s' (%s: %s)\n", id, ruleType, pattern)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "unique rule ID (optional, generated if omitted)")
	cmd.Flags().StringVar(&commandPat, "command", "", "command line pattern to whitelist")
	cmd.Flags().StringVar(&toolPat, "tool", "", "tool call name/pattern to whitelist")
	cmd.Flags().StringVar(&matchMode, "mode", "exact", "matching mode (exact, prefix, glob, regex)")
	cmd.Flags().StringVar(&decision, "decision", "allow", "decision (allow, deny, ask)")
	cmd.Flags().StringVar(&description, "desc", "", "description / rationale for rule")
	cmd.Flags().StringSliceVar(&agents, "agents", []string{"all"}, "agents this rule applies to (antigravity, claude, all)")

	return cmd
}

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [rule_id]",
		Short: "Remove a whitelist rule by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			path := getPolicyPath()
			p, err := loadOrCreatePolicy(path)
			if err != nil {
				return err
			}

			newRules := make([]policy.PolicyRule, 0, len(p.Rules))
			found := false
			for _, r := range p.Rules {
				if r.ID == id {
					found = true
					continue
				}
				newRules = append(newRules, r)
			}

			if !found {
				return fmt.Errorf("rule with ID '%s' not found", id)
			}

			p.Rules = newRules
			if err := policy.SaveToFile(p, path); err != nil {
				return fmt.Errorf("failed to save policy: %w", err)
			}

			cmd.Printf("Successfully removed rule '%s'\n", id)
			return nil
		},
	}

	return cmd
}

func newCheckCmd() *cobra.Command {
	var (
		agent    string
		toolName string
	)

	cmd := &cobra.Command{
		Use:   "check [command]",
		Short: "Evaluate a candidate command or tool call against active policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			p, err := loadOrCreatePolicy(path)
			if err != nil {
				return err
			}

			var res policy.EvaluationResult
			if toolName != "" {
				res = policy.EvaluateTool(p, agent, toolName)
				cmd.Printf("Evaluating Tool: '%s' (Agent: %s)\n", toolName, agent)
			} else {
				if len(args) == 0 {
					return fmt.Errorf("must specify a command argument or --tool flag")
				}
				candidateCmd := strings.Join(args, " ")
				res = policy.EvaluateCommand(p, agent, candidateCmd)
				cmd.Printf("Evaluating Command: '%s' (Agent: %s)\n", candidateCmd, agent)
			}

			cmd.Printf("DECISION: %s\n", strings.ToUpper(string(res.Decision)))
			cmd.Printf("Reason:   %s\n", res.Reason)
			if res.MatchedRule != nil {
				cmd.Printf("Matched Rule ID: %s (Mode: %s, Pattern: %s)\n", res.MatchedRule.ID, res.MatchedRule.MatchMode, res.MatchedRule.Pattern)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "all", "evaluating agent (antigravity, claude, all)")
	cmd.Flags().StringVar(&toolName, "tool", "", "evaluate tool name instead of command")

	return cmd
}
