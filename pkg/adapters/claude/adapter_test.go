package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/pkg/adapters/claude"
	"github.com/twoBoots/antimatter/pkg/policy"
)

func TestClaudeAdapter_ExportAndImport(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, ".claude.json")
	backupDir := filepath.Join(tmpDir, "backups")

	// Pre-populate with existing settings to verify safe merge
	initialJSON := `{
		"theme": "dark",
		"allowedTools": ["Bash"],
		"approvedCommands": ["ls"]
	}`
	require.NoError(t, os.WriteFile(configFile, []byte(initialJSON), 0644))

	adapter := claude.NewAdapter(
		claude.WithConfigPath(configFile),
		claude.WithBackupDir(backupDir),
	)

	assert.Equal(t, "claude", adapter.Name())
	assert.Equal(t, "Claude Code", adapter.DisplayName())
	assert.Equal(t, configFile, adapter.ConfigPath())

	// Step 1: Export a canonical policy
	p := policy.NewPolicy("claude-test-profile", "Test")
	p.Rules = []policy.PolicyRule{
		{
			ID:          "claude-git-status",
			Type:        policy.RuleTypeCommand,
			Pattern:     "git status",
			MatchMode:   policy.MatchExact,
			Decision:    policy.DecisionAllow,
			Description: "Allow git status in Claude",
			Agents:      []string{"claude"},
		},
		{
			ID:        "claude-mcp-tool",
			Type:      policy.RuleTypeTool,
			Pattern:   "mcp_jira_*",
			MatchMode: policy.MatchGlob,
			Decision:  policy.DecisionAllow,
			Agents:    []string{"all"},
		},
		{
			ID:        "ag-exclusive-rule",
			Type:      policy.RuleTypeCommand,
			Pattern:   "gemini run",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
			Agents:    []string{"antigravity"}, // Should be filtered out from Claude export
		},
	}

	err := adapter.Export(p)
	require.NoError(t, err)

	// Verify existing fields ("theme") are preserved
	content, err := os.ReadFile(configFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"theme": "dark"`)
	assert.Contains(t, string(content), "git status")

	// Step 2: Import from written file
	imported, err := adapter.Import()
	require.NoError(t, err)
	assert.Equal(t, "imported-claude", imported.Name)

	var hasGitStatus, hasJiraTool bool
	for _, r := range imported.Rules {
		if r.Pattern == "git status" && r.Type == policy.RuleTypeCommand {
			hasGitStatus = true
		}
		if r.Pattern == "mcp_jira_*" && r.Type == policy.RuleTypeTool {
			hasJiraTool = true
		}
	}
	assert.True(t, hasGitStatus)
	assert.True(t, hasJiraTool)

	// Step 3: Verify backup was created
	backups, err := os.ReadDir(filepath.Join(backupDir, "claude"))
	require.NoError(t, err)
	assert.NotEmpty(t, backups)
}
