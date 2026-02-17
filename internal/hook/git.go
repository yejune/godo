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
	// Filter out untracked files (??) — only report tracked file changes
	var tracked []string
	for _, line := range strings.Split(output, "\n") {
		if !strings.HasPrefix(line, "?? ") {
			tracked = append(tracked, line)
		}
	}
	if len(tracked) == 0 {
		return false, ""
	}
	return true, strings.Join(tracked, "\n")
}
