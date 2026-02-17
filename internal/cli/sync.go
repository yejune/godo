package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Assemble and deploy persona files to .claude/ directory",
	Long: `Sync runs the assembler pipeline and copies the output to .claude/.
Equivalent to 'make dev'. Also ensures settings.local.json exists.`,
	RunE: runSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	// Find godo binary path (self)
	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("find executable: %w", err)
	}

	distDir := "dist"

	// Step 1: Run assembler
	assembleArgs := []string{
		"assemble",
		"--core", "core",
		"--persona", "personas/do/manifest.yaml",
		"--out", distDir,
	}

	assembleCmd := exec.Command(self, assembleArgs...)
	assembleCmd.Stdout = cmd.OutOrStdout()
	assembleCmd.Stderr = cmd.ErrOrStderr()
	if err := assembleCmd.Run(); err != nil {
		return fmt.Errorf("assemble: %w", err)
	}

	// Step 2: Copy dist/ to .claude/
	if err := os.RemoveAll(".claude"); err != nil {
		return fmt.Errorf("clean .claude: %w", err)
	}
	if err := os.MkdirAll(".claude", 0755); err != nil {
		return fmt.Errorf("create .claude: %w", err)
	}

	// Copy directories
	dirs := []string{"agents", "commands", "rules", "skills", "styles", "characters", "spinners"}
	for _, dir := range dirs {
		src := filepath.Join(distDir, dir)
		if info, err := os.Stat(src); err == nil && info.IsDir() {
			if err := copyDir(src, filepath.Join(".claude", dir)); err != nil {
				return fmt.Errorf("copy %s: %w", dir, err)
			}
		}
	}

	// Copy files
	files := []string{"settings.json", "registry.yaml"}
	for _, f := range files {
		src := filepath.Join(distDir, f)
		dst := filepath.Join(".claude", f)
		if data, err := os.ReadFile(src); err == nil {
			if err := os.WriteFile(dst, data, 0644); err != nil {
				return fmt.Errorf("copy %s: %w", f, err)
			}
		}
	}

	// Copy CLAUDE.md to project root
	if data, err := os.ReadFile(filepath.Join(distDir, "CLAUDE.md")); err == nil {
		if err := os.WriteFile("CLAUDE.md", data, 0644); err != nil {
			return fmt.Errorf("copy CLAUDE.md: %w", err)
		}
	}

	// Step 3: Ensure settings.local.json
	EnsureSettingsLocal()

	fmt.Fprintln(cmd.OutOrStdout(), "Sync complete: .claude/ + CLAUDE.md")
	return nil
}

// copyDir recursively copies a directory tree.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path and destination
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}
