package theme

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Common symbols
const (
	Check = "✓"
	Cross = "✗"
	Dot   = "•"
	Arrow = "→"
)

// Prefixed output helpers - API compatible with old style package

// OK returns a success checkmark with message
func OK(msg string) string {
	return SuccessStyle.Render(Check) + " " + msg
}

// Fail returns a failure X with message
func Fail(msg string) string {
	return ErrorStyle.Render(Cross) + " " + msg
}

// Item returns a bullet point with message
func Item(msg string) string {
	return CommandStyle.Render(Dot) + " " + msg
}

// Note returns a friendly note prefix with message
func Note(msg string) string {
	return WarningStyle.Render("Note:") + " " + msg
}

// Print helpers

// PrintHeader prints a styled header
func PrintHeader(s string) {
	fmt.Println(Header(s))
}

// PrintOK prints a success message
func PrintOK(s string) {
	fmt.Println(OK(s))
}

// PrintFail prints a failure message
func PrintFail(s string) {
	fmt.Println(Fail(s))
}

// PrintItem prints a bullet item
func PrintItem(s string) {
	fmt.Println(Item(s))
}

// PrintCmd prints a labeled command
func PrintCmd(label, c string) {
	fmt.Printf("  %s %s\n", Faint(label), Cmd(c))
}

// Section helpers for structured output

// Section prints a section header with spacing
func Section(title string) {
	fmt.Printf("\n%s\n", Header(title))
}

// SubSection prints a subsection header
func SubSection(title string) {
	fmt.Printf("  %s\n", Faint(title))
}

// Row helpers for two-column layouts

// RowStyle returns a style for fixed-width rows
func RowStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width)
}

// ToolCheck formats a tool check result (checkmark/X with name)
func ToolCheck(name string, found bool) string {
	if found {
		return SuccessStyle.Render(Check) + " " + name
	}
	return MutedStyle.Render(Cross + " " + name)
}

// ToolCheckWidth returns the visible width of a ToolCheck string
func ToolCheckWidth(name string) int {
	return 2 + len(name) // symbol + space + name
}

// KeyValue formats a key-value pair for config display
func KeyValue(key, value string) string {
	return fmt.Sprintf("  %s %s", MutedStyle.Render(key), value)
}

// StatusTag returns a colored status indicator
func StatusTag(status string, isGood bool) string {
	if isGood {
		return SuccessStyle.Render(status)
	}
	return WarningStyle.Render(status)
}
