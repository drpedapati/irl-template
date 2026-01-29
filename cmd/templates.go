package cmd

import (
	"fmt"

	"github.com/drpedapati/irl-template/pkg/style"
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
			fmt.Printf("%sNo templates available.%s Run %sirl update%s to fetch.\n",
				style.Yellow, style.Reset, style.Cyan, style.Reset)
			return nil
		}

		fmt.Printf("%sTemplates%s\n", style.BoldCyan, style.Reset)
		fmt.Println()
		for _, t := range list {
			fmt.Printf("  %s%s%s %s%s%s\n",
				style.Cyan, style.Dot, style.Reset,
				style.Bold, t.Name, style.Reset)
			fmt.Printf("    %s%s%s\n", style.Dim, t.Description, style.Reset)
		}
		fmt.Println()
		fmt.Printf("%sUsage:%s irl init -t %s<template>%s \"purpose\"\n",
			style.Dim, style.Reset, style.Cyan, style.Reset)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
