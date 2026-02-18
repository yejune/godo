package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/yejune/godo/internal/profile"
)

var claudeCmd = &cobra.Command{
	Use:   "claude [flags] [-- claude-args...]",
	Short: "Launch Claude Code with configured flags",
	Long: `Launch Claude Code with flags configured via 'godo setup'.
Reads DO_CLAUDE_* settings from .claude/settings.local.json.
Pass additional arguments to Claude Code after --.`,
	Aliases:            []string{"cc"},
	RunE:               runClaude,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(claudeCmd)
}

func runClaude(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && args[0] == "profile" {
		return runClaudeProfileCompat(cmd, args[1:])
	}

	profileName, filteredArgs := parseClaudeProfileFlag(args)
	if profileName != "" && profileName != "default" {
		if err := profile.EnsureDir(profileName); err != nil {
			return fmt.Errorf("set profile: %w", err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Profile: %s\n", profileName)
	}

	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude not found in PATH. Install Claude Code first")
	}

	settings := readSettingsLocal()
	env := getEnvMap(settings)

	getString := func(key string) string {
		if v, ok := env[key].(string); ok {
			return v
		}
		return ""
	}

	bypass := getString("DO_CLAUDE_BYPASS") == "true"
	chrome := getString("DO_CLAUDE_CHROME") == "true"
	cont := getString("DO_CLAUDE_CONTINUE") == "true"
	autoSync := getString("DO_CLAUDE_AUTO_SYNC") == "true"
	model := getString("DO_CLAUDE_MODEL")

	// Default model to opus[1m] if not specified
	if model == "" {
		model = "opus[1m]"
	}

	var passThrough []string
	for _, arg := range filteredArgs {
		switch arg {
		case "--chrome":
			chrome = true
		case "--no-chrome":
			chrome = false
		case "-b", "--bypass":
			bypass = true
		case "-c", "--continue":
			cont = true
		default:
			passThrough = append(passThrough, arg)
		}
	}

	if autoSync {
		fmt.Fprintln(cmd.ErrOrStderr(), "Auto-syncing...")
		self, _ := os.Executable()
		syncExec := exec.Command(self, "sync")
		syncExec.Stdout = cmd.OutOrStdout()
		syncExec.Stderr = cmd.ErrOrStderr()
		if err := syncExec.Run(); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: sync failed: %v\n", err)
		}
	}

	claudeArgs := []string{"claude"}
	if bypass {
		claudeArgs = append(claudeArgs, "--dangerously-skip-permissions")
	}
	if !chrome {
		claudeArgs = append(claudeArgs, "--no-chrome")
	}
	if cont {
		claudeArgs = append(claudeArgs, "--continue")
	}
	if model != "" {
		claudeArgs = append(claudeArgs, "--model", model)
	}
	claudeArgs = append(claudeArgs, passThrough...)

	return syscall.Exec(claudeBin, claudeArgs, os.Environ())
}

func parseClaudeProfileFlag(args []string) (string, []string) {
	var profileName string
	filtered := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		if (args[i] == "--profile" || args[i] == "-p") && i+1 < len(args) {
			profileName = args[i+1]
			i++
			continue
		}
		filtered = append(filtered, args[i])
	}

	return profileName, filtered
}

func runClaudeProfileCompat(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: godo claude profile <list|current|delete>")
	}

	switch args[0] {
	case "list", "ls":
		return runProfileList(cmd, nil)
	case "current":
		return runProfileCurrent(cmd, nil)
	case "delete", "rm":
		if len(args) < 2 {
			return fmt.Errorf("usage: godo claude profile delete <name>")
		}
		return runProfileDelete(cmd, []string{args[1]})
	default:
		return fmt.Errorf("unknown profile command %q", args[0])
	}
}
