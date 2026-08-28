package cmd_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/cmd"
	"github.com/twoBoots/antimatter/pkg/adapters"
	"github.com/twoBoots/antimatter/pkg/adapters/antigravity"
	"github.com/twoBoots/antimatter/pkg/adapters/claude"
	"github.com/twoBoots/antimatter/pkg/policy"
)

func TestExportAndImportCmd_Execution(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "policy.yaml")
	agDir := filepath.Join(tmpDir, "antigravity-cli")
	claudeFile := filepath.Join(tmpDir, ".claude.json")
	backupDir := filepath.Join(tmpDir, "backups")

	// Set up manager with custom test adapters
	mgr := adapters.NewManager()
	mgr.Register(antigravity.NewAdapter(
		antigravity.WithConfigDir(agDir),
		antigravity.WithBackupDir(backupDir),
	))
	mgr.Register(claude.NewAdapter(
		claude.WithConfigPath(claudeFile),
		claude.WithBackupDir(backupDir),
	))
	cmd.SetDefaultManager(mgr)

	// Save initial policy
	p := policy.NewPolicy("test-export", "Export Test")
	p.Rules = []policy.PolicyRule{
		{
			ID:        "rule-echo",
			Type:      policy.RuleTypeCommand,
			Pattern:   "echo hello",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
		},
	}
	require.NoError(t, policy.SaveToFile(p, configFile))

	t.Run("export to all agents", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "export", "--target", "all"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Exported policy to Google Antigravity")
		assert.Contains(t, buf.String(), "Exported policy to Claude Code")
		assert.FileExists(t, filepath.Join(agDir, "whitelist.json"))
		assert.FileExists(t, claudeFile)
	})

	t.Run("import from antigravity", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "import", "--from", "antigravity"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Successfully imported rules from antigravity")
	})

	t.Run("import from unknown agent error", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "import", "--from", "unknown_agent"})

		err := rootCmd.Execute()
		assert.Error(t, err)
	})
}
