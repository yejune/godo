package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup for .claude/settings.local.json",
	Long: `Setup interactively configures user preferences and saves them
to .claude/settings.local.json. Existing keys not managed by setup are preserved.`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

// option represents a selectable choice for numbered prompts.
type option struct {
	label string
	value string
}

// settingsFile represents the structure of settings.local.json.
type settingsFile map[string]interface{}

// readSettingsLocal reads and parses .claude/settings.local.json.
// Returns an empty map if the file does not exist or cannot be parsed.
func readSettingsLocal() settingsFile {
	data, err := os.ReadFile(".claude/settings.local.json")
	if err != nil {
		return settingsFile{}
	}
	var settings settingsFile
	if err := json.Unmarshal(data, &settings); err != nil {
		return settingsFile{}
	}
	return settings
}

// getEnvMap extracts the "env" map from settings, creating it if absent.
func getEnvMap(settings settingsFile) map[string]interface{} {
	env, ok := settings["env"].(map[string]interface{})
	if !ok {
		env = make(map[string]interface{})
	}
	return env
}

// writeSettingsLocal writes settings to .claude/settings.local.json,
// creating the .claude directory if needed.
func writeSettingsLocal(settings settingsFile) error {
	if err := os.MkdirAll(".claude", 0755); err != nil {
		return fmt.Errorf("create .claude directory: %w", err)
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}
	if err := os.WriteFile(".claude/settings.local.json", append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("write settings: %w", err)
	}
	return nil
}

// currentUser returns the current OS username via whoami.
func currentUser() string {
	out, err := exec.Command("whoami").Output()
	if err != nil {
		return "user"
	}
	return strings.TrimSpace(string(out))
}

// promptText asks the user for free-text input with a default value.
func promptText(reader *bufio.Reader, label, defaultVal string) string {
	fmt.Printf("%s [%s]: ", label, defaultVal)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

// promptChoice asks the user to pick from numbered options.
// Returns the value of the selected option, or defaultVal on empty input.
func promptChoice(reader *bufio.Reader, label string, opts []option, defaultIdx int) string {
	parts := make([]string, len(opts))
	for i, o := range opts {
		parts[i] = fmt.Sprintf("%d:%s", i+1, o.label)
	}
	fmt.Printf("%s (%s) [%d]: ", label, strings.Join(parts, " "), defaultIdx+1)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return opts[defaultIdx].value
	}
	// Parse number
	var idx int
	if _, err := fmt.Sscanf(line, "%d", &idx); err == nil && idx >= 1 && idx <= len(opts) {
		return opts[idx-1].value
	}
	// Invalid input, use default
	return opts[defaultIdx].value
}

// promptYN asks a yes/no question and returns a string value.
// yesVal is the value when Y, noVal when N.
func promptYN(reader *bufio.Reader, label string, defaultYes bool, yesVal, noVal string) string {
	defaultStr := "N"
	if defaultYes {
		defaultStr = "Y"
	}
	fmt.Printf("%s (Y/N) [%s]: ", label, defaultStr)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToUpper(line))
	if line == "" {
		if defaultYes {
			return yesVal
		}
		return noVal
	}
	if line == "Y" || line == "YES" {
		return yesVal
	}
	return noVal
}

// findDefault returns the index of the option matching currentVal,
// or fallback if not found.
func findDefault(opts []option, currentVal string, fallback int) int {
	for i, o := range opts {
		if o.value == currentVal {
			return i
		}
	}
	return fallback
}

