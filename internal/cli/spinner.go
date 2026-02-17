package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	if err := applySpinnerToSettings(verbs); err != nil {
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
	settings, err := readClaudeSettings(getClaudeSettingsPath())
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

func getClaudeSettingsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "settings.json")
	}
	return filepath.Join(homeDir, ".claude", "settings.json")
}

func readClaudeSettings(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]interface{}), nil
		}
		return nil, err
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	return settings, nil
}

func writeClaudeSettings(path string, settings map[string]interface{}) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func applySpinnerToSettings(verbs []string) error {
	settingsPath := getClaudeSettingsPath()
	settings, err := readClaudeSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("reading settings: %w", err)
	}

	iverbs := make([]interface{}, len(verbs))
	for i, v := range verbs {
		iverbs[i] = v
	}

	settings["spinnerVerbs"] = map[string]interface{}{
		"mode":  "replace",
		"verbs": iverbs,
	}

	if err := writeClaudeSettings(settingsPath, settings); err != nil {
		return fmt.Errorf("writing settings: %w", err)
	}
	return nil
}

func removeSpinnerFromSettings() error {
	settingsPath := getClaudeSettingsPath()
	settings, err := readClaudeSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("reading settings: %w", err)
	}
	delete(settings, "spinnerVerbs")
	if err := writeClaudeSettings(settingsPath, settings); err != nil {
		return fmt.Errorf("writing settings: %w", err)
	}
	return nil
}
