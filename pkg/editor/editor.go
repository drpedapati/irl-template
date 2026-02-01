// Package editor provides plan file editor detection and launching.
package editor

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drpedapati/irl-template/pkg/config"
)

// EditorType distinguishes terminal vs GUI editors
type EditorType string

const (
	EditorTypeTerminal EditorType = "terminal"
	EditorTypeGUI      EditorType = "gui"
)

// Editor represents an editor for plan files
type Editor struct {
	Name    string
	Command string
	Type    EditorType
}

// TerminalEditors lists available terminal-based editors
var TerminalEditors = []Editor{
	{Name: "nano", Command: "nano", Type: EditorTypeTerminal},
	{Name: "vim", Command: "vim", Type: EditorTypeTerminal},
	{Name: "vi", Command: "vi", Type: EditorTypeTerminal},
	{Name: "Fresh", Command: "fresh", Type: EditorTypeTerminal},
}

// GUIEditors lists available GUI-based editors
var GUIEditors = []Editor{
	{Name: "VS Code", Command: "code", Type: EditorTypeGUI},
	{Name: "Cursor", Command: "cursor", Type: EditorTypeGUI},
	{Name: "Zed", Command: "zed", Type: EditorTypeGUI},
}

// EditorFinishedMsg is sent when a terminal editor closes
type EditorFinishedMsg struct {
	Err error
}

// EditorOpenedMsg is sent when a GUI editor is launched
type EditorOpenedMsg struct {
	Err error
}

// IsAvailable checks if an editor command is available
func IsAvailable(command string) bool {
	// Check command line
	if _, err := exec.LookPath(command); err == nil {
		return true
	}

	// Check macOS .app for GUI editors
	if runtime.GOOS == "darwin" {
		appNames := map[string]string{
			"cursor": "Cursor",
			"code":   "Visual Studio Code",
			"zed":    "Zed",
		}
		if appName, ok := appNames[command]; ok {
			paths := []string{
				filepath.Join("/Applications", appName+".app"),
				filepath.Join(os.Getenv("HOME"), "Applications", appName+".app"),
			}
			for _, p := range paths {
				if _, err := os.Stat(p); err == nil {
					return true
				}
			}
		}
	}

	return false
}

// GetAvailableTerminal returns list of available terminal editors
func GetAvailableTerminal() []Editor {
	var available []Editor
	for _, e := range TerminalEditors {
		if IsAvailable(e.Command) {
			available = append(available, e)
		}
	}
	return available
}

// GetAvailableGUI returns list of available GUI editors
func GetAvailableGUI() []Editor {
	var available []Editor
	for _, e := range GUIEditors {
		if IsAvailable(e.Command) {
			available = append(available, e)
		}
	}
	return available
}

// GetPreferred returns the user's preferred editor based on config
// Falls back to auto-detection if not configured or configured editor not found
func GetPreferred() (Editor, bool) {
	cfg, err := config.Load()
	if err != nil {
		return autoDetect()
	}

	editorCmd := cfg.PlanEditor
	editorType := cfg.PlanEditorType

	// If "auto" or empty, auto-detect
	if editorCmd == "" || editorCmd == "auto" {
		return autoDetect()
	}

	// Check if configured editor is available
	if !IsAvailable(editorCmd) {
		return Editor{}, false
	}

	// Determine type
	eType := EditorTypeTerminal
	if editorType == "gui" {
		eType = EditorTypeGUI
	}

	// Find the editor name
	name := editorCmd
	for _, e := range TerminalEditors {
		if e.Command == editorCmd {
			name = e.Name
			break
		}
	}
	for _, e := range GUIEditors {
		if e.Command == editorCmd {
			name = e.Name
			break
		}
	}

	return Editor{
		Name:    name,
		Command: editorCmd,
		Type:    eType,
	}, true
}

// autoDetect finds the best available editor automatically
func autoDetect() (Editor, bool) {
	// 1. Check $EDITOR environment variable
	if envEditor := os.Getenv("EDITOR"); envEditor != "" {
		if IsAvailable(envEditor) {
			return Editor{
				Name:    envEditor,
				Command: envEditor,
				Type:    EditorTypeTerminal, // $EDITOR is typically terminal-based
			}, true
		}
	}

	// 2. Check $VISUAL environment variable
	if envVisual := os.Getenv("VISUAL"); envVisual != "" {
		if IsAvailable(envVisual) {
			// VISUAL can be terminal or GUI
			eType := EditorTypeTerminal
			for _, e := range GUIEditors {
				if e.Command == envVisual {
					eType = EditorTypeGUI
					break
				}
			}
			return Editor{
				Name:    envVisual,
				Command: envVisual,
				Type:    eType,
			}, true
		}
	}

	// 3. Check terminal editors in order: nano, vim
	for _, e := range TerminalEditors[:2] { // nano, vim (not vi as last resort)
		if IsAvailable(e.Command) {
			return e, true
		}
	}

	// 4. Check GUI editors: code, cursor
	for _, e := range GUIEditors[:2] {
		if IsAvailable(e.Command) {
			return e, true
		}
	}

	// 5. Fall back to vi (always available on Unix)
	if IsAvailable("vi") {
		return Editor{Name: "vi", Command: "vi", Type: EditorTypeTerminal}, true
	}

	return Editor{}, false
}

// Open launches the editor with the given file path
// For terminal editors: uses tea.ExecProcess (suspends TUI)
// For GUI editors: runs in background
func Open(editor Editor, path string) tea.Cmd {
	if editor.Type == EditorTypeTerminal {
		return openTerminalEditor(editor.Command, path)
	}
	return openGUIEditor(editor.Command, path)
}

// openTerminalEditor suspends the TUI and opens the editor
func openTerminalEditor(command, path string) tea.Cmd {
	c := exec.Command(command, path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return EditorFinishedMsg{Err: err}
	})
}

// openGUIEditor launches the editor in the background
func openGUIEditor(command, path string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		// On macOS, use 'open -a' for GUI apps
		if runtime.GOOS == "darwin" {
			appNames := map[string]string{
				"cursor": "Cursor",
				"code":   "Visual Studio Code",
				"zed":    "Zed",
			}
			if appName, ok := appNames[command]; ok {
				cmd = exec.Command("open", "-a", appName, path)
			} else {
				cmd = exec.Command(command, path)
			}
		} else {
			cmd = exec.Command(command, path)
		}

		err := cmd.Start() // Don't wait
		return EditorOpenedMsg{Err: err}
	}
}

// GetPlanPath returns the path to main-plan.md in a project
func GetPlanPath(projectPath string) string {
	// Check root level first
	rootPlan := filepath.Join(projectPath, "main-plan.md")
	if _, err := os.Stat(rootPlan); err == nil {
		return rootPlan
	}

	// Check IRL structure
	irlPlan := filepath.Join(projectPath, "01-plans", "main-plan.md")
	if _, err := os.Stat(irlPlan); err == nil {
		return irlPlan
	}

	// Default to root level (for new projects)
	return rootPlan
}
