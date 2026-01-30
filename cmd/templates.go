package cmd

import (
	"fmt"

	"github.com/drpedapati/irl-template/pkg/templates"
	"github.com/drpedapati/irl-template/pkg/theme"
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
			fmt.Printf("%s Run %s to fetch.\n",
				theme.Warn("No templates available."),
				theme.Cmd("irl update"))
			return nil
		}

		theme.Section("Templates")
		fmt.Println()
		for _, t := range list {
			fmt.Printf("  %s %s\n",
				theme.Cmd(theme.Dot),
				theme.B(t.Name))
			fmt.Printf("    %s\n", theme.Faint(t.Description))
		}
		fmt.Println()
		fmt.Printf("%s irl init -t %s \"purpose\"\n",
			theme.Faint("Usage:"),
			theme.Cmd("<template>"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
