package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/extractor"
	"github.com/do-focus/convert/internal/model"
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

	// Print summary
	slotCount := len(registry.Slots)
	personaFiles := countPersonaFiles(manifest)
	fmt.Fprintf(cmd.OutOrStdout(), "Extraction complete: %d core slots, %d persona files extracted\n", slotCount, personaFiles)
	fmt.Fprintf(cmd.OutOrStdout(), "  Core:    %s\n", coreDir)
	fmt.Fprintf(cmd.OutOrStdout(), "  Persona: %s\n", personaDir)

	return nil
}

// countPersonaFiles counts the total number of persona asset references in a manifest.
func countPersonaFiles(m *model.PersonaManifest) int {
	n := len(m.Agents) + len(m.Skills) + len(m.Rules) + len(m.Styles) +
		len(m.Commands) + len(m.HookScripts) + len(m.SlotContent) + len(m.AgentPatches)
	if m.ClaudeMD != "" {
		n++
	}
	return n
}
