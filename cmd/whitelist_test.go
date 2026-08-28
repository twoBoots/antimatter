package cmd_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/cmd"
	"github.com/twoBoots/antimatter/pkg/policy"
)

func TestWhitelistCmd_Hierarchy(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "policy.yaml")

	p := policy.NewPolicy("test-profile", "Test")
	p.Rules = []policy.PolicyRule{
		{
			ID:        "rule-1",
			Type:      policy.RuleTypeCommand,
			Pattern:   "git status",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
		},
	}
	require.NoError(t, policy.SaveToFile(p, configFile))

	t.Run("list rules", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "whitelist", "list"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "rule-1")
		assert.Contains(t, buf.String(), "git status")
	})

	t.Run("add rule", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "whitelist", "add", "--id", "rule-2", "--command", "npm test*", "--mode", "glob", "--desc", "Allow npm test"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Successfully added rule 'rule-2'")

		// Verify persisted
		loaded, err := policy.LoadFromFile(configFile)
		require.NoError(t, err)
		require.Len(t, loaded.Rules, 2)
	})

	t.Run("check command - allowed", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "whitelist", "check", "git status"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "DECISION: ALLOW")
		assert.Contains(t, buf.String(), "rule-1")
	})

	t.Run("check command - ask", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "whitelist", "check", "curl https://evil.com"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "DECISION: ASK")
	})

	t.Run("remove rule", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "whitelist", "remove", "rule-1"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Successfully removed rule 'rule-1'")

		loaded, err := policy.LoadFromFile(configFile)
		require.NoError(t, err)
		require.Len(t, loaded.Rules, 1)
		assert.Equal(t, "rule-2", loaded.Rules[0].ID)
	})
}
