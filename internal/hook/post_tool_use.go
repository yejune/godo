package hook

// HandlePostToolUse handles the PostToolUse hook event.
// Keep it non-intrusive for baseline compatibility.
func HandlePostToolUse(input *Input) *Output {
	return &Output{}
}
