package hook

import (
	"encoding/json"
	"os"
	"strings"
)

// HandlePreTool handles the PreToolUse hook event.
// It checks file paths, bash commands, and enforces checklist requirements.
func HandlePreTool(input *Input) *Output {
	policy := DefaultSecurityPolicy()
	toolName := input.ToolName

	// [HARD] All code modification tools require valid checklist (only if .do/jobs exists)
	jobsDir := ".do/jobs"
	if _, err := os.Stat(jobsDir); err == nil {
		// Project has Do framework - enforce checklist
		switch toolName {
		case "Task", "Write", "Edit":
			if output := checkChecklistRequirement(input); output != nil {
				return output
			}
		}
	}

	switch toolName {
	case "Write", "Edit", "Read", "Glob":
		return checkFileAccess(policy, input)
	case "Bash":
		return checkBashCommand(policy, input)
	default:
		return NewAllowOutput()
	}
}

// checkChecklistRequirement enforces that Task tool requires a valid sub-checklist path.
// [HARD] Proceeding with development while checklist is unwritten is a VIOLATION
func checkChecklistRequirement(input *Input) *Output {
	// Extract prompt from Task tool input
	prompt := extractPrompt(input.ToolInput)
	if prompt == "" {
		return NewDenyOutput("체크리스트 없음: Task 프롬프트에 체크리스트 경로가 없습니다. 먼저 체크리스트를 생성하세요.")
	}

	// Check if prompt contains a valid jobs path (.do/jobs/YY/MM/DD/title/checklists/XX_agent.md)
	if !strings.Contains(prompt, ".do/jobs/") || !strings.Contains(prompt, "/checklists/") {
		return NewDenyOutput("체크리스트 없음: 작업 지시에 서브 체크리스트 경로(.do/jobs/YY/MM/DD/title/checklists/XX_agent.md)가 없습니다. 먼저 체크리스트를 생성하세요.")
	}

	// Extract and verify the checklist path exists
	checklistPath := extractChecklistPath(prompt)
	if checklistPath == "" {
		return NewDenyOutput("체크리스트 없음: 유효한 체크리스트 경로를 찾을 수 없습니다.")
	}

	if _, err := os.Stat(checklistPath); os.IsNotExist(err) {
		return NewDenyOutput("체크리스트 없음: " + checklistPath + " 파일이 없습니다. 먼저 체크리스트를 생성하세요.")
	}

	return nil // Allow - valid checklist exists
}

// extractPrompt extracts the prompt from Task tool input JSON.
func extractPrompt(toolInput json.RawMessage) string {
	if len(toolInput) == 0 {
		return ""
	}
	var data map[string]interface{}
	if json.Unmarshal(toolInput, &data) != nil {
		return ""
	}
	if prompt, ok := data["prompt"].(string); ok {
		return prompt
	}
	return ""
}

// extractChecklistPath extracts .do/jobs/.../checklists/...md path from prompt.
func extractChecklistPath(prompt string) string {
	// Find .do/jobs path
	idx := strings.Index(prompt, ".do/jobs/")
	if idx == -1 {
		return ""
	}

	// Extract from .do/jobs to end of line or space
	rest := prompt[idx:]
	end := strings.IndexAny(rest, " \n\t")
	if end == -1 {
		end = len(rest)
	}

	path := rest[:end]
	// Verify it's a checklist path
	if strings.Contains(path, "/checklists/") && strings.HasSuffix(path, ".md") {
		return path
	}
	return ""
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
