package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/twoBoots/antimatter/pkg/adapters"
)

// StatusOutput models structured JSON output for the status command
type StatusOutput struct {
	Policy   PolicyStatusSummary      `json:"policy"`
	Adapters []adapters.AdapterStatus `json:"adapters"`
}

// PolicyStatusSummary provides summary stats of the active policy
type PolicyStatusSummary struct {
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	RuleCount   int       `json:"ruleCount"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// NewStatusCmd returns the status command
func NewStatusCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display sync and configuration status across agent adapters",
		Long:  "Inspects local policy status, rule counts, and configuration state for Antigravity and Claude.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getPolicyPath()
			p, err := loadOrCreatePolicy(path)
			if err != nil {
				return fmt.Errorf("failed to load policy: %w", err)
			}

			mgr := GetDefaultManager()
			adapterStatuses := mgr.GetStatuses()

			summary := PolicyStatusSummary{
				Path:        path,
				Name:        p.Name,
				RuleCount:   len(p.Rules),
				LastUpdated: p.UpdatedAt,
			}

			if jsonOutput {
				out := StatusOutput{
					Policy:   summary,
					Adapters: adapterStatuses,
				}
				data, err := json.MarshalIndent(out, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(data))
				return nil
			}

			cmd.Println("==================================================")
			cmd.Println("          Antimatter Policy Status                ")
			cmd.Println("==================================================")
			cmd.Printf("Active Policy File: %s\n", path)
			cmd.Printf("Profile Name:       %s\n", p.Name)
			cmd.Printf("Total Rules:        %d\n", len(p.Rules))
			cmd.Println()

			cmd.Println("Configured Agent Adapters:")
			cmd.Println(strings.Repeat("-", 60))
			for _, a := range adapterStatuses {
				statusLabel := "[Not Found]"
				if a.Detected {
					statusLabel = "[Active/Detected]"
				}
				cmd.Printf("• %-20s %s\n", a.DisplayName, statusLabel)
				cmd.Printf("  Config Path: %s\n", a.ConfigPath)
				if a.LastBackup != "" {
					cmd.Printf("  Last Backup: %s\n", a.LastBackup)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	return cmd
}
