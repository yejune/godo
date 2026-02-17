package hook

// HandleCompact handles the PreCompact hook event.
// Keep compact hook non-intrusive to match original godo behavior.
func HandleCompact(input *Input) *Output {
	return &Output{Continue: true}
}
