package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/style"
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
			fmt.Fprintf(os.Stderr, "%sError:%s %v\n", style.Red, style.Reset, err)
			os.Exit(1)
		}
		status := style.Green + "exists" + style.Reset
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			status = style.Yellow + "will be created" + style.Reset
		}
		fmt.Printf("%s%s%s Set default directory: %s (%s)\n",
			style.Green, style.Check, style.Reset, dir, status)
		return
	}

	// Show current config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", style.Red, style.Reset, err)
		os.Exit(1)
	}

	fmt.Printf("%sConfiguration%s\n", style.BoldCyan, style.Reset)
	fmt.Println()

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".irl", "config.json")
	fmt.Printf("  %sConfig file%s       %s\n", style.Dim, style.Reset, configPath)

	if cfg.DefaultDirectory != "" {
		status := style.Green + "exists" + style.Reset
		if _, err := os.Stat(cfg.DefaultDirectory); os.IsNotExist(err) {
			status = style.Yellow + "will be created" + style.Reset
		}
		fmt.Printf("  %sDefault directory%s %s (%s)\n",
			style.Dim, style.Reset, cfg.DefaultDirectory, status)
	} else {
		defaultDir := filepath.Join(home, "Documents", "irl_projects")
		fmt.Printf("  %sDefault directory%s %snot set%s %s(will use %s)%s\n",
			style.Dim, style.Reset, style.Yellow, style.Reset, style.Dim, defaultDir, style.Reset)
	}

	fmt.Println()
	fmt.Printf("%sTo change:%s irl config --dir %s~/path%s\n",
		style.Dim, style.Reset, style.Cyan, style.Reset)
}
