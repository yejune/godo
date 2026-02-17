package hook

// HandleSessionEnd handles the SessionEnd hook event.
// Minimal implementation - allows session to end normally.
func HandleSessionEnd(input *Input) *Output {
	return &Output{Continue: true}
}
