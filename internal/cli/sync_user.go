package cli

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const releaseURL = "https://github.com/yejune/godo/releases/latest/download/do-release.tar.gz"

var userSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Install or update Do framework from GitHub releases",
	Long: `Downloads the latest Do framework release and installs it into the current project.
Backs up existing framework files before updating.`,
	RunE: runUserSync,
}

func init() {
	rootCmd.AddCommand(userSyncCmd)
	userSyncCmd.Flags().Bool("dry-run", false, "Show what would be changed without making changes")
	userSyncCmd.Flags().Bool("init-custom", false, "Scaffold claude/ directory for custom overrides")
}

// fileExists returns true if the path exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirExists returns true if the path exists and is a directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// frameworkDirs lists the framework directories to backup/delete.
var frameworkDirs = []string{
	filepath.Join(".claude", "agents", "do"),
	filepath.Join(".claude", "commands", "do"),
	filepath.Join(".claude", "styles"),
}

// frameworkRuleGlob is the pattern for framework rule files.
const frameworkRuleGlob = ".claude/rules/dev-*.md"

// frameworkFiles lists standalone framework files to backup/delete.
var frameworkFiles = []string{
	"CLAUDE.md",
}

func runUserSync(cmd *cobra.Command, args []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	initCustom, _ := cmd.Flags().GetBool("init-custom")

	// --init-custom: scaffold claude/ directory and exit
	if initCustom {
		return scaffoldCustomDir(cmd)
	}

	// Step 1: Dev folder warning (stop with warning, not error)
	if fileExists("tobrew.yaml") && fileExists(filepath.Join("cmd", "godo", "main.go")) {
		fmt.Fprintln(cmd.ErrOrStderr(), "Warning: running in framework development directory")
		return nil
	}

	// --dry-run: show what would happen and exit
	if dryRun {
		return runDryRun(cmd)
	}

	// Step 2: Initialize global dir (~/.do/)
	if err := initGlobalDir(cmd); err != nil {
		return fmt.Errorf("init global dir: %w", err)
	}

	// Step 3: Register project
	if err := registerProject(cmd); err != nil {
		return fmt.Errorf("register project: %w", err)
	}

	// Step 4: Create settings.local.json
	EnsureSettingsLocal()

	// Step 5: Backup framework files
	backupDir, err := backupFrameworkFiles(cmd)
	if err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	// Set up restore-on-failure
	success := false
	defer func() {
		if !success && backupDir != "" {
			fmt.Fprintln(cmd.ErrOrStderr(), "Restoring from backup due to failure...")
			if restoreErr := restoreFromBackup(backupDir); restoreErr != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Warning: restore failed: %v\n", restoreErr)
			}
		}
	}()

	// Step 6: Delete framework files
	if err := deleteFrameworkFiles(cmd); err != nil {
		return fmt.Errorf("delete framework files: %w", err)
	}

	// Step 7: Download and extract framework
	fmt.Fprintf(cmd.ErrOrStderr(), "Downloading from %s ...\n", releaseURL)
	if err := downloadAndExtract(cmd, releaseURL); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	// Ensure required directories exist
	os.MkdirAll(".claude", 0755)
	os.MkdirAll(filepath.Join(".do", "config", "sections"), 0755)

	// Step 8: Verify installation
	if err := verifyInstallation(); err != nil {
		return fmt.Errorf("verify: %w", err)
	}

	// Step 9: Merge custom overrides
	if err := mergeCustomOverrides(cmd); err != nil {
		return fmt.Errorf("merge custom overrides: %w", err)
	}

	// Step 10: Sync product assets
	if err := syncProductAssets(cmd); err != nil {
		return fmt.Errorf("sync product assets: %w", err)
	}

	// Step 11: Clean old backups (keep 3 most recent)
	cleanOldBackups(cmd, 3)

	success = true
	fmt.Fprintln(cmd.ErrOrStderr(), "Sync complete.")
	return nil
}

