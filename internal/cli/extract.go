package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yejune/godo/internal/detector"
	"github.com/yejune/godo/internal/extractor"
	"github.com/yejune/godo/internal/model"
	"github.com/yejune/godo/internal/template"
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

	// Create BrandSlotifier from detected persona name.
	// This will be used to strip brand prefixes from core skill paths
	// and replace brand references in core file content with slot variables.
	personaName := extractPersona
	if personaName == "" {
		personaName = manifest.Name
	}
	slotifier := extractor.NewBrandSlotifier(personaName)

	// Remap FoundIn paths in registry slots before saving.
	// This ensures registry.yaml references the stripped skill names
	// (e.g., "skills/lang-python/SKILL.md" instead of "skills/moai-lang-python/SKILL.md").
	if slotifier != nil {
		remapRegistryPaths(registry, slotifier)
	}

	// Save core templates (registry.yaml) to <out>/core/
	coreDir := filepath.Join(extractOut, "core")
	if err := registry.Save(coreDir); err != nil {
		return fmt.Errorf("save core templates: %w", err)
	}

	// Copy core files to <out>/core/ with brand slotification.
	// For each core file:
	//   - RemapCorePath: strip brand prefix from skill dir names in the output path
	//   - SlotifyContent: replace brand references with {{slot:BRAND}} etc. in text files
	coreFileCount := 0
	for _, relPath := range manifest.CoreFiles {
		src := filepath.Join(manifest.SourceDir, relPath)
		dstRelPath := relPath
		if slotifier != nil {
			dstRelPath = slotifier.RemapCorePath(relPath)
			dstRelPath = slotifier.StripBrandSubdir(dstRelPath)
		}
		dst := filepath.Join(coreDir, dstRelPath)

		if slotifier != nil && isTextFile(relPath) {
			if err := copyFileSlotified(src, dst, slotifier); err != nil {
				return fmt.Errorf("copy core file %s: %w", relPath, err)
			}
		} else {
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("copy core file %s: %w", relPath, err)
			}
		}
		coreFileCount++
	}

	// Determine persona name for output directory
	outputPersonaName := extractPersona
	if outputPersonaName == "" {
		outputPersonaName = manifest.Name
	}
	if outputPersonaName == "" {
		outputPersonaName = "default"
	}

	// Populate brand identity fields on the manifest.
	// These are used by the assembler to deslotify core templates
	// (replacing {{slot:BRAND}} etc. with the persona's brand values).
	if outputPersonaName != "default" {
		manifest.Brand = outputPersonaName
		manifest.BrandDir = outputPersonaName
		manifest.BrandCmd = outputPersonaName
	}

	// Save persona manifest to <out>/personas/<name>/manifest.yaml
	personaDir := filepath.Join(extractOut, "personas", outputPersonaName)
	if err := os.MkdirAll(personaDir, 0o755); err != nil {
		return fmt.Errorf("create persona dir %s: %w", personaDir, err)
	}

	// Remap persona file paths: strip brand subdirectory for cleaner storage.
	// e.g., agents/moai/manager-ddd.md → agents/manager-ddd.md
	// The assembler will add the brand subdir back during assembly.
	if slotifier != nil {
		remapManifestPersonaPaths(manifest, slotifier)
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

// remapRegistryPaths updates all FoundIn paths in registry slot entries
// to use stripped skill names (e.g., "skills/lang-python/" instead of "skills/moai-lang-python/").
func remapRegistryPaths(reg *template.Registry, slotifier *extractor.BrandSlotifier) {
	for _, entry := range reg.Slots {
		entry.Description = slotifier.SlotifyContent(entry.Description)
		for i := range entry.FoundIn {
			entry.FoundIn[i].Path = slotifier.RemapCorePath(entry.FoundIn[i].Path)
			entry.FoundIn[i].Path = slotifier.StripBrandSubdir(entry.FoundIn[i].Path)
		}
	}
}

// textExtensions lists file extensions that should have brand references slotified.
var textExtensions = map[string]bool{
	".md":   true,
	".yaml": true,
	".yml":  true,
	".json": true,
	".py":   true,
	".sh":   true,
	".go":   true,
	".ts":   true,
	".js":   true,
	".tsx":  true,
	".jsx":  true,
	".txt":  true,
	".toml": true,
	".cfg":  true,
	".ini":  true,
	".html": true,
	".css":  true,
}

// isTextFile returns true if the file extension indicates a text file
// whose content should be slotified.
func isTextFile(relPath string) bool {
	ext := strings.ToLower(filepath.Ext(relPath))
	return textExtensions[ext]
}

// copyFileSlotified reads a text file, applies brand slotification to its content,
// and writes the result to the destination path.
func copyFileSlotified(src, dst string, slotifier *extractor.BrandSlotifier) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read source %s: %w", src, err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source %s: %w", src, err)
	}

	content := slotifier.SlotifyContent(string(data))

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create parent dir for %s: %w", dst, err)
	}

	if err := os.WriteFile(dst, []byte(content), srcInfo.Mode().Perm()); err != nil {
		return fmt.Errorf("write destination %s: %w", dst, err)
	}

	return nil
}

// remapManifestPersonaPaths strips brand subdirectories from all persona file paths
// in the manifest. This ensures persona files are stored without brand nesting.
// The assembler will add the brand subdir back during assembly.
func remapManifestPersonaPaths(manifest *model.PersonaManifest, slotifier *extractor.BrandSlotifier) {
	// Remap PersonaFiles map keys.
	newPersonaFiles := make(map[string]string, len(manifest.PersonaFiles))
	for relPath, srcPath := range manifest.PersonaFiles {
		stripped := slotifier.StripBrandSubdir(relPath)
		newPersonaFiles[stripped] = srcPath
	}
	manifest.PersonaFiles = newPersonaFiles

	// Remap named persona file lists.
	manifest.Agents = remapPaths(manifest.Agents, slotifier)
	manifest.Rules = remapPaths(manifest.Rules, slotifier)
	manifest.Commands = remapPaths(manifest.Commands, slotifier)
	manifest.HookScripts = remapPaths(manifest.HookScripts, slotifier)
	manifest.Styles = remapPaths(manifest.Styles, slotifier)
	manifest.Characters = remapPaths(manifest.Characters, slotifier)
	manifest.Spinners = remapPaths(manifest.Spinners, slotifier)
	manifest.Skills = remapPaths(manifest.Skills, slotifier)

	// Remap AgentPatches map keys.
	if len(manifest.AgentPatches) > 0 {
		newPatches := make(map[string]*model.AgentPatch, len(manifest.AgentPatches))
		for relPath, patch := range manifest.AgentPatches {
			stripped := slotifier.StripBrandSubdir(relPath)
			newPatches[stripped] = patch
		}
		manifest.AgentPatches = newPatches
	}
}

// remapPaths applies StripBrandSubdir to each path in the slice.
func remapPaths(paths []string, slotifier *extractor.BrandSlotifier) []string {
	if len(paths) == 0 {
		return paths
	}
	result := make([]string, len(paths))
	for i, p := range paths {
		result[i] = slotifier.StripBrandSubdir(p)
	}
	return result
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
