package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twoBoots/antimatter/pkg/version"
)

func TestGetInfo_Defaults(t *testing.T) {
	info := version.Get()
	assert.NotEmpty(t, info.Version)
	assert.NotEmpty(t, info.String())
}

func TestGetInfo_FormattedString(t *testing.T) {
	version.Version = "v1.2.3"
	version.Commit = "abcdef1"
	version.Date = "2026-08-28T00:00:00Z"

	info := version.Get()
	assert.Equal(t, "v1.2.3", info.Version)
	assert.Equal(t, "abcdef1", info.Commit)
	assert.Equal(t, "2026-08-28T00:00:00Z", info.Date)
	assert.Contains(t, info.String(), "antimatter v1.2.3 (commit: abcdef1, built at: 2026-08-28T00:00:00Z)")
}
