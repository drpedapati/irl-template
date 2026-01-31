package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/theme"
)

// HelpSection represents which section of help we're viewing
type HelpSection int

const (
	HelpSectionMenu       HelpSection = iota // The help menu
	HelpSectionWhatIsIRL                     // "What is IRL?" slide deck
	// Future sections:
	// HelpSectionFirstProject
	// HelpSectionTemplates
)

// Menu items for the help menu
type helpMenuItem struct {
	title       string
	description string
	duration    string // e.g., "5 slides · 1 min"
	section     HelpSection
}

var helpMenuItems = []helpMenuItem{
	{
		title:       "What is IRL?",
		description: "Learn the core concept",
		duration:    "5 slides · 1 min",
		section:     HelpSectionWhatIsIRL,
	},
	// Future items:
	// {
	// 	title:       "Your First Project",
	// 	description: "Step-by-step tutorial",
	// 	duration:    "8 slides · 3 min",
	// 	section:     HelpSectionFirstProject,
	// },
}

const whatIsIRLSlides = 5

// HelpLoadedMsg is sent when help content is loaded (kept for compatibility)
type HelpLoadedMsg struct {
	Err error
}

// HelpModel is a hierarchical help view with menu and slide decks
type HelpModel struct {
	section      HelpSection
	menuCursor   int
	currentSlide int
	width        int
	height       int
}

// NewHelpModel creates a new help view starting at the menu
func NewHelpModel() HelpModel {
	return HelpModel{
		section:      HelpSectionMenu,
		menuCursor:   0,
		currentSlide: 0,
	}
}

// SetSize sets the view dimensions
func (m *HelpModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// LoadHelp returns a command (no-op for slides)
func (m *HelpModel) LoadHelp() tea.Cmd {
	return nil
}

// IsAtMenu returns true if we're at the help menu (not in a slide deck)
func (m HelpModel) IsAtMenu() bool {
	return m.section == HelpSectionMenu
}

// IsFirstSlide returns true if on the first slide of current section
func (m HelpModel) IsFirstSlide() bool {
	if m.section == HelpSectionMenu {
		return true
	}
	return m.currentSlide == 0
}

// CurrentSlide returns the current slide index
func (m HelpModel) CurrentSlide() int {
	return m.currentSlide
}

// TotalSlides returns total slides in current section
func (m HelpModel) TotalSlides() int {
	switch m.section {
	case HelpSectionWhatIsIRL:
		return whatIsIRLSlides
	default:
		return 0
	}
}

// Update handles messages
func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m HelpModel) handleKey(msg tea.KeyMsg) (HelpModel, tea.Cmd) {
	if m.section == HelpSectionMenu {
		return m.handleMenuKey(msg)
	}
	return m.handleSlideKey(msg)
}

func (m HelpModel) handleMenuKey(msg tea.KeyMsg) (HelpModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuCursor > 0 {
			m.menuCursor--
		}
	case "down", "j":
		if m.menuCursor < len(helpMenuItems)-1 {
			m.menuCursor++
		}
	case "enter", "right", "l":
		// Enter the selected section
		m.section = helpMenuItems[m.menuCursor].section
		m.currentSlide = 0
	case "1":
		if len(helpMenuItems) >= 1 {
			m.section = helpMenuItems[0].section
			m.currentSlide = 0
		}
	case "2":
		if len(helpMenuItems) >= 2 {
			m.section = helpMenuItems[1].section
			m.currentSlide = 0
		}
	}
	return m, nil
}

func (m HelpModel) handleSlideKey(msg tea.KeyMsg) (HelpModel, tea.Cmd) {
	totalSlides := m.TotalSlides()

	switch msg.String() {
	case "esc":
		// Go back to menu
		m.section = HelpSectionMenu
		m.currentSlide = 0
	case "left", "h":
		if m.currentSlide > 0 {
			m.currentSlide--
		} else {
			// At first slide, go back to menu
			m.section = HelpSectionMenu
		}
	case "right", "l", "enter", " ":
		if m.currentSlide < totalSlides-1 {
			m.currentSlide++
		}
	case "home", "g":
		m.currentSlide = 0
	case "end", "G":
		m.currentSlide = totalSlides - 1
	case "1":
		m.currentSlide = 0
	case "2":
		if totalSlides > 1 {
			m.currentSlide = 1
		}
	case "3":
		if totalSlides > 2 {
			m.currentSlide = 2
		}
	case "4":
		if totalSlides > 3 {
			m.currentSlide = 3
		}
	case "5":
		if totalSlides > 4 {
			m.currentSlide = 4
		}
	}
	return m, nil
}

