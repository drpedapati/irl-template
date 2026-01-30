package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or set configuration",
	Long: `View or set IRL configuration.

Examples:
  irl config                    # Show current config
  irl config --dir ~/Research   # Set default directory`,
	Run: runConfig,
}

var configDirFlag string

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVar(&configDirFlag, "dir", "", "Set default directory for new projects")
}

func runConfig(cmd *cobra.Command, args []string) {
	if configDirFlag != "" {
		// Set new directory
		dir := expandPath(configDirFlag)
		if err := config.SetDefaultDirectory(dir); err != nil {
			fmt.Fprintf(os.Stderr, "%s %v\n", theme.Err("Error:"), err)
			os.Exit(1)
		}
		status := theme.StatusTag("exists", true)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			status = theme.StatusTag("will be created", false)
		}
		fmt.Printf("%s Set default directory: %s (%s)\n",
			theme.OK(""), dir, status)
		return
	}

	// Show current config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", theme.Err("Error:"), err)
		os.Exit(1)
	}

	theme.Section("Configuration")
	fmt.Println()

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".irl", "config.json")
	fmt.Println(theme.KeyValue("Config file      ", configPath))

	if cfg.DefaultDirectory != "" {
		status := theme.StatusTag("exists", true)
		if _, err := os.Stat(cfg.DefaultDirectory); os.IsNotExist(err) {
			status = theme.StatusTag("will be created", false)
		}
		fmt.Printf("%s (%s)\n",
			theme.KeyValue("Default directory", cfg.DefaultDirectory),
			status)
	} else {
		defaultDir := filepath.Join(home, "Documents", "irl_projects")
		fmt.Printf("%s %s\n",
			theme.KeyValue("Default directory", theme.Warn("not set")),
			theme.Faint("(will use "+defaultDir+")"))
	}

	fmt.Println()
	fmt.Printf("%s irl config --dir %s\n",
		theme.Faint("To change:"),
		theme.Cmd("~/path"))
}
