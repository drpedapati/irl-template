package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/doctor"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// DoctorModel is the environment check view
type DoctorModel struct {
	systemInfo string
	results    []ToolResult
	width      int
	height     int
	loaded     bool
}

// NewDoctorModel creates a new doctor view
func NewDoctorModel() DoctorModel {
	return DoctorModel{}
}

// SetSize sets the view dimensions
func (m *DoctorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// RunChecks returns a command that runs environment checks
func (m *DoctorModel) RunChecks() tea.Cmd {
	return func() tea.Msg {
		sysInfo := doctor.GetSystemInfo()
		results := doctor.CheckAllTools()
		return DoctorResultMsg{
			SystemInfo: sysInfo.String(),
			Results:    convertResults(results),
		}
	}
}

func convertResults(dr []doctor.ToolResult) []ToolResult {
	results := make([]ToolResult, len(dr))
	for i, r := range dr {
		results[i] = ToolResult{
			Category: r.Tool.Category,
			Name:     r.Tool.Name,
			Found:    r.Found,
		}
	}
	return results
}

// Update handles messages
func (m DoctorModel) Update(msg tea.Msg) (DoctorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case DoctorResultMsg:
		m.loaded = true
		m.systemInfo = msg.SystemInfo
		m.results = msg.Results
	}
	return m, nil
}

// View renders the doctor view
func (m DoctorModel) View() string {
	var b strings.Builder

	if !m.loaded {
		return b.String()
	}

	// System info on first line
	infoStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	b.WriteString(infoStyle.Render(m.systemInfo))
	b.WriteString("\n")

	// Group results by category
	grouped := make(map[string][]ToolResult)
	for _, r := range m.results {
		grouped[r.Category] = append(grouped[r.Category], r)
	}

	// Two-column layout
	colWidth := 30

	// Row 1: Core Tools | AI Assistants
	b.WriteString(m.renderColumnHeaders("Core Tools", "AI Assistants", colWidth))
	b.WriteString(m.renderTwoColumns(grouped["Core Tools"], grouped["AI Assistants"], colWidth))

	// Row 2: IDEs | Sandbox
	b.WriteString("\n")
	b.WriteString(m.renderColumnHeaders("IDEs", "Sandbox", colWidth))
	b.WriteString(m.renderTwoColumns(grouped["IDEs"], grouped["Sandbox"], colWidth))

	// Docker hint
	if doctor.HasDocker() {
		b.WriteString("\n")
		hintStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
		cmdStyle := lipgloss.NewStyle().Foreground(theme.Accent)
		b.WriteString(hintStyle.Render("Tip: ") + cmdStyle.Render("docker sandbox run claude"))
		b.WriteString("\n")
	}

	return b.String()
}

func (m DoctorModel) renderColumnHeaders(left, right string, colWidth int) string {
	headerStyle := lipgloss.NewStyle().Foreground(theme.Muted).MarginLeft(2)
	leftCol := lipgloss.NewStyle().Width(colWidth)
	return headerStyle.Render(leftCol.Render(left) + "  " + right) + "\n"
}

func (m DoctorModel) renderTwoColumns(left, right []ToolResult, colWidth int) string {
	var b strings.Builder

	maxRows := len(left)
	if len(right) > maxRows {
		maxRows = len(right)
	}

	leftCol := lipgloss.NewStyle().Width(colWidth)

	for i := 0; i < maxRows; i++ {
		leftStr := ""
		rightStr := ""

		if i < len(left) {
			leftStr = formatToolCheck(left[i])
		}
		if i < len(right) {
			rightStr = formatToolCheck(right[i])
		}

		b.WriteString("  " + leftCol.Render(leftStr) + "  " + rightStr + "\n")
	}

	return b.String()
}

func formatToolCheck(r ToolResult) string {
	check := lipgloss.NewStyle().Foreground(theme.Success).Render("✓")
	cross := lipgloss.NewStyle().Foreground(theme.Muted).Render("✗")

	if r.Found {
		return check + " " + r.Name
	}
	return cross + " " + lipgloss.NewStyle().Foreground(theme.Muted).Render(r.Name)
}
