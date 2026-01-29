package cmd

import (
	"fmt"

	"github.com/drpedapati/irl-template/pkg/style"
	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update templates from GitHub",
	Long:  "Fetch the latest templates from the IRL template repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%sFetching templates...%s\n", style.Dim, style.Reset)

		if err := templates.Update(); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		// Show available templates
		list, _ := templates.ListTemplates()
		fmt.Printf("%s%s%s Updated %s%d templates%s\n",
			style.Green, style.Check, style.Reset,
			style.Cyan, len(list), style.Reset)

		fmt.Println()
		for _, t := range list {
			fmt.Printf("  %s%s%s %s%s%s %s%s%s\n",
				style.Cyan, style.Dot, style.Reset,
				style.Bold, t.Name, style.Reset,
				style.Dim, t.Description, style.Reset)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
