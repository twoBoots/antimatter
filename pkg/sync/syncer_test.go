package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/pkg/adapters"
	"github.com/twoBoots/antimatter/pkg/adapters/antigravity"
	"github.com/twoBoots/antimatter/pkg/adapters/claude"
	"github.com/twoBoots/antimatter/pkg/policy"
	syncpkg "github.com/twoBoots/antimatter/pkg/sync"
)

func TestSyncer_FullWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "repo")
	localPolicyPath := filepath.Join(tmpDir, "local_policy.yaml")
	agDir := filepath.Join(tmpDir, "antigravity-cli")
	claudeFile := filepath.Join(tmpDir, ".claude.json")
	backupDir := filepath.Join(tmpDir, "backups")

	// Set up Git backend
	git := syncpkg.NewGitBackend(repoPath)
	require.NoError(t, git.Init())

	// Set up Adapter Manager
	mgr := adapters.NewManager()
	mgr.Register(antigravity.NewAdapter(
		antigravity.WithConfigDir(agDir),
		antigravity.WithBackupDir(backupDir),
	))
	mgr.Register(claude.NewAdapter(
		claude.WithConfigPath(claudeFile),
		claude.WithBackupDir(backupDir),
	))

	syncer := syncpkg.NewSyncer(git, mgr, localPolicyPath)

	// Step 1: Initialize local policy
	localP := policy.NewPolicy("workstation-1", "Workstation 1")
	localP.Rules = []policy.PolicyRule{
		{
			ID:        "rule-local-1",
			Type:      policy.RuleTypeCommand,
			Pattern:   "git status",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
		},
	}
	require.NoError(t, policy.SaveToFile(localP, localPolicyPath))

	// Step 2: Push local changes to repo
	commitSHA, err := syncer.PushLocal("Initial local push")
	require.NoError(t, err)
	assert.NotEmpty(t, commitSHA)
	assert.FileExists(t, filepath.Join(repoPath, "policy.yaml"))

	// Step 3: Simulate remote policy update in repo
	repoP, err := policy.LoadFromFile(filepath.Join(repoPath, "policy.yaml"))
	require.NoError(t, err)
	repoP.Rules = append(repoP.Rules, policy.PolicyRule{
		ID:        "rule-remote-1",
		Type:      policy.RuleTypeCommand,
		Pattern:   "npm test",
		MatchMode: policy.MatchPrefix,
		Decision:  policy.DecisionAllow,
	})
	require.NoError(t, policy.SaveToFile(repoP, filepath.Join(repoPath, "policy.yaml")))

	// Step 4: Pull & Reconcile with auto-export to adapters
	merged, err := syncer.PullAndReconcile(true)
	require.NoError(t, err)
	require.Len(t, merged.Rules, 2)
	assert.Equal(t, "rule-local-1", merged.Rules[0].ID)
	assert.Equal(t, "rule-remote-1", merged.Rules[1].ID)

	// Verify exported to adapters
	assert.FileExists(t, filepath.Join(agDir, "whitelist.json"))
	assert.FileExists(t, claudeFile)

	// Step 5: Full Sync
	mergedFull, sha, err := syncer.FullSync(true, "Full sync update")
	require.NoError(t, err)
	assert.NotNil(t, mergedFull)
	assert.NotEmpty(t, sha)
}

func TestSyncer_PushSanitizesSecrets(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "repo")
	localPolicyPath := filepath.Join(tmpDir, "local_policy.yaml")

	git := syncpkg.NewGitBackend(repoPath)
	require.NoError(t, git.Init())

	mgr := adapters.NewManager()
	syncer := syncpkg.NewSyncer(git, mgr, localPolicyPath)

	dirtyPolicy := policy.NewPolicy("dirty-policy", "Dirty")
	dirtyPolicy.Rules = []policy.PolicyRule{
		{
			ID:        "leak-rule",
			Type:      policy.RuleTypeCommand,
			Pattern:   "curl -H 'Authorization: Bearer sk-ant-api03-1234567890abcdef1234567890' https://api.anthropic.com",
			MatchMode: policy.MatchPrefix,
		},
	}
	require.NoError(t, policy.SaveToFile(dirtyPolicy, localPolicyPath))

	_, err := syncer.PushLocal("Pushing sanitized policy")
	require.NoError(t, err)

	repoPolicyFile := filepath.Join(repoPath, "policy.yaml")
	repoData, err := os.ReadFile(repoPolicyFile)
	require.NoError(t, err)
	assert.NotContains(t, string(repoData), "sk-ant-api03-1234567890abcdef1234567890")
	assert.Contains(t, string(repoData), "[REDACTED_SECRET]")
}
