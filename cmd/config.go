package cmd

import (
	"encoding/json"
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
  irl config                        # Show current config
  irl config --json                 # JSON output
  irl config --dir ~/Research       # Set default directory
  irl config --editor cursor        # Set preferred editor`,
	RunE: runConfig,
}

var (
	configDirFlag    string
	configEditorFlag string
	configJSONFlag   bool
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVar(&configDirFlag, "dir", "", "Set default directory for new projects")
	configCmd.Flags().StringVar(&configEditorFlag, "editor", "", "Set preferred editor (e.g., cursor, code, vim)")
	configCmd.Flags().BoolVar(&configJSONFlag, "json", false, "Output as JSON")
}

func runConfig(cmd *cobra.Command, args []string) error {
	changed := false

	// Set directory
	if configDirFlag != "" {
		dir := expandPath(configDirFlag)
		if err := config.SetDefaultDirectory(dir); err != nil {
			return fmt.Errorf("failed to set directory: %w", err)
		}
		status := theme.StatusTag("exists", true)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			status = theme.StatusTag("will be created", false)
		}
		fmt.Printf("%s Set default directory: %s (%s)\n",
			theme.OK(""), dir, status)
		changed = true
	}

	// Set editor
	if configEditorFlag != "" {
		editorType := "gui"
		switch configEditorFlag {
		case "vim", "nvim", "nano", "helix", "emacs":
			editorType = "terminal"
		}
		if err := config.SetPlanEditor(configEditorFlag, editorType); err != nil {
			return fmt.Errorf("failed to set editor: %w", err)
		}
		fmt.Printf("%s Set editor: %s (%s)\n",
			theme.OK(""), theme.Cmd(configEditorFlag), editorType)
		changed = true
	}

	if changed && !configJSONFlag {
		return nil
	}

	// Show current config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if configJSONFlag {
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
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

	if cfg.PlanEditor != "" {
		fmt.Println(theme.KeyValue("Editor          ", cfg.PlanEditor+" ("+cfg.PlanEditorType+")"))
	} else {
		fmt.Printf("%s\n", theme.KeyValue("Editor          ", theme.Faint("auto-detect")))
	}

	if config.HasProfile() {
		p := cfg.Profile
		label := p.Name
		if p.Institution != "" {
			label += ", " + p.Institution
		}
		fmt.Println(theme.KeyValue("Profile         ", label))
	}

	fmt.Println()
	fmt.Printf("%s irl config --dir %s\n",
		theme.Faint("To change:"),
		theme.Cmd("~/path"))

	return nil
}
