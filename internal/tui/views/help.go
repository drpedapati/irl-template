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
	HelpSectionMenu          HelpSection = iota // The help menu
	HelpSectionWhatIsIRL                        // "What is IRL?" slide deck
	HelpSectionHowIRLHelps                      // "How IRL helps" slide deck
	HelpSectionWhoIsIRLFor                      // "Who is IRL for?" slide deck
	HelpSectionWhatCanYouBuild                  // "What can you build?" slide deck
	HelpSectionWhatYouNeed                      // "What you need" slide deck
	HelpSectionSeeItInAction                    // "See it in action" slide deck
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
	{
		title:       "How IRL helps",
		description: "A concrete example",
		duration:    "11 slides · 3 min",
		section:     HelpSectionHowIRLHelps,
	},
	{
		title:       "Who is IRL for?",
		description: "Find your use case",
		duration:    "6 slides · 2 min",
		section:     HelpSectionWhoIsIRLFor,
	},
	{
		title:       "What can you build?",
		description: "See example projects",
		duration:    "8 slides · 3 min",
		section:     HelpSectionWhatCanYouBuild,
	},
	{
		title:       "What you need",
		description: "Prerequisites and setup",
		duration:    "3 slides · 1 min",
		section:     HelpSectionWhatYouNeed,
	},
	{
		title:       "See it in action",
		description: "Watch a complete workflow",
		duration:    "6 slides · 2 min",
		section:     HelpSectionSeeItInAction,
	},
}

const whatIsIRLSlides = 5
const howIRLHelpsSlides = 11
const whoIsIRLForSlides = 6
const whatCanYouBuildSlides = 8
const whatYouNeedSlides = 3
const seeItInActionSlides = 6

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
	case HelpSectionHowIRLHelps:
		return howIRLHelpsSlides
	case HelpSectionWhoIsIRLFor:
		return whoIsIRLForSlides
	case HelpSectionWhatCanYouBuild:
		return whatCanYouBuildSlides
	case HelpSectionWhatYouNeed:
		return whatYouNeedSlides
	case HelpSectionSeeItInAction:
		return seeItInActionSlides
	default:
		return 0
	}
}

