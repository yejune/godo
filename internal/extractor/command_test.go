package extractor

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestExtractCommands_MixedFiles(t *testing.T) {
	sourceDir := t.TempDir()
	coreCmdDir := t.TempDir()
	personaCmdDir := t.TempDir()

	// Create persona commands (prefixed names)
	personaFiles := []string{
		"moai-plan.md",
		"do:setup.md",
		"do-init.md",
	}
	for _, f := range personaFiles {
		if err := os.WriteFile(filepath.Join(sourceDir, f), []byte("# "+f), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create core commands (generic names)
	coreFiles := []string{
		"build.md",
		"test.md",
		"deploy.md",
	}
	for _, f := range coreFiles {
		if err := os.WriteFile(filepath.Join(sourceDir, f), []byte("# "+f), 0644); err != nil {
			t.Fatal(err)
		}
	}

	gotCore, gotPersona, err := ExtractCommands(sourceDir, coreCmdDir, personaCmdDir)
	if err != nil {
		t.Fatalf("ExtractCommands() error: %v", err)
	}

	sort.Strings(gotCore)
	sort.Strings(gotPersona)

	if len(gotCore) != 3 {
		t.Errorf("core count = %d, want 3; got %v", len(gotCore), gotCore)
	}
	if len(gotPersona) != 3 {
		t.Errorf("persona count = %d, want 3; got %v", len(gotPersona), gotPersona)
	}

	// Verify core files were copied
	for _, f := range coreFiles {
		dst := filepath.Join(coreCmdDir, f)
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			t.Errorf("core file %q not copied to %s", f, coreCmdDir)
		}
	}

	// Verify persona files were copied
	for _, f := range personaFiles {
		dst := filepath.Join(personaCmdDir, f)
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			t.Errorf("persona file %q not copied to %s", f, personaCmdDir)
		}
	}
}

func TestExtractCommands_SubdirectoryPreserved(t *testing.T) {
	sourceDir := t.TempDir()
	coreCmdDir := t.TempDir()
	personaCmdDir := t.TempDir()

	// Create a subdirectory with a persona command
	subDir := filepath.Join(sourceDir, "workflows")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "moai-deploy.md"), []byte("# deploy"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a subdirectory with a core command
	if err := os.WriteFile(filepath.Join(subDir, "lint.md"), []byte("# lint"), 0644); err != nil {
		t.Fatal(err)
	}

	_, _, err := ExtractCommands(sourceDir, coreCmdDir, personaCmdDir)
	if err != nil {
		t.Fatalf("ExtractCommands() error: %v", err)
	}

	// Check subdirectory structure is preserved
	personaDst := filepath.Join(personaCmdDir, "workflows", "moai-deploy.md")
	if _, err := os.Stat(personaDst); os.IsNotExist(err) {
		t.Errorf("persona file not copied with subdirectory structure: %s", personaDst)
	}

	coreDst := filepath.Join(coreCmdDir, "workflows", "lint.md")
	if _, err := os.Stat(coreDst); os.IsNotExist(err) {
		t.Errorf("core file not copied with subdirectory structure: %s", coreDst)
	}
}

func TestExtractCommands_EmptyDir(t *testing.T) {
	sourceDir := t.TempDir()
	coreCmdDir := t.TempDir()
	personaCmdDir := t.TempDir()

	gotCore, gotPersona, err := ExtractCommands(sourceDir, coreCmdDir, personaCmdDir)
	if err != nil {
		t.Fatalf("ExtractCommands() error: %v", err)
	}
	if len(gotCore) != 0 {
		t.Errorf("core count = %d, want 0", len(gotCore))
	}
	if len(gotPersona) != 0 {
		t.Errorf("persona count = %d, want 0", len(gotPersona))
	}
}

func TestExtractCommands_NonexistentSourceDir(t *testing.T) {
	_, _, err := ExtractCommands("/nonexistent/path", t.TempDir(), t.TempDir())
	if err == nil {
		t.Error("ExtractCommands() expected error for nonexistent source, got nil")
	}
}

func TestExtractHookScripts_PersonaBinaryReference(t *testing.T) {
	sourceDir := t.TempDir()
	coreHookDir := t.TempDir()
	personaHookDir := t.TempDir()

	// Script that references godo binary
	godoScript := "#!/bin/bash\nset -euo pipefail\ngodo hook pre-tool \"$@\"\n"
	if err := os.WriteFile(filepath.Join(sourceDir, "pre-tool.sh"), []byte(godoScript), 0755); err != nil {
		t.Fatal(err)
	}

	// Script that references moai binary
	moaiScript := "#!/bin/bash\nmoai validate --strict\n"
	if err := os.WriteFile(filepath.Join(sourceDir, "validate.sh"), []byte(moaiScript), 0755); err != nil {
		t.Fatal(err)
	}

	_, gotPersona, err := ExtractHookScripts(sourceDir, coreHookDir, personaHookDir)
	if err != nil {
		t.Fatalf("ExtractHookScripts() error: %v", err)
	}

	if len(gotPersona) != 2 {
		t.Errorf("persona count = %d, want 2; got %v", len(gotPersona), gotPersona)
	}

	// Verify files were copied to persona dir
	for _, name := range []string{"pre-tool.sh", "validate.sh"} {
		dst := filepath.Join(personaHookDir, name)
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			t.Errorf("persona hook script %q not copied to %s", name, personaHookDir)
		}
	}
}

