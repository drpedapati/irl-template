package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/projects"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

var listJSONFlag bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List IRL projects in workspace",
	Long: `List all IRL projects in the configured workspace directory.

A folder is considered a project if it contains a main-plan.md file.

Examples:
  irl list          # Table output
  irl list --json   # JSON for agents`,
	Aliases: []string{"ls"},
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&listJSONFlag, "json", false, "Output as JSON")
}

func runList(cmd *cobra.Command, args []string) error {
	baseDir := config.GetDefaultDirectory()
	if baseDir == "" {
		return fmt.Errorf("no default directory configured (run 'irl config --dir ~/path' to set one)")
	}

	list, err := projects.ScanDir(baseDir)
	if err != nil {
		return fmt.Errorf("failed to scan projects: %w", err)
	}

	if listJSONFlag {
		data, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	if len(list) == 0 {
		fmt.Println(theme.Faint("No projects found in " + baseDir))
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n",
		theme.Faint("NAME"), theme.Faint("MODIFIED"), theme.Faint("PATH"))

	for _, p := range list {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			p.Name, smartDate(p.Modified), theme.Faint(p.Path))
	}
	w.Flush()

	fmt.Printf("\n%s %d projects in %s\n", theme.Faint("Total:"), len(list), baseDir)

	return nil
}

func smartDate(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < 24*time.Hour:
		return "today"
	case diff < 48*time.Hour:
		return "yesterday"
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	default:
		return t.Format("Jan 2, 2006")
	}
}
