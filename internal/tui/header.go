package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Header represents the top header component
type Header struct {
	width   int
	version string
}

// NewHeader creates a new header component
func NewHeader(version string) Header {
	return Header{
		version: version,
	}
}

// SetWidth sets the header width
func (h *Header) SetWidth(width int) {
	h.width = width
}

// View renders the header
func (h Header) View() string {
	logo := LogoStyle.Render("IRL")
	title := TitleStyle.Render("Idempotent Research Loop")
	version := VersionStyle.Render("v" + h.version)

	left := logo + title

	// Calculate padding to right-align version
	padding := h.width - lipgloss.Width(left) - lipgloss.Width(version) - 4
	if padding < 1 {
		padding = 1
	}

	return left + strings.Repeat(" ", padding) + version
}

// Divider returns a horizontal divider line
func Divider(width int) string {
	if width < 2 {
		width = 2
	}
	return DividerStyle.Render(strings.Repeat(DividerCh, width-2))
}
