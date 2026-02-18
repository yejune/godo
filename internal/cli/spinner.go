package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yejune/godo/internal/persona"
)

var spinnerCmd = &cobra.Command{
	Use:   "spinner [apply|restore|status|persona-type]",
	Short: "Manage or print spinner verbs",
	Long: `Compatibility mode:
  godo spinner apply|restore|status
Legacy mode:
  godo spinner [young-f|young-m|senior-f|senior-m]`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSpinner,
}

func init() {
	rootCmd.AddCommand(spinnerCmd)
}

func runSpinner(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return spinnerStatus(cmd)
	}

	switch args[0] {
	case "apply":
		return spinnerApply(cmd)
	case "restore":
		return spinnerRestore(cmd)
	case "status":
		return spinnerStatus(cmd)
	default:
		verbs := persona.GetSpinnerVerbs(args[0])
		if len(verbs) == 0 {
			return fmt.Errorf("no spinner verbs found for persona %q", args[0])
		}
		fmt.Fprintln(cmd.OutOrStdout(), strings.Join(verbs, "\n"))
		return nil
	}
}

func spinnerApply(cmd *cobra.Command) error {
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

	verbs := persona.GetSpinnerVerbs(personaType)
	if len(verbs) == 0 {
		return fmt.Errorf("no spinner verbs found for persona %q", personaType)
	}

	if err := persona.ApplySpinnerToSettings(verbs); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Applied %d spinner verbs (persona: %s)\n", len(verbs), personaType)
	return nil
}

func spinnerRestore(cmd *cobra.Command) error {
	if err := removeSpinnerFromSettings(); err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Removed spinner verbs (restored English defaults)")
	return nil
}

func spinnerStatus(cmd *cobra.Command) error {
	settings, err := persona.ReadClaudeSettings(persona.GetClaudeSettingsPath())
	if err != nil {
		return err
	}

	spinnerVerbs, exists := settings["spinnerVerbs"]
	if !exists {
		fmt.Fprintln(cmd.OutOrStdout(), "Spinner: English (default)")
		fmt.Fprintln(cmd.OutOrStdout(), "  Run 'godo spinner apply' to switch to Korean")
		return nil
	}

	if section, ok := spinnerVerbs.(map[string]interface{}); ok {
		if raw, ok := section["verbs"]; ok {
			if list, ok := raw.([]interface{}); ok {
				fmt.Fprintf(cmd.OutOrStdout(), "Spinner: Korean (%d verbs)\n", len(list))
				if len(list) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "  Sample: %v\n", list[0])
				}
				fmt.Fprintln(cmd.OutOrStdout(), "  Run 'godo spinner restore' to switch back to English")
				return nil
			}
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Spinner: custom (unknown format)")
	return nil
}

func removeSpinnerFromSettings() error {
	settingsPath := persona.GetClaudeSettingsPath()
	settings, err := persona.ReadClaudeSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("reading settings: %w", err)
	}
	delete(settings, "spinnerVerbs")
	if err := persona.WriteClaudeSettings(settingsPath, settings); err != nil {
		return fmt.Errorf("writing settings: %w", err)
	}
	return nil
}
