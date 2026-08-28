package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twoBoots/antimatter/pkg/adapters"
	"github.com/twoBoots/antimatter/pkg/adapters/antigravity"
	"github.com/twoBoots/antimatter/pkg/adapters/claude"
	"github.com/twoBoots/antimatter/pkg/policy"
)

var defaultManager *adapters.Manager

// GetDefaultManager returns the shared adapter manager (lazily initializing default adapters if nil)
func GetDefaultManager() *adapters.Manager {
	if defaultManager == nil {
		mgr := adapters.NewManager()
		mgr.Register(antigravity.NewAdapter())
		mgr.Register(claude.NewAdapter())
		defaultManager = mgr
	}
	return defaultManager
}

// SetDefaultManager overrides the default manager (used for testing with custom dirs)
func SetDefaultManager(mgr *adapters.Manager) {
	defaultManager = mgr
}

// NewExportCmd returns the export command
func NewExportCmd() *cobra.Command {
	var target string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export canonical whitelist policies to active AI agent configurations",
		Long:  "Translates and exports canonical policies to Google Antigravity, Claude Code, and other configured assistants.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			p, err := loadOrCreatePolicy(path)
			if err != nil {
				return fmt.Errorf("failed to load policy: %w", err)
			}

			mgr := GetDefaultManager()
			target = strings.ToLower(strings.TrimSpace(target))

			if target == "all" || target == "" {
				errs := mgr.ExportAll(p)
				for _, a := range mgr.List() {
					if err, failed := errs[a.Name()]; failed {
						cmd.Printf("✖ Failed to export to %s: %v\n", a.DisplayName(), err)
					} else {
						cmd.Printf("✔ Exported policy to %s (%s)\n", a.DisplayName(), a.ConfigPath())
					}
				}
				if len(errs) > 0 {
					return fmt.Errorf("export completed with %d failure(s)", len(errs))
				}
				return nil
			}

			adapter, ok := mgr.Get(target)
			if !ok {
				return fmt.Errorf("unknown target agent '%s'", target)
			}

			if err := adapter.Export(p); err != nil {
				return fmt.Errorf("export to %s failed: %w", adapter.DisplayName(), err)
			}

			cmd.Printf("✔ Exported policy to %s (%s)\n", adapter.DisplayName(), adapter.ConfigPath())
			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "all", "target agent to export to (antigravity, claude, all)")
	return cmd
}

// NewImportCmd returns the import command
func NewImportCmd() *cobra.Command {
	var (
		from    string
		replace bool
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import whitelist rules from an AI agent configuration into canonical policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if from == "" {
				return fmt.Errorf("must specify --from <agent> (antigravity, claude)")
			}

			mgr := GetDefaultManager()
			from = strings.ToLower(strings.TrimSpace(from))

			importedPolicy, err := mgr.ImportFrom(from)
			if err != nil {
				return fmt.Errorf("import from '%s' failed: %w", from, err)
			}

			path := getPolicyPath()
			var finalPolicy *policy.Policy

			if replace {
				finalPolicy = importedPolicy
				finalPolicy.Name = "imported-" + from
			} else {
				basePolicy, err := loadOrCreatePolicy(path)
				if err != nil {
					return fmt.Errorf("failed to load local policy: %w", err)
				}
				finalPolicy, err = policy.MergePolicies(basePolicy, importedPolicy)
				if err != nil {
					return fmt.Errorf("failed to merge imported policy: %w", err)
				}
			}

			if err := policy.SaveToFile(finalPolicy, path); err != nil {
				return fmt.Errorf("failed to save imported policy: %w", err)
			}

			cmd.Printf("✔ Successfully imported rules from %s into %s (total: %d rules)\n", from, path, len(finalPolicy.Rules))
			return nil
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "source agent to import from (antigravity, claude)")
	cmd.Flags().BoolVar(&replace, "replace", false, "replace current policy entirely instead of merging")

	return cmd
}
