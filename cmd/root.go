package cmd

import (
	"fmt"
	"os"

	"github.com/drpedapati/irl-template/pkg/style"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "irl",
	Short: "Idempotent Research Loop",
	Long:  `irl - Idempotent Research Loop`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Hide plumbing commands
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Custom help for root command only
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "irl" {
			printRootHelp()
		} else {
			defaultHelp(cmd, args)
		}
	})
}

func printRootHelp() {
	c := style.Cyan
	g := style.Green
	d := style.Dim
	b := style.Bold
	r := style.Reset

	fmt.Printf("%sirl%s %s- Idempotent Research Loop%s\n", b+c, r, d, r)
	fmt.Println()
	fmt.Printf("%sQuick start:%s\n", b, r)
	fmt.Printf("  %sirl init \"your research purpose\"%s\n", g, r)
	fmt.Println()
	fmt.Printf("%sCommands:%s\n", d, r)
	fmt.Printf("  %sinit%s        Create a new IRL project\n", c, r)
	fmt.Println()
	fmt.Printf("%sInfo:%s\n", d, r)
	fmt.Printf("  %stemplates%s   List available templates\n", c, r)
	fmt.Printf("  %sdoctor%s      Check environment setup\n", c, r)
	fmt.Println()
	fmt.Printf("%sSettings:%s\n", d, r)
	fmt.Printf("  %sconfig%s      View or set configuration\n", c, r)
	fmt.Printf("  %supdate%s      Update templates from GitHub\n", c, r)
	fmt.Println()
	fmt.Printf("%sirl <command> --help%s for details\n", d, r)
}

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Print version",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("irl version %s\n", Version)
	},
}
