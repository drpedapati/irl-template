// Package tui provides the interactive terminal UI for IRL.
package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// TUI-specific styles built on theme colors
var (
	// Box styles for layout
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Muted).
			Padding(0, 1)

	// Header area
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Padding(0, 1)

	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#CC785C", Dark: "#D97757"}) // Claude orange

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#1F1F1F", Dark: "#F5F5F5"}). // High contrast text
			PaddingLeft(1)

	VersionStyle = lipgloss.NewStyle().
			Foreground(theme.Muted)

	// Menu styles
	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	MenuSelectedStyle = lipgloss.NewStyle().
				Foreground(theme.Accent).
				Bold(true).
				PaddingLeft(2)

	MenuCursorStyle = lipgloss.NewStyle().
			Foreground(theme.Accent).
			Bold(true)

	MenuDescStyle = lipgloss.NewStyle().
			Foreground(theme.Muted).
			PaddingLeft(4)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(theme.Muted).
			Padding(0, 1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(theme.Accent).
			Bold(true)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(theme.Muted)

	// Content area
	ContentStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// View titles
	ViewTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			MarginBottom(1)

	// List items in views
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	ListSelectedStyle = lipgloss.NewStyle().
				Foreground(theme.Accent).
				Bold(true).
				PaddingLeft(2)

	// Status indicators
	SuccessStyle = lipgloss.NewStyle().
			Foreground(theme.Success)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(theme.Error)

	WarningStyle = lipgloss.NewStyle().
			Foreground(theme.Warning)

	// Spinner/loading
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(theme.Accent)

	// Input fields
	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Muted).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(theme.Accent).
				Padding(0, 1)

	// Help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(theme.Muted)

	// Divider
	DividerStyle = lipgloss.NewStyle().
			Foreground(theme.Muted)
)

// Symbols used in TUI
const (
	Cursor     = "●"
	CursorDim  = " "
	Check      = "✓"
	Cross      = "✗"
	Arrow      = "→"
	Dot        = "•"
	Spinner    = "◐"
	DividerCh  = "─"
)

// FormatKey formats a keyboard shortcut for display
func FormatKey(key, desc string) string {
	return KeyStyle.Render(key) + " " + KeyDescStyle.Render(desc)
}

// FormatKeyCompact formats a key without description styling
func FormatKeyCompact(key string) string {
	return KeyStyle.Render(key)
}
