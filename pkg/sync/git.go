package sync

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RepoStatus represents the sync repository state
type RepoStatus struct {
	RepoPath       string `json:"repoPath"`
	RemoteURL      string `json:"remoteUrl,omitempty"`
	Branch         string `json:"branch"`
	Clean          bool   `json:"clean"`
	HasUncommitted bool   `json:"hasUncommitted"`
	HasUntracked   bool   `json:"hasUntracked"`
}

// GitBackend wraps git CLI operations targeting the sync store
type GitBackend struct {
	repoPath string
}

// NewGitBackend returns a new GitBackend instance
func NewGitBackend(repoPath string) *GitBackend {
	return &GitBackend{repoPath: repoPath}
}

// DefaultRepoPath returns ~/.antimatter/repo
func DefaultRepoPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".antimatter", "repo")
	}
	return filepath.Join(home, ".antimatter", "repo")
}

func (g *GitBackend) runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.repoPath
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	outStr := strings.TrimSpace(stdout.String())
	errStr := strings.TrimSpace(stderr.String())
	if err != nil {
		combined := errStr
		if outStr != "" {
			if combined != "" {
				combined += ": " + outStr
			} else {
				combined = outStr
			}
		}
		return outStr, fmt.Errorf("git %s failed: %w: %s", strings.Join(args, " "), err, combined)
	}
	return outStr, nil
}

// Init creates and initializes the git repository if not present
func (g *GitBackend) Init() error {
	if err := os.MkdirAll(g.repoPath, 0755); err != nil {
		return fmt.Errorf("failed to create repo directory '%s': %w", g.repoPath, err)
	}

	gitDir := filepath.Join(g.repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		if _, err := g.runGit("init", "-b", "main"); err != nil {
			// Fallback for older git versions
			if _, fErr := g.runGit("init"); fErr != nil {
				return fErr
			}
		}

		// Configure safe local fallback user if needed
		_, _ = g.runGit("config", "user.name", "Antimatter")
		_, _ = g.runGit("config", "user.email", "antimatter@twoboots.local")
	}
	return nil
}

// SetRemote configures the origin remote URL
func (g *GitBackend) SetRemote(url string) error {
	remotes, _ := g.runGit("remote")
	hasOrigin := false
	for _, r := range strings.Split(remotes, "\n") {
		if strings.TrimSpace(r) == "origin" {
			hasOrigin = true
			break
		}
	}

	if hasOrigin {
		_, err := g.runGit("remote", "set-url", "origin", url)
		return err
	}
	_, err := g.runGit("remote", "add", "origin", url)
	return err
}

// GetRemote returns the current origin URL
func (g *GitBackend) GetRemote() (string, error) {
	url, err := g.runGit("remote", "get-url", "origin")
	if err != nil {
		return "", nil // No remote configured
	}
	return url, nil
}

// Status inspects the repository status
func (g *GitBackend) Status() (*RepoStatus, error) {
	out, err := g.runGit("status", "--porcelain")
	if err != nil {
		return nil, err
	}

	remote, _ := g.GetRemote()
	branch, _ := g.runGit("branch", "--show-current")
	if branch == "" {
		branch = "main"
	}

	status := &RepoStatus{
		RepoPath:  g.repoPath,
		RemoteURL: remote,
		Branch:    branch,
		Clean:     true,
	}

	if strings.TrimSpace(out) != "" {
		status.Clean = false
		lines := strings.Split(strings.TrimSpace(out), "\n")
		for _, l := range lines {
			if strings.HasPrefix(l, "??") {
				status.HasUntracked = true
			} else {
				status.HasUncommitted = true
			}
		}
	}

	return status, nil
}

// StageAndCommit stages all changes and creates a commit
func (g *GitBackend) StageAndCommit(msg string) (string, error) {
	if _, err := g.runGit("add", "-A"); err != nil {
		return "", err
	}

	out, err := g.runGit("commit", "-m", msg)
	if err != nil {
		errStr := strings.ToLower(err.Error() + " " + out)
		if strings.Contains(errStr, "nothing to commit") || strings.Contains(errStr, "clean") {
			return "", nil
		}
		return "", err
	}

	sha, err := g.runGit("rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	return sha, nil
}

// Diff returns unstaged or uncommitted diffs
func (g *GitBackend) Diff() (string, error) {
	out, err := g.runGit("diff", "HEAD")
	if err != nil {
		// If HEAD doesn't exist yet, return empty diff
		return "", nil
	}
	return out, nil
}

// Pull pulls changes from origin
func (g *GitBackend) Pull() error {
	remote, err := g.GetRemote()
	if err != nil || remote == "" {
		return fmt.Errorf("no remote origin configured")
	}

	branch, _ := g.runGit("branch", "--show-current")
	if branch == "" {
		branch = "main"
	}

	_, err = g.runGit("pull", "--rebase", "origin", branch)
	return err
}

// Push pushes local commits to origin
func (g *GitBackend) Push() error {
	remote, err := g.GetRemote()
	if err != nil || remote == "" {
		return fmt.Errorf("no remote origin configured")
	}

	branch, _ := g.runGit("branch", "--show-current")
	if branch == "" {
		branch = "main"
	}

	_, err = g.runGit("push", "-u", "origin", branch)
	return err
}
