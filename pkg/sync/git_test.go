package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	syncpkg "github.com/twoBoots/antimatter/pkg/sync"
)

func TestGitBackend_Lifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "sync_repo")

	backend := syncpkg.NewGitBackend(repoPath)

	// Step 1: Init repository
	err := backend.Init()
	require.NoError(t, err)
	assert.DirExists(t, filepath.Join(repoPath, ".git"))

	// Step 2: Set and Get Remote
	testRemote := "git@github.com:twoBoots/my-antimatter-store.git"
	err = backend.SetRemote(testRemote)
	require.NoError(t, err)

	remoteURL, err := backend.GetRemote()
	require.NoError(t, err)
	assert.Equal(t, testRemote, remoteURL)

	// Step 3: Write file and check status/diff
	policyFile := filepath.Join(repoPath, "policy.yaml")
	require.NoError(t, os.WriteFile(policyFile, []byte("version: antimatter/v1\n"), 0644))

	status, err := backend.Status()
	require.NoError(t, err)
	assert.False(t, status.Clean)
	assert.True(t, status.HasUntracked || status.HasUncommitted)

	// Step 4: Stage & Commit
	commitSHA, err := backend.StageAndCommit("Initial policy commit")
	require.NoError(t, err)
	assert.NotEmpty(t, commitSHA)

	statusAfter, err := backend.Status()
	require.NoError(t, err)
	assert.True(t, statusAfter.Clean)

	// Step 5: Modify file and check Diff
	require.NoError(t, os.WriteFile(policyFile, []byte("version: antimatter/v1\nname: my-profile\n"), 0644))
	diff, err := backend.Diff()
	require.NoError(t, err)
	assert.Contains(t, diff, "+name: my-profile")

	// Step 6: Test nothing to commit
	emptySHA, err := backend.StageAndCommit("No-op commit")
	require.NoError(t, err)
	assert.NotEmpty(t, emptySHA) // committed the modification
	noOpSHA, err := backend.StageAndCommit("Empty commit")
	require.NoError(t, err)
	assert.Empty(t, noOpSHA) // nothing to commit

	// Step 7: Push & Pull error handling when remote cannot be reached
	err = backend.Push()
	assert.Error(t, err) // cannot push to dummy remote

	err = backend.Pull()
	assert.Error(t, err) // cannot pull from dummy remote
}

func TestGitBackend_NoRemoteErrors(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "empty_repo")
	backend := syncpkg.NewGitBackend(repoPath)
	require.NoError(t, backend.Init())

	// Push without remote configured
	err := backend.Push()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no remote origin configured")

	// Pull without remote configured
	err = backend.Pull()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no remote origin configured")

	// DefaultRepoPath helper
	defPath := syncpkg.DefaultRepoPath()
	assert.NotEmpty(t, defPath)
	assert.Contains(t, defPath, ".antimatter")
}
