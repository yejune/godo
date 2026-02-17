package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yejune/godo/internal/hook"
)

var hookCmd = &cobra.Command{
	Use:   "hook [event-type]",
	Short: "Execute hook logic for Claude Code lifecycle events",
	Long: `Hook reads JSON from stdin and dispatches to the appropriate hook handler
based on the event type. Used by Claude Code's hooks system.`,
	Args: cobra.ExactArgs(1),
	RunE: runHook,
}

func init() {
	rootCmd.AddCommand(hookCmd)
}

// handlerFunc is the signature for all hook handlers.
type handlerFunc func(*hook.Input) *hook.Output

// hookHandlers maps CLI event names (kebab-case) to handler functions.
var hookHandlers = map[string]handlerFunc{
	"session-start":      hook.HandleSessionStart,
	"pre-tool":           hook.HandlePreTool,
	"post-tool-use":      hook.HandlePostToolUse,
	"compact":            hook.HandleCompact,
	"stop":               hook.HandleStop,
	"session-end":        hook.HandleSessionEnd,
	"subagent-stop":      hook.HandleSubagentStop,
	"user-prompt-submit": hook.HandleUserPromptSubmit,
}

func runHook(cmd *cobra.Command, args []string) error {
	cliEventType := args[0]

	handler, ok := hookHandlers[cliEventType]
	if !ok {
		return fmt.Errorf("unknown event type: %s (valid: session-start, pre-tool, post-tool-use, compact, stop, session-end, subagent-stop, user-prompt-submit)", cliEventType)
	}

	// Read structured input from stdin.
	input := hook.ReadInput()

	// Validate contract.
	workDir, _ := os.Getwd()
	contract := hook.NewContract(workDir)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := contract.Validate(ctx); err != nil {
		output := hook.NewDenyOutput(fmt.Sprintf("contract violation: %v", err))
		hook.WriteOutput(output)
		return nil
	}

	// Dispatch to handler.
	output := handler(input)
	if output != nil {
		hook.WriteOutput(output)
	}

	return nil
}
