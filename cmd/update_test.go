package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/bender/pkg/updater"
)

func TestNewUpdateCmd_Flags(t *testing.T) {
	cmd := NewUpdateCmd()
	require.NotNil(t, cmd)

	assert.Equal(t, "update", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	checkFlag := cmd.Flags().Lookup("check")
	require.NotNil(t, checkFlag)
	assert.Equal(t, "false", checkFlag.DefValue)

	targetFlag := cmd.Flags().Lookup("target-version")
	require.NotNil(t, targetFlag)
	assert.Equal(t, "", targetFlag.DefValue)

	forceFlag := cmd.Flags().Lookup("force")
	require.NotNil(t, forceFlag)
	assert.Equal(t, "false", forceFlag.DefValue)

	// Verify help output contains flags
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "--check")
	assert.Contains(t, buf.String(), "--target-version")
	assert.Contains(t, buf.String(), "--force")
}

func TestUpdater_SemVerStandardCompliance(t *testing.T) {
	// Tests SemVer 2.0.0 compliance provided by Bender v1.1.0 (golang.org/x/mod/semver)
	// In v1.0.1, beta.2 vs beta.11 failed due to lexicographical string sorting of parts.
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.0.0-beta.2", "1.0.0-beta.11", -1},
		{"1.0.0-beta.11", "1.0.0-beta.2", 1},
		{"1.0.0-rc.1", "1.0.0", -1},
		{"1.0.0+20130313", "1.0.0+20120313", 0},
		{"0.1", "0.1.0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			res := updater.CompareVersions(tt.v1, tt.v2)
			assert.Equal(t, tt.expected, res, "comparing %s and %s", tt.v1, tt.v2)
		})
	}
}
