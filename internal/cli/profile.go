package cli

import (
	"fmt"

	"github.com/do-focus/convert/internal/profile"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Claude configuration profiles",
	Long: `Profile lists, switches, and deletes Claude configuration profiles
stored in ~/.claude-profiles/.`,
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
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
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := profile.Delete(args[0]); err != nil {
			return fmt.Errorf("delete profile: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "profile %q deleted\n", args[0])
		return nil
	},
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	rootCmd.AddCommand(profileCmd)
}
