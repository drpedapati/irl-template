package cmd

import (
	"fmt"
	"os"

	"github.com/drpedapati/irl-template/internal/tui"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "irl",
	Short: "Idempotent Research Loop",
	Long:  `irl - Idempotent Research Loop`,
	Run: func(cmd *cobra.Command, args []string) {
		// Launch TUI when no subcommand is provided
		if err := tui.Run(Version); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Add --version flag
	rootCmd.Flags().BoolP("version", "v", false, "Print version")
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Printf("irl version %s\n", Version)
			os.Exit(0)
		}
	}

	// Hide plumbing commands
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Custom help for root command only
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "irl" {
			printRootHelp()
		} else {
			defaultHelp(cmd, args)
		}
	})
}

func printRootHelp() {
	fmt.Printf("%s %s\n", theme.B(theme.Cmd("irl")), theme.Faint("- Idempotent Research Loop"))
	fmt.Println()
	fmt.Printf("%s\n", theme.B("Quick start:"))
	fmt.Printf("  %s\n", theme.Succ("irl init \"your research purpose\""))
	fmt.Println()
	fmt.Printf("%s\n", theme.Faint("Commands:"))
	fmt.Printf("  %s        Create a new IRL project\n", theme.Cmd("init"))
	fmt.Println()
	fmt.Printf("%s\n", theme.Faint("Info:"))
	fmt.Printf("  %s   List available templates\n", theme.Cmd("templates"))
	fmt.Printf("  %s      Check environment setup\n", theme.Cmd("doctor"))
	fmt.Println()
	fmt.Printf("%s\n", theme.Faint("Settings:"))
	fmt.Printf("  %s      View or set configuration\n", theme.Cmd("config"))
	fmt.Printf("  %s      Update templates from GitHub\n", theme.Cmd("update"))
	fmt.Println()
	fmt.Printf("%s for details\n", theme.Faint("irl <command> --help"))
}

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Print version",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("irl version %s\n", Version)
	},
}