// SectionTitle returns the title of the current section
func (m HelpModel) SectionTitle() string {
	switch m.section {
	case HelpSectionWhatIsIRL:
		return "What is IRL?"
	case HelpSectionHowIRLHelps:
		return "How IRL helps"
	case HelpSectionWhoIsIRLFor:
		return "Who is IRL for?"
	case HelpSectionWhatCanYouBuild:
		return "What can you build?"
	case HelpSectionWhatYouNeed:
		return "What you need"
	case HelpSectionSeeItInAction:
		return "See it in action"
	default:
		return ""
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
	case "3":
		if len(helpMenuItems) >= 3 {
			m.section = helpMenuItems[2].section
			m.currentSlide = 0
		}
	case "4":
		if len(helpMenuItems) >= 4 {
			m.section = helpMenuItems[3].section
			m.currentSlide = 0
		}
	case "5":
		if len(helpMenuItems) >= 5 {
			m.section = helpMenuItems[4].section
			m.currentSlide = 0
		}
	case "6":
		if len(helpMenuItems) >= 6 {
			m.section = helpMenuItems[5].section
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
	case "6":
		if totalSlides > 5 {
			m.currentSlide = 5
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

	// Menu content (outer container provides left margin)
	menuContent := content.String()

	// Footer
	footer := m.renderMenuFooter()

	return menuContent + "\n" + footer
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

	// Left-justified (outer container provides margin)
	return hint + "\n"
}

// ============================================================================
// SLIDE DECK RENDERER
// ============================================================================

func (m HelpModel) renderSlides() string {
	var content string

	switch m.section {
	case HelpSectionWhatIsIRL:
		content = m.renderWhatIsIRLSlide()
	case HelpSectionHowIRLHelps:
		content = m.renderHowIRLHelpsSlide()
	case HelpSectionWhoIsIRLFor:
		content = m.renderWhoIsIRLForSlide()
	case HelpSectionWhatCanYouBuild:
		content = m.renderWhatCanYouBuildSlide()
	case HelpSectionWhatYouNeed:
		content = m.renderWhatYouNeedSlide()
	case HelpSectionSeeItInAction:
		content = m.renderSeeItInActionSlide()
	}

	// Header with lesson title and slide counter
	header := m.renderSlideHeader()

	// Fixed heights
	const headerHeight = 2
	const footerHeight = 2

	// Content area fills everything except header and footer
	contentHeight := m.height - headerHeight - footerHeight
	if contentHeight < 8 {
		contentHeight = 8
	}

	// Center content in its dedicated area
	centeredContent := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	// Footer with navigation
	footer := m.renderSlideFooterFixed()

	return header + centeredContent + footer
}

func (m HelpModel) renderSlideHeader() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	counterStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	title := titleStyle.Render(m.SectionTitle())
	counter := counterStyle.Render("(" + string('0'+byte(m.currentSlide+1)) + "/" + string('0'+byte(m.TotalSlides())) + ")")

	// Handle slide counts > 9
	if m.TotalSlides() > 9 || m.currentSlide+1 > 9 {
		counter = counterStyle.Render("(" + slideItoa(m.currentSlide+1) + "/" + slideItoa(m.TotalSlides()) + ")")
	}

	header := title + "  " + counter

	// No extra padding - outer container (tui.go) provides left margin
	return header + "\n"
}

// slideItoa converts small ints to string for slide counters
func slideItoa(n int) string {
	if n < 10 {
		return string('0' + byte(n))
	}
	return string('0'+byte(n/10)) + string('0'+byte(n%10))
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

// renderSlideFooterFixed returns footer with exactly 2 lines height
func (m HelpModel) renderSlideFooterFixed() string {
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

	footer := dots + "    " + navHint

	// No extra padding - outer container (tui.go) provides left margin
	return footer + "\n"
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

// ============================================================================
// "HOW IRL HELPS" SLIDES
// ============================================================================

func (m HelpModel) renderHowIRLHelpsSlide() string {
	switch m.currentSlide {
	case 0:
		return m.renderHowIRLHelpsTitleSlide()
	case 1:
		return m.renderTheTaskSlide()
	case 2:
		return m.renderTraditionalWaySlide()
	case 3:
		return m.renderChatWaySlide()
	case 4:
		return m.renderWriteYourPlanSlide()
	case 5:
		return m.renderYourMethodologySlide()
	case 6:
		return m.renderPointToSourcesSlide()
	case 7:
		return m.renderAIDoesTheWorkSlide()
	case 8:
		return m.renderSeeYourOutputsSlide()
	case 9:
		return m.renderIterateSlide()
	case 10:
		return m.renderTheDifferenceSlide()
	}
	return ""
}

func (m HelpModel) renderHowIRLHelpsTitleSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("How IRL Helps")
	subtitle := subtitleStyle.Render("A real example, three approaches")
	hint := mutedStyle.Render("See why document-centric matters")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		subtitle,
		"",
		hint,
	)
}

func (m HelpModel) renderTheTaskSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Accent).
		Padding(1, 3)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("The Task")

	taskContent := lipgloss.JoinVertical(
		lipgloss.Center,
		accentStyle.Render("Compare 2 papers on sleep & memory"),
		"",
		mutedStyle.Render("One PDF you have, one article online"),
		mutedStyle.Render("Write a short synthesis"),
	)

	box := boxStyle.Render(taskContent)

	subtitle := mutedStyle.Render("Let's see three ways to do this...")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		subtitle,
	)
}

func (m HelpModel) renderTraditionalWaySlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Warning).
		Padding(1, 2)

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Your Method Is Sound")
	label := labelStyle.Render("The challenge isn't the method—it's scale")

	points := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render("You have a systematic approach"),
		mutedStyle.Render("  ...but applying it to 20 papers is exhausting"),
		"",
		accentStyle.Render("You know exactly what to extract"),
		mutedStyle.Render("  ...but consistency fades under fatigue"),
		"",
		accentStyle.Render("Deadlines arrive"),
		mutedStyle.Render("  ...and your careful method gets compromised"),
	)

	box := boxStyle.Render(points)

	insight := mutedStyle.Render("The method works. Scaling it is the hard part.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		label,
		"",
		box,
		"",
		insight,
	)
}