// scaffoldCustomDir creates the claude/ directory structure for custom overrides.
func scaffoldCustomDir(cmd *cobra.Command) error {
	dirs := []string{
		"claude/agents",
		"claude/commands",
		"claude/styles",
		"claude/rules",
		"claude/skills",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("create %s: %w", d, err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Created %s/\n", d)
	}
	fmt.Fprintln(cmd.ErrOrStderr(), "Custom override directory scaffolded. Place your overrides in claude/.")
	return nil
}

// runDryRun shows what would be changed without making changes.
func runDryRun(cmd *cobra.Command) error {
	fmt.Fprintln(cmd.OutOrStdout(), "=== Dry Run ===")
	fmt.Fprintf(cmd.OutOrStdout(), "Release URL: %s\n", releaseURL)
	fmt.Fprintln(cmd.OutOrStdout(), "\nFramework directories to delete:")
	for _, d := range frameworkDirs {
		if dirExists(d) {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s (exists)\n", d)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s (not present)\n", d)
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), "\nSkill directories to delete:")
	fmt.Fprintf(cmd.OutOrStdout(), "  - .claude/skills/do-* patterns\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  - .claude/skills/moai-* patterns\n")

	fmt.Fprintln(cmd.OutOrStdout(), "\nFramework files to delete:")
	for _, f := range frameworkFiles {
		if fileExists(f) {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s (exists)\n", f)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s (not present)\n", f)
		}
	}

	ruleFiles, _ := filepath.Glob(frameworkRuleGlob)
	if len(ruleFiles) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "\nFramework rule files to delete:")
		for _, f := range ruleFiles {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", f)
		}
	}

	if dirExists("claude") {
		fmt.Fprintln(cmd.OutOrStdout(), "\nCustom overrides: claude/ directory found, will merge")
	}
	if dirExists(".product") {
		fmt.Fprintln(cmd.OutOrStdout(), "\nProduct assets: .product/ directory found, will sync")
	}

	fmt.Fprintln(cmd.OutOrStdout(), "\nNo changes made.")
	return nil
}

// initGlobalDir creates the ~/.do/ directory structure and default config files.
func initGlobalDir(cmd *cobra.Command) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	doDir := filepath.Join(home, ".do")
	dirs := []string{
		doDir,
		filepath.Join(doDir, "bin"),
		filepath.Join(doDir, "config"),
		filepath.Join(doDir, "logs"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("create %s: %w", d, err)
		}
	}

	// Create config.json if not exists
	configPath := filepath.Join(doDir, "config.json")
	if !fileExists(configPath) {
		data := []byte(`{"version": "1.0.0"}` + "\n")
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return fmt.Errorf("write config.json: %w", err)
		}
		fmt.Fprintln(cmd.ErrOrStderr(), "Created ~/.do/config.json")
	}

	// Create projects.json if not exists
	projectsPath := filepath.Join(doDir, "projects.json")
	if !fileExists(projectsPath) {
		data := []byte(`{"projects":[]}` + "\n")
		if err := os.WriteFile(projectsPath, data, 0644); err != nil {
			return fmt.Errorf("write projects.json: %w", err)
		}
		fmt.Fprintln(cmd.ErrOrStderr(), "Created ~/.do/projects.json")
	}

	return nil
}

// projectEntry represents a project registration in projects.json.
type projectEntry struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	RegisteredAt int64  `json:"registered_at"`
}

// projectsFile represents the structure of projects.json.
type projectsFile struct {
	Projects []projectEntry `json:"projects"`
}

// registerProject adds the current directory to ~/.do/projects.json.
func registerProject(cmd *cobra.Command) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get cwd: %w", err)
	}

	projectsPath := filepath.Join(home, ".do", "projects.json")
	data, err := os.ReadFile(projectsPath)
	if err != nil {
		return fmt.Errorf("read projects.json: %w", err)
	}

	var pf projectsFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return fmt.Errorf("parse projects.json: %w", err)
	}

	// Check if already registered
	for _, p := range pf.Projects {
		if p.Path == cwd {
			return nil // already registered
		}
	}

	// Add new project
	pf.Projects = append(pf.Projects, projectEntry{
		Path:         cwd,
		Name:         filepath.Base(cwd),
		RegisteredAt: time.Now().Unix(),
	})

	out, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal projects.json: %w", err)
	}
	if err := os.WriteFile(projectsPath, append(out, '\n'), 0644); err != nil {
		return fmt.Errorf("write projects.json: %w", err)
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Registered project: %s\n", cwd)
	return nil
}

// collectFrameworkSkillDirs returns skill directories matching do-* and moai-* patterns.
func collectFrameworkSkillDirs() []string {
	var dirs []string
	patterns := []string{
		filepath.Join(".claude", "skills", "do-*"),
		filepath.Join(".claude", "skills", "moai-*"),
	}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(pattern)
		for _, m := range matches {
			if dirExists(m) {
				dirs = append(dirs, m)
			}
		}
	}
	return dirs
}

