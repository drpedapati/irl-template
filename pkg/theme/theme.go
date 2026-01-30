// Package theme provides warm, Claude-inspired terminal styling using Lip Gloss.
package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - warm, approachable colors inspired by Claude's personality.
// These adapt to light/dark terminal themes via lipgloss.AdaptiveColor.
var (
	// Primary: Coral/Terracotta - headers, highlights - warm and inviting
	Primary = lipgloss.AdaptiveColor{Light: "#C75146", Dark: "#E07A5F"}

	// Accent: Soft Purple - commands, interactive elements
	Accent = lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#A78BFA"}

	// Success: Sage Green - checkmarks, completion - calming
	Success = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#6EE7B7"}

	// Warning: Warm Amber - warnings - friendly alert
	Warning = lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"}

	// Error: Soft Rose - errors - noticeable but not harsh
	Error = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#F87171"}

	// Muted: Warm Gray - descriptions, secondary text - readable, not cold
	Muted = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}
)

// Reusable styles
var (
	// HeaderStyle for section headers
	HeaderStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	// CommandStyle for command names and paths
	CommandStyle = lipgloss.NewStyle().
			Foreground(Accent)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	// WarningStyle for warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error)

	// MutedStyle for dim/secondary text
	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// BoldStyle for emphasized text
	BoldStyle = lipgloss.NewStyle().
			Bold(true)
)

// Semantic styling functions

// Header returns a styled header string
func Header(s string) string {
	return HeaderStyle.Render(s)
}

// Cmd returns a styled command/path string
func Cmd(s string) string {
	return CommandStyle.Render(s)
}

// Succ returns a styled success string
func Succ(s string) string {
	return SuccessStyle.Render(s)
}

// Warn returns a styled warning string
func Warn(s string) string {
	return WarningStyle.Render(s)
}

// Err returns a styled error string
func Err(s string) string {
	return ErrorStyle.Render(s)
}

// Faint returns a styled muted/dim string
func Faint(s string) string {
	return MutedStyle.Render(s)
}

// B returns a bold string
func B(s string) string {
	return BoldStyle.Render(s)
}
