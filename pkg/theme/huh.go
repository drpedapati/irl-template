package theme

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// HuhTheme returns a custom Huh theme matching our warm color palette.
func HuhTheme() *huh.Theme {
	t := huh.ThemeBase()

	// Title styling (prompt questions)
	t.Focused.Title = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)
	t.Blurred.Title = lipgloss.NewStyle().
		Foreground(Muted)

	// Description styling
	t.Focused.Description = lipgloss.NewStyle().
		Foreground(Muted)
	t.Blurred.Description = lipgloss.NewStyle().
		Foreground(Muted)

	// Text input styling
	t.Focused.TextInput.Cursor = lipgloss.NewStyle().
		Foreground(Accent)
	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().
		Foreground(Muted)
	t.Focused.TextInput.Prompt = lipgloss.NewStyle().
		Foreground(Accent)

	// Select styling
	t.Focused.SelectSelector = lipgloss.NewStyle().
		Foreground(Accent).
		SetString(Arrow + " ")
	t.Focused.SelectedOption = lipgloss.NewStyle().
		Foreground(Accent)
	t.Focused.UnselectedOption = lipgloss.NewStyle().
		Foreground(Muted)

	// Confirm styling
	t.Focused.FocusedButton = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(Accent).
		Padding(0, 1)
	t.Focused.BlurredButton = lipgloss.NewStyle().
		Foreground(Muted).
		Padding(0, 1)

	// Base styling
	t.Focused.Base = lipgloss.NewStyle().
		PaddingLeft(0)

	// Next/Error styling
	t.Focused.Next = lipgloss.NewStyle().
		Foreground(Success)
	t.Focused.ErrorMessage = lipgloss.NewStyle().
		Foreground(Error)
	t.Focused.ErrorIndicator = lipgloss.NewStyle().
		Foreground(Error)

	return t
}

// NewForm creates a new Huh form with our custom theme.
func NewForm(groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).WithTheme(HuhTheme()).WithKeyMap(appKeyMap())
}

func appKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()

	km.Select.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "left"),
		key.WithHelp("←", "back"),
	)
	km.Select.Next = key.NewBinding(
		key.WithKeys("enter", "tab", "right"),
		key.WithHelp("→", "next"),
	)

	km.MultiSelect.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "left"),
		key.WithHelp("←", "back"),
	)
	km.MultiSelect.Next = key.NewBinding(
		key.WithKeys("enter", "tab", "right"),
		key.WithHelp("→", "next"),
	)

	km.Note.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "left"),
		key.WithHelp("←", "back"),
	)
	km.Note.Next = key.NewBinding(
		key.WithKeys("enter", "tab", "right"),
		key.WithHelp("→", "next"),
	)

	return km
}
