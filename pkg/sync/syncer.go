package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/twoBoots/antimatter/pkg/adapters"
	"github.com/twoBoots/antimatter/pkg/policy"
)

// Syncer orchestrates sync operations across the Git repo, policy engine, and agent adapters
type Syncer struct {
	git             *GitBackend
	manager         *adapters.Manager
	localPolicyPath string
}

// NewSyncer returns a new Syncer instance
func NewSyncer(git *GitBackend, mgr *adapters.Manager, localPolicyPath string) *Syncer {
	return &Syncer{
		git:             git,
		manager:         mgr,
		localPolicyPath: localPolicyPath,
	}
}

// RepoPolicyPath returns the path to policy.yaml inside the sync repository
func (s *Syncer) RepoPolicyPath() string {
	return filepath.Join(s.git.repoPath, "policy.yaml")
}

// PullAndReconcile reads from the sync repo, merges with local policy, saves both, and optionally exports to adapters
func (s *Syncer) PullAndReconcile(autoExport bool) (*policy.Policy, error) {
	if err := s.git.Init(); err != nil {
		return nil, fmt.Errorf("failed to init git repository: %w", err)
	}

	repoPolicyFile := s.RepoPolicyPath()
	var repoP *policy.Policy
	if _, err := os.Stat(repoPolicyFile); err == nil {
		p, err := policy.LoadFromFile(repoPolicyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load repo policy: %w", err)
		}
		repoP = p
	}

	var localP *policy.Policy
	if _, err := os.Stat(s.localPolicyPath); err == nil {
		p, err := policy.LoadFromFile(s.localPolicyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load local policy: %w", err)
		}
		localP = p
	} else {
		localP = policy.NewPolicy("default", "Default Policy")
	}

	merged, err := policy.MergePolicies(localP, repoP)
	if err != nil {
		return nil, fmt.Errorf("failed to merge policies: %w", err)
	}

	// Persist to local policy file
	if err := policy.SaveToFile(merged, s.localPolicyPath); err != nil {
		return nil, fmt.Errorf("failed to save local policy: %w", err)
	}

	// Persist to repo policy file
	if err := policy.SaveToFile(merged, repoPolicyFile); err != nil {
		return nil, fmt.Errorf("failed to save repo policy: %w", err)
	}

	// Auto export to adapters if requested
	if autoExport && s.manager != nil {
		s.manager.ExportAll(merged)
	}

	return merged, nil
}

// PushLocal sanitizes local policy, writes it to the sync repo, commits, and pushes upstream
func (s *Syncer) PushLocal(commitMsg string) (string, error) {
	if err := s.git.Init(); err != nil {
		return "", fmt.Errorf("failed to init git repository: %w", err)
	}

	localP, err := policy.LoadFromFile(s.localPolicyPath)
	if err != nil {
		if os.IsNotExist(err) {
			localP = policy.NewPolicy("default", "Default Policy")
		} else {
			return "", fmt.Errorf("failed to read local policy: %w", err)
		}
	}

	// Sanitize secrets before staging
	cleanP, _ := SanitizePolicy(localP)

	repoPolicyFile := s.RepoPolicyPath()
	if err := policy.SaveToFile(cleanP, repoPolicyFile); err != nil {
		return "", fmt.Errorf("failed to write sanitized policy to repo: %w", err)
	}

	if commitMsg == "" {
		commitMsg = "sync(policy): Update Antimatter whitelist policy"
	}

	sha, err := s.git.StageAndCommit(commitMsg)
	if err != nil {
		return "", fmt.Errorf("failed to commit policy to repo: %w", err)
	}

	// Push if remote origin exists
	remote, _ := s.git.GetRemote()
	if remote != "" {
		_ = s.git.Push() // non-fatal in local/offline test mode
	}

	return sha, nil
}

// FullSync orchestrates Pull -> Merge -> Export -> Sanitize -> Commit -> Push
func (s *Syncer) FullSync(autoExport bool, commitMsg string) (*policy.Policy, string, error) {
	// Attempt git pull from remote if configured
	remote, _ := s.git.GetRemote()
	if remote != "" {
		_ = s.git.Pull()
	}

	merged, err := s.PullAndReconcile(autoExport)
	if err != nil {
		return nil, "", err
	}

	sha, err := s.PushLocal(commitMsg)
	if err != nil {
		return nil, "", err
	}

	return merged, sha, nil
}
