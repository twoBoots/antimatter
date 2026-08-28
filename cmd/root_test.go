package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/cmd"
	"github.com/twoBoots/antimatter/pkg/version"
)

func TestRootCmd_Help(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "antimatter is a high-performance CLI utility")
	assert.Contains(t, output, "Usage:")
	assert.Contains(t, output, "version")
}

func TestVersionCmd_Output(t *testing.T) {
	version.Version = "v1.0.0"
	version.Commit = "1234567"
	version.Date = "2026-08-28T00:00:00Z"

	t.Run("plain text", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"version"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "antimatter v1.0.0")
	})

	t.Run("json output", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"version", "--json"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), `"version": "v1.0.0"`)
		assert.Contains(t, buf.String(), `"commit": "1234567"`)
	})
}