func runSetup(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)
	settings := readSettingsLocal()
	env := getEnvMap(settings)

	// Helper to get current env value as string
	cur := func(key string) string {
		if v, ok := env[key].(string); ok {
			return v
		}
		return ""
	}

	fmt.Fprintln(cmd.OutOrStdout(), "godo setup")
	fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 40))

	// 1. User name
	nameDefault := cur("DO_USER_NAME")
	if nameDefault == "" {
		nameDefault = currentUser()
	}
	userName := promptText(reader, "User name", nameDefault)

	// 2. Language
	langOpts := []option{
		{"ko", "ko"}, {"en", "en"}, {"ja", "ja"}, {"zh", "zh"},
	}
	langDefault := findDefault(langOpts, cur("DO_LANGUAGE"), 0)
	lang := promptChoice(reader, "Language", langOpts, langDefault)

	// 3. Commit language
	commitLangOpts := []option{
		{"en", "en"}, {"ko", "ko"},
	}
	commitLangDefault := findDefault(commitLangOpts, cur("DO_COMMIT_LANGUAGE"), 0)
	commitLang := promptChoice(reader, "Commit language", commitLangOpts, commitLangDefault)

	// 4. Persona
	personaOpts := []option{
		{"young-f", "young-f"}, {"young-m", "young-m"},
		{"senior-f", "senior-f"}, {"senior-m", "senior-m"},
	}
	personaDefault := findDefault(personaOpts, cur("DO_PERSONA"), 0)
	persona := promptChoice(reader, "Persona", personaOpts, personaDefault)

	// 5. Execution mode
	modeOpts := []option{
		{"do", "do"}, {"focus", "focus"}, {"team", "team"},
	}
	modeDefault := findDefault(modeOpts, cur("DO_MODE"), 0)
	execMode := promptChoice(reader, "Mode", modeOpts, modeDefault)

	// 6. Style
	styleOpts := []option{
		{"sprint", "sprint"}, {"pair", "pair"}, {"direct", "direct"},
	}
	styleDefault := findDefault(styleOpts, cur("DO_STYLE"), 1)
	style := promptChoice(reader, "Style", styleOpts, styleDefault)

	// 7. AI footer
	aiFooterCur := cur("DO_AI_FOOTER") == "true"
	aiFooter := promptYN(reader, "AI footer?", aiFooterCur, "true", "false")

	// 8. Jobs language
	jobsLangOpts := []option{
		{"en", "en"}, {"ko", "ko"},
	}
	jobsLangDefault := findDefault(jobsLangOpts, cur("DO_JOBS_LANGUAGE"), 0)
	jobsLang := promptChoice(reader, "Jobs language", jobsLangOpts, jobsLangDefault)

	// --- Claude Launch Flags ---
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), "Claude Launch Flags")
	fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 40))

	// 9. Bypass permissions
	bypassCur := cur("DO_CLAUDE_BYPASS") == "true"
	bypass := promptYN(reader, "Bypass permissions?", bypassCur, "true", "false")

	// 10. Chrome MCP
	chromeCur := cur("DO_CLAUDE_CHROME") == "true"
	chrome := promptYN(reader, "Enable Chrome MCP?", chromeCur, "true", "false")

	// 11. Continue session
	continueCur := cur("DO_CLAUDE_CONTINUE") == "true"
	continueSession := promptYN(reader, "Continue previous session?", continueCur, "true", "false")

	// 12. Auto sync
	autoSyncCur := cur("DO_CLAUDE_AUTO_SYNC")
	autoSyncDefault := autoSyncCur != "false" // default true
	autoSync := promptYN(reader, "Auto sync before launch?", autoSyncDefault, "true", "false")

	// Merge into env (preserves existing keys not managed here)
	env["DO_USER_NAME"] = userName
	env["DO_LANGUAGE"] = lang
	env["DO_COMMIT_LANGUAGE"] = commitLang
	env["DO_PERSONA"] = persona
	env["DO_MODE"] = execMode
	env["DO_STYLE"] = style
	env["DO_AI_FOOTER"] = aiFooter
	env["DO_JOBS_LANGUAGE"] = jobsLang
	env["DO_CLAUDE_BYPASS"] = bypass
	env["DO_CLAUDE_CHROME"] = chrome
	env["DO_CLAUDE_CONTINUE"] = continueSession
	env["DO_CLAUDE_AUTO_SYNC"] = autoSync
	// CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS is auto-injected via EnsureSettingsLocal (sync time)

	settings["env"] = env

	if err := writeSettingsLocal(settings); err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), "Saved to .claude/settings.local.json")
	return nil
}

// EnsureSettingsLocal creates .claude/settings.local.json with defaults if it
// does not exist. If the file exists but is missing CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS,
// it adds the key (migration). Call this during sync or other init paths.
func EnsureSettingsLocal() {
	settings := readSettingsLocal()

	settingsPath := ".claude/settings.local.json"
	if _, err := os.Stat(settingsPath); err != nil {
		// File does not exist: create with defaults
		env := map[string]interface{}{
			"DO_USER_NAME":                         currentUser(),
			"DO_LANGUAGE":                          "ko",
			"DO_COMMIT_LANGUAGE":                   "en",
			"DO_PERSONA":                           "young-f",
			"DO_MODE":                              "do",
			"DO_STYLE":                             "pair",
			"DO_AI_FOOTER":                         "false",
			"DO_JOBS_LANGUAGE":                     "en",
			"DO_CLAUDE_BYPASS":                     "false",
			"DO_CLAUDE_CHROME":                     "false",
			"DO_CLAUDE_CONTINUE":                   "false",
			"DO_CLAUDE_AUTO_SYNC":                  "true",
			"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1",
		}
		settings["env"] = env
		if err := writeSettingsLocal(settings); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to create settings.local.json: %v\n", err)
		}
		return
	}

	// File exists: migration - ensure CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS is present
	env := getEnvMap(settings)
	if _, ok := env["CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS"]; !ok {
		env["CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS"] = "1"
		settings["env"] = env
		if err := writeSettingsLocal(settings); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to migrate settings.local.json: %v\n", err)
		}
	}
}
