package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var claudeCmd = &cobra.Command{
	Use:   "claude [flags] [-- claude-args...]",
	Short: "Launch Claude Code with configured flags",
	Long: `Launch Claude Code with flags configured via 'godo setup'.
Reads DO_CLAUDE_* settings from .claude/settings.local.json.
Pass additional arguments to Claude Code after --.`,
	RunE:               runClaude,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(claudeCmd)
}

func runClaude(cmd *cobra.Command, args []string) error {
	// Find claude binary
	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude not found in PATH. Install Claude Code first")
	}

	// Read settings
	settings := readSettingsLocal()
	env := getEnvMap(settings)

	getString := func(key string) string {
		if v, ok := env[key].(string); ok {
			return v
		}
		return ""
	}

	// Auto sync if enabled
	if getString("DO_CLAUDE_AUTO_SYNC") == "true" {
		fmt.Fprintln(cmd.ErrOrStderr(), "Auto-syncing...")
		self, _ := os.Executable()
		syncExec := exec.Command(self, "sync")
		syncExec.Stdout = cmd.OutOrStdout()
		syncExec.Stderr = cmd.ErrOrStderr()
		if err := syncExec.Run(); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: sync failed: %v\n", err)
		}
	}

	// Build claude args
	claudeArgs := []string{"claude"}

	if getString("DO_CLAUDE_BYPASS") == "true" {
		claudeArgs = append(claudeArgs, "--dangerously-skip-permissions")
	}

	if getString("DO_CLAUDE_CONTINUE") == "true" {
		claudeArgs = append(claudeArgs, "--continue")
	}

	// Append any extra args passed after --
	claudeArgs = append(claudeArgs, args...)

	// Launch claude via exec (replaces current process)
	return syscall.Exec(claudeBin, claudeArgs, os.Environ())
}
