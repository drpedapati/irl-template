package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/naming"
	"github.com/drpedapati/irl-template/pkg/scaffold"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

var (
	adoptRenameFlag   bool
	adoptTemplateFlag string
	adoptDirFlag      string
)

var adoptCmd = &cobra.Command{
	Use:   "adopt <folder-path>",
	Short: "Adopt an existing folder as an IRL project",
	Long: `Copy an existing folder into the IRL workspace and add IRL scaffolding.

The folder is copied to your IRL workspace directory and given a plans/main-plan.md
file so it appears in the project list.

Examples:
  irl adopt ~/Downloads/my-research       Copy to workspace, keep name
  irl adopt ./experiment-data --rename     Copy with YYMMDD prefix
  irl adopt ~/paper -t irl-basic           Use specific template
  irl adopt ~/analysis -d ~/Research       Specify workspace directory`,
	Args: cobra.ExactArgs(1),
	RunE: runAdopt,
}

func init() {
	rootCmd.AddCommand(adoptCmd)
	adoptCmd.Flags().BoolVar(&adoptRenameFlag, "rename", false,
		"Rename with YYMMDD prefix (e.g., 260210-my-folder)")
	adoptCmd.Flags().StringVarP(&adoptTemplateFlag, "template", "t", "",
		"Template to use for main-plan.md")
	adoptCmd.Flags().StringVarP(&adoptDirFlag, "dir", "d", "",
		"Workspace directory (overrides default)")
}

func runAdopt(cmd *cobra.Command, args []string) error {
	// Resolve and validate source folder
	sourcePath := expandPath(args[0])

	info, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("folder not found: %s", sourcePath)
		}
		return fmt.Errorf("cannot access folder: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", sourcePath)
	}

	// Determine workspace directory
	var baseDir string
	if adoptDirFlag != "" {
		baseDir = expandPath(adoptDirFlag)
	} else {
		baseDir = config.GetDefaultDirectory()
		if baseDir == "" {
			return fmt.Errorf("no default directory configured (run 'irl init' to set one)")
		}
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("cannot create workspace directory %s: %w", baseDir, err)
	}

	// Determine target name
	folderName := filepath.Base(sourcePath)
	if adoptRenameFlag {
		folderName = naming.Timestamp() + "-" + naming.Slugify(folderName)
	}
	destPath := filepath.Join(baseDir, folderName)

	// Check destination doesn't already exist
	if _, err := os.Stat(destPath); !os.IsNotExist(err) {
		return fmt.Errorf("destination already exists: %s", destPath)
	}

	// Check source isn't already inside workspace
	absSource, _ := filepath.Abs(sourcePath)
	absBase, _ := filepath.Abs(baseDir)
	if strings.HasPrefix(absSource, absBase+string(filepath.Separator)) {
		return fmt.Errorf("folder is already inside the workspace: %s", absSource)
	}

	// Copy folder to workspace
	fmt.Println()
	fmt.Println(theme.Faint("Copying folder..."))

	cpCmd := exec.Command("cp", "-r", sourcePath, destPath)
	if err := cpCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy folder: %w", err)
	}

	// Add .gitignore if not present
	gitignorePath := filepath.Join(destPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		if err := scaffold.Create(destPath); err != nil {
			fmt.Println(theme.Note(fmt.Sprintf("couldn't create .gitignore: %v", err)))
		}
	}

	// Add plans/main-plan.md if no plan file exists
	hasPlan := fileExists(filepath.Join(destPath, "plans", "main-plan.md")) ||
		fileExists(filepath.Join(destPath, "main-plan.md")) ||
		fileExists(filepath.Join(destPath, "01-plans", "main-plan.md"))

	if !hasPlan {
		var planContent string
		if adoptTemplateFlag != "" {
			tmpl, err := templates.GetTemplate(adoptTemplateFlag)
			if err != nil {
				fmt.Println(theme.Note(fmt.Sprintf("%v, using basic template", err)))
				tmpl = templates.EmbeddedTemplates["irl-basic"]
			}
			planContent = tmpl.Content
		} else {
			planContent = templates.EmbeddedTemplates["irl-basic"].Content
		}

		if err := scaffold.WritePlan(destPath, planContent); err != nil {
			fmt.Println(theme.Note(fmt.Sprintf("couldn't create plan: %v", err)))
		}
	}

	// Git init if not already a repo, otherwise commit IRL files
	gitDir := filepath.Join(destPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		if err := scaffold.GitInit(destPath); err != nil {
			fmt.Println(theme.Note(fmt.Sprintf(
				"couldn't set up git: %v (no worries, you can do it later)", err)))
		}
	} else if !hasPlan {
		adoptCommitIRLFiles(destPath)
	}

	// Success output
	fmt.Printf("\n%s Adopted %s\n",
		theme.OK("You're all set!"),
		theme.Cmd(folderName))

	fmt.Printf("  %s %s\n", theme.Faint("Source:"), sourcePath)
	fmt.Printf("  %s %s\n", theme.Faint("Path:"), destPath)

	if adoptRenameFlag {
		fmt.Printf("  %s %s\n", theme.Faint("Renamed:"), folderName)
	}

	if !hasPlan {
		tmplName := "irl-basic"
		if adoptTemplateFlag != "" {
			tmplName = adoptTemplateFlag
		}
		fmt.Printf("  %s %s\n", theme.Faint("Template:"), tmplName)
	}

	fmt.Printf("\n%s\n", theme.B("Next steps:"))
	fmt.Printf("  %s\n", theme.Faint("Edit plans/main-plan.md to define your workflow"))

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func adoptCommitIRLFiles(projectPath string) {
	var filesToAdd []string
	for _, f := range []string{"plans/main-plan.md", ".gitignore"} {
		if fileExists(filepath.Join(projectPath, f)) {
			filesToAdd = append(filesToAdd, f)
		}
	}
	if len(filesToAdd) == 0 {
		return
	}

	addArgs := append([]string{"add"}, filesToAdd...)
	addCmd := exec.Command("git", addArgs...)
	addCmd.Dir = projectPath
	if err := addCmd.Run(); err != nil {
		return
	}

	commitCmd := exec.Command("git", "commit", "-q", "-m", "Add IRL scaffolding (adopted project)")
	commitCmd.Dir = projectPath
	_ = commitCmd.Run()
}
