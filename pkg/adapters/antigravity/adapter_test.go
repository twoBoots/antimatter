package antigravity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/pkg/adapters/antigravity"
	"github.com/twoBoots/antimatter/pkg/policy"
)

func TestAntigravityAdapter_ExportAndImport(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "antigravity-cli")
	backupDir := filepath.Join(tmpDir, "backups")

	adapter := antigravity.NewAdapter(
		antigravity.WithConfigDir(configDir),
		antigravity.WithBackupDir(backupDir),
	)

	assert.Equal(t, "antigravity", adapter.Name())
	assert.Equal(t, "Google Antigravity", adapter.DisplayName())
	assert.Equal(t, filepath.Join(configDir, "whitelist.json"), adapter.ConfigPath())

	// Step 1: Export a canonical policy
	p := policy.NewPolicy("dev-profile", "Dev Profile")
	p.Rules = []policy.PolicyRule{
		{
			ID:          "ag-git-status",
			Type:        policy.RuleTypeCommand,
			Pattern:     "git status",
			MatchMode:   policy.MatchExact,
			Decision:    policy.DecisionAllow,
			Description: "Allow checking git status",
			Agents:      []string{"antigravity"},
		},
		{
			ID:        "ag-view-file",
			Type:      policy.RuleTypeTool,
			Pattern:   "view_file",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
			Agents:    []string{"all"},
		},
		{
			ID:        "claude-only-cmd",
			Type:      policy.RuleTypeCommand,
			Pattern:   "pytest",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
			Agents:    []string{"claude"}, // Should be filtered out from Antigravity export
		},
	}

	err := adapter.Export(p)
	require.NoError(t, err)
	assert.FileExists(t, adapter.ConfigPath())

	// Step 2: Import from written file
	imported, err := adapter.Import()
	require.NoError(t, err)
	assert.Equal(t, "imported-antigravity", imported.Name)
	require.Len(t, imported.Rules, 2) // only the 2 applicable rules

	// Verify imported command rule
	cmdRule := imported.Rules[0]
	assert.Equal(t, policy.RuleTypeCommand, cmdRule.Type)
	assert.Equal(t, "git status", cmdRule.Pattern)
	assert.Equal(t, policy.MatchExact, cmdRule.MatchMode)
	assert.Contains(t, cmdRule.Agents, "antigravity")

	// Step 3: Export again to verify backup creation
	err = adapter.Export(p)
	require.NoError(t, err)

	backups, err := os.ReadDir(filepath.Join(backupDir, "antigravity"))
	require.NoError(t, err)
	assert.NotEmpty(t, backups)
}
