package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// UpdateModel handles template updates with progress
type UpdateModel struct {
	width    int
	height   int
	progress progress.Model
	status   string
	done     bool
	count    int
	err      error
	percent  float64
}

// UpdateProgressMsg updates the progress bar
type UpdateProgressMsg struct {
	Percent float64
	Status  string
}

// UpdateCompleteMsg signals update completion
type UpdateCompleteMsg struct {
	Count int
	Err   error
}

// NewUpdateModel creates a new update view
func NewUpdateModel() UpdateModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
		progress.WithoutPercentage(),
	)
	return UpdateModel{
		progress: p,
		status:   "Connecting to GitHub...",
	}
}

// SetSize sets the view dimensions
func (m *UpdateModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// StartUpdate begins the update process
func (m *UpdateModel) StartUpdate() tea.Cmd {
	return tea.Batch(
		m.animateProgress(),
		m.fetchTemplates(),
	)
}

func (m *UpdateModel) animateProgress() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return UpdateProgressMsg{Percent: -1} // Signal to animate
	})
}

func (m *UpdateModel) fetchTemplates() tea.Cmd {
	return func() tea.Msg {
		list, err := templates.FetchTemplates()
		if err != nil {
			return UpdateCompleteMsg{Err: err}
		}
		return UpdateCompleteMsg{Count: len(list)}
	}
}

// IsDone returns true if update is complete
func (m UpdateModel) IsDone() bool {
	return m.done
}

// Update handles messages
func (m UpdateModel) Update(msg tea.Msg) (UpdateModel, tea.Cmd) {
	switch msg := msg.(type) {
	case UpdateProgressMsg:
		if m.done {
			return m, nil
		}
		// Animate progress up to 90% while waiting
		if m.percent < 0.9 {
			m.percent += 0.02
			if m.percent > 0.3 && m.percent < 0.35 {
				m.status = "Fetching template list..."
			} else if m.percent > 0.5 && m.percent < 0.55 {
				m.status = "Downloading templates..."
			} else if m.percent > 0.7 && m.percent < 0.75 {
				m.status = "Caching templates..."
			}
		}
		return m, m.animateProgress()

	case UpdateCompleteMsg:
		m.done = true
		m.percent = 1.0
		m.count = msg.Count
		m.err = msg.Err
		if m.err != nil {
			m.status = "Update failed"
		} else {
			m.status = fmt.Sprintf("Found %d templates", m.count)
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

// View renders the update view
func (m UpdateModel) View() string {
	var b strings.Builder

	// Center the content
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary).
		MarginBottom(1)

	statusStyle := lipgloss.NewStyle().
		Foreground(theme.Muted).
		MarginBottom(1)

	b.WriteString("\n\n")
	b.WriteString("  " + titleStyle.Render("Updating Templates"))
	b.WriteString("\n\n")

	// Progress bar
	b.WriteString("  " + m.progress.ViewAs(m.percent))
	b.WriteString("\n\n")

	// Status
	if m.done {
		if m.err != nil {
			errStyle := lipgloss.NewStyle().Foreground(theme.Error)
			b.WriteString("  " + errStyle.Render("✗ " + m.err.Error()))
		} else {
			successStyle := lipgloss.NewStyle().Foreground(theme.Success)
			b.WriteString("  " + successStyle.Render(fmt.Sprintf("✓ %s", m.status)))
		}
	} else {
		b.WriteString("  " + statusStyle.Render(m.status))
	}
	b.WriteString("\n")

	return b.String()
}
