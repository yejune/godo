package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "convert",
	Short: "Extract and assemble .claude/ directory layers",
	Long: `Convert extracts moai-adk's .claude/ directory into core (methodology-agnostic)
and persona layers, then reassembles core + any persona into a deployable .claude/ output.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
