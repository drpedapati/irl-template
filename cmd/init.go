package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/naming"
	"github.com/drpedapati/irl-template/pkg/scaffold"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/spf13/cobra"
)

var (
	templateFlag string
	nameFlag     string
	dirFlag      string
)

var initCmd = &cobra.Command{
	Use:   "init [purpose]",
	Short: "Create a new IRL project",
	Long: `Create a new IRL project with automatic naming.

Examples:
  irl init                              # Interactive mode
  irl init "ERP analysis study"         # Auto-generates: 260129-erp-analysis-study
  irl init -n my-project                # Use exact name: my-project
  irl init -t meeting-abstract          # With specific template
  irl init -d ~/Research "APA poster"   # Create in specific directory`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&templateFlag, "template", "t", "", "Template to use")
	initCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "Exact project name (skip auto-naming)")
	initCmd.Flags().StringVarP(&dirFlag, "dir", "d", "", "Directory to create project in (overrides default)")
}

func runInit(cmd *cobra.Command, args []string) error {
	var projectName string
	var purpose string
	var baseDir string

	// Determine if interactive mode
	isInteractive := len(args) == 0 && nameFlag == ""

	// Get base directory
	if dirFlag != "" {
		// Flag overrides everything
		baseDir = expandPath(dirFlag)
	} else if isInteractive {
		// Interactive: check config or ask
		baseDir = getOrAskDefaultDirectory()
	} else {
		// Non-interactive: use config or current directory
		baseDir = config.GetDefaultDirectory()
		if baseDir == "" {
			baseDir, _ = os.Getwd()
		}
	}

	// Ensure base directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("cannot create directory %s: %w", baseDir, err)
	}

	// Determine project name
	if nameFlag != "" {
		projectName = nameFlag
	} else if len(args) > 0 {
		purpose = strings.Join(args, " ")
		projectName = naming.GenerateName(purpose)
	} else {
		// Interactive mode
		prompt := &survey.Input{
			Message: "What's this project for?",
			Help:    "Brief description (e.g., 'ERP correlation analysis')",
		}
		if err := survey.AskOne(prompt, &purpose); err != nil {
			return err
		}
		if purpose == "" {
			return fmt.Errorf("project purpose is required")
		}
		projectName = naming.GenerateName(purpose)

		// Confirm or allow edit
		confirmPrompt := &survey.Input{
			Message: "Project folder name:",
			Default: projectName,
		}
		if err := survey.AskOne(confirmPrompt, &projectName); err != nil {
			return err
		}
	}

	// Full project path
	projectPath := filepath.Join(baseDir, projectName)

	// Check if directory exists
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", projectPath)
	}

	// Select template (prompt unless -t flag provided)
	var selectedTemplate string
	if templateFlag != "" {
		selectedTemplate = templateFlag
	} else {
		// Always offer template selection if no -t flag
		templateList, err := templates.ListTemplates()
		if err != nil {
			fmt.Println("Warning: could not fetch templates, using basic")
		}

		if len(templateList) > 0 {
			options := []string{"None (empty plan)"}
			for _, t := range templateList {
				options = append(options, fmt.Sprintf("%s - %s", t.Name, t.Description))
			}

			prompt := &survey.Select{
				Message: "Select template:",
				Options: options,
			}

			var result string
			if err := survey.AskOne(prompt, &result); err != nil {
				return err
			}

			if result != "None (empty plan)" {
				for _, t := range templateList {
					if strings.HasPrefix(result, t.Name+" - ") {
						selectedTemplate = t.Name
						break
					}
				}
			}
		}
	}

	// Create project
	fmt.Printf("\nCreating: %s\n", projectPath)

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Println("  Setting up structure...")
	if err := scaffold.Create(projectPath); err != nil {
		return err
	}

	// Apply template
	if selectedTemplate != "" {
		fmt.Printf("  Applying template: %s\n", selectedTemplate)
		tmpl, err := templates.GetTemplate(selectedTemplate)
		if err != nil {
			fmt.Printf("  Warning: %v, using basic template\n", err)
			tmpl = templates.EmbeddedTemplates["basic"]
		}
		if err := scaffold.WritePlan(projectPath, tmpl.Content); err != nil {
			return err
		}
	} else {
		// Write minimal plan
		minimalPlan := "# IRL Plan\n\n[Edit this file to define your research plan]\n"
		if err := scaffold.WritePlan(projectPath, minimalPlan); err != nil {
			return err
		}
	}

	fmt.Println("  Initializing git...")
	if err := scaffold.GitInit(projectPath); err != nil {
		fmt.Printf("  Warning: git init failed: %v\n", err)
	}

	fmt.Printf("\n✓ Created: %s\n", projectPath)
	fmt.Printf("\nNext:\n")
	fmt.Printf("  cd %s\n", projectPath)
	fmt.Printf("  # Edit 01-plans/main-plan.md\n")

	return nil
}

func getOrAskDefaultDirectory() string {
	// Check if already configured
	defaultDir := config.GetDefaultDirectory()
	if defaultDir != "" {
		// Show current setting and offer to change
		fmt.Printf("Default directory: %s\n", defaultDir)
		var changeDir bool
		prompt := &survey.Confirm{
			Message: "Use this directory?",
			Default: true,
		}
		survey.AskOne(prompt, &changeDir)
		if changeDir {
			return defaultDir
		}
	}

	// Ask for directory - use Documents/irl_projects as default (cross-platform, POSIX-compliant)
	home, _ := os.UserHomeDir()
	suggestion := filepath.Join(home, "Documents", "irl_projects")
	if defaultDir != "" {
		suggestion = defaultDir
	}

	var newDir string
	prompt := &survey.Input{
		Message: "Where should IRL projects be created?",
		Default: suggestion,
		Help:    "This will be saved as your default directory",
	}
	if err := survey.AskOne(prompt, &newDir); err != nil {
		return suggestion
	}

	newDir = expandPath(newDir)

	// Save to config
	if err := config.SetDefaultDirectory(newDir); err != nil {
		fmt.Printf("Warning: could not save config: %v\n", err)
	} else {
		fmt.Printf("Saved default directory: %s\n", newDir)
	}

	return newDir
}

func expandPath(path string) string {
	path = strings.TrimSpace(path)
	home, _ := os.UserHomeDir()

	// Handle bare ~ or ~/path
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}

	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err == nil {
			return abs
		}
	}
	return path
}