// backupFrameworkFiles backs up framework files to .do/backup/<timestamp>/.
func backupFrameworkFiles(cmd *cobra.Command) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupBase := filepath.Join(".do", "backup", timestamp)

	hasContent := false

	// Backup directories
	for _, d := range frameworkDirs {
		if dirExists(d) {
			dst := filepath.Join(backupBase, d)
			if err := copyDir(d, dst); err != nil {
				return "", fmt.Errorf("backup dir %s: %w", d, err)
			}
			hasContent = true
		}
	}

	// Backup skill directories
	for _, d := range collectFrameworkSkillDirs() {
		dst := filepath.Join(backupBase, d)
		if err := copyDir(d, dst); err != nil {
			return "", fmt.Errorf("backup skill dir %s: %w", d, err)
		}
		hasContent = true
	}

	// Backup standalone files
	for _, f := range frameworkFiles {
		if fileExists(f) {
			dst := filepath.Join(backupBase, f)
			if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
				return "", err
			}
			if err := copyFile(f, dst); err != nil {
				return "", fmt.Errorf("backup file %s: %w", f, err)
			}
			hasContent = true
		}
	}

	// Backup rule files
	ruleFiles, _ := filepath.Glob(frameworkRuleGlob)
	for _, f := range ruleFiles {
		dst := filepath.Join(backupBase, f)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return "", err
		}
		if err := copyFile(f, dst); err != nil {
			return "", fmt.Errorf("backup rule %s: %w", f, err)
		}
		hasContent = true
	}

	if hasContent {
		fmt.Fprintf(cmd.ErrOrStderr(), "Backed up to %s\n", backupBase)
		return backupBase, nil
	}

	return "", nil
}

// restoreFromBackup restores framework files from a backup directory.
func restoreFromBackup(backupDir string) error {
	return filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(backupDir, path)
		if err != nil {
			return err
		}
		target := rel // restore to original relative location

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return copyFile(path, target)
	})
}

