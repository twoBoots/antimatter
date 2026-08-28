package adapters_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/pkg/adapters"
	"github.com/twoBoots/antimatter/pkg/policy"
)

func TestBackupEngine_CreateAndList(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "settings.json")
	require.NoError(t, os.WriteFile(sourceFile, []byte(`{"key": "value"}`), 0644))

	backupBaseDir := filepath.Join(tmpDir, "backups")

	backupPath, err := adapters.CreateBackupWithBaseDir(sourceFile, "antigravity", backupBaseDir)
	require.NoError(t, err)
	assert.FileExists(t, backupPath)
	assert.Contains(t, backupPath, "settings.json")

	content, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, `{"key": "value"}`, string(content))

	backups, err := adapters.ListBackupsWithBaseDir("antigravity", backupBaseDir)
	require.NoError(t, err)
	require.Len(t, backups, 1)
	assert.Equal(t, backupPath, backups[0])
}

// MockAdapter for testing manager registry and dispatch
type MockAdapter struct {
	name        string
	displayName string
	detected    bool
	configPath  string
	exported    bool
	imported    bool
}

func (m *MockAdapter) Name() string        { return m.name }
func (m *MockAdapter) DisplayName() string { return m.displayName }
func (m *MockAdapter) Detected() bool      { return m.detected }
func (m *MockAdapter) ConfigPath() string  { return m.configPath }
func (m *MockAdapter) Backup() (string, error) {
	return "/tmp/mock_backup", nil
}
func (m *MockAdapter) Export(p *policy.Policy) error {
	m.exported = true
	return nil
}
func (m *MockAdapter) Import() (*policy.Policy, error) {
	m.imported = true
	p := policy.NewPolicy("imported-mock", "Mock Imported")
	p.Rules = []policy.PolicyRule{
		{ID: "mock-rule", Type: policy.RuleTypeCommand, Pattern: "echo 1", MatchMode: policy.MatchExact},
	}
	return p, nil
}

func TestManager_RegistryAndDispatch(t *testing.T) {
	mgr := adapters.NewManager()
	mockAG := &MockAdapter{name: "antigravity", displayName: "Google Antigravity", detected: true, configPath: "/path/ag"}
	mockClaude := &MockAdapter{name: "claude", displayName: "Claude Code", detected: false, configPath: "/path/claude"}

	mgr.Register(mockAG)
	mgr.Register(mockClaude)

	require.Len(t, mgr.List(), 2)

	ag, found := mgr.Get("antigravity")
	require.True(t, found)
	assert.Equal(t, "antigravity", ag.Name())

	_, found = mgr.Get("unknown")
	assert.False(t, found)

	// ExportAll
	p := policy.NewPolicy("test-p", "Test")
	errs := mgr.ExportAll(p)
	assert.Empty(t, errs)
	assert.True(t, mockAG.exported)
	assert.True(t, mockClaude.exported)

	// ImportFrom
	imported, err := mgr.ImportFrom("antigravity")
	require.NoError(t, err)
	assert.Equal(t, "imported-mock", imported.Name)
	assert.True(t, mockAG.imported)

	// ImportFrom unknown
	_, err = mgr.ImportFrom("unknown")
	assert.Error(t, err)

	// Statuses
	statuses := mgr.GetStatuses()
	require.Len(t, statuses, 2)
	assert.Equal(t, "antigravity", statuses[0].Name)
	assert.True(t, statuses[0].Detected)
}