func TestExtractHookScripts_NoPersonaReference(t *testing.T) {
	sourceDir := t.TempDir()
	coreHookDir := t.TempDir()
	personaHookDir := t.TempDir()

	// Generic script with no persona binary reference
	genericScript := "#!/bin/bash\nset -euo pipefail\necho \"Running pre-commit checks\"\nnpm run lint\n"
	if err := os.WriteFile(filepath.Join(sourceDir, "pre-commit.sh"), []byte(genericScript), 0644); err != nil {
		t.Fatal(err)
	}

	gotCore, gotPersona, err := ExtractHookScripts(sourceDir, coreHookDir, personaHookDir)
	if err != nil {
		t.Fatalf("ExtractHookScripts() error: %v", err)
	}

	if len(gotCore) != 1 {
		t.Errorf("core count = %d, want 1; got %v", len(gotCore), gotCore)
	}
	if len(gotPersona) != 0 {
		t.Errorf("persona count = %d, want 0; got %v", len(gotPersona), gotPersona)
	}

	// Verify file was copied to core dir
	dst := filepath.Join(coreHookDir, "pre-commit.sh")
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Errorf("core hook script not copied to %s", dst)
	}
}

func TestExtractHookScripts_Subdirectory(t *testing.T) {
	sourceDir := t.TempDir()
	coreHookDir := t.TempDir()
	personaHookDir := t.TempDir()

	// Create subdirectory (e.g., .claude/hooks/moai/)
	subDir := filepath.Join(sourceDir, "moai")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	personaScript := "#!/bin/bash\ngodo hook session-start\n"
	if err := os.WriteFile(filepath.Join(subDir, "session-start.sh"), []byte(personaScript), 0755); err != nil {
		t.Fatal(err)
	}

	_, gotPersona, err := ExtractHookScripts(sourceDir, coreHookDir, personaHookDir)
	if err != nil {
		t.Fatalf("ExtractHookScripts() error: %v", err)
	}

	if len(gotPersona) != 1 {
		t.Errorf("persona count = %d, want 1", len(gotPersona))
	}

	// Verify subdirectory structure preserved
	dst := filepath.Join(personaHookDir, "moai", "session-start.sh")
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Errorf("persona hook not copied with subdirectory: %s", dst)
	}
}

func TestExtractHookScripts_EmptyDir(t *testing.T) {
	sourceDir := t.TempDir()
	coreHookDir := t.TempDir()
	personaHookDir := t.TempDir()

	gotCore, gotPersona, err := ExtractHookScripts(sourceDir, coreHookDir, personaHookDir)
	if err != nil {
		t.Fatalf("ExtractHookScripts() error: %v", err)
	}
	if len(gotCore) != 0 {
		t.Errorf("core count = %d, want 0", len(gotCore))
	}
	if len(gotPersona) != 0 {
		t.Errorf("persona count = %d, want 0", len(gotPersona))
	}
}

func TestExtractHookScripts_FalsePositiveAvoidance(t *testing.T) {
	sourceDir := t.TempDir()
	coreHookDir := t.TempDir()
	personaHookDir := t.TempDir()

	// Script that mentions "godo" without a space after (not a binary call)
	script := "#!/bin/bash\n# This references godocker, not the godo binary\necho \"starting godocker container\"\n"
	if err := os.WriteFile(filepath.Join(sourceDir, "docker.sh"), []byte(script), 0644); err != nil {
		t.Fatal(err)
	}

	gotCore, gotPersona, err := ExtractHookScripts(sourceDir, coreHookDir, personaHookDir)
	if err != nil {
		t.Fatalf("ExtractHookScripts() error: %v", err)
	}

	if len(gotCore) != 1 {
		t.Errorf("core count = %d, want 1 (false positive avoided)", len(gotCore))
	}
	if len(gotPersona) != 0 {
		t.Errorf("persona count = %d, want 0 (false positive avoided)", len(gotPersona))
	}
}

func TestCopyFile(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	content := []byte("hello world\n")
	srcPath := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Copy to a nested destination (parent dirs should be created)
	dstPath := filepath.Join(dstDir, "sub", "dir", "test.txt")
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile() error: %v", err)
	}

	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("reading copied file: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("copied content = %q, want %q", got, content)
	}
}