func (m HelpModel) renderChatWaySlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Warning).
		Padding(1, 2)

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.Warning).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	errorStyle := lipgloss.NewStyle().
		Foreground(theme.Error)

	title := titleStyle.Render("Chat Interfaces")
	label := labelStyle.Render("You adapt to it—not the other way around")

	points := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render("Your method doesn't transfer"),
		mutedStyle.Render("  You work the way the chat wants"),
		"",
		accentStyle.Render("Each conversation starts fresh"),
		mutedStyle.Render("  No memory of your approach"),
		"",
		accentStyle.Render("Results vanish"),
		mutedStyle.Render("  Can't re-run, can't share, can't cite"),
	)

	box := boxStyle.Render(points)

	insight := errorStyle.Render("You lose your methodology") + mutedStyle.Render(" when you use chat.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		label,
		"",
		box,
		"",
		insight,
	)
}

func (m HelpModel) renderWriteYourPlanSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	// Code-like box for the plan
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Success).
		Padding(0, 2)

	commentStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	headingStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	title := titleStyle.Render("IRL: Write Your Plan")
	label := labelStyle.Render("It's just English in a markdown file")

	// Show actual plan content - emphasize it's natural language
	plan := lipgloss.JoinVertical(
		lipgloss.Left,
		headingStyle.Render("## Goal"),
		textStyle.Render("Compare these 2 papers on sleep and"),
		textStyle.Render("memory. Find what they agree on."),
		"",
		headingStyle.Render("## What I Want"),
		textStyle.Render("A 1-page synthesis with:"),
		textStyle.Render("- Summary of each paper"),
		textStyle.Render("- Common themes"),
		textStyle.Render("- Key differences"),
	)

	box := boxStyle.Render(plan)

	hint := commentStyle.Render("← No special syntax. Just describe what you want.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		label,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderYourMethodologySlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent).
		Padding(0, 2)

	headingStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	title := titleStyle.Render("Go As Deep As You Need")
	label := labelStyle.Render("\"Summary of each paper\" can become...")

	// Show expanded, detailed instructions
	plan := lipgloss.JoinVertical(
		lipgloss.Left,
		headingStyle.Render("## For each paper, extract:"),
		detailStyle.Render("- Full citation (APA 7th)"),
		detailStyle.Render("- Study design & sample size"),
		detailStyle.Render("- Statistical methods used"),
		detailStyle.Render("- Effect sizes with 95% CI"),
		detailStyle.Render("- Limitations noted by authors"),
		detailStyle.Render("- Direct quotes with page numbers"),
	)

	box := boxStyle.Render(plan)

	hint := mutedStyle.Render("This is ") + successStyle.Render("your") + mutedStyle.Render(" methodology. The AI follows ") + successStyle.Render("your") + mutedStyle.Render(" standards.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		label,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderPointToSourcesSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent).
		Padding(0, 2)

	headingStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	pathStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	urlStyle := lipgloss.NewStyle().
		Foreground(theme.Warning)

	commentStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("IRL: Point to Your Sources")

	// Show how to reference files and URLs
	plan := lipgloss.JoinVertical(
		lipgloss.Left,
		headingStyle.Render("## Sources"),
		"",
		commentStyle.Render("# A PDF in your project:"),
		pathStyle.Render("02-data/papers/walker-2017.pdf"),
		"",
		commentStyle.Render("# An article online:"),
		urlStyle.Render("https://nature.com/articles/sleep-study"),
	)

	box := boxStyle.Render(plan)

	hint := commentStyle.Render("Files, URLs, or both. The AI reads them all.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderAIDoesTheWorkSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Success).
		Padding(1, 3)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	title := titleStyle.Render("IRL: You Direct, AI Executes")
	label := labelStyle.Render("Your plan is the control surface")

	steps := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render("Your plan says what to do"),
		mutedStyle.Render("  AI follows your instructions—nothing more"),
		"",
		accentStyle.Render("Your sources define the scope"),
		mutedStyle.Render("  AI only sees what you provide"),
		"",
		accentStyle.Render("Your outputs prove compliance"),
		mutedStyle.Render("  Verify the AI did exactly what you asked"),
	)

	box := boxStyle.Render(steps)

	result := successStyle.Render("You stay in control. The outputs prove it.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		label,
		"",
		box,
		"",
		result,
	)
}

func (m HelpModel) renderSeeYourOutputsSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Success).
		Padding(1, 2)

	folderStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	fileStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	familiarStyle := lipgloss.NewStyle().
		Foreground(theme.Warning)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("IRL: See Your Outputs")

	tree := lipgloss.JoinVertical(
		lipgloss.Left,
		folderStyle.Render("03-outputs/"),
		fileStyle.Render("  synthesis.md")+"      "+mutedStyle.Render("← Working draft"),
		fileStyle.Render("  paper1-notes.md")+"   "+mutedStyle.Render("← Extracted notes"),
		"",
		mutedStyle.Render("  Renders to formats you know:"),
		familiarStyle.Render("  synthesis.docx")+"     "+mutedStyle.Render("← Word"),
		familiarStyle.Render("  synthesis.pdf")+"      "+mutedStyle.Render("← PDF"),
		familiarStyle.Render("  slides.pptx")+"        "+mutedStyle.Render("← PowerPoint"),
	)

	box := boxStyle.Render(tree)

	hint := mutedStyle.Render("Quarto renders to Word, PDF, PowerPoint, and more.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderIterateSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent).
		Padding(1, 2)

	headingStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	oldStyle := lipgloss.NewStyle().
		Foreground(theme.Error).
		Strikethrough(true)

	newStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	title := titleStyle.Render("IRL: Need Changes? Just Edit")

	// Show editing the plan
	plan := lipgloss.JoinVertical(
		lipgloss.Left,
		headingStyle.Render("## Sources"),
		mutedStyle.Render("02-data/papers/walker-2017.pdf"),
		mutedStyle.Render("https://nature.com/articles/..."),
		newStyle.Render("+ 02-data/papers/diekelmann-2010.pdf"),
		"",
		headingStyle.Render("## What I Want"),
		oldStyle.Render("A 1-page synthesis"),
		newStyle.Render("A 2-page synthesis with methodology"),
	)

	box := boxStyle.Render(plan)

	hint := successStyle.Render("Run again → New outputs → Same structure")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderTheDifferenceSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	// Comparison table with borders
	headerStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	errorStyle := lipgloss.NewStyle().
		Foreground(theme.Error)

	warningStyle := lipgloss.NewStyle().
		Foreground(theme.Warning)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("The Difference")

	// Simple comparison
	header := "              " + headerStyle.Render("Traditional") + "    " + headerStyle.Render("Chat") + "    " + headerStyle.Render("IRL")

	row1 := mutedStyle.Render("Repeatable    ") +
		errorStyle.Render("    no") + "         " +
		warningStyle.Render("no") + "      " +
		successStyle.Render("yes")

	row2 := mutedStyle.Render("Versioned     ") +
		errorStyle.Render("    no") + "         " +
		warningStyle.Render("no") + "      " +
		successStyle.Render("yes")

	row3 := mutedStyle.Render("Documented    ") +
		errorStyle.Render("    no") + "         " +
		warningStyle.Render("no") + "      " +
		successStyle.Render("yes")

	row4 := mutedStyle.Render("Shareable     ") +
		warningStyle.Render("   hard") + "       " +
		warningStyle.Render("no") + "      " +
		successStyle.Render("yes")

	takeaway := successStyle.Render("IRL = Reproducible AI workflows")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		header,
		"",
		row1,
		row2,
		row3,
		row4,
		"",
		"",
		takeaway,
	)
}

