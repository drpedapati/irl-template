package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/drpedapati/irl-template/pkg/config"
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
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		exists := "exists"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			exists = "will be created"
		}
		fmt.Printf("Default directory: %s (%s)\n", dir, exists)
		return
	}

	// Show current config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("IRL Configuration")
	fmt.Println("─────────────────")

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".irl", "config.json")
	fmt.Printf("Config file:       %s\n", configPath)

	if cfg.DefaultDirectory != "" {
		exists := "exists"
		if _, err := os.Stat(cfg.DefaultDirectory); os.IsNotExist(err) {
			exists = "will be created"
		}
		fmt.Printf("Default directory: %s (%s)\n", cfg.DefaultDirectory, exists)
	} else {
		defaultDir := filepath.Join(home, "Documents", "irl_projects")
		fmt.Printf("Default directory: (not set, will prompt or use %s)\n", defaultDir)
	}

	fmt.Println("\nTo change: irl config --dir ~/path/to/directory")
}
