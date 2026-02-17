package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yejune/godo/internal/lint"
)

var lintCmd = &cobra.Command{
	Use:   "lint [files...]",
	Short: "Run language-aware linters on changed or specified files",
	Long: `Lint detects languages from file extensions and runs the appropriate linter
(go vet, ruff, tsc, eslint, cargo clippy). Without arguments, lints git-changed files.`,
	RunE: runLint,
}

var lintSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive linter installation",
	RunE:  runLintSetup,
}

var lintAll bool

func init() {
	lintCmd.Flags().BoolVar(&lintAll, "all", false, "lint all project files instead of only changed files")
	lintCmd.AddCommand(lintSetupCmd)
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

func runLintSetup(cmd *cobra.Command, args []string) error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Scanning project for languages...")
	fmt.Fprintln(cmd.OutOrStdout())

	langs := lint.ScanProjectLanguages(projectDir)
	if len(langs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No supported code files found in project.")
		return nil
	}

	var langNames []string
	for _, l := range langs {
		langNames = append(langNames, string(l))
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Detected languages: %s\n\n", strings.Join(langNames, ", "))

	status := lint.CheckSetupStatus(langs)
	lint.PrintSetupStatus(status)
	if len(status.Missing) == 0 {
		return nil
	}

	managers := lint.DetectPackageManagers()
	scanner := bufio.NewScanner(cmd.InOrStdin())

	for _, info := range status.Missing {
		fmt.Fprintf(cmd.OutOrStdout(), "Install %s?\n", info.DisplayName)
		options := lint.GetInstallOptions(info, managers)
		if len(options) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  No package manager available.\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Install manually: %s\n\n", info.DownloadURL)
			continue
		}

		for i, opt := range options {
			fmt.Fprintf(cmd.OutOrStdout(), "  %d) %s\n", i+1, opt.Label)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  s) Skip\n")
		fmt.Fprintf(cmd.OutOrStdout(), "Choose [1-%d, s]: ", len(options))

		if !scanner.Scan() {
			return nil
		}
		choice := strings.TrimSpace(scanner.Text())
		if choice == "" || strings.EqualFold(choice, "s") {
			fmt.Fprintln(cmd.OutOrStdout(), "  Skipped.")
			fmt.Fprintln(cmd.OutOrStdout())
			continue
		}

		idx := 0
		if _, err := fmt.Sscanf(choice, "%d", &idx); err != nil || idx < 1 || idx > len(options) {
			fmt.Fprintln(cmd.OutOrStdout(), "  Invalid choice, skipping.")
			fmt.Fprintln(cmd.OutOrStdout())
			continue
		}

		selected := options[idx-1]
		fmt.Fprintf(cmd.OutOrStdout(), "  Running: %s\n", selected.Label)
		if err := lint.RunInstall(selected.Args); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "  Failed: %v\n", err)
			fmt.Fprintf(cmd.OutOrStdout(), "  Try manual install: %s\n\n", info.DownloadURL)
			continue
		}
		fmt.Fprintln(cmd.OutOrStdout(), "  Installed!")
		fmt.Fprintln(cmd.OutOrStdout())
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Setup complete.")
	return nil
}
