package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godo",
	Short: "Do framework CLI — extract, assemble, and manage .claude/ personas",
	Long: `Godo is the Do framework CLI. It extracts moai-adk's .claude/ directory into
core and persona layers, reassembles them into deployable output, and provides
runtime utilities for hooks, mode switching, linting, scaffolding, and more.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
