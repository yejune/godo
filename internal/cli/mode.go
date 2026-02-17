package cli

import (
	"fmt"

	"github.com/yejune/godo/internal/mode"
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

var modePermissionCmd = &cobra.Command{
	Use:   "permission [bypass|accept|default|plan]",
	Short: "Set the Claude Code permission mode",
	Args:  cobra.ExactArgs(1),
	RunE:  runModePermission,
}

func init() {
	rootCmd.AddCommand(modeCmd)
	modeCmd.AddCommand(modePermissionCmd)
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

func runModePermission(cmd *cobra.Command, args []string) error {
	name := args[0]
	ccMode, ok := mode.PermissionModes[name]
	if !ok {
		return fmt.Errorf("invalid permission mode %q (valid: bypass, accept, default, plan)", name)
	}

	if err := mode.SetDefaultMode(ccMode); err != nil {
		return fmt.Errorf("set permission mode: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "permission mode set to %s (%s)\n", name, ccMode)
	return nil
}
