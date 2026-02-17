package cli

import (
	"fmt"

	"github.com/yejune/godo/internal/scaffold"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [agent|skill] [name]",
	Short: "Scaffold new agent or skill definitions",
	Long: `Create generates boilerplate files for new agents or skills
following the MoAI-ADK conventions.`,
	Args: cobra.ExactArgs(2),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	kind := args[0]
	name := args[1]

	switch kind {
	case "agent":
		if err := scaffold.CreateAgent(name); err != nil {
			return fmt.Errorf("create agent: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "agent %q created\n", name)
	case "skill":
		if err := scaffold.CreateSkill(name); err != nil {
			return fmt.Errorf("create skill: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "skill %q created\n", name)
	default:
		return fmt.Errorf("unknown kind %q (valid: agent, skill)", kind)
	}

	return nil
}
