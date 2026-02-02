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
		keys: DefaultMenuKeys(),
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
	return DefaultMenuKeysWithUpdate(0, false)
}

// DefaultMenuKeysWithUpdate returns menu keys with update count
func DefaultMenuKeysWithUpdate(newTemplates int, checking bool) []KeyBinding {
	updateDesc := "Update"
	if checking {
		updateDesc = "Update (...)"
	} else if newTemplates > 0 {
		updateDesc = "Update (" + itoa(newTemplates) + " new)"
	}
	return []KeyBinding{
		{Key: "i", Desc: "Profile"},
		{Key: "u", Desc: updateDesc},
		{Key: "↑↓", Desc: "Navigate"},
		{Key: "→", Desc: "Select"},
		{Key: "q", Desc: "Quit"},
	}
}

// Simple int to string for status bar
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
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

// ProjectsViewKeys returns key bindings for the projects view (selection mode)
func ProjectsViewKeys() []KeyBinding {
	return []KeyBinding{
		{Key: "↑↓", Desc: "Navigate"},
		{Key: "e", Desc: "Edit"},
		{Key: "b", Desc: "Backup"},
		{Key: "x", Desc: "Delete"},
		{Key: "←", Desc: "Back"},
	}
}

// ProjectsFilterKeys returns key bindings for the projects view (filter mode)
func ProjectsFilterKeys() []KeyBinding {
	return []KeyBinding{
		{Key: "↓", Desc: "Select"},
		{Key: "Enter", Desc: "Open"},
		{Key: "Esc", Desc: "Clear"},
		{Key: "←", Desc: "Back"},
	}
}