// deleteFrameworkFiles removes all framework files and directories.
func deleteFrameworkFiles(cmd *cobra.Command) error {
	// Delete directories
	for _, d := range frameworkDirs {
		if dirExists(d) {
			if err := os.RemoveAll(d); err != nil {
				return fmt.Errorf("remove %s: %w", d, err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Removed %s\n", d)
		}
	}

	// Delete skill directories
	for _, d := range collectFrameworkSkillDirs() {
		if err := os.RemoveAll(d); err != nil {
			return fmt.Errorf("remove skill %s: %w", d, err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Removed %s\n", d)
	}

	// Delete standalone files
	for _, f := range frameworkFiles {
		if fileExists(f) {
			if err := os.Remove(f); err != nil {
				return fmt.Errorf("remove %s: %w", f, err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Removed %s\n", f)
		}
	}

	// Delete rule files
	ruleFiles, _ := filepath.Glob(frameworkRuleGlob)
	for _, f := range ruleFiles {
		if err := os.Remove(f); err != nil {
			return fmt.Errorf("remove rule %s: %w", f, err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Removed %s\n", f)
	}

	return nil
}

// downloadAndExtract downloads a tar.gz from url and extracts it to the current directory.
func downloadAndExtract(cmd *cobra.Command, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	count := 0

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar read: %w", err)
		}

		// Skip install.sh
		name := header.Name
		if filepath.Base(name) == "install.sh" {
			continue
		}

		// Clean the path: strip leading ./ or /
		name = filepath.Clean(name)
		if name == "." {
			continue
		}

		// Security: prevent path traversal
		if strings.Contains(name, "..") {
			continue
		}

		target := name

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("mkdir %s: %w", target, err)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("mkdir parent %s: %w", target, err)
			}

			mode := os.FileMode(header.Mode)
			if mode == 0 {
				mode = 0644
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
			if err != nil {
				return fmt.Errorf("create %s: %w", target, err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("write %s: %w", target, err)
			}
			f.Close()
			count++
		}
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Extracted %d files\n", count)
	return nil
}

// verifyInstallation checks that critical files exist after download.
func verifyInstallation() error {
	required := []string{
		"CLAUDE.md",
		filepath.Join(".claude", "settings.json"),
	}
	for _, f := range required {
		if !fileExists(f) {
			return fmt.Errorf("missing required file after download: %s", f)
		}
	}
	return nil
}

// mergeCustomOverrides copies files from claude/ (dot-less) on top of .claude/.
func mergeCustomOverrides(cmd *cobra.Command) error {
	if !dirExists("claude") {
		return nil
	}

	fmt.Fprintln(cmd.ErrOrStderr(), "Merging custom overrides from claude/ ...")

	// Directory mappings: source -> destination
	dirMappings := map[string]string{
		"claude/agents":   filepath.Join(".claude", "agents", "do"),
		"claude/commands": filepath.Join(".claude", "commands", "do"),
		"claude/styles":   filepath.Join(".claude", "styles"),
		"claude/rules":    filepath.Join(".claude", "rules"),
		"claude/skills":   filepath.Join(".claude", "skills"),
	}

	for src, dst := range dirMappings {
		if dirExists(src) {
			if err := copyDir(src, dst); err != nil {
				return fmt.Errorf("merge %s -> %s: %w", src, dst, err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "  Merged %s -> %s\n", src, dst)
		}
	}

	// Handle CLAUDE.md override
	claudeFullOverride := filepath.Join("claude", "CLAUDE.md")
	claudeAppend := filepath.Join("claude", "CLAUDE.append.md")

	if fileExists(claudeFullOverride) {
		// Full override: replace CLAUDE.md entirely
		if err := copyFile(claudeFullOverride, "CLAUDE.md"); err != nil {
			return fmt.Errorf("override CLAUDE.md: %w", err)
		}
		fmt.Fprintln(cmd.ErrOrStderr(), "  Replaced CLAUDE.md with claude/CLAUDE.md")
	} else if fileExists(claudeAppend) {
		// Append mode: add content with markers
		appendData, err := os.ReadFile(claudeAppend)
		if err != nil {
			return fmt.Errorf("read CLAUDE.append.md: %w", err)
		}

		existing, err := os.ReadFile("CLAUDE.md")
		if err != nil {
			return fmt.Errorf("read CLAUDE.md: %w", err)
		}

		// Remove existing custom section if present
		marker := "<!-- DO:CUSTOM -->"
		if idx := strings.Index(string(existing), marker); idx >= 0 {
			existing = []byte(strings.TrimRight(string(existing[:idx]), "\n"))
		}

		// Append with markers
		combined := fmt.Sprintf("%s\n\n%s\n%s\n", string(existing), marker, string(appendData))
		if err := os.WriteFile("CLAUDE.md", []byte(combined), 0644); err != nil {
			return fmt.Errorf("write CLAUDE.md with append: %w", err)
		}
		fmt.Fprintln(cmd.ErrOrStderr(), "  Appended claude/CLAUDE.append.md to CLAUDE.md")
	}

	return nil
}

// syncProductAssets copies files from .product/ to .claude/.
func syncProductAssets(cmd *cobra.Command) error {
	if !dirExists(".product") {
		return nil
	}

	fmt.Fprintln(cmd.ErrOrStderr(), "Syncing product assets from .product/ ...")

	// Directory mappings
	dirMappings := map[string]string{
		".product/agents":   filepath.Join(".claude", "agents", "product"),
		".product/skills":   filepath.Join(".claude", "skills", "product"),
		".product/commands": filepath.Join(".claude", "commands", "product"),
		".product/rules":    filepath.Join(".claude", "rules", "product"),
		".product/hooks":    filepath.Join(".claude", "hooks", "product"),
	}

	for src, dst := range dirMappings {
		if dirExists(src) {
			if err := copyDir(src, dst); err != nil {
				return fmt.Errorf("sync %s -> %s: %w", src, dst, err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "  Synced %s -> %s\n", src, dst)
		}
	}

	// Handle .product/CLAUDE.md -> .claude/rules/product.md
	productClaude := filepath.Join(".product", "CLAUDE.md")
	if fileExists(productClaude) {
		dst := filepath.Join(".claude", "rules", "product.md")
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return fmt.Errorf("create rules dir: %w", err)
		}
		if err := copyFile(productClaude, dst); err != nil {
			return fmt.Errorf("sync product CLAUDE.md: %w", err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "  Synced %s -> %s\n", productClaude, dst)
	}

	return nil
}

// cleanOldBackups removes old backup directories, keeping the most recent 'keep' count.
func cleanOldBackups(cmd *cobra.Command, keep int) {
	backupBase := filepath.Join(".do", "backup")
	if !dirExists(backupBase) {
		return
	}

	entries, err := os.ReadDir(backupBase)
	if err != nil {
		return
	}

	// Filter to directories only
	var dirs []os.DirEntry
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e)
		}
	}

	if len(dirs) <= keep {
		return
	}

	// Sort by name (timestamp format ensures lexicographic = chronological)
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	// Remove oldest, keeping the last 'keep' entries
	toRemove := dirs[:len(dirs)-keep]
	for _, d := range toRemove {
		path := filepath.Join(backupBase, d.Name())
		if err := os.RemoveAll(path); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to remove old backup %s: %v\n", path, err)
			continue
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Cleaned old backup: %s\n", d.Name())
	}
}
