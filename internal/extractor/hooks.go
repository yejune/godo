package extractor

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// personaBinaryRe matches persona binary invocations at the start of a line
// (possibly preceded by whitespace). The trailing space or end-of-line ensures
// we match actual commands, not substrings like "godocker".
var personaBinaryRe = regexp.MustCompile(`(?m)^\s*(godo|moai)\s`)

// hasPersonaBinaryReference returns true if content contains a persona binary
// invocation as a command (at line start), not as a substring in prose.
func hasPersonaBinaryReference(content string) bool {
	return personaBinaryRe.MatchString(content)
}

// ExtractHookScripts classifies and copies hook scripts from sourceDir.
// Scripts are classified based on content (binary references).
//
// Persona scripts: contain references to persona binaries (godo, moai)
// Core scripts: no persona binary references
func ExtractHookScripts(sourceDir, coreHookDir, personaHookDir string) (coreFiles, personaFiles []string, err error) {
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		relPath, relErr := filepath.Rel(sourceDir, path)
		if relErr != nil {
			return relErr
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		if hasPersonaBinaryReference(string(content)) {
			dst := filepath.Join(personaHookDir, relPath)
			if cpErr := copyFile(path, dst); cpErr != nil {
				return cpErr
			}
			personaFiles = append(personaFiles, relPath)
		} else {
			dst := filepath.Join(coreHookDir, relPath)
			if cpErr := copyFile(path, dst); cpErr != nil {
				return cpErr
			}
			coreFiles = append(coreFiles, relPath)
		}

		return nil
	})
	return
}

// copyFile copies a file from src to dst, creating parent directories as needed.
func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
