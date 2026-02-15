package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/drpedapati/irl-template/pkg/config"
	"github.com/drpedapati/irl-template/pkg/theme"
	"github.com/spf13/cobra"
)

var (
	profileNameFlag         string
	profileTitleFlag        string
	profileInstitutionFlag  string
	profileDepartmentFlag   string
	profileEmailFlag        string
	profileInstructionsFlag string
	profileClearFlag        bool
	profileJSONFlag         bool
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "View or set your profile",
	Long: `View or set your profile information for plan injection.

Profile fields are added as YAML front matter when creating new projects.

Examples:
  irl profile                                # Show current profile
  irl profile --json                         # JSON output
  irl profile --name "Jane Doe"              # Set name
  irl profile --institution "UCSF"           # Set institution
  irl profile --name "Jane Doe" --title "MD" --institution "UCSF"
  irl profile --instructions "Always cite sources"
  irl profile --clear                        # Clear all fields`,
	RunE: runProfile,
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.Flags().StringVar(&profileNameFlag, "name", "", "Set name")
	profileCmd.Flags().StringVar(&profileTitleFlag, "title", "", "Set title (e.g., MD, PhD)")
	profileCmd.Flags().StringVar(&profileInstitutionFlag, "institution", "", "Set institution")
	profileCmd.Flags().StringVar(&profileDepartmentFlag, "department", "", "Set department")
	profileCmd.Flags().StringVar(&profileEmailFlag, "email", "", "Set email")
	profileCmd.Flags().StringVar(&profileInstructionsFlag, "instructions", "", "Set AI instructions")
	profileCmd.Flags().BoolVar(&profileClearFlag, "clear", false, "Clear all profile fields")
	profileCmd.Flags().BoolVar(&profileJSONFlag, "json", false, "Output as JSON")
}

func runProfile(cmd *cobra.Command, args []string) error {
	// Clear
	if profileClearFlag {
		if err := config.ClearProfile(); err != nil {
			return fmt.Errorf("failed to clear profile: %w", err)
		}
		fmt.Println(theme.OK("Profile cleared"))
		return nil
	}

	// Check if any set flags were provided
	setting := cmd.Flags().Changed("name") || cmd.Flags().Changed("title") ||
		cmd.Flags().Changed("institution") || cmd.Flags().Changed("department") ||
		cmd.Flags().Changed("email") || cmd.Flags().Changed("instructions")

	if setting {
		// Merge with existing profile
		profile := config.GetProfile()
		if cmd.Flags().Changed("name") {
			profile.Name = profileNameFlag
		}
		if cmd.Flags().Changed("title") {
			profile.Title = profileTitleFlag
		}
		if cmd.Flags().Changed("institution") {
			profile.Institution = profileInstitutionFlag
		}
		if cmd.Flags().Changed("department") {
			profile.Department = profileDepartmentFlag
		}
		if cmd.Flags().Changed("email") {
			profile.Email = profileEmailFlag
		}
		if cmd.Flags().Changed("instructions") {
			profile.Instructions = profileInstructionsFlag
		}

		if err := config.SetProfile(profile); err != nil {
			return fmt.Errorf("failed to save profile: %w", err)
		}
		fmt.Println(theme.OK("Profile updated"))
		fmt.Println()
	}

	// Show current profile
	profile := config.GetProfile()

	if profileJSONFlag {
		data, err := json.MarshalIndent(profile, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	if !config.HasProfile() && !setting {
		fmt.Println(theme.Faint("No profile set"))
		fmt.Println()
		fmt.Printf("%s irl profile --name \"Your Name\" --institution \"Your Institution\"\n",
			theme.Faint("Set one:"))
		return nil
	}

	theme.Section("Profile")
	fmt.Println()
	if profile.Name != "" {
		fmt.Println(theme.KeyValue("Name        ", profile.Name))
	}
	if profile.Title != "" {
		fmt.Println(theme.KeyValue("Title       ", profile.Title))
	}
	if profile.Institution != "" {
		fmt.Println(theme.KeyValue("Institution ", profile.Institution))
	}
	if profile.Department != "" {
		fmt.Println(theme.KeyValue("Department  ", profile.Department))
	}
	if profile.Email != "" {
		fmt.Println(theme.KeyValue("Email       ", profile.Email))
	}
	if profile.Instructions != "" {
		fmt.Println(theme.KeyValue("Instructions", profile.Instructions))
	}
	fmt.Println()

	return nil
}
