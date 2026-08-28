package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/twoBoots/antimatter/pkg/adapters"
	"github.com/twoBoots/antimatter/pkg/policy"
)

const (
	AdapterName        = "claude"
	AdapterDisplayName = "Claude Code"
)

// Adapter implements the adapters.Adapter interface for Claude Code
type Adapter struct {
	configPath string
	backupDir  string
}

// Option configures the Claude Adapter
type Option func(*Adapter)

// WithConfigPath sets a custom path for .claude.json
func WithConfigPath(path string) Option {
	return func(a *Adapter) {
		a.configPath = path
	}
}

// WithBackupDir sets a custom backup base directory
func WithBackupDir(dir string) Option {
	return func(a *Adapter) {
		a.backupDir = dir
	}
}

// NewAdapter creates a Claude adapter instance
func NewAdapter(opts ...Option) *Adapter {
	home, _ := os.UserHomeDir()
	a := &Adapter{
		configPath: filepath.Join(home, ".claude.json"),
		backupDir:  adapters.DefaultBackupDir(),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *Adapter) Name() string {
	return AdapterName
}

func (a *Adapter) DisplayName() string {
	return AdapterDisplayName
}

func (a *Adapter) Detected() bool {
	_, err := os.Stat(a.configPath)
	return err == nil
}

func (a *Adapter) ConfigPath() string {
	return a.configPath
}

func (a *Adapter) Backup() (string, error) {
	return adapters.CreateBackupWithBaseDir(a.configPath, AdapterName, a.backupDir)
}

func (a *Adapter) Export(p *policy.Policy) error {
	if p == nil {
		return fmt.Errorf("policy cannot be nil")
	}

	// Step 1: Create backup of existing config
	if _, err := a.Backup(); err != nil {
		return fmt.Errorf("failed to create pre-export backup: %w", err)
	}

	// Step 2: Read existing config map to preserve other keys
	configMap := make(map[string]interface{})
	if data, err := os.ReadFile(a.configPath); err == nil {
		_ = json.Unmarshal(data, &configMap)
	}

	// Step 3: Extract tools and commands
	allowedTools := make([]string, 0)
	approvedCommands := make([]string, 0)

	for _, rule := range p.Rules {
		if !isApplicableToClaude(rule.Agents) {
			continue
		}

		if rule.Type == policy.RuleTypeCommand {
			approvedCommands = append(approvedCommands, rule.Pattern)
		} else if rule.Type == policy.RuleTypeTool {
			allowedTools = append(allowedTools, rule.Pattern)
		}
	}

	configMap["allowedTools"] = allowedTools
	configMap["approvedCommands"] = approvedCommands

	dir := filepath.Dir(a.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory '%s': %w", dir, err)
	}

	data, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Claude config: %w", err)
	}

	tmpPath := a.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tmpPath, a.configPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to atomically write Claude config: %w", err)
	}

	return nil
}

func (a *Adapter) Import() (*policy.Policy, error) {
	data, err := os.ReadFile(a.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Claude config: %w", err)
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(data, &configMap); err != nil {
		return nil, fmt.Errorf("failed to parse Claude config: %w", err)
	}

	p := policy.NewPolicy("imported-claude", "Imported from Claude Code configuration")
	p.Rules = make([]policy.PolicyRule, 0)

	// Extract commands
	if rawCmds, ok := configMap["approvedCommands"].([]interface{}); ok {
		for i, cmdItem := range rawCmds {
			if cmdStr, ok := cmdItem.(string); ok && cmdStr != "" {
				p.Rules = append(p.Rules, policy.PolicyRule{
					ID:          fmt.Sprintf("claude-cmd-%d", i+1),
					Type:        policy.RuleTypeCommand,
					Pattern:     cmdStr,
					MatchMode:   policy.MatchExact,
					Decision:    policy.DecisionAllow,
					Description: "Imported from Claude Code approvedCommands",
					Agents:      []string{"claude"},
					CreatedAt:   time.Now().UTC(),
				})
			}
		}
	}

	// Extract tools
	if rawTools, ok := configMap["allowedTools"].([]interface{}); ok {
		for i, toolItem := range rawTools {
			if toolStr, ok := toolItem.(string); ok && toolStr != "" {
				p.Rules = append(p.Rules, policy.PolicyRule{
					ID:          fmt.Sprintf("claude-tool-%d", i+1),
					Type:        policy.RuleTypeTool,
					Pattern:     toolStr,
					MatchMode:   policy.MatchGlob,
					Decision:    policy.DecisionAllow,
					Description: "Imported from Claude Code allowedTools",
					Agents:      []string{"claude"},
					CreatedAt:   time.Now().UTC(),
				})
			}
		}
	}

	return p, nil
}

func isApplicableToClaude(agents []string) bool {
	if len(agents) == 0 {
		return true
	}
	for _, a := range agents {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "all" || a == "*" || a == "claude" || a == "claude-code" || a == "anthropic" {
			return true
		}
	}
	return false
}
