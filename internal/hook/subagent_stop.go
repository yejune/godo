package hook

// HandleSubagentStop handles the SubagentStop hook event.
// Baseline compatibility: continue without injecting extra warnings.
func HandleSubagentStop(input *Input) *Output {
	return &Output{Continue: true}
}
