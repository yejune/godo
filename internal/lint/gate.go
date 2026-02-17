package lint

import (
	"fmt"
	"sort"
	"strings"
)

const (
	ExitSuccess = 0
	ExitError   = 2
)

// EvaluateResults checks diagnostics and prints formatted output.
// Returns exit code: 0 for success/warnings-only, 2 for errors.
func EvaluateResults(diags []Diagnostic) int {
	if len(diags) == 0 {
		fmt.Println("Lint: all clean")
		return ExitSuccess
	}

	errors := 0
	warnings := 0
	for _, d := range diags {
		if d.Severity == "error" {
			errors++
		} else {
			warnings++
		}
	}

	fmt.Print(FormatDiagnostics(diags))

	fileCount := CountUniqueFiles(diags)
	fmt.Printf("\n%d errors, %d warnings in %d files\n", errors, warnings, fileCount)

	if errors > 0 {
		return ExitError
	}
	return ExitSuccess
}

// FormatDiagnostics formats diagnostics grouped by file.
func FormatDiagnostics(diags []Diagnostic) string {
	if len(diags) == 0 {
		return ""
	}

	byFile := make(map[string][]Diagnostic)
	var fileOrder []string
	for _, d := range diags {
		if _, seen := byFile[d.File]; !seen {
			fileOrder = append(fileOrder, d.File)
		}
		byFile[d.File] = append(byFile[d.File], d)
	}
	sort.Strings(fileOrder)

	var sb strings.Builder
	for _, file := range fileOrder {
		fileDiags := byFile[file]
		sort.Slice(fileDiags, func(i, j int) bool {
			if fileDiags[i].Line != fileDiags[j].Line {
				return fileDiags[i].Line < fileDiags[j].Line
			}
			return fileDiags[i].Column < fileDiags[j].Column
		})

		sb.WriteString(file)
		sb.WriteString("\n")
		for _, d := range fileDiags {
			prefix := "warning"
			if d.Severity == "error" {
				prefix = "error"
			}
			ruleStr := ""
			if d.Rule != "" {
				ruleStr = " [" + d.Rule + "]"
			}
			sb.WriteString(fmt.Sprintf("  %d:%d  %s  %s%s  (%s)\n",
				d.Line, d.Column, prefix, d.Message, ruleStr, d.Source))
		}
	}
	return sb.String()
}

// CountUniqueFiles counts unique file paths in diagnostics.
func CountUniqueFiles(diags []Diagnostic) int {
	seen := make(map[string]bool)
	for _, d := range diags {
		seen[d.File] = true
	}
	return len(seen)
}
