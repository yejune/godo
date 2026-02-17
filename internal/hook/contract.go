package hook

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// ErrContractFail indicates the hook execution contract was violated.
var ErrContractFail = errors.New("hook: execution contract violated")

// Contract validates the hook execution environment.
type Contract struct {
	WorkDir string
}

// NewContract creates a new contract for the given working directory.
func NewContract(workDir string) *Contract {
	return &Contract{WorkDir: workDir}
}

// Validate checks that the execution environment meets contract requirements.
// It verifies:
//   - The context is not already cancelled or expired
//   - The working directory is specified and accessible
func (c *Contract) Validate(ctx context.Context) error {
	// Check context is still valid
	if ctx.Err() != nil {
		return fmt.Errorf("%w: context already done: %v", ErrContractFail, ctx.Err())
	}

	// Check working directory is specified
	if c.WorkDir == "" {
		return fmt.Errorf("%w: working directory not specified", ErrContractFail)
	}

	// Check working directory exists and is accessible
	info, err := os.Stat(c.WorkDir)
	if err != nil {
		return fmt.Errorf("%w: working directory not accessible: %v", ErrContractFail, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: path is not a directory: %s", ErrContractFail, c.WorkDir)
	}

	return nil
}

// Guarantees returns the list of guaranteed execution conditions.
func (c *Contract) Guarantees() []string {
	return []string{
		"stdin: valid JSON conforming to Claude Code hook protocol",
		"exit code: 0 (allow/success), 2 (block), other (non-blocking error)",
		"timeout: configurable via context.WithTimeout (default 30s)",
		"working directory: project root ($CLAUDE_PROJECT_DIR)",
	}
}

// NonGuarantees returns the list of non-guaranteed execution conditions.
func (c *Contract) NonGuarantees() []string {
	return []string{
		"User PATH: binary must be in system PATH",
		"shell environment variables: .bashrc/.zshrc are not loaded",
		"shell functions or alias definitions",
		"Python/Node.js/uv runtime availability",
	}
}