// ============================================================================
// "WHO IS IRL FOR?" SLIDES
// ============================================================================

func (m HelpModel) renderWhoIsIRLForSlide() string {
	switch m.currentSlide {
	case 0:
		return m.renderWhoIsIRLForTitleSlide()
	case 1:
		return m.renderTaskContinuumSlide()
	case 2:
		return m.renderResearchersSlide()
	case 3:
		return m.renderStudentsSlide()
	case 4:
		return m.renderProfessionalsSlide()
	case 5:
		return m.renderWhoIsIRLForSummarySlide()
	}
	return ""
}

func (m HelpModel) renderWhoIsIRLForTitleSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Who is IRL for?")
	subtitle := subtitleStyle.Render("A responsible way to bring AI into your work")

	personas := mutedStyle.Render("Researchers  ·  Students  ·  Professionals")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		subtitle,
		"",
		personas,
	)
}

func (m HelpModel) renderTaskContinuumSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent).
		Padding(1, 2)

	smallStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	medStyle := lipgloss.NewStyle().
		Foreground(theme.Warning)

	largeStyle := lipgloss.NewStyle().
		Foreground(theme.Primary)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	arrowStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	title := titleStyle.Render("From Small Tasks to Large Projects")

	// Visual continuum
	continuum := lipgloss.JoinVertical(
		lipgloss.Center,
		smallStyle.Render("    Small    ")+"     "+medStyle.Render("    Medium    ")+"     "+largeStyle.Render("    Large    "),
		arrowStyle.Render("      │       ")+"     "+arrowStyle.Render("       │       ")+"     "+arrowStyle.Render("      │      "),
		arrowStyle.Render("      ▼       ")+"     "+arrowStyle.Render("       ▼       ")+"     "+arrowStyle.Render("      ▼      "),
		"",
		smallStyle.Render("  Abstract   ")+"     "+medStyle.Render("  Lit Review   ")+"     "+largeStyle.Render("   Grant    "),
		smallStyle.Render("   Email     ")+"     "+medStyle.Render("   Analysis    ")+"     "+largeStyle.Render("  Report    "),
		smallStyle.Render("   Summary   ")+"     "+medStyle.Render("  Data Viz     ")+"     "+largeStyle.Render("   Thesis   "),
	)

	box := boxStyle.Render(continuum)

	hint := mutedStyle.Render("Same responsible harness. Any size task.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderResearchersSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		Padding(1, 2)

	smallStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	medStyle := lipgloss.NewStyle().
		Foreground(theme.Warning)

	largeStyle := lipgloss.NewStyle().
		Foreground(theme.Primary)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	header := titleStyle.Render("Researchers & Scientists")

	cardContent := lipgloss.JoinVertical(
		lipgloss.Left,
		smallStyle.Render("Small")+"    "+mutedStyle.Render("Conference abstract, paper summary"),
		medStyle.Render("Medium")+"   "+mutedStyle.Render("Literature review, methods section"),
		largeStyle.Render("Large")+"    "+mutedStyle.Render("Grant proposal, systematic review"),
	)

	card := cardStyle.Render(cardContent)

	hint := mutedStyle.Render("Your methodology, applied consistently at any scale")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		card,
		"",
		hint,
	)
}

