package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twoBoots/antimatter/pkg/policy"
	syncpkg "github.com/twoBoots/antimatter/pkg/sync"
)

var (
	defaultSyncer     *syncpkg.Syncer
	defaultGitBackend *syncpkg.GitBackend
)

// GetDefaultGitBackend returns the git backend
func GetDefaultGitBackend() *syncpkg.GitBackend {
	if defaultGitBackend == nil {
		defaultGitBackend = syncpkg.NewGitBackend(syncpkg.DefaultRepoPath())
	}
	return defaultGitBackend
}

// SetDefaultGitBackend overrides git backend for testing
func SetDefaultGitBackend(git *syncpkg.GitBackend) {
	defaultGitBackend = git
}

// GetDefaultSyncer returns the syncer
func GetDefaultSyncer(policyPath string) *syncpkg.Syncer {
	if defaultSyncer == nil {
		git := GetDefaultGitBackend()
		mgr := GetDefaultManager()
		defaultSyncer = syncpkg.NewSyncer(git, mgr, policyPath)
	}
	return defaultSyncer
}

// SetDefaultSyncer overrides syncer for testing
func SetDefaultSyncer(s *syncpkg.Syncer) {
	defaultSyncer = s
}

// NewSyncCmd returns the sync command
func NewSyncCmd() *cobra.Command {
	var (
		autoExport bool
		message    string
	)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize local policy with remote Git repository and active agent adapters",
		Long:  "Executes a complete sync lifecycle: pulls remote changes, reconciles with local policy, updates agent adapters, and pushes local additions upstream.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			syncer := GetDefaultSyncer(path)

			merged, sha, err := syncer.FullSync(autoExport, message)
			if err != nil {
				return fmt.Errorf("full sync failed: %w", err)
			}

			cmd.Println("==================================================")
			cmd.Println("          Antimatter Full Sync Complete           ")
			cmd.Println("==================================================")
			cmd.Printf("✔ Policy File:      %s\n", path)
			cmd.Printf("✔ Total Rules:      %d\n", len(merged.Rules))
			if sha != "" {
				cmd.Printf("✔ Commit Hash:      %s\n", sha)
			}
			if autoExport {
				cmd.Println("✔ Agent Adapters:   Synchronized (Antigravity, Claude)")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&autoExport, "auto-export", true, "automatically export merged policy to active agent adapters")
	cmd.Flags().StringVarP(&message, "message", "m", "", "commit message for sync updates")

	return cmd
}

// NewPushCmd returns the push command
func NewPushCmd() *cobra.Command {
	var message string

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Sanitize, commit, and push local whitelist policy to the sync repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			syncer := GetDefaultSyncer(path)

			sha, err := syncer.PushLocal(message)
			if err != nil {
				return fmt.Errorf("push failed: %w", err)
			}

			if sha == "" {
				cmd.Println("✔ Sync store is already up to date (nothing to commit)")
			} else {
				cmd.Printf("✔ Successfully committed and synced policy to local store (commit: %s)\n", sha)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "commit message for local push")
	return cmd
}

// NewPullCmd returns the pull command
func NewPullCmd() *cobra.Command {
	var autoExport bool

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull latest policies from sync repository and reconcile with local settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			syncer := GetDefaultSyncer(path)

			merged, err := syncer.PullAndReconcile(autoExport)
			if err != nil {
				return fmt.Errorf("pull failed: %w", err)
			}

			cmd.Printf("✔ Successfully pulled and reconciled policy '%s' (total: %d rules)\n", merged.Name, len(merged.Rules))
			if autoExport {
				cmd.Println("✔ Automatically exported rules to active agent adapters")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&autoExport, "auto-export", true, "automatically export reconciled policy to agent adapters")
	return cmd
}

// NewDiffCmd returns the diff command
func NewDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Inspect differences between local policy and sync repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			git := GetDefaultGitBackend()

			// Check local uncommitted changes
			localP, err := loadOrCreatePolicy(path)
			if err != nil {
				return err
			}

			repoPolicyFile := syncpkg.NewSyncer(git, nil, path).RepoPolicyPath()
			var repoP *policy.Policy
			if _, err := os.Stat(repoPolicyFile); err == nil {
				repoP, _ = policy.LoadFromFile(repoPolicyFile)
			}

			cmd.Println("Local uncommitted policy changes:")
			cmd.Println(strings.Repeat("-", 60))

			if repoP == nil {
				cmd.Printf("• Initializing new policy store with %d rules:\n", len(localP.Rules))
				for _, r := range localP.Rules {
					cmd.Printf("  + [%s] %s (%s: %s)\n", r.Type, r.ID, r.MatchMode, r.Pattern)
				}
				return nil
			}

			// Compute rule diff
			repoRules := make(map[string]policy.PolicyRule)
			for _, r := range repoP.Rules {
				repoRules[r.ID] = r
			}

			diffCount := 0
			for _, r := range localP.Rules {
				oldR, exists := repoRules[r.ID]
				if !exists {
					cmd.Printf("  + [%s] %s (%s: %s)\n", r.Type, r.ID, r.MatchMode, r.Pattern)
					diffCount++
				} else if oldR.Pattern != r.Pattern || oldR.MatchMode != r.MatchMode || oldR.Decision != r.Decision {
					cmd.Printf("  ~ [%s] %s\n      Old: %s (%s)\n      New: %s (%s)\n", r.Type, r.ID, oldR.MatchMode, oldR.Pattern, r.MatchMode, r.Pattern)
					diffCount++
				}
			}

			if diffCount == 0 {
				cmd.Println("✔ No local changes; policy is in sync with repository")
			}

			return nil
		},
	}
}

// NewRemoteCmd returns the remote configuration command
func NewRemoteCmd() *cobra.Command {
	remoteCmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage remote Git synchronization repository",
	}

	remoteCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Display configured remote origin URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			git := GetDefaultGitBackend()
			url, err := git.GetRemote()
			if err != nil {
				return err
			}
			if url == "" {
				cmd.Println("No remote repository configured. Set one using: antimatter remote set <url>")
				return nil
			}
			cmd.Printf("Remote origin: %s\n", url)
			return nil
		},
	})

	remoteCmd.AddCommand(&cobra.Command{
		Use:   "set [url]",
		Short: "Set or update remote origin repository URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			git := GetDefaultGitBackend()
			if err := git.Init(); err != nil {
				return err
			}
			url := args[0]
			if err := git.SetRemote(url); err != nil {
				return fmt.Errorf("failed to set remote: %w", err)
			}
			cmd.Printf("✔ Configured remote origin: %s\n", url)
			return nil
		},
	})

	return remoteCmd
}
