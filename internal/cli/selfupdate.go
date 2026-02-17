package cli

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var selfupdateCmd = &cobra.Command{
	Use:     "selfupdate",
	Aliases: []string{"self-update"},
	Short:   "Update godo to the latest version via Homebrew",
	RunE:    runSelfupdate,
}

func init() {
	rootCmd.AddCommand(selfupdateCmd)
}

func runSelfupdate(cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew not found. Install godo manually or install Homebrew first")
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Updating godo via Homebrew...")

	brewCmd := exec.Command("brew", "upgrade", "yejune/tap/godo")
	brewCmd.Stdout = cmd.OutOrStdout()
	brewCmd.Stderr = cmd.ErrOrStderr()

	if err := brewCmd.Run(); err != nil {
		return fmt.Errorf("brew upgrade failed: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Update complete!")
	return nil
}
