package cli

import (
	"github.com/do-focus/convert/internal/mode"
	"github.com/do-focus/convert/internal/profile"
	"github.com/do-focus/convert/internal/statusline"
	"github.com/spf13/cobra"
)

var statuslineCmd = &cobra.Command{
	Use:   "statusline",
	Short: "Render the Claude Code status line from stdin JSON",
	Long: `Statusline reads Claude Code's JSON status payload from stdin and prints
a formatted status line with mode, model, context usage, and cost info.`,
	Run: func(cmd *cobra.Command, args []string) {
		statusline.Render(statusline.Config{
			Version:        rootCmd.Version,
			ReadModeState:  mode.ReadState,
			GetProfileName: profile.GetCurrentName,
		})
	},
}

func init() {
	rootCmd.AddCommand(statuslineCmd)
}
