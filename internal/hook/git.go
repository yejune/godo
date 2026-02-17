package hook

import (
	"os/exec"
	"strings"
)

// GitStatus checks for uncommitted changes in the current working directory.
// Returns true if there are uncommitted changes, and a summary string.
var GitStatus = func() (bool, string) {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		// Not a git repo or git not available — skip check
		return false, ""
	}
	output := strings.TrimSpace(string(out))
	if output == "" {
		return false, ""
	}
	// Count changed files
	lines := strings.Split(output, "\n")
	return true, strings.Join(lines, "\n")
}
