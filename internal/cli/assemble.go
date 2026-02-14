package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/do-focus/convert/internal/assembler"
	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/template"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var assembleCmd = &cobra.Command{
	Use:   "assemble",
	Short: "Assemble core templates + persona into deployable .claude/ directory",
	Long: `Assemble merges core templates with a persona manifest to produce
a complete .claude/ directory. Core template slots are filled with
persona-specific content, and persona-only files are copied to the output.`,
	RunE: runAssemble,
}

var (
	assembleCoreDir    string
	assemblePersona    string
	assembleOutputDir  string
)

func init() {
	assembleCmd.Flags().StringVar(&assembleCoreDir, "core", "", "path to core templates directory")
	assembleCmd.Flags().StringVar(&assemblePersona, "persona", "", "path to persona manifest file (persona.yaml)")
	assembleCmd.Flags().StringVar(&assembleOutputDir, "out", "", "output directory for assembled .claude/")

	assembleCmd.MarkFlagRequired("core")
	assembleCmd.MarkFlagRequired("persona")
	assembleCmd.MarkFlagRequired("out")

	rootCmd.AddCommand(assembleCmd)
}

func runAssemble(cmd *cobra.Command, args []string) error {
	// Load template registry from core directory.
	registry, err := template.LoadRegistry(assembleCoreDir)
	if err != nil {
		return fmt.Errorf("load registry: %w", err)
	}

	// Load persona manifest.
	manifestData, err := os.ReadFile(assemblePersona)
	if err != nil {
		return fmt.Errorf("read persona manifest %s: %w", assemblePersona, err)
	}

	var manifest model.PersonaManifest
	if err := yaml.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("parse persona manifest %s: %w", assemblePersona, err)
	}

	// Persona directory is the parent of the manifest file.
	personaDir := filepath.Dir(assemblePersona)

	// Create assembler and run.
	asm := assembler.NewAssembler(assembleCoreDir, personaDir, assembleOutputDir, &manifest, registry)
	result, err := asm.Assemble()
	if err != nil {
		return fmt.Errorf("assemble: %w", err)
	}

	// Print summary.
	fmt.Fprintf(cmd.OutOrStdout(), "%d files assembled to %s/\n", result.FilesWritten, assembleOutputDir)

	if len(result.Warnings) > 0 {
		for _, w := range result.Warnings {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: %s\n", w)
		}
	}

	return nil
}