// View renders the current view (menu or slides)
func (m HelpModel) View() string {
	if m.section == HelpSectionMenu {
		return m.renderMenu()
	}
	return m.renderSlides()
}

// ============================================================================
// HELP MENU
// ============================================================================

func (m HelpModel) renderMenu() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	// Build menu content
	var content strings.Builder

	content.WriteString(titleStyle.Render("Help & Tutorials"))
	content.WriteString("\n\n")

	// Render each menu item
	for i, item := range helpMenuItems {
		content.WriteString(m.renderMenuItem(i, item))
		if i < len(helpMenuItems)-1 {
			content.WriteString("\n\n")
		}
	}

	// Center the menu
	menuContent := content.String()
	contentHeight := m.height - 3
	if contentHeight < 10 {
		contentHeight = 10
	}

	centered := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Center,
		lipgloss.Center,
		menuContent,
	)

	// Footer
	footer := m.renderMenuFooter()

	return centered + footer
}

func (m HelpModel) renderMenuItem(index int, item helpMenuItem) string {
	isSelected := index == m.menuCursor

	// Styles
	keyStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	titleStyle := lipgloss.NewStyle().Bold(true)
	selectedTitleStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	durationStyle := lipgloss.NewStyle().Foreground(theme.Success)
	cursorStyle := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true)

	var line strings.Builder

	// Line 1: cursor [key] Title                    duration
	if isSelected {
		line.WriteString(cursorStyle.Render("●"))
	} else {
		line.WriteString(" ")
	}
	line.WriteString(" ")
	line.WriteString(keyStyle.Render("[" + string('1'+rune(index)) + "]"))
	line.WriteString(" ")

	if isSelected {
		line.WriteString(selectedTitleStyle.Render(item.title))
	} else {
		line.WriteString(titleStyle.Render(item.title))
	}

	line.WriteString("  ")
	line.WriteString(durationStyle.Render(item.duration))
	line.WriteString("\n")

	// Line 2: description (indented)
	line.WriteString("      ")
	line.WriteString(descStyle.Render(item.description))

	return line.String()
}

func (m HelpModel) renderMenuFooter() string {
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	accentStyle := lipgloss.NewStyle().Foreground(theme.Accent)

	hint := accentStyle.Render("↑↓") + mutedStyle.Render(" navigate  ") +
		accentStyle.Render("Enter") + mutedStyle.Render(" select  ") +
		accentStyle.Render("esc") + mutedStyle.Render(" back")

	return lipgloss.Place(m.width, 2, lipgloss.Center, lipgloss.Bottom, hint)
}

// ============================================================================
// SLIDE DECK RENDERER
// ============================================================================

func (m HelpModel) renderSlides() string {
	var content string

	switch m.section {
	case HelpSectionWhatIsIRL:
		content = m.renderWhatIsIRLSlide()
	}

	// Build footer with navigation
	footer := m.renderSlideFooter()
	footerHeight := lipgloss.Height(footer)

	// Calculate available height for content
	contentHeight := m.height - footerHeight - 1
	if contentHeight < 10 {
		contentHeight = 10
	}

	// Center content both horizontally and vertically
	centeredContent := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	return centeredContent + footer
}

