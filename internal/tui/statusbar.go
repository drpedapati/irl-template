package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar represents the bottom status bar
type StatusBar struct {
	width int
	keys  []KeyBinding
}

// KeyBinding represents a keyboard shortcut
type KeyBinding struct {
	Key  string
	Desc string
}

// NewStatusBar creates a new status bar
func NewStatusBar() StatusBar {
	return StatusBar{
		keys: []KeyBinding{
			{Key: "↑↓", Desc: "Navigate"},
			{Key: "Enter", Desc: "Select"},
			{Key: "q", Desc: "Quit"},
		},
	}
}

// SetWidth sets the status bar width
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// SetKeys sets custom key bindings for the current view
func (s *StatusBar) SetKeys(keys []KeyBinding) {
	s.keys = keys
}

// View renders the status bar (centered)
func (s StatusBar) View() string {
	var parts []string
	for _, k := range s.keys {
		parts = append(parts, FormatKey(k.Key, k.Desc))
	}

	content := strings.Join(parts, "  ")

	// Center the content
	contentWidth := lipgloss.Width(content)
	padding := (s.width - contentWidth) / 2
	if padding < 0 {
		padding = 0
	}

	return StatusBarStyle.Render(strings.Repeat(" ", padding) + content)
}

// ViewCompact renders the status bar keys left-aligned without centering
func (s StatusBar) ViewCompact() string {
	var parts []string
	for _, k := range s.keys {
		parts = append(parts, FormatKey(k.Key, k.Desc))
	}

	return strings.Join(parts, "  ")
}

// DefaultMenuKeys returns the default key bindings for the main menu
func DefaultMenuKeys() []KeyBinding {
	return []KeyBinding{
		{Key: "↑↓", Desc: "Navigate"},
		{Key: "→", Desc: "Select"},
		{Key: "q", Desc: "Quit"},
	}
}

// ViewKeys returns key bindings for sub-views
func ViewKeys() []KeyBinding {
	return []KeyBinding{
		{Key: "↑↓", Desc: "Navigate"},
		{Key: "←", Desc: "Back"},
		{Key: "q", Desc: "Quit"},
	}
}

// InitViewKeys returns key bindings for the init view
func InitViewKeys() []KeyBinding {
	return []KeyBinding{
		{Key: "→", Desc: "Next"},
		{Key: "←", Desc: "Back"},
		{Key: "Enter", Desc: "Confirm"},
	}
}

// TemplateViewKeys returns key bindings for the templates view
func TemplateViewKeys() []KeyBinding {
	return []KeyBinding{
		{Key: "↑↓", Desc: "Navigate"},
		{Key: "←", Desc: "Back"},
		{Key: "r", Desc: "Refresh"},
	}
}
