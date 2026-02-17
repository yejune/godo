package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/do-focus/convert/internal/persona"
	"github.com/spf13/cobra"
)

var spinnerCmd = &cobra.Command{
	Use:   "spinner [persona-type]",
	Short: "Print spinner verb list for a persona type",
	Long: `Spinner outputs the list of spinner verbs for the given persona type
(young-f, young-m, senior-f, senior-m). Used by Claude Code's status hooks.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSpinner,
}

func init() {
	rootCmd.AddCommand(spinnerCmd)
}

func runSpinner(cmd *cobra.Command, args []string) error {
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}
	if len(args) > 0 {
		personaType = args[0]
	}

	verbs := persona.GetSpinnerVerbs(personaType)
	if len(verbs) == 0 {
		return fmt.Errorf("no spinner verbs found for persona %q", personaType)
	}

	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(verbs, "\n"))
	return nil
}
