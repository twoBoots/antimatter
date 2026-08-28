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

func TestStatusCmd_Execution(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "policy.yaml")

	p := policy.NewPolicy("status-test-profile", "Status Test")
	p.Rules = []policy.PolicyRule{
		{ID: "r1", Type: policy.RuleTypeCommand, Pattern: "ls", MatchMode: policy.MatchExact},
	}
	require.NoError(t, policy.SaveToFile(p, configFile))

	t.Run("plain status output", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "status"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Antimatter Policy Status")
		assert.Contains(t, buf.String(), "Configured Agent Adapters:")
		assert.Contains(t, buf.String(), "Google Antigravity")
		assert.Contains(t, buf.String(), "Claude Code")
	})

	t.Run("json status output", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", configFile, "status", "--json"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), `"policy"`)
		assert.Contains(t, buf.String(), `"adapters"`)
	})
}
