package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "irl",
	Short: "Idempotent Research Loop",
	Long:  `irl - Idempotent Research Loop`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Hide plumbing commands
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Custom help for root command only
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "irl" {
			fmt.Print(`irl - Idempotent Research Loop

Quick start:
  irl init "your research purpose"

Commands:
  init        Create a new IRL project

Info:
  templates   List available templates
  doctor      Check environment setup

Settings:
  config      View or set configuration
  update      Update templates from GitHub

Use "irl <command> --help" for details
Use "irl version" for version info
`)
		} else {
			defaultHelp(cmd, args)
		}
	})
}

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Print version",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("irl version %s\n", Version)
	},
}
