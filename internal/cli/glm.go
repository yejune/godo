package cli

import (
	"fmt"

	"github.com/do-focus/convert/internal/glm"
	"github.com/spf13/cobra"
)

var glmCmd = &cobra.Command{
	Use:   "glm",
	Short: "Manage GLM API credentials",
	Long:  `GLM provides credential management for the GLM (Generative Language Model) API.`,
}

var glmSetupCmd = &cobra.Command{
	Use:   "setup [api-key]",
	Short: "Store a GLM API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := glm.SetupCredentials(args[0]); err != nil {
			return fmt.Errorf("setup credentials: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "GLM API key stored (%s)\n", glm.MaskAPIKey(args[0]))
		return nil
	},
}

var glmStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current GLM credential status",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := glm.LoadCredentials()
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "no GLM credentials configured")
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "GLM API key: %s\n", glm.MaskAPIKey(creds.APIKey))
		return nil
	},
}

func init() {
	glmCmd.AddCommand(glmSetupCmd)
	glmCmd.AddCommand(glmStatusCmd)
	rootCmd.AddCommand(glmCmd)
}
