package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/do-focus/convert/internal/hook"
	"github.com/spf13/cobra"
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

func runHook(cmd *cobra.Command, args []string) error {
	eventType := hook.EventType(args[0])
	if !hook.IsValidEventType(eventType) {
		return fmt.Errorf("unknown event type: %s (valid: %v)", eventType, hook.ValidEventTypes())
	}

	// Read hook input from stdin.
	input := hook.ReadStdin()

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

	// For now, echo the input back as JSON for debugging.
	// Individual hook handlers will be wired up as they are implemented.
	data, _ := json.MarshalIndent(input, "", "  ")
	fmt.Fprintf(cmd.OutOrStdout(), "event=%s input=%s\n", eventType, string(data))

	return nil
}
