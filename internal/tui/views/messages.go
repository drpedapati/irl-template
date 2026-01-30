// Package views contains the sub-views for the TUI.
package views

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// itoa is a helper to convert int to string
func itoa(i int) string {
	return strconv.Itoa(i)
}

// Message types for view communication

// TemplatesLoadedMsg is sent when templates are loaded
type TemplatesLoadedMsg struct {
	Templates []TemplateItem
	Err       error
}

// TemplateItem represents a template for display
type TemplateItem struct {
	Name        string
	Description string
	Content     string
	Source      string // "default" (GitHub) or "custom" (local _templates folder)
}

// DoctorResultMsg is sent when doctor checks complete
type DoctorResultMsg struct {
	SystemInfo string
	Results    []ToolResult
}

// ToolResult represents a tool check result
type ToolResult struct {
	Category string
	Name     string
	Found    bool
}

// InitCompleteMsg is sent when project creation completes
type InitCompleteMsg struct {
	ProjectPath string
	Err         error
}

// ConfigLoadedMsg is sent when config is loaded
type ConfigLoadedMsg struct {
	DefaultDir string
	Err        error
}

// BackToMenuMsg signals to return to the main menu
type BackToMenuMsg struct{}

// BackToMenu returns a command that sends BackToMenuMsg
func BackToMenu() tea.Cmd {
	return func() tea.Msg {
		return BackToMenuMsg{}
	}
}
