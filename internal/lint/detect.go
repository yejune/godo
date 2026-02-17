package lint

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Language represents a detected programming language.
type Language string

const (
	LangGo         Language = "go"
	LangPython     Language = "python"
	LangTypeScript Language = "typescript"
	LangJavaScript Language = "javascript"
	LangRust       Language = "rust"
	LangUnknown    Language = "unknown"
)

// LinterInfo holds the linter command name and display name for a language.
type LinterInfo struct {
	Language    Language
	Command     string // binary name for exec.LookPath
	DisplayName string
}

// AllLinters returns the linter info for each supported language.
func AllLinters() []LinterInfo {
	return []LinterInfo{
		{LangGo, "go", "go vet"},
		{LangPython, "ruff", "ruff"},
		{LangTypeScript, "tsc", "tsc (TypeScript)"},
		{LangJavaScript, "eslint", "eslint"},
		{LangRust, "cargo", "cargo clippy"},
	}
}

// LinterForLanguage returns the linter info for a given language.
func LinterForLanguage(lang Language) (LinterInfo, bool) {
	for _, l := range AllLinters() {
		if l.Language == lang {
			return l, true
		}
	}
	return LinterInfo{}, false
}

// DetectLanguage maps a file extension to a language.
func DetectLanguage(filePath string) Language {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return LangGo
	case ".py", ".pyi":
		return LangPython
	case ".ts", ".tsx", ".mts", ".cts":
		return LangTypeScript
	case ".js", ".jsx", ".mjs", ".cjs":
		return LangJavaScript
	case ".rs":
		return LangRust
	default:
		return LangUnknown
	}
}

// IsCodeFile returns true if the file has a recognized code extension.
func IsCodeFile(filePath string) bool {
	return DetectLanguage(filePath) != LangUnknown
}

// CheckLinterInstalled returns true if the linter for the given language is installed.
func CheckLinterInstalled(lang Language) bool {
	info, ok := LinterForLanguage(lang)
	if !ok {
		return false
	}
	_, err := exec.LookPath(info.Command)
	return err == nil
}

// GetChangedFiles returns files changed in git (staged + unstaged).
// If all is true, returns all tracked files instead.
func GetChangedFiles(projectDir string, all bool) []string {
	var cmd *exec.Cmd
	if all {
		cmd = exec.Command("git", "ls-files")
	} else {
		cmd = exec.Command("git", "diff", "--name-only", "HEAD")
	}
	cmd.Dir = projectDir
	out, err := cmd.Output()
	if err != nil {
		if !all {
			cmd = exec.Command("git", "diff", "--name-only")
			cmd.Dir = projectDir
			out, err = cmd.Output()
			if err != nil {
				return nil
			}
		} else {
			return nil
		}
	}

	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && IsCodeFile(line) {
			files = append(files, line)
		}
	}
	return files
}

// GroupFilesByLanguage groups file paths by their detected language.
func GroupFilesByLanguage(files []string) map[Language][]string {
	groups := make(map[Language][]string)
	for _, f := range files {
		lang := DetectLanguage(f)
		if lang != LangUnknown {
			groups[lang] = append(groups[lang], f)
		}
	}
	return groups
}

// RunForHook runs lint on a specific file and returns diagnostics as a string.
// Returns empty string if no issues found or linter not installed.
func RunForHook(filePath string, projectDir string) string {
	lang := DetectLanguage(filePath)
	if lang == LangUnknown {
		return ""
	}

	if !CheckLinterInstalled(lang) {
		info, _ := LinterForLanguage(lang)
		return "Lint skipped: " + info.DisplayName + " not installed."
	}

	diags := RunLinter(lang, []string{filePath}, projectDir)
	if len(diags) == 0 {
		return ""
	}

	return FormatDiagnostics(diags)
}