func (m HelpModel) renderStudentsSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent).
		Padding(1, 2)

	smallStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	medStyle := lipgloss.NewStyle().
		Foreground(theme.Warning)

	largeStyle := lipgloss.NewStyle().
		Foreground(theme.Primary)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	header := titleStyle.Render("Graduate Students")

	cardContent := lipgloss.JoinVertical(
		lipgloss.Left,
		smallStyle.Render("Small")+"    "+mutedStyle.Render("Reading notes, weekly summary"),
		medStyle.Render("Medium")+"   "+mutedStyle.Render("Course paper, qualifying exam prep"),
		largeStyle.Render("Large")+"    "+mutedStyle.Render("Thesis chapter, dissertation"),
	)

	card := cardStyle.Render(cardContent)

	hint := mutedStyle.Render("Build your portfolio with a clear audit trail")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		card,
		"",
		hint,
	)
}

func (m HelpModel) renderProfessionalsSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Success).
		Padding(1, 2)

	smallStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	medStyle := lipgloss.NewStyle().
		Foreground(theme.Warning)

	largeStyle := lipgloss.NewStyle().
		Foreground(theme.Primary)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	header := titleStyle.Render("Industry Professionals")

	cardContent := lipgloss.JoinVertical(
		lipgloss.Left,
		smallStyle.Render("Small")+"    "+mutedStyle.Render("Status update, meeting notes"),
		medStyle.Render("Medium")+"   "+mutedStyle.Render("Technical report, analysis deck"),
		largeStyle.Render("Large")+"    "+mutedStyle.Render("White paper, regulatory filing"),
	)

	card := cardStyle.Render(cardContent)

	hint := mutedStyle.Render("Transparent AI use for compliance and handoff")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		card,
		"",
		hint,
	)
}

