package hook

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupGitRepo(t *testing.T) (string, func()) {
	t.Helper()
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Init git repo
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	exec.Command("git", "config", "user.name", "test").Run()

	// Create initial commit
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("init"), 0644)
	exec.Command("git", "add", "README.md").Run()
	exec.Command("git", "commit", "-m", "init").Run()

	return tmpDir, func() { os.Chdir(origDir) }
}

func Test_GitStatus_clean_repo(t *testing.T) {
	_, cleanup := setupGitRepo(t)
	defer cleanup()

	hasChanges, _ := GitStatus()
	if hasChanges {
		t.Error("expected no changes in clean repo")
	}
}

func Test_GitStatus_with_uncommitted_changes(t *testing.T) {
	tmpDir, cleanup := setupGitRepo(t)
	defer cleanup()

	// Modify a tracked file to create uncommitted changes
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("modified"), 0644)

	hasChanges, summary := GitStatus()
	if !hasChanges {
		t.Error("expected changes with uncommitted file")
	}
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func Test_GitStatus_with_staged_changes(t *testing.T) {
	tmpDir, cleanup := setupGitRepo(t)
	defer cleanup()

	os.WriteFile(filepath.Join(tmpDir, "staged.go"), []byte("package main"), 0644)
	exec.Command("git", "add", "staged.go").Run()

	hasChanges, _ := GitStatus()
	if !hasChanges {
		t.Error("expected changes with staged file")
	}
}

func Test_GitStatus_not_a_git_repo(t *testing.T) {
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	hasChanges, _ := GitStatus()
	if hasChanges {
		t.Error("expected no changes in non-git directory")
	}
}

func Test_GitStatus_ignores_untracked_files(t *testing.T) {
	tmpDir, cleanup := setupGitRepo(t)
	defer cleanup()

	// Create an untracked file (not staged, not committed)
	os.WriteFile(filepath.Join(tmpDir, "untracked.txt"), []byte("hello"), 0644)

	hasChanges, _ := GitStatus()
	if hasChanges {
		t.Error("expected no changes when only untracked files exist")
	}
}

func Test_GitStatus_mockable(t *testing.T) {
	// Verify GitStatus is a var that can be overridden
	original := GitStatus
	defer func() { GitStatus = original }()

	GitStatus = func() (bool, string) {
		return true, "M fake.go"
	}
	has, summary := GitStatus()
	if !has || summary != "M fake.go" {
		t.Error("mock override did not work")
	}
}
