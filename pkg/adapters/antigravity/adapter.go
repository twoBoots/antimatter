package antigravity

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
	AdapterName        = "antigravity"
	AdapterDisplayName = "Google Antigravity"
)

// WhitelistedCommand defines Antigravity's JSON structure for command rules
type WhitelistedCommand struct {
	Pattern     string `json:"pattern"`
	Mode        string `json:"mode"`
	Description string `json:"description,omitempty"`
	Decision    string `json:"decision,omitempty"`
}

// ConfigFile represents Antigravity's whitelist storage structure
type ConfigFile struct {
	WhitelistedCommands []WhitelistedCommand `json:"whitelistedCommands"`
	WhitelistedTools    []string             `json:"whitelistedTools"`
}

// Adapter implements the adapters.Adapter interface for Google Antigravity
type Adapter struct {
	configDir string
	backupDir string
}

// Option configures the Antigravity Adapter
type Option func(*Adapter)

// WithConfigDir sets a custom configuration directory
func WithConfigDir(dir string) Option {
	return func(a *Adapter) {
		a.configDir = dir
	}
}

// WithBackupDir sets a custom backup directory
func WithBackupDir(dir string) Option {
	return func(a *Adapter) {
		a.backupDir = dir
	}
}

// NewAdapter creates an Antigravity adapter instance
func NewAdapter(opts ...Option) *Adapter {
	home, _ := os.UserHomeDir()
	a := &Adapter{
		configDir: filepath.Join(home, ".gemini", "antigravity-cli"),
		backupDir: adapters.DefaultBackupDir(),
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
	info, err := os.Stat(a.configDir)
	return err == nil && info.IsDir()
}

func (a *Adapter) ConfigPath() string {
	return filepath.Join(a.configDir, "whitelist.json")
}

func (a *Adapter) Backup() (string, error) {
	return adapters.CreateBackupWithBaseDir(a.ConfigPath(), AdapterName, a.backupDir)
}

func (a *Adapter) Export(p *policy.Policy) error {
	if p == nil {
		return fmt.Errorf("policy cannot be nil")
	}

	// Step 1: Create backup if target file exists
	if _, err := a.Backup(); err != nil {
		return fmt.Errorf("failed to create pre-export backup: %w", err)
	}

	// Step 2: Transform canonical policy to Antigravity config
	var config ConfigFile
	config.WhitelistedCommands = make([]WhitelistedCommand, 0)
	config.WhitelistedTools = make([]string, 0)

	for _, rule := range p.Rules {
		if !isApplicableToAntigravity(rule.Agents) {
			continue
		}

		if rule.Type == policy.RuleTypeCommand {
			config.WhitelistedCommands = append(config.WhitelistedCommands, WhitelistedCommand{
				Pattern:     rule.Pattern,
				Mode:        string(rule.MatchMode),
				Description: rule.Description,
				Decision:    string(rule.Decision),
			})
		} else if rule.Type == policy.RuleTypeTool {
			config.WhitelistedTools = append(config.WhitelistedTools, rule.Pattern)
		}
	}

	if err := os.MkdirAll(a.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory '%s': %w", a.configDir, err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Antigravity config: %w", err)
	}

	targetPath := a.ConfigPath()
	tmpPath := targetPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tmpPath, targetPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to atomically replace config file: %w", err)
	}

	return nil
}

func (a *Adapter) Import() (*policy.Policy, error) {
	data, err := os.ReadFile(a.ConfigPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read Antigravity config: %w", err)
	}

	var config ConfigFile
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse Antigravity config: %w", err)
	}

	p := policy.NewPolicy("imported-antigravity", "Imported from Google Antigravity configuration")
	p.Rules = make([]policy.PolicyRule, 0)

	for i, cmd := range config.WhitelistedCommands {
		mode := policy.MatchMode(cmd.Mode)
		if mode == "" {
			mode = policy.MatchExact
		}
		dec := policy.Decision(cmd.Decision)
		if dec == "" {
			dec = policy.DecisionAllow
		}
		p.Rules = append(p.Rules, policy.PolicyRule{
			ID:          fmt.Sprintf("ag-cmd-%d", i+1),
			Type:        policy.RuleTypeCommand,
			Pattern:     cmd.Pattern,
			MatchMode:   mode,
			Decision:    dec,
			Description: cmd.Description,
			Agents:      []string{"antigravity"},
			CreatedAt:   time.Now().UTC(),
		})
	}

	for i, tool := range config.WhitelistedTools {
		p.Rules = append(p.Rules, policy.PolicyRule{
			ID:        fmt.Sprintf("ag-tool-%d", i+1),
			Type:      policy.RuleTypeTool,
			Pattern:   tool,
			MatchMode: policy.MatchGlob,
			Decision:  policy.DecisionAllow,
			Agents:    []string{"antigravity"},
			CreatedAt: time.Now().UTC(),
		})
	}

	return p, nil
}

func isApplicableToAntigravity(agents []string) bool {
	if len(agents) == 0 {
		return true
	}
	for _, a := range agents {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "all" || a == "*" || a == "antigravity" || a == "agy" || a == "gemini" {
			return true
		}
	}
	return false
}
