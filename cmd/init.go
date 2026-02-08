package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/naming"
	"github.com/drpedapati/irl-template/pkg/scaffold"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
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
  irl init -t irl-basic                 # With specific template
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

	// Only prompt for the base directory when no purpose/name was provided.
	shouldPromptForDir := len(args) == 0 && nameFlag == ""

	// Get base directory
	if dirFlag != "" {
		// Flag overrides everything
		baseDir = expandPath(dirFlag)
	} else if shouldPromptForDir {
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
		// Interactive mode - ask for purpose
		form := theme.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("What's this project for?").
					Description("Brief description (e.g., 'ERP correlation analysis')").
					Placeholder("Enter your project purpose...").
					Value(&purpose),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		if purpose == "" {
			return fmt.Errorf("project purpose is required")
		}

		projectName = naming.GenerateName(purpose)

		// Confirm or allow edit
		form = theme.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Project folder name").
					Description("Edit if you'd like a different name").
					Value(&projectName),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}
	}

	// Full project path
	projectPath := filepath.Join(baseDir, projectName)

	// Check if directory exists
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("heads up: '%s' already exists", projectPath)
	}

	// Select template (prompt unless -t flag provided)
	var selectedTemplate string
	if templateFlag != "" {
		selectedTemplate = templateFlag
	} else {
		// Always offer template selection if no -t flag
		templateList, err := templates.ListTemplates()
		if err != nil {
			fmt.Println(theme.Note("couldn't fetch templates, using basic"))
		}

		if len(templateList) > 0 {
			// Build options for Huh select - type-safe
			options := []huh.Option[string]{
				huh.NewOption("None (empty plan)", ""),
			}
			for _, t := range templateList {
				label := fmt.Sprintf("%s - %s", t.Name, t.Description)
				options = append(options, huh.NewOption(label, t.Name))
			}

			form := theme.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Pick a template").
						Description("Templates give your project a head start").
						Options(options...).
						Value(&selectedTemplate),
				),
			)

			if err := form.Run(); err != nil {
				return err
			}
		}
	}

	// Create project
	fmt.Println()
	fmt.Println(theme.Faint("Setting things up..."))

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := scaffold.Create(projectPath); err != nil {
		return err
	}

	// Apply template
	if selectedTemplate != "" {
		tmpl, err := templates.GetTemplate(selectedTemplate)
		if err != nil {
			fmt.Println(theme.Note(fmt.Sprintf("%v, using basic template", err)))
			tmpl = templates.EmbeddedTemplates["irl-basic"]
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

	if err := scaffold.GitInit(projectPath); err != nil {
		fmt.Println(theme.Note(fmt.Sprintf("couldn't set up git: %v (no worries, you can do it later)", err)))
	}

	// Success output - warm and friendly
	fmt.Printf("\n%s Created %s\n",
		theme.OK("You're all set!"),
		theme.Cmd(projectPath))

	if selectedTemplate != "" {
		fmt.Printf("  %s %s\n", theme.Faint("Template:"), selectedTemplate)
	}

	fmt.Printf("\n%s\n", theme.B("Next steps:"))
	fmt.Printf("  %s %s\n", theme.Cmd("cd"), projectPath)
	fmt.Printf("  %s\n", theme.Faint("# Edit main-plan.md"))

	fmt.Printf("\n%s\n", theme.B("Open in IDE:"))
	fmt.Printf("  %s %s\n", theme.Cmd("cursor"), projectPath)
	fmt.Printf("  %s %s\n", theme.Cmd("code"), projectPath)
	fmt.Printf("  %s %s\n", theme.Cmd("positron"), projectPath)

	fmt.Printf("\n%s\n",
		theme.Faint("No shell command? Open IDE → Cmd+Shift+P → \"Install ... in PATH\""))

	return nil
}

func getOrAskDefaultDirectory() string {
	// Check if already configured
	defaultDir := config.GetDefaultDirectory()
	if defaultDir != "" {
		// Show current setting and offer to change
		fmt.Printf("Default directory: %s\n", defaultDir)

		var useExisting bool = true
		form := theme.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Use this directory?").
					Affirmative("Yes").
					Negative("Change it").
					Value(&useExisting),
			),
		)
		form.Run()

		if useExisting {
			return defaultDir
		}
	}

	// Ask for directory (browse or type a path)
	home, _ := os.UserHomeDir()
	suggestion := filepath.Join(home, "Documents", "irl_projects")
	if defaultDir != "" {
		suggestion = defaultDir
	}
	suggestion = expandPath(suggestion)

	// Start near the suggestion if possible, otherwise fall back to home
	startDir := suggestion
	if info, err := os.Stat(startDir); err != nil || !info.IsDir() {
		startDir = filepath.Dir(suggestion)
	}
	if info, err := os.Stat(startDir); err != nil || !info.IsDir() {
		startDir = home
	}
	if startDir == "" {
		startDir = "."
	}

	for {
		method := "browse"
		if err := theme.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose a project directory").
					Description("Browse folders or type a path. Use ←/→ to go back/next.").
					Options(
						huh.NewOption("Browse folders", "browse"),
						huh.NewOption("Type a path", "type"),
					).
					Value(&method),
			),
		).Run(); err != nil {
			return suggestion
		}

		newDir := suggestion
		switch method {
		case "browse":
			var picked string
			if err := theme.NewForm(
				huh.NewGroup(
					huh.NewFilePicker().
						Title("Select a folder").
						Description("↑/↓ move • → open • ← up • enter select • esc back").
						DirAllowed(true).
						FileAllowed(false).
						ShowHidden(false).
						ShowSize(false).
						ShowPermissions(false).
						Height(18).
						Picking(true).
						CurrentDirectory(startDir).
						Value(&picked),
				),
			).Run(); err != nil {
				return suggestion
			}
			if picked == "" {
				continue
			}
			newDir = picked

			var folderName string
			if err := theme.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Create a new folder here?").
						Description("Leave blank to keep the selected folder.").
						Placeholder("New folder name (optional)").
						Value(&folderName),
				),
			).Run(); err == nil {
				folderName = strings.TrimSpace(folderName)
				if folderName != "" {
					newDir = filepath.Join(newDir, folderName)
				}
			}
		case "type":
			// newDir already set to suggestion
		default:
			return suggestion
		}

		inputTitle := "Project directory"
		inputDescription := "You can type a new folder; it can be created if needed."
		if method == "browse" {
			inputTitle = "Confirm directory"
			inputDescription = "Edit to create a new folder or refine the path."
		}
		if err := theme.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(inputTitle).
					Description(inputDescription).
					Placeholder(suggestion).
					Value(&newDir),
			),
		).Run(); err != nil {
			return suggestion
		}

		newDir = expandPath(newDir)
		if newDir == "" {
			return suggestion
		}

		if info, err := os.Stat(newDir); err == nil && !info.IsDir() {
			fmt.Println(theme.Note("That path is a file; using its parent folder instead."))
			newDir = filepath.Dir(newDir)
		}

		status := "exists"
		if _, err := os.Stat(newDir); os.IsNotExist(err) {
			status = "will be created"
		}

		useDir := true
		if err := theme.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Use this directory?").
					Description(fmt.Sprintf("%s (%s)", newDir, status)).
					Affirmative("Use").
					Negative("Start over").
					Value(&useDir),
			),
		).Run(); err != nil {
			return suggestion
		}
		if !useDir {
			continue
		}

		if status == "will be created" {
			if err := os.MkdirAll(newDir, 0755); err != nil {
				fmt.Println(theme.Note(fmt.Sprintf("couldn't create folder: %v", err)))
			}
		}

		// Save to config
		if err := config.SetDefaultDirectory(newDir); err != nil {
			fmt.Println(theme.Note(fmt.Sprintf("couldn't save config: %v", err)))
		} else {
			fmt.Printf("Saved default directory: %s\n", newDir)
		}

		return newDir
	}
}

func expandPath(path string) string {
	path = strings.TrimSpace(path)
	home, _ := os.UserHomeDir()

	// Handle bare ~ or ~/<path> (plus ~\<path> on Windows)
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~"+string(filepath.Separator)) {
		return filepath.Join(home, path[2:])
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
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
