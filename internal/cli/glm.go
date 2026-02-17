package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yejune/godo/internal/glm"
)

var glmCmd = &cobra.Command{
	Use:                "glm [flags] [-- claude-args...]",
	Short:              "Launch Claude Code with GLM backend",
	Long:               `Without subcommands, runs Claude Code using stored GLM credentials. Supports -c/--continue and other Claude flags.`,
	RunE:               runGLM,
	DisableFlagParsing: true,
}

var glmSetupCmd = &cobra.Command{
	Use:   "setup [api-key]",
	Short: "Store a GLM API key",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGLMSetup,
}

var glmStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current GLM credential status",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := glm.LoadCredentials()
		if err != nil || creds == nil || creds.APIKey == "" {
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

func runGLM(cmd *cobra.Command, args []string) error {
	creds, err := glm.LoadCredentials()
	if err != nil {
		return fmt.Errorf("load GLM credentials: %w", err)
	}
	if creds == nil || creds.APIKey == "" {
		return fmt.Errorf("GLM API key is not configured. Run 'godo glm setup'")
	}

	glm.SetGLMEnv(creds.APIKey)
	fmt.Fprintln(cmd.OutOrStdout(), "Launching Claude Code with GLM backend")
	return runClaude(cmd, args)
}

func runGLMSetup(cmd *cobra.Command, args []string) error {
	apiKey := ""
	if len(args) == 1 {
		apiKey = strings.TrimSpace(args[0])
	} else {
		fmt.Fprint(cmd.OutOrStdout(), "GLM API key: ")
		scanner := bufio.NewScanner(cmd.InOrStdin())
		if !scanner.Scan() {
			return nil
		}
		apiKey = strings.TrimSpace(scanner.Text())
	}

	if apiKey == "" {
		return fmt.Errorf("empty API key")
	}

	if err := glm.SetupCredentials(apiKey); err != nil {
		return fmt.Errorf("setup credentials: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "GLM API key stored (%s)\n", glm.MaskAPIKey(apiKey))
	return nil
}