func (m HelpModel) renderWhoIsIRLForSummarySlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Success).
		Padding(1, 3)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("The IRL Promise")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render("Responsible"),
		mutedStyle.Render("  AI that follows your instructions"),
		"",
		accentStyle.Render("Transparent"),
		mutedStyle.Render("  See exactly what was done"),
		"",
		accentStyle.Render("Auditable"),
		mutedStyle.Render("  Every version preserved in git"),
	)

	box := boxStyle.Render(content)

	tagline := successStyle.Render("Reproducible AI-assisted research.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		tagline,
	)
}

// ============================================================================
// "WHAT CAN YOU BUILD?" SLIDES
// ============================================================================

func (m HelpModel) renderWhatCanYouBuildSlide() string {
	switch m.currentSlide {
	case 0:
		return m.renderWhatCanYouBuildTitleSlide()
	case 1:
		return m.renderUniversalWorkflowSlide()
	case 2:
		return m.renderBeforeAfterRunSlide()
	case 3:
		return m.renderTemplateSetupSlide()
	case 4:
		return m.renderLitReviewSlide()
	case 5:
		return m.renderDataAnalysisSlide()
	case 6:
		return m.renderCodeProjectSlide()
	case 7:
		return m.renderWritingProjectSlide()
	}
	return ""
}

func (m HelpModel) renderWhatCanYouBuildTitleSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("What can you build?")
	subtitle := subtitleStyle.Render("Templates for every workflow")

	hint := mutedStyle.Render("All projects share a common structure")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		subtitle,
		"",
		hint,
	)
}

func (m HelpModel) renderUniversalWorkflowSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent).
		Padding(1, 2)

	phaseStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	arrowStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	title := titleStyle.Render("How Professionals Work")
	label := labelStyle.Render("Every project follows this pattern")

	workflow := lipgloss.JoinVertical(
		lipgloss.Center,
		phaseStyle.Render("    PLAN    ")+"  "+arrowStyle.Render("→")+"  "+phaseStyle.Render("   EXECUTE   ")+"  "+arrowStyle.Render("→")+"  "+phaseStyle.Render("   REVIEW   "),
		"",
		mutedStyle.Render(" Define goals ")+"     "+mutedStyle.Render("  Do the work  ")+"     "+mutedStyle.Render(" Check quality "),
		mutedStyle.Render("  Set method  ")+"     "+mutedStyle.Render(" Follow method ")+"     "+mutedStyle.Render(" Format output "),
		mutedStyle.Render("List sources  ")+"     "+mutedStyle.Render("Create outputs ")+"     "+mutedStyle.Render("   Finalize   "),
	)

	box := boxStyle.Render(workflow)

	hint := mutedStyle.Render("IRL templates encode this pattern.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		label,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderBeforeAfterRunSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Success).
		Padding(1, 2)

	beforeStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	afterStyle := lipgloss.NewStyle().
		Foreground(theme.Warning).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Before Each Run & After Each Run")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		beforeStyle.Render("Before Each Run")+"  "+mutedStyle.Render("(start of each loop)"),
		mutedStyle.Render("  \"Read each paper's methods section first\""),
		mutedStyle.Render("  \"Use APA 7th citation format\""),
		mutedStyle.Render("  \"Focus on quantitative findings only\""),
		"",
		afterStyle.Render("After Each Run")+"   "+mutedStyle.Render("(end of each loop)"),
		mutedStyle.Render("  \"Summarize key themes in a table\""),
		mutedStyle.Render("  \"Flag any contradictory findings\""),
		mutedStyle.Render("  \"Export to Word format\""),
	)

	box := boxStyle.Render(content)

	hint := mutedStyle.Render("Your standards, applied every loop.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderTemplateSetupSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Accent).
		Padding(1, 2)

	sectionStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("You Define the Setup")
	label := labelStyle.Render("Every plan has these sections—you fill in the details")

	sections := lipgloss.JoinVertical(
		lipgloss.Left,
		sectionStyle.Render("## Setup")+"          "+mutedStyle.Render("← Folders you need (you choose)"),
		sectionStyle.Render("## Before Each Run")+" "+mutedStyle.Render("← What happens before each loop"),
		sectionStyle.Render("## Approach")+"        "+mutedStyle.Render("← Your methodology and goals"),
		sectionStyle.Render("## Tasks")+"           "+mutedStyle.Render("← The work to be done"),
		sectionStyle.Render("## After Each Run")+"  "+mutedStyle.Render("← What happens after each loop"),
	)

	box := boxStyle.Render(sections)

	hint := mutedStyle.Render("Lit review needs abstracts/. Code project needs repo/. You decide.")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		label,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderLitReviewSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	// Double border for emphasis
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Primary).
		Padding(1, 2)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Literature Review")

	diagram := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render("02-data/"),
		mutedStyle.Render("  papers/      ← PDFs and sources"),
		mutedStyle.Render("  extracted/   ← Key quotes"),
		"",
		accentStyle.Render("03-outputs/"),
		mutedStyle.Render("  synthesis.md ← Your analysis"),
		mutedStyle.Render("  themes.md    ← Common threads"),
	)

	box := boxStyle.Render(diagram)

	hint := successStyle.Render("Template: lit-review")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderDataAnalysisSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Accent).
		Padding(1, 2)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Data Analysis")

	diagram := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render("02-data/"),
		mutedStyle.Render("  raw/         ← Original data"),
		mutedStyle.Render("  cleaned/     ← Processed data"),
		"",
		accentStyle.Render("03-outputs/"),
		mutedStyle.Render("  figures/     ← Visualizations"),
		mutedStyle.Render("  report.md    ← Findings"),
	)

	box := boxStyle.Render(diagram)

	hint := successStyle.Render("Template: data-analysis")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderCodeProjectSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Success).
		Padding(1, 2)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Code Project")

	diagram := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render("repo/"),
		mutedStyle.Render("  src/         ← Your code"),
		mutedStyle.Render("  tests/       ← Test suite"),
		"",
		accentStyle.Render("notes/"),
		mutedStyle.Render("  arch.md      ← Architecture"),
		mutedStyle.Render("  decisions.md ← ADRs"),
	)

	box := boxStyle.Render(diagram)

	hint := successStyle.Render("Template: github-clone")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

