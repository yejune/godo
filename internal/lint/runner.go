package lint

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Diagnostic represents a single lint finding.
type Diagnostic struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"` // "error" or "warning"
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Source   string `json:"source"` // linter name
}

// RunLinter dispatches to the appropriate linter for the given language.
func RunLinter(lang Language, files []string, projectDir string) []Diagnostic {
	switch lang {
	case LangGo:
		return RunGoVet(files, projectDir)
	case LangPython:
		return RunRuff(files, projectDir)
	case LangTypeScript:
		return RunTsc(files, projectDir)
	case LangJavaScript:
		return RunESLint(files, projectDir)
	case LangRust:
		return RunCargoClippy(projectDir)
	default:
		return nil
	}
}

// RunGoVet runs `go vet ./...` and parses the output.
func RunGoVet(files []string, projectDir string) []Diagnostic {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = projectDir
	out, _ := cmd.CombinedOutput()
	return ParseGoVetOutput(string(out))
}

// goVetPattern matches: file.go:line:column: message
var goVetPattern = regexp.MustCompile(`^\.?/?([^:]+\.go):(\d+):(\d+):\s*(.+)$`)

// ParseGoVetOutput parses go vet stderr output.
func ParseGoVetOutput(output string) []Diagnostic {
	var diags []Diagnostic
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		matches := goVetPattern.FindStringSubmatch(line)
		if len(matches) != 5 {
			continue
		}
		lineNum, err := strconv.Atoi(matches[2])
		if err != nil {
			continue
		}
		colNum, err := strconv.Atoi(matches[3])
		if err != nil {
			continue
		}
		diags = append(diags, Diagnostic{
			File:     matches[1],
			Line:     lineNum,
			Column:   colNum,
			Severity: "warning",
			Message:  matches[4],
			Source:   "go vet",
		})
	}
	return diags
}

// RunRuff runs `ruff check --output-format=json` and parses the JSON output.
func RunRuff(files []string, projectDir string) []Diagnostic {
	args := []string{"check", "--output-format=json"}
	args = append(args, files...)
	cmd := exec.Command("ruff", args...)
	cmd.Dir = projectDir
	out, _ := cmd.Output()
	return ParseRuffJSON(out)
}

// ParseRuffJSON parses ruff JSON output.
func ParseRuffJSON(data []byte) []Diagnostic {
	if len(data) == 0 {
		return nil
	}

	var ruffDiags []struct {
		Code     string `json:"code"`
		Message  string `json:"message"`
		Location struct {
			Row    int `json:"row"`
			Column int `json:"column"`
		} `json:"location"`
		Filename string `json:"filename"`
	}

	if err := json.Unmarshal(data, &ruffDiags); err != nil {
		return nil
	}

	var diags []Diagnostic
	for _, d := range ruffDiags {
		severity := "warning"
		if strings.HasPrefix(d.Code, "E") || strings.HasPrefix(d.Code, "F") {
			severity = "error"
		}
		diags = append(diags, Diagnostic{
			File:     d.Filename,
			Line:     d.Location.Row,
			Column:   d.Location.Column,
			Severity: severity,
			Message:  d.Message,
			Rule:     d.Code,
			Source:   "ruff",
		})
	}
	return diags
}

// RunTsc runs `tsc --noEmit --pretty false` and parses the output.
func RunTsc(files []string, projectDir string) []Diagnostic {
	args := []string{"--noEmit", "--pretty", "false"}
	cmd := exec.Command("tsc", args...)
	cmd.Dir = projectDir
	out, _ := cmd.CombinedOutput()
	return ParseTscOutput(string(out))
}

// tscPattern matches: file.ts(line,column): error TS1234: message
var tscPattern = regexp.MustCompile(`^([^(]+)\((\d+),(\d+)\):\s*(error|warning)\s+(TS\d+):\s*(.+)$`)

// ParseTscOutput parses tsc text output.
func ParseTscOutput(output string) []Diagnostic {
	var diags []Diagnostic
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		matches := tscPattern.FindStringSubmatch(line)
		if len(matches) != 7 {
			continue
		}
		lineNum, err := strconv.Atoi(matches[2])
		if err != nil {
			continue
		}
		colNum, err := strconv.Atoi(matches[3])
		if err != nil {
			continue
		}
		diags = append(diags, Diagnostic{
			File:     strings.TrimSpace(matches[1]),
			Line:     lineNum,
			Column:   colNum,
			Severity: matches[4],
			Message:  matches[6],
			Rule:     matches[5],
			Source:   "tsc",
		})
	}
	return diags
}

// RunESLint runs `eslint --format json` and parses the JSON output.
func RunESLint(files []string, projectDir string) []Diagnostic {
	args := []string{"--format", "json"}
	args = append(args, files...)
	cmd := exec.Command("eslint", args...)
	cmd.Dir = projectDir
	out, _ := cmd.Output()
	return ParseESLintJSON(out)
}

// ParseESLintJSON parses eslint JSON output.
func ParseESLintJSON(data []byte) []Diagnostic {
	if len(data) == 0 {
		return nil
	}

	var eslintResults []struct {
		FilePath string `json:"filePath"`
		Messages []struct {
			RuleID   string `json:"ruleId"`
			Severity int    `json:"severity"`
			Message  string `json:"message"`
			Line     int    `json:"line"`
			Column   int    `json:"column"`
		} `json:"messages"`
	}

	if err := json.Unmarshal(data, &eslintResults); err != nil {
		return nil
	}

	var diags []Diagnostic
	for _, file := range eslintResults {
		for _, msg := range file.Messages {
			severity := "warning"
			if msg.Severity == 2 {
				severity = "error"
			}
			diags = append(diags, Diagnostic{
				File:     file.FilePath,
				Line:     msg.Line,
				Column:   msg.Column,
				Severity: severity,
				Message:  msg.Message,
				Rule:     msg.RuleID,
				Source:   "eslint",
			})
		}
	}
	return diags
}

// RunCargoClippy runs `cargo clippy --message-format=json` and parses the output.
func RunCargoClippy(projectDir string) []Diagnostic {
	cmd := exec.Command("cargo", "clippy", "--message-format=json")
	cmd.Dir = projectDir
	out, _ := cmd.Output()
	return ParseClippyJSON(out)
}

// ParseClippyJSON parses cargo clippy JSON output (one JSON per line).
func ParseClippyJSON(data []byte) []Diagnostic {
	if len(data) == 0 {
		return nil
	}

	var diags []Diagnostic
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var msg struct {
			Reason  string `json:"reason"`
			Message *struct {
				Code *struct {
					Code string `json:"code"`
				} `json:"code"`
				Level   string `json:"level"`
				Message string `json:"message"`
				Spans   []struct {
					FileName    string `json:"file_name"`
					LineStart   int    `json:"line_start"`
					ColumnStart int    `json:"column_start"`
				} `json:"spans"`
			} `json:"message"`
		}

		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		if msg.Reason != "compiler-message" || msg.Message == nil {
			continue
		}

		severity := "warning"
		if msg.Message.Level == "error" {
			severity = "error"
		}

		rule := ""
		if msg.Message.Code != nil {
			rule = msg.Message.Code.Code
		}

		for _, span := range msg.Message.Spans {
			diags = append(diags, Diagnostic{
				File:     span.FileName,
				Line:     span.LineStart,
				Column:   span.ColumnStart,
				Severity: severity,
				Message:  msg.Message.Message,
				Rule:     rule,
				Source:   "clippy",
			})
			break // Only use first span
		}
	}
	return diags
}
