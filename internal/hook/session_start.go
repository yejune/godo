package hook

import (
	"fmt"

	"github.com/yejune/godo/internal/mode"
)

// HandleSessionStart handles the SessionStart hook event.
// Keep startup message minimal and compatible with the original godo behavior.
func HandleSessionStart(input *Input) *Output {
	currentMode := mode.ReadState()
	message := fmt.Sprintf("current_mode: %s", currentMode)
	return NewSessionOutput(true, message)
}
