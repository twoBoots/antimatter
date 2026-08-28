package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

// NewRootCmd initializes and returns the base root cobra command
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "antimatter",
		Short: "antimatter - AI Agent Configuration & Whitelist Sync Engine",
		Long: `antimatter is a high-performance CLI utility and embedded MCP server 
designed to persist, version control, and synchronize agent configurations, 
whitelisted tool calls, and approved CLI commands across Google Antigravity 
and Claude Code.`,
		SilenceUsage: true,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.antimatter/policy.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose debug logging")

	// Register subcommands
	rootCmd.AddCommand(NewVersionCmd())

	return rootCmd
}

// Execute runs the root command
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
