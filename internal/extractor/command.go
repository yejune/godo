package extractor

import (
	"os"
	"path/filepath"
	"strings"
)

// personaCommandPrefixes lists filename prefixes that indicate a persona command.
var personaCommandPrefixes = []string{
	"moai",
	"do:",
	"do-",
}

// isPersonaCommand returns true if the filename matches a persona command prefix.
func isPersonaCommand(filename string) bool {
	lower := strings.ToLower(filename)
	for _, prefix := range personaCommandPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

// ExtractCommands classifies and copies command files from sourceDir.
// Commands are classified as core or persona based on filename patterns.
//
// Persona commands: files starting with persona prefixes (moai, do:, do-, godo)
// Core commands: everything else
func ExtractCommands(sourceDir, coreCmdDir, personaCmdDir string) (coreFiles, personaFiles []string, err error) {
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

		filename := filepath.Base(relPath)

		if isPersonaCommand(filename) {
			dst := filepath.Join(personaCmdDir, relPath)
			if cpErr := copyFile(path, dst); cpErr != nil {
				return cpErr
			}
			personaFiles = append(personaFiles, relPath)
		} else {
			dst := filepath.Join(coreCmdDir, relPath)
			if cpErr := copyFile(path, dst); cpErr != nil {
				return cpErr
			}
			coreFiles = append(coreFiles, relPath)
		}

		return nil
	})
	return
}
