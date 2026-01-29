package cmd

import (
	"fmt"

	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update templates from GitHub",
	Long:  "Fetch the latest templates from the IRL template repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching latest templates...")

		if err := templates.Update(); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		fmt.Println("✓ Templates updated")

		// Show available templates
		list, _ := templates.ListTemplates()
		fmt.Printf("\nAvailable templates (%d):\n", len(list))
		for _, t := range list {
			fmt.Printf("  • %s - %s\n", t.Name, t.Description)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
