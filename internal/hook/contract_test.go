package hook

import (
	"context"
	"errors"
	"os"
	"testing"
)

func Test_NewContract_creates_contract(t *testing.T) {
	c := NewContract("/tmp")
	if c == nil {
		t.Fatal("expected non-nil Contract")
	}
	if c.WorkDir != "/tmp" {
		t.Errorf("WorkDir: got %q, want %q", c.WorkDir, "/tmp")
	}
}

func Test_Contract_Validate_valid_directory(t *testing.T) {
	dir := t.TempDir()
	c := NewContract(dir)
	err := c.Validate(context.Background())
	if err != nil {
		t.Errorf("expected no error for valid directory, got: %v", err)
	}
}

func Test_Contract_Validate_empty_workdir(t *testing.T) {
	c := NewContract("")
	err := c.Validate(context.Background())
	if err == nil {
		t.Fatal("expected error for empty working directory")
	}
	if !errors.Is(err, ErrContractFail) {
		t.Errorf("expected ErrContractFail, got: %v", err)
	}
}

func Test_Contract_Validate_nonexistent_directory(t *testing.T) {
	c := NewContract("/nonexistent/path/that/does/not/exist")
	err := c.Validate(context.Background())
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
	if !errors.Is(err, ErrContractFail) {
		t.Errorf("expected ErrContractFail, got: %v", err)
	}
}

func Test_Contract_Validate_cancelled_context(t *testing.T) {
	dir := t.TempDir()
	c := NewContract(dir)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := c.Validate(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	if !errors.Is(err, ErrContractFail) {
		t.Errorf("expected ErrContractFail, got: %v", err)
	}
}

func Test_Contract_Validate_path_is_file_not_directory(t *testing.T) {
	// Create a temp file (not a directory)
	f, err := os.CreateTemp("", "contract-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())

	c := NewContract(f.Name())
	err = c.Validate(context.Background())
	if err == nil {
		t.Fatal("expected error when path is a file, not a directory")
	}
	if !errors.Is(err, ErrContractFail) {
		t.Errorf("expected ErrContractFail, got: %v", err)
	}
}

func Test_Contract_Guarantees_returns_non_empty(t *testing.T) {
	c := NewContract("/tmp")
	guarantees := c.Guarantees()
	if len(guarantees) == 0 {
		t.Error("expected non-empty guarantees list")
	}
}

func Test_Contract_NonGuarantees_returns_non_empty(t *testing.T) {
	c := NewContract("/tmp")
	nonGuarantees := c.NonGuarantees()
	if len(nonGuarantees) == 0 {
		t.Error("expected non-empty non-guarantees list")
	}
}
