package hook

import (
	"encoding/json"
	"strings"
)

// HandlePreTool handles the PreToolUse hook event.
// It checks file paths and bash commands against security policies.
func HandlePreTool(input *Input) *Output {
	policy := DefaultSecurityPolicy()
	toolName := input.ToolName

	switch toolName {
	case "Write", "Edit", "Read", "Glob":
		return checkFileAccess(policy, input)
	case "Bash":
		return checkBashCommand(policy, input)
	default:
		return NewAllowOutput()
	}
}

// checkFileAccess validates file tool access against security patterns.
func checkFileAccess(policy *SecurityPolicy, input *Input) *Output {
	filePath := extractFilePath(input.ToolInput)
	if filePath == "" {
		return NewAllowOutput()
	}

	// Check deny patterns
	for _, re := range policy.DenyFilePatterns {
		if re.MatchString(filePath) {
			return NewDenyOutput("Blocked: file matches security deny pattern: " + filePath)
		}
	}

	// Check ask patterns (only for write operations)
	if input.ToolName == "Write" || input.ToolName == "Edit" {
		for _, re := range policy.AskFilePatterns {
			if re.MatchString(filePath) {
				return NewAskOutput("File requires confirmation: " + filePath)
			}
		}
	}

	return NewAllowOutput()
}

// checkBashCommand validates bash commands against security patterns.
func checkBashCommand(policy *SecurityPolicy, input *Input) *Output {
	command := extractCommand(input.ToolInput)
	if command == "" {
		return NewAllowOutput()
	}

	// Check deny patterns
	for _, re := range policy.DenyBashPatterns {
		if re.MatchString(command) {
			return NewDenyOutput("Blocked: command matches security deny pattern: " + command)
		}
	}

	// Check ask patterns
	for _, re := range policy.AskBashPatterns {
		if re.MatchString(command) {
			return NewAskOutput("Command requires confirmation: " + command)
		}
	}

	// Check for sensitive content in command
	for _, re := range policy.SensitiveContentPatterns {
		if re.MatchString(command) {
			return NewDenyOutput("Blocked: command contains sensitive content")
		}
	}

	return NewAllowOutput()
}

// extractFilePath extracts the file_path from tool input JSON.
func extractFilePath(toolInput json.RawMessage) string {
	if len(toolInput) == 0 {
		return ""
	}
	var data map[string]interface{}
	if json.Unmarshal(toolInput, &data) != nil {
		return ""
	}
	// Try file_path first, then path
	if fp, ok := data["file_path"].(string); ok {
		return fp
	}
	if fp, ok := data["path"].(string); ok {
		return fp
	}
	return ""
}

// extractCommand extracts the command from Bash tool input JSON.
func extractCommand(toolInput json.RawMessage) string {
	if len(toolInput) == 0 {
		return ""
	}
	var data map[string]interface{}
	if json.Unmarshal(toolInput, &data) != nil {
		return ""
	}
	if cmd, ok := data["command"].(string); ok {
		return strings.TrimSpace(cmd)
	}
	return ""
}