func (m HelpModel) renderSlideFooter() string {
	totalSlides := m.TotalSlides()

	// Slide indicator dots
	var dots string
	for i := 0; i < totalSlides; i++ {
		if i == m.currentSlide {
			dots += lipgloss.NewStyle().Foreground(theme.Primary).Render("●")
		} else {
			dots += lipgloss.NewStyle().Foreground(theme.Muted).Render("○")
		}
		if i < totalSlides-1 {
			dots += " "
		}
	}

	// Navigation hints
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)
	accentStyle := lipgloss.NewStyle().Foreground(theme.Accent)

	var navHint string
	if m.currentSlide == 0 {
		navHint = accentStyle.Render("←") + mutedStyle.Render(" menu  ") +
			accentStyle.Render("→") + mutedStyle.Render(" next")
	} else if m.currentSlide == totalSlides-1 {
		navHint = accentStyle.Render("←") + mutedStyle.Render(" back  ") +
			accentStyle.Render("esc") + mutedStyle.Render(" menu")
	} else {
		navHint = accentStyle.Render("← →") + mutedStyle.Render(" navigate  ") +
			accentStyle.Render("esc") + mutedStyle.Render(" menu")
	}

	footer := lipgloss.JoinHorizontal(
		lipgloss.Center,
		dots,
		"    ",
		navHint,
	)

	return lipgloss.Place(m.width, 2, lipgloss.Center, lipgloss.Bottom, footer)
}

// ============================================================================
// "WHAT IS IRL?" SLIDES
// ============================================================================

func (m HelpModel) renderWhatIsIRLSlide() string {
	switch m.currentSlide {
	case 0:
		return m.renderTitleSlide()
	case 1:
		return m.renderProblemSlide()
	case 2:
		return m.renderSolutionSlide()
	case 3:
		return m.renderHowItWorksSlide()
	case 4:
		return m.renderGetStartedSlide()
	}
	return ""
}

func (m HelpModel) renderTitleSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	taglineStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	title := titleStyle.Render(`
 ██╗██████╗ ██╗
 ██║██╔══██╗██║
 ██║██████╔╝██║
 ██║██╔══██╗██║
 ██║██║  ██║███████╗
 ╚═╝╚═╝  ╚═╝╚══════╝`)

	tagline := taglineStyle.Render("Idempotent Research Loop")
	subtitle := accentStyle.Render("Document-centric AI workflows")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		tagline,
		"",
		subtitle,
	)
}

func (m HelpModel) renderProblemSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("The Problem")

	bullet1 := accentStyle.Render("Conversations disappear")
	desc1 := mutedStyle.Render("Can't rerun last week's analysis")

	bullet2 := accentStyle.Render("Results are unrepeatable")
	desc2 := mutedStyle.Render("Same question, different answer")

	bullet3 := accentStyle.Render("No documentation")
	desc3 := mutedStyle.Render("No audit trail, no provenance")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		bullet1,
		desc1,
		"",
		bullet2,
		desc2,
		"",
		bullet3,
		desc3,
	)
}

func (m HelpModel) renderSolutionSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	boxStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	title := titleStyle.Render("The Solution")
	intro := "Write a " + accentStyle.Render("plan") + ", not a chat"

	diagram := boxStyle.Render(`
  ╭───────────────────────╮
  │                       │
  │     main-plan.md      │
  │                       │
  │   Your single source  │
  │       of truth        │
  │                       │
  ╰───────────────────────╯`)

	tagline := mutedStyle.Render("The AI reads your plan. Every time.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		intro,
		diagram,
		"",
		tagline,
	)
}

func (m HelpModel) renderHowItWorksSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("How It Works")
	subtitle := mutedStyle.Render("The three-step loop")

	step1 := successStyle.Render(" 1 ") + accentStyle.Render("WRITE   ") + mutedStyle.Render("Edit your plan")
	step2 := successStyle.Render(" 2 ") + accentStyle.Render("EXECUTE ") + mutedStyle.Render("AI runs the plan")
	step3 := successStyle.Render(" 3 ") + accentStyle.Render("REVIEW  ") + mutedStyle.Render("Check outputs, commit")

	arrow := mutedStyle.Render("↻ repeat")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		subtitle,
		"",
		"",
		step1,
		"",
		step2,
		"",
		step3,
		"",
		"",
		arrow,
	)
}

func (m HelpModel) renderGetStartedSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Get Started")

	cta := successStyle.Render("Press ") + accentStyle.Render("n") + successStyle.Render(" to create your first project")

	tip1 := mutedStyle.Render("Pick a template to start from")
	tip2 := mutedStyle.Render("Edit main-plan.md with your goals")
	tip3 := mutedStyle.Render("Run the AI and review the results")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		cta,
		"",
		"",
		tip1,
		tip2,
		tip3,
	)
}
