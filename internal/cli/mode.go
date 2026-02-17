package cli

import (
	"fmt"

	"github.com/do-focus/convert/internal/mode"
	"github.com/spf13/cobra"
)

var modeCmd = &cobra.Command{
	Use:   "mode [do|focus|team]",
	Short: "Get or set the current execution mode",
	Long: `Mode manages the Do framework execution mode (do, focus, or team).
Without arguments, prints the current mode. With an argument, switches to the specified mode.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMode,
}

func init() {
	rootCmd.AddCommand(modeCmd)
}

func runMode(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), mode.ReadState())
		return nil
	}

	newMode := args[0]
	switch newMode {
	case "do", "focus", "team":
		mode.WriteState(newMode)
		if err := mode.SetDefaultMode(newMode); err != nil {
			return fmt.Errorf("set default mode: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "mode switched to %s\n", newMode)
	default:
		return fmt.Errorf("invalid mode %q (valid: do, focus, team)", newMode)
	}

	return nil
}
