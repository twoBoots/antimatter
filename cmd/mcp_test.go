package cmd_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/cmd"
	"github.com/twoBoots/antimatter/pkg/policy"
)

func TestUpdateCmd_Help(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"update", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Queries GitHub Releases API")
	assert.Contains(t, buf.String(), "--check")
}

func TestMCPCmd_HelpAndInstall(t *testing.T) {
	t.Run("mcp help", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"mcp", "--help"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Runs a stdio JSON-RPC 2.0 MCP server")
	})

	t.Run("mcp install list clients", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"mcp", "install", "--list"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Supported AI Clients")
	})
}

func TestMCPServer_ToolsExecution(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "policy.yaml")

	p := policy.NewPolicy("test-mcp-policy", "MCP Test")
	p.Rules = []policy.PolicyRule{
		{
			ID:        "rule-git",
			Type:      policy.RuleTypeCommand,
			Pattern:   "git status",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
		},
	}
	require.NoError(t, policy.SaveToFile(p, configFile))

	srv := cmd.BuildMCPServer(configFile)
	require.NotNil(t, srv)

	// Test antimatter_get_whitelist tool handler
	getToolHandler, ok := cmd.GetToolHandler(srv, "antimatter_get_whitelist")
	require.True(t, ok)

	getRes, err := getToolHandler(context.Background(), nil)
	require.NoError(t, err)
	assert.False(t, getRes.IsError)
	assert.Contains(t, getRes.Content[0].Text, "rule-git")

	// Test antimatter_check_command tool handler directly
	checkToolHandler, ok := cmd.GetToolHandler(srv, "antimatter_check_command")
	require.True(t, ok)

	res, err := checkToolHandler(context.Background(), map[string]interface{}{
		"command": "git status",
	})
	require.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Contains(t, res.Content[0].Text, "ALLOW")

	// Test tool check
	resTool, err := checkToolHandler(context.Background(), map[string]interface{}{
		"tool": "view_file",
	})
	require.NoError(t, err)
	assert.False(t, resTool.IsError)
	assert.Contains(t, resTool.Content[0].Text, "DECISION:")

	// Test missing arguments
	resErr, err := checkToolHandler(context.Background(), map[string]interface{}{})
	require.NoError(t, err)
	assert.True(t, resErr.IsError)

	// Test antimatter_add_whitelist_rule tool handler
	addToolHandler, ok := cmd.GetToolHandler(srv, "antimatter_add_whitelist_rule")
	require.True(t, ok)

	addRes, err := addToolHandler(context.Background(), map[string]interface{}{
		"id":      "mcp-rule-1",
		"command": "go test ./...",
		"mode":    "prefix",
	})
	require.NoError(t, err)
	assert.False(t, addRes.IsError)
	assert.Contains(t, addRes.Content[0].Text, "Successfully added")

	// Test add missing pattern
	addResErr, err := addToolHandler(context.Background(), map[string]interface{}{})
	require.NoError(t, err)
	assert.True(t, addResErr.IsError)

	// Verify persistence
	loaded, err := policy.LoadFromFile(configFile)
	require.NoError(t, err)
	require.Len(t, loaded.Rules, 2)
}
