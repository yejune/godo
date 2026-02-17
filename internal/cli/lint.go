package cli

import (
	"fmt"
	"os"

	"github.com/yejune/godo/internal/lint"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint [files...]",
	Short: "Run language-aware linters on changed or specified files",
	Long: `Lint detects languages from file extensions and runs the appropriate linter
(go vet, ruff, tsc, eslint, cargo clippy). Without arguments, lints git-changed files.`,
	RunE: runLint,
}

var lintAll bool

func init() {
	lintCmd.Flags().BoolVar(&lintAll, "all", false, "lint all project files instead of only changed files")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	var files []string
	if len(args) > 0 {
		files = args
	} else {
		files = lint.GetChangedFiles(projectDir, lintAll)
	}

	if len(files) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no files to lint")
		return nil
	}

	grouped := lint.GroupFilesByLanguage(files)
	var allDiags []lint.Diagnostic

	for lang, langFiles := range grouped {
		diags := lint.RunLinter(lang, langFiles, projectDir)
		allDiags = append(allDiags, diags...)
	}

	if len(allDiags) > 0 {
		fmt.Fprint(cmd.OutOrStdout(), lint.FormatDiagnostics(allDiags))
	}

	exitCode := lint.EvaluateResults(allDiags)
	if exitCode != 0 {
		os.Exit(exitCode)
	}

	return nil
}
