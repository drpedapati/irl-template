package cmd

import (
	"fmt"

	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := templates.ListTemplates()
		if err != nil {
			return err
		}

		if len(list) == 0 {
			fmt.Println("No templates available. Run 'irl update' to fetch templates.")
			return nil
		}

		fmt.Println("Available templates:")
		for _, t := range list {
			fmt.Printf("  • %s\n    %s\n", t.Name, t.Description)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
