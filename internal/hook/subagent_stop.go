package hook

// HandleSubagentStop handles the SubagentStop hook event.
// Minimal implementation - allows subagent to stop normally.
func HandleSubagentStop(input *Input) *Output {
	return &Output{Continue: true}
}
