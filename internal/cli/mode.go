package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yejune/godo/internal/mode"
)

var modeCmd = &cobra.Command{
	Use:   "mode [get|set <mode>|<mode>]",
	Short: "Get or set the current execution/permission mode",
	Long: `Compatible forms:
  godo mode
  godo mode get
  godo mode set <do|focus|team|bypass|accept|default|plan>
  godo mode <do|focus|team|bypass|accept|default|plan>`,
	Args: cobra.MaximumNArgs(2),
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

	arg0 := strings.ToLower(args[0])

	if arg0 == "get" {
		fmt.Fprintln(cmd.OutOrStdout(), mode.ReadState())
		return nil
	}

	if arg0 == "set" {
		if len(args) < 2 {
			return fmt.Errorf("usage: godo mode set <do|focus|team|bypass|accept|default|plan>")
		}
		return applyMode(cmd, strings.ToLower(args[1]))
	}

	return applyMode(cmd, arg0)
}

func applyMode(cmd *cobra.Command, value string) error {
	if ccMode, ok := mode.PermissionModes[value]; ok {
		if err := mode.SetDefaultMode(ccMode); err != nil {
			return fmt.Errorf("set permission mode: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Permission mode set: %s (%s)\n", value, ccMode)
		return nil
	}

	if mode.ExecutionModes[value] {
		mode.WriteState(value)
		fmt.Fprintf(cmd.OutOrStdout(), "Mode set: %s\n", value)
		return nil
	}

	return fmt.Errorf("invalid mode %q (valid execution: do/focus/team, permission: bypass/accept/default/plan)", value)
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
