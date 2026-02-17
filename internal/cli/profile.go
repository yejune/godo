package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yejune/godo/internal/profile"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Claude configuration profiles",
	Long: `Profile lists, switches, and deletes Claude configuration profiles
stored in ~/.claude-profiles/.`,
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all available profiles",
	RunE:    runProfileList,
}

var profileCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current profile name",
	RunE:  runProfileCurrent,
}

var profileDeleteCmd = &cobra.Command{
	Use:     "delete [name]",
	Aliases: []string{"rm"},
	Short:   "Delete a profile",
	Args:    cobra.ExactArgs(1),
	RunE:    runProfileDelete,
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileCurrentCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	rootCmd.AddCommand(profileCmd)
}

func runProfileList(cmd *cobra.Command, args []string) error {
	entries := profile.List()
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no profiles found")
		return nil
	}
	for _, e := range entries {
		marker := "  "
		if e.Current {
			marker = "* "
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", marker, e.Name)
	}
	return nil
}

func runProfileCurrent(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(cmd.OutOrStdout(), profile.GetCurrentName())
	return nil
}

func runProfileDelete(cmd *cobra.Command, args []string) error {
	if err := profile.Delete(args[0]); err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}
	return nil
}
