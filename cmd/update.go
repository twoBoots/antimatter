package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/twoBoots/antimatter/pkg/version"
	"github.com/twoBoots/bender/pkg/updater"
)

const (
	defaultRepo       = "twoBoots/antimatter"
	defaultBinaryName = "antimatter"
)

// NewUpdateCmd returns the self-updater subcommand
func NewUpdateCmd() *cobra.Command {
	var (
		checkOnly     bool
		targetVersion string
		force         bool
	)

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update Antimatter in-place to latest or specified version",
		Long:  "Queries GitHub Releases API, compares semantic versions, and atomically replaces the local binary in-place.",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := updater.SelfUpdate(updater.Options{
				Repo:           defaultRepo,
				BinaryName:     defaultBinaryName,
				CurrentVersion: version.Version,
				TargetVersion:  targetVersion,
				CheckOnly:      checkOnly,
				Force:          force,
			})
			if err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			cmd.Println(res.Message)
			return nil
		},
	}

	updateCmd.Flags().BoolVar(&checkOnly, "check", false, "check for updates without applying")
	updateCmd.Flags().StringVar(&targetVersion, "target-version", "", "specific release tag to install (e.g. v1.2.0)")
	updateCmd.Flags().BoolVar(&force, "force", false, "force reinstall even if already up to date")

	return updateCmd
}
