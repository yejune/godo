package cli

import (
	"github.com/spf13/cobra"
)

var moaiCmd = &cobra.Command{
	Use:   "moai",
	Short: "Framework developer tools",
	Long:  `Developer tools for building and managing the Do framework itself.`,
}

func init() {
	rootCmd.AddCommand(moaiCmd)
}