func (m HelpModel) renderWritingProjectSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Warning).
		Padding(1, 2)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Writing Projects")

	diagram := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render("01-plans/"),
		mutedStyle.Render("  outline.md   ← Setup"),
		mutedStyle.Render("  notes.md     ← Research"),
		"",
		accentStyle.Render("03-outputs/"),
		mutedStyle.Render("  draft.md     ← Working draft"),
		mutedStyle.Render("  final.pdf    ← Rendered output"),
	)

	box := boxStyle.Render(diagram)

	hint := successStyle.Render("Template: meeting-abstract")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
		"",
		hint,
	)
}

// ============================================================================
// "WHAT YOU NEED" SLIDES
// ============================================================================

func (m HelpModel) renderWhatYouNeedSlide() string {
	switch m.currentSlide {
	case 0:
		return m.renderWhatYouNeedTitleSlide()
	case 1:
		return m.renderRequirementsSlide()
	case 2:
		return m.renderReadyToStartSlide()
	}
	return ""
}

func (m HelpModel) renderWhatYouNeedTitleSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("What you need")
	subtitle := subtitleStyle.Render("Simple requirements")
	hint := mutedStyle.Render("Most things are already on your Mac")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		subtitle,
		"",
		hint,
	)
}

func (m HelpModel) renderRequirementsSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	// Thick border for requirements
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(theme.Success).
		Padding(1, 3)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Requirements")

	check := successStyle.Render("[✓]")

	reqs := lipgloss.JoinVertical(
		lipgloss.Left,
		check+" "+accentStyle.Render("Git"),
		mutedStyle.Render("    For version control"),
		"",
		check+" "+accentStyle.Render("Text editor"),
		mutedStyle.Render("    VS Code, Cursor, or any editor"),
		"",
		check+" "+accentStyle.Render("AI assistant"),
		mutedStyle.Render("    Claude Code, Cursor, or similar"),
	)

	box := boxStyle.Render(reqs)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		box,
	)
}

