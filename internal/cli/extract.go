package cli

import (
	"fmt"
	"os"
	"path/filepath"

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
  - Persona manifest to <out>/personas/<name>/manifest.yaml`,
	RunE: runExtract,
}

var (
	extractSrc     string
	extractOut     string
	extractPersona string
)

func init() {
	extractCmd.Flags().StringVar(&extractSrc, "src", "", "source .claude/ directory path (required)")
	extractCmd.Flags().StringVar(&extractOut, "out", "", "output directory for core templates and persona manifest (required)")
	extractCmd.Flags().StringVar(&extractPersona, "persona", "", "persona name (default: auto-detect from source)")
	_ = extractCmd.MarkFlagRequired("src")
	_ = extractCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(extractCmd)
}

func runExtract(cmd *cobra.Command, args []string) error {
	// Validate source directory exists
	info, err := os.Stat(extractSrc)
	if err != nil {
		return fmt.Errorf("source directory %q: %w", extractSrc, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("source %q is not a directory", extractSrc)
	}

	// Set up detector and pattern registry
	patternReg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(patternReg)
	if err != nil {
		return fmt.Errorf("create persona detector: %w", err)
	}

	// Create orchestrator and run extraction
	orch := extractor.NewExtractorOrchestrator(det, patternReg)
	registry, manifest, err := orch.Extract(extractSrc)
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
