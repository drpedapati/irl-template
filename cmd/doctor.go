package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

type tool struct {
	name    string
	cmd     string
	install string
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment and show recommendations",
	Long:  "Check for required tools, AI assistants, IDEs, and system info",
	Run:   runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) {
	theme.Section("Environment")

	// System info - single line
	printSystemInfoCompact()

	// Docs link
	fmt.Printf("  %s %s\n", theme.Faint("Docs:"), theme.Cmd("https://www.irloop.org"))

	// Two-column layout for tools
	coreTools := []tool{
		{"Git", "git", "brew install git"},
		{"Quarto", "quarto", "brew install --cask quarto"},
		{"R", "R", "brew install r"},
		{"Python", "python3", "brew install python"},
	}

	aiTools := []tool{
		{"Claude Code", "claude", "npm i -g @anthropic-ai/claude-code"},
		{"Aider", "aider", "pip install aider-chat"},
		{"Copilot CLI", "gh copilot", "gh extension install github/gh-copilot"},
		{"Ollama", "ollama", "brew install ollama"},
	}

	ideTools := []tool{
		{"Positron", "positron", "brew install --cask positron"},
		{"VS Code", "code", "brew install --cask visual-studio-code"},
		{"Cursor", "cursor", "brew install --cask cursor"},
		{"RStudio", "rstudio", "brew install --cask rstudio"},
	}

	sandboxTools := []tool{
		{"Docker", "docker", "brew install --cask docker"},
	}

	// Print two columns: Core Tools | AI Assistants
	fmt.Println()
	printColumnHeaders("Core Tools", "AI Assistants")
	printTwoColumns(coreTools, aiTools)

	// Print two columns: IDEs | Sandbox
	fmt.Println()
	printColumnHeaders("IDEs", "Sandbox")
	printTwoColumns(ideTools, sandboxTools)

	// Sandbox hint
	fmt.Println()
	if checkCmd("docker") {
		fmt.Printf("  %s %s\n",
			theme.Faint("Tip:"),
			theme.Cmd("docker sandbox run claude"))
	}
	fmt.Println()
}

func printSystemInfoCompact() {
	var parts []string

	// Platform
	parts = append(parts, fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))

	// CPU
	parts = append(parts, fmt.Sprintf("%d cores", runtime.NumCPU()))

	// Memory
	if runtime.GOOS == "darwin" {
		if out, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
			memBytes := strings.TrimSpace(string(out))
			if len(memBytes) > 9 {
				parts = append(parts, memBytes[:len(memBytes)-9]+" GB")
			}
		}
	} else if runtime.GOOS == "linux" {
		if out, err := exec.Command("free", "-g").Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 1 {
				fields := strings.Fields(lines[1])
				if len(fields) > 1 {
					parts = append(parts, fields[1]+" GB")
				}
			}
		}
	}

	// Disk
	if out, err := exec.Command("df", "-h", ".").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 4 {
				parts = append(parts, fields[3]+" free")
			}
		}
	}

	fmt.Printf("  %s\n", theme.Faint(strings.Join(parts, " · ")))
}

func printColumnHeaders(left, right string) {
	// Use Lip Gloss for fixed-width columns
	leftStyle := lipgloss.NewStyle().Width(30).Foreground(theme.Muted)
	rightStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	fmt.Printf("  %s  %s\n", leftStyle.Render(left), rightStyle.Render(right))
}

func printTwoColumns(left, right []tool) {
	maxRows := len(left)
	if len(right) > maxRows {
		maxRows = len(right)
	}

	// Create a fixed-width style for the left column
	leftColStyle := lipgloss.NewStyle().Width(30)

	for i := 0; i < maxRows; i++ {
		leftStr := ""
		rightStr := ""

		if i < len(left) {
			leftStr = formatToolCheck(left[i])
		}
		if i < len(right) {
			rightStr = formatToolCheck(right[i])
		}

		// Use Lip Gloss width for proper alignment
		fmt.Printf("  %s  %s\n", leftColStyle.Render(leftStr), rightStr)
	}
}

func formatToolCheck(t tool) string {
	return theme.ToolCheck(t.name, checkTool(t))
}

func checkTool(t tool) bool {
	switch t.cmd {
	case "gh copilot":
		return checkGhCopilot()
	case "positron":
		return checkApp("Positron")
	case "cursor":
		return checkApp("Cursor")
	case "rstudio":
		return checkApp("RStudio")
	default:
		return checkCmd(t.cmd)
	}
}

func checkCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func checkGhCopilot() bool {
	out, err := exec.Command("gh", "extension", "list").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "copilot")
}

func checkApp(name string) bool {
	if runtime.GOOS != "darwin" {
		return checkCmd(strings.ToLower(name))
	}
	paths := []string{
		filepath.Join("/Applications", name+".app"),
		filepath.Join(os.Getenv("HOME"), "Applications", name+".app"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}