func (m HelpModel) renderReadyToStartSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	// Glowing effect with nested borders
	outerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Muted).
		Padding(0, 1)

	innerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Success).
		Padding(1, 3)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	title := titleStyle.Render("Ready to start!")

	innerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		successStyle.Render("You're all set"),
		"",
		"Press "+accentStyle.Render("n")+" to create",
		"your first project",
	)

	inner := innerStyle.Render(innerContent)
	outer := outerStyle.Render(inner)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		outer,
	)
}

// ============================================================================
// "SEE IT IN ACTION" SLIDES
// ============================================================================

func (m HelpModel) renderSeeItInActionSlide() string {
	switch m.currentSlide {
	case 0:
		return m.renderSeeItInActionTitleSlide()
	case 1:
		return m.renderSeeItInActionStep1Slide()
	case 2:
		return m.renderSeeItInActionStep2Slide()
	case 3:
		return m.renderSeeItInActionStep3Slide()
	case 4:
		return m.renderSeeItInActionStep4Slide()
	case 5:
		return m.renderSeeItInActionStep5Slide()
	}
	return ""
}

func (m HelpModel) renderSeeItInActionTitleSlide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	title := titleStyle.Render("See it in action")
	subtitle := subtitleStyle.Render("A complete workflow in 5 steps")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		subtitle,
	)
}

func (m HelpModel) renderSeeItInActionStep1Slide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	title := titleStyle.Render("Step 1: Write your plan")

	diagram := boxStyle.Render(`
┌─ main-plan.md ─────────────┐
│ ## Goal                    │
│ Summarize 3 papers on X    │
│                            │
│ ## Tasks                   │
│ - [ ] Read abstracts       │
│ - [ ] Extract key points   │
└────────────────────────────┘`)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		diagram,
	)
}

func (m HelpModel) renderSeeItInActionStep2Slide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Step 2: AI executes")

	bullet1 := accentStyle.Render("AI reads your plan")
	bullet2 := accentStyle.Render("Creates files in outputs/")
	bullet3 := accentStyle.Render("Follows your instructions")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		bullet1,
		mutedStyle.Render("Your plan is the control surface"),
		"",
		bullet2,
		mutedStyle.Render("Artifacts are generated automatically"),
		"",
		bullet3,
		mutedStyle.Render("Exactly as you specified"),
	)
}

func (m HelpModel) renderSeeItInActionStep3Slide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success)

	warningStyle := lipgloss.NewStyle().
		Foreground(theme.Accent)

	title := titleStyle.Render("Step 3: Review the diff")

	diagram := successStyle.Render("+ Added: summary.md") + "\n" +
		successStyle.Render("+ Added: key_themes.md") + "\n" +
		warningStyle.Render("M Modified: main-plan.md")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		diagram,
	)
}

func (m HelpModel) renderSeeItInActionStep4Slide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	title := titleStyle.Render("Step 4: Commit your work")

	bullet1 := accentStyle.Render("Save a checkpoint")
	bullet2 := accentStyle.Render("Can always go back")
	bullet3 := accentStyle.Render("Your audit trail")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		bullet1,
		mutedStyle.Render("git commit -m \"iteration 1\""),
		"",
		bullet2,
		mutedStyle.Render("Every version is preserved"),
		"",
		bullet3,
		mutedStyle.Render("Complete history of your work"),
	)
}

func (m HelpModel) renderSeeItInActionStep5Slide() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	accentStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Muted)

	successStyle := lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	title := titleStyle.Render("Step 5: Repeat")

	bullet1 := accentStyle.Render("Refine your plan")
	bullet2 := accentStyle.Render("Run again")
	bullet3 := accentStyle.Render("Build iteratively")

	loop := successStyle.Render("↻")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		"",
		bullet1,
		mutedStyle.Render("Add details, fix issues"),
		"",
		bullet2,
		mutedStyle.Render("AI picks up where you left off"),
		"",
		bullet3,
		mutedStyle.Render("Each iteration improves the output"),
		"",
		"",
		loop,
	)
}
