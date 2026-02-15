package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/projects"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

var openEditorFlag string

var openCmd = &cobra.Command{
	Use:   "open <project>",
	Short: "Open a project in an editor",
	Long: `Open an IRL project directory in your preferred editor.

Uses the configured plan editor, or specify one with --editor.

Examples:
  irl open my-project                # Preferred editor
  irl open my-project --editor code  # VS Code
  irl open my-project --editor cursor`,
	Args: cobra.ExactArgs(1),
	RunE: runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().StringVar(&openEditorFlag, "editor", "", "Editor command (e.g., code, cursor, vim)")
}

func runOpen(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// Find the project
	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return fmt.Errorf("no default directory configured (run 'irl config --dir ~/path' to set one)")
	}

	// Try exact path first
	projectPath := filepath.Join(baseDir, projectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		// Try fuzzy match from project list
		list, scanErr := projects.ScanDir(baseDir)
		if scanErr != nil {
			return fmt.Errorf("project %q not found", projectName)
		}
		for _, p := range list {
			if p.Name == projectName {
				projectPath = p.Path
				break
			}
		}
		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			return fmt.Errorf("project %q not found in %s", projectName, baseDir)
		}
	}

	// Determine editor
	editorCmd := openEditorFlag
	if editorCmd == "" {
		editorCmd = config.GetPlanEditor()
	}
	if editorCmd == "" || editorCmd == "auto" {
		editorCmd = detectEditor()
	}
	if editorCmd == "" {
		return fmt.Errorf("no editor configured (run 'irl config --editor <cmd>' to set one)")
	}

	// Open project
	if err := launchEditor(editorCmd, projectPath); err != nil {
		return fmt.Errorf("failed to open %s with %s: %w", projectName, editorCmd, err)
	}

	fmt.Printf("%s Opened %s in %s\n",
		theme.OK(""),
		theme.Cmd(projectName),
		editorCmd)

	return nil
}

func detectEditor() string {
	// Try common editors in preference order
	for _, cmd := range []string{"cursor", "code", "positron", "vim"} {
		if _, err := exec.LookPath(cmd); err == nil {
			return cmd
		}
	}
	return ""
}

func launchEditor(editorCmd, projectPath string) error {
	var cmd *exec.Cmd

	// GUI editors on macOS: use 'open -a' if it's an app name
	switch editorCmd {
	case "vim", "nvim", "nano", "helix", "emacs":
		// Terminal editors: open in a new terminal window
		if runtime.GOOS == "darwin" {
			script := fmt.Sprintf(`tell application "Terminal" to do script "cd '%s' && %s"`,
				projectPath, editorCmd)
			cmd = exec.Command("osascript", "-e", script)
		} else {
			cmd = exec.Command("x-terminal-emulator", "-e", "sh", "-c",
				fmt.Sprintf("cd '%s' && %s", projectPath, editorCmd))
		}
	default:
		// GUI editors: launch directly
		cmd = exec.Command(editorCmd, projectPath)
	}

	return cmd.Start()
}
