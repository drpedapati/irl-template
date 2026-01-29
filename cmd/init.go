package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/drpedapati/irl-template/pkg/naming"
	"github.com/drpedapati/irl-template/pkg/scaffold"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/spf13/cobra"
)

var (
	templateFlag string
	nameFlag     string
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
  irl init "APA poster" -t meeting-abstract`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&templateFlag, "template", "t", "", "Template to use")
	initCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "Exact project name (skip auto-naming)")
}

func runInit(cmd *cobra.Command, args []string) error {
	var projectName string
	var purpose string

	// Determine project name
	if nameFlag != "" {
		// Exact name provided
		projectName = nameFlag
	} else if len(args) > 0 {
		// Purpose provided, generate name
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

	// Check if directory exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", projectName)
	}

	// Select template
	var selectedTemplate string
	if templateFlag != "" {
		selectedTemplate = templateFlag
	} else {
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
	fmt.Printf("\nCreating: %s\n", projectName)

	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Println("  Setting up structure...")
	if err := scaffold.Create(projectName); err != nil {
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
		if err := scaffold.WritePlan(projectName, tmpl.Content); err != nil {
			return err
		}
	} else {
		// Write minimal plan
		minimalPlan := "# IRL Plan\n\n[Edit this file to define your research plan]\n"
		if err := scaffold.WritePlan(projectName, minimalPlan); err != nil {
			return err
		}
	}

	fmt.Println("  Initializing git...")
	if err := scaffold.GitInit(projectName); err != nil {
		fmt.Printf("  Warning: git init failed: %v\n", err)
	}

	fmt.Printf("\n✓ Created: %s\n", projectName)
	fmt.Printf("\nNext:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  # Edit 01-plans/main-plan.md\n")

	return nil
}
