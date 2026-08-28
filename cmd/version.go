package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/twoBoots/antimatter/pkg/version"
)

// NewVersionCmd returns the version subcommand
func NewVersionCmd() *cobra.Command {
	var jsonOutput bool

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version, commit, and build information",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := version.Get()
			if jsonOutput {
				data, err := json.MarshalIndent(info, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(data))
				return nil
			}
			cmd.Println(info.String())
			return nil
		},
	}

	versionCmd.Flags().BoolVar(&jsonOutput, "json", false, "output version metadata in JSON format")

	return versionCmd
}
