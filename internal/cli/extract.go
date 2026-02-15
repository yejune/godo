package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/extractor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract .claude/ directory into core templates and persona manifest",
	Long: `Extract walks a .claude/ source directory, separates methodology-agnostic
core templates from persona-specific content, and writes:
  - Core templates to <out>/core/
  - Persona manifest to <out>/personas/<name>/manifest.yaml

Source can be a local directory (--src) or a GitHub repository (--repo).
The --src and --repo flags are mutually exclusive.`,
	RunE: runExtract,
}

var (
	extractSrc     string
	extractOut     string
	extractPersona string
	extractRepo    string
	extractBranch  string
)

func init() {
	extractCmd.Flags().StringVar(&extractSrc, "src", "", "source .claude/ directory path")
	extractCmd.Flags().StringVar(&extractRepo, "repo", "", "GitHub repository URL or shorthand (e.g., org/repo)")
	extractCmd.Flags().StringVar(&extractBranch, "branch", "", "branch or tag to clone (default: repository default branch)")
	extractCmd.Flags().StringVar(&extractOut, "out", "", "output directory for core templates and persona manifest (required)")
	extractCmd.Flags().StringVar(&extractPersona, "persona", "", "persona name (default: auto-detect from source)")
	_ = extractCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(extractCmd)
}

// expandRepoURL normalizes a repository reference into a full git clone URL.
//   - "org/repo"                   -> "https://github.com/org/repo.git"
//   - "github.com/org/repo"        -> "https://github.com/org/repo.git"
//   - "https://github.com/o/r"     -> used as-is
//   - "https://github.com/o/r.git" -> used as-is
func expandRepoURL(raw string) string {
	// Already a full URL — use as-is
	if strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "git@") {
		return raw
	}

	// "github.com/org/repo" -> strip domain prefix, treat as shorthand
	if strings.HasPrefix(raw, "github.com/") {
		raw = strings.TrimPrefix(raw, "github.com/")
	}

	// "org/repo" shorthand
	raw = strings.TrimSuffix(raw, ".git")
	return "https://github.com/" + raw + ".git"
}

// cloneRepo performs a shallow git clone and returns the temp directory path.
// The caller is responsible for removing the temp directory.
func cloneRepo(repoURL, branch string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "convert-clone-*")
	if err != nil {
		return "", fmt.Errorf("create temp directory: %w", err)
	}

	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, repoURL, tmpDir)

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Clean up on clone failure
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("git clone %s: %w", repoURL, err)
	}

	return tmpDir, nil
}

func runExtract(cmd *cobra.Command, args []string) error {
	// Validate mutually exclusive flags
	hasSrc := extractSrc != ""
	hasRepo := extractRepo != ""

	if hasSrc && hasRepo {
		return fmt.Errorf("--src and --repo are mutually exclusive; provide one, not both")
	}
	if !hasSrc && !hasRepo {
		return fmt.Errorf("either --src or --repo is required")
	}

	srcDir := extractSrc

	// Handle --repo: clone and resolve .claude/ directory
	if hasRepo {
		repoURL := expandRepoURL(extractRepo)
		fmt.Fprintf(cmd.OutOrStdout(), "Cloning %s ...\n", repoURL)

		tmpDir, err := cloneRepo(repoURL, extractBranch)
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)

		claudeDir := filepath.Join(tmpDir, ".claude")
		info, err := os.Stat(claudeDir)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("repository does not contain a .claude/ directory")
		}

		srcDir = claudeDir
	}

	// Validate source directory exists
	info, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("source directory %q: %w", srcDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("source %q is not a directory", srcDir)
	}

	// Set up detector and pattern registry
	patternReg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(patternReg)
	if err != nil {
		return fmt.Errorf("create persona detector: %w", err)
	}

	// Create orchestrator and run extraction
	orch := extractor.NewExtractorOrchestrator(det, patternReg)
	registry, manifest, err := orch.Extract(srcDir)
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	// Save core templates (registry.yaml) to <out>/core/
	coreDir := filepath.Join(extractOut, "core")
	if err := registry.Save(coreDir); err != nil {
		return fmt.Errorf("save core templates: %w", err)
	}

	// Copy core files to <out>/core/
	coreFileCount := 0
	for _, relPath := range manifest.CoreFiles {
		src := filepath.Join(manifest.SourceDir, relPath)
		dst := filepath.Join(coreDir, relPath)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("copy core file %s: %w", relPath, err)
		}
		coreFileCount++
	}

	// Determine persona name
	personaName := extractPersona
	if personaName == "" {
		personaName = manifest.Name
	}
	if personaName == "" {
		personaName = "default"
	}

	// Save persona manifest to <out>/personas/<name>/manifest.yaml
	personaDir := filepath.Join(extractOut, "personas", personaName)
	if err := os.MkdirAll(personaDir, 0o755); err != nil {
		return fmt.Errorf("create persona dir %s: %w", personaDir, err)
	}

	manifestData, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("marshal persona manifest: %w", err)
	}

	manifestPath := filepath.Join(personaDir, "manifest.yaml")
	if err := os.WriteFile(manifestPath, manifestData, 0o644); err != nil {
		return fmt.Errorf("write persona manifest %s: %w", manifestPath, err)
	}

	// Copy persona files to <out>/personas/<name>/
	personaFileCount := 0
	for relPath, srcPath := range manifest.PersonaFiles {
		dst := filepath.Join(personaDir, relPath)
		if err := copyFile(srcPath, dst); err != nil {
			return fmt.Errorf("copy persona file %s: %w", relPath, err)
		}
		personaFileCount++
	}

	// Print summary
	fmt.Fprintf(cmd.OutOrStdout(), "Extracted %d core files, %d persona files to %s\n", coreFileCount, personaFileCount, extractOut)
	fmt.Fprintf(cmd.OutOrStdout(), "  Core:    %s (%d files + registry.yaml)\n", coreDir, coreFileCount)
	fmt.Fprintf(cmd.OutOrStdout(), "  Persona: %s (%d files + manifest.yaml)\n", personaDir, personaFileCount)

	return nil
}

// copyFile copies a file from src to dst, preserving the source file's permissions.
// Parent directories are created as needed.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source %s: %w", src, err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source %s: %w", src, err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create parent dir for %s: %w", dst, err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode().Perm())
	if err != nil {
		return fmt.Errorf("create destination %s: %w", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy data to %s: %w", dst, err)
	}

	return nil
}
