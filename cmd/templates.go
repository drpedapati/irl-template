package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage templates",
	Long: `List, show, create, or delete templates.

Examples:
  irl templates                          # List all templates
  irl templates list                     # Same as above
  irl templates show <name>              # Print template content
  irl templates create <name>            # Create custom template from irl-basic
  irl templates create <name> --from X   # Create custom template from existing
  irl templates delete <name>            # Delete a custom template`,
	RunE: runTemplatesList,
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	RunE:  runTemplatesList,
}

var templatesShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Print template content",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplatesShow,
}

var templateCreateFromFlag string

var templatesCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a custom template",
	Long: `Create a custom template in your workspace's _templates/ folder.

By default, copies from irl-basic. Use --from to copy from another template.`,
	Args: cobra.ExactArgs(1),
	RunE: runTemplatesCreate,
}

var templatesDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a custom template",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplatesDelete,
}

func init() {
	rootCmd.AddCommand(templatesCmd)
	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesShowCmd)
	templatesCmd.AddCommand(templatesCreateCmd)
	templatesCmd.AddCommand(templatesDeleteCmd)
	templatesCreateCmd.Flags().StringVar(&templateCreateFromFlag, "from", "irl-basic",
		"Source template to copy from")
}

func runTemplatesList(cmd *cobra.Command, args []string) error {
	list, err := templates.ListTemplates()
	if err != nil {
		return err
	}

	// Also load custom templates
	customList := loadCustomTemplatesCLI()

	if len(list) == 0 && len(customList) == 0 {
		fmt.Printf("%s Run %s to fetch.\n",
			theme.Warn("No templates available."),
			theme.Cmd("irl update"))
		return nil
	}

	if len(list) > 0 {
		theme.Section("Templates")
		fmt.Println()
		for _, t := range list {
			fmt.Printf("  %s %s\n",
				theme.Cmd(theme.Dot),
				theme.B(t.Name))
			fmt.Printf("    %s\n", theme.Faint(t.Description))
		}
	}

	if len(customList) > 0 {
		fmt.Println()
		theme.Section("Custom Templates")
		fmt.Println()
		for _, ct := range customList {
			fmt.Printf("  %s %s\n",
				theme.Cmd(theme.Dot),
				theme.B(ct.name))
			fmt.Printf("    %s\n", theme.Faint(ct.description))
		}
	}

	fmt.Println()
	fmt.Printf("%s irl init -t %s \"purpose\"\n",
		theme.Faint("Usage:"),
		theme.Cmd("<template>"))

	return nil
}

func runTemplatesShow(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Try standard templates first
	tmpl, err := templates.GetTemplate(name)
	if err == nil {
		fmt.Print(tmpl.Content)
		return nil
	}

	// Try custom templates
	content, err := readCustomTemplate(name)
	if err != nil {
		return fmt.Errorf("template %q not found", name)
	}

	fmt.Print(content)
	return nil
}

func runTemplatesCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return fmt.Errorf("no default directory configured (run 'irl config --dir ~/path' to set one)")
	}

	// Get source template content
	var content string
	// Try custom templates first
	customContent, err := readCustomTemplate(templateCreateFromFlag)
	if err == nil {
		content = customContent
	} else {
		// Try standard templates
		tmpl, err := templates.GetTemplate(templateCreateFromFlag)
		if err != nil {
			return fmt.Errorf("source template %q not found", templateCreateFromFlag)
		}
		content = tmpl.Content
	}

	// Create _templates/<name>/main-plan.md
	templatesDir := filepath.Join(baseDir, "_templates")
	templateDir := filepath.Join(templatesDir, name)

	if _, err := os.Stat(templateDir); !os.IsNotExist(err) {
		return fmt.Errorf("template %q already exists at %s", name, templateDir)
	}

	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}

	planPath := filepath.Join(templateDir, "main-plan.md")
	if err := os.WriteFile(planPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}

	fmt.Printf("%s Created template %s\n",
		theme.OK(""),
		theme.Cmd(name))
	fmt.Printf("  %s %s\n", theme.Faint("Path:"), templateDir)
	fmt.Printf("  %s %s\n", theme.Faint("Edit:"), planPath)

	return nil
}

func runTemplatesDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return fmt.Errorf("no default directory configured")
	}

	templateDir := filepath.Join(baseDir, "_templates", name)
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("custom template %q not found", name)
	}

	if err := os.RemoveAll(templateDir); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	fmt.Printf("%s Deleted template %s\n",
		theme.OK(""),
		theme.Cmd(name))

	return nil
}

// customTemplate represents a custom template found in _templates/
type customTemplate struct {
	name        string
	description string
	path        string
}

func loadCustomTemplatesCLI() []customTemplate {
	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return nil
	}

	templatesDir := filepath.Join(baseDir, "_templates")
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil
	}

	var result []customTemplate
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		planPath := filepath.Join(templatesDir, entry.Name(), "main-plan.md")
		content, err := os.ReadFile(planPath)
		if err != nil {
			continue
		}

		desc := extractDescriptionCLI(string(content))
		if desc == "" {
			desc = "Custom template"
		}

		result = append(result, customTemplate{
			name:        entry.Name(),
			description: desc,
			path:        filepath.Join(templatesDir, entry.Name()),
		})
	}
	return result
}

func readCustomTemplate(name string) (string, error) {
	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return "", fmt.Errorf("no workspace directory")
	}

	planPath := filepath.Join(baseDir, "_templates", name, "main-plan.md")
	content, err := os.ReadFile(planPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func extractDescriptionCLI(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}
