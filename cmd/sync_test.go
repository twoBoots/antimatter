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
	syncpkg "github.com/twoBoots/antimatter/pkg/sync"
)

func TestSyncCommands_Lifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "sync_repo")
	localPolicyFile := filepath.Join(tmpDir, "policy.yaml")
	agDir := filepath.Join(tmpDir, "antigravity-cli")
	claudeFile := filepath.Join(tmpDir, ".claude.json")
	backupDir := filepath.Join(tmpDir, "backups")

	// Set custom backend and manager for CLI tests
	git := syncpkg.NewGitBackend(repoPath)
	mgr := adapters.NewManager()
	mgr.Register(antigravity.NewAdapter(
		antigravity.WithConfigDir(agDir),
		antigravity.WithBackupDir(backupDir),
	))
	mgr.Register(claude.NewAdapter(
		claude.WithConfigPath(claudeFile),
		claude.WithBackupDir(backupDir),
	))

	syncer := syncpkg.NewSyncer(git, mgr, localPolicyFile)
	cmd.SetDefaultSyncer(syncer)
	cmd.SetDefaultGitBackend(git)

	// Save initial policy
	p := policy.NewPolicy("test-sync-cli", "Sync CLI Test")
	p.Rules = []policy.PolicyRule{
		{
			ID:        "rule-1",
			Type:      policy.RuleTypeCommand,
			Pattern:   "git status",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
		},
	}
	require.NoError(t, policy.SaveToFile(p, localPolicyFile))

	t.Run("remote set and show", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"remote", "set", "git@github.com:twoBoots/my-repo.git"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Configured remote origin")

		// remote show
		buf.Reset()
		rootCmd.SetArgs([]string{"remote", "show"})
		err = rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "git@github.com:twoBoots/my-repo.git")
	})

	t.Run("push command", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", localPolicyFile, "push", "--message", "CLI push"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Successfully committed and synced policy to local store")
	})

	t.Run("diff command", func(t *testing.T) {
		// Add new rule to local policy
		p.Rules = append(p.Rules, policy.PolicyRule{
			ID:        "rule-2",
			Type:      policy.RuleTypeCommand,
			Pattern:   "go test",
			MatchMode: policy.MatchPrefix,
		})
		require.NoError(t, policy.SaveToFile(p, localPolicyFile))

		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", localPolicyFile, "diff"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Local uncommitted policy changes")
		assert.Contains(t, buf.String(), "rule-2")
	})

	t.Run("pull command", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", localPolicyFile, "pull"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Successfully pulled and reconciled policy")
	})

	t.Run("sync command", func(t *testing.T) {
		rootCmd := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"--config", localPolicyFile, "sync", "--message", "End to end sync"})

		err := rootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Antimatter Full Sync Complete")
	})
}
