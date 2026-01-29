package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│                     IRL Doctor                              │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")

	// System Info
	printSystemInfo()

	// Core tools
	printSection("Core Tools")
	coreTools := []tool{
		{"Git", "git", "brew install git"},
		{"Quarto", "quarto", "brew install --cask quarto"},
		{"R", "R", "brew install r"},
		{"Python", "python3", "brew install python"},
	}
	printToolTable(coreTools)

	// AI tools
	printSection("AI Assistants")
	aiTools := []tool{
		{"Claude Code", "claude", "npm i -g @anthropic-ai/claude-code"},
		{"Aider", "aider", "pip install aider-chat"},
		{"Copilot CLI", "gh copilot", "gh extension install github/gh-copilot"},
		{"Ollama", "ollama", "brew install ollama"},
	}
	printToolTable(aiTools)

	// IDEs
	printSection("IDEs")
	ideTools := []tool{
		{"Positron", "positron", "brew install --cask positron"},
		{"VS Code", "code", "brew install --cask visual-studio-code"},
		{"Cursor", "cursor", "brew install --cask cursor"},
		{"RStudio", "rstudio", "brew install --cask rstudio"},
	}
	printToolTable(ideTools)

	// Sandbox
	printSection("Sandbox")
	sandboxTools := []tool{
		{"Docker", "docker", "brew install --cask docker"},
	}
	printToolTable(sandboxTools)

	// Recommendations box
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Sandbox Commands                                            │")
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	if checkCmd("docker") {
		fmt.Println("│  docker sandbox run claude     # Claude in container        │")
		fmt.Println("│  docker sandbox run aider      # Aider in container         │")
	} else {
		fmt.Println("│  Install Docker, then:  docker sandbox run claude          │")
	}
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
}

func printSystemInfo() {
	fmt.Println()
	fmt.Println("System")
	fmt.Println(strings.Repeat("─", 61))

	// OS & Arch
	fmt.Printf("%-20s %s/%s\n", "Platform", runtime.GOOS, runtime.GOARCH)

	// CPU cores
	fmt.Printf("%-20s %d cores\n", "CPU", runtime.NumCPU())

	// Memory (macOS specific)
	if runtime.GOOS == "darwin" {
		if out, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
			memBytes := strings.TrimSpace(string(out))
			if len(memBytes) > 9 {
				gb := memBytes[:len(memBytes)-9]
				fmt.Printf("%-20s %s GB\n", "Memory", gb)
			}
		}
	} else if runtime.GOOS == "linux" {
		if out, err := exec.Command("free", "-g").Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 1 {
				fields := strings.Fields(lines[1])
				if len(fields) > 1 {
					fmt.Printf("%-20s %s GB\n", "Memory", fields[1])
				}
			}
		}
	}

	// Disk space of current directory
	if out, err := exec.Command("df", "-h", ".").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 4 {
				fmt.Printf("%-20s %s free of %s\n", "Disk", fields[3], fields[1])
			}
		}
	}

	// Current directory
	if wd, err := os.Getwd(); err == nil {
		if len(wd) > 40 {
			wd = "..." + wd[len(wd)-37:]
		}
		fmt.Printf("%-20s %s\n", "Directory", wd)
	}
}

func printSection(title string) {
	fmt.Println()
	fmt.Printf("%-20s %-6s %s\n", title, "", "Install")
	fmt.Println(strings.Repeat("─", 61))
}

func printToolTable(tools []tool) {
	for _, t := range tools {
		status := "✗"
		hint := t.install
		if checkTool(t) {
			status = "✓"
			hint = ""
		}
		if len(hint) > 36 {
			hint = hint[:33] + "..."
		}
		fmt.Printf("%-20s %-6s %s\n", t.name, status, hint)
	}
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
