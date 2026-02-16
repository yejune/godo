package validator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/do-focus/convert/internal/model"
)

// ValidationResult holds the outcome of validating all dependencies.
// OK is true when no required (blocking) dependencies failed.
type ValidationResult struct {
	OK       bool
	Blocked  []string // required deps that failed
	Warnings []string // soft deps that failed (required=false)
}

// DependencyValidator checks whether agent dependencies are satisfied
// at runtime. JobDir points to the .do/jobs/{id}/ directory that contains
// state.json, checklist.md, and the checklists/ sub-directory.
type DependencyValidator struct {
	JobDir string
}

// stateJSON mirrors the structure of state.json used for phase tracking.
type stateJSON struct {
	Phases map[string]phaseEntry `json:"phases"`
}

type phaseEntry struct {
	Status string `json:"status"`
}

// dockerPSEntry represents one service line from `docker compose ps --format json`.
type dockerPSEntry struct {
	Name    string `json:"Name"`
	Service string `json:"Service"`
	State   string `json:"State"`
	Health  string `json:"Health"`
}

// statusRegex matches checklist status symbols: [ ], [~], [*], [!], [o], [x]
var statusRegex = regexp.MustCompile(`\[([ ~*!ox])\]`)

// itemIDRegex matches checklist item IDs like #1, #12, #123
var itemIDRegex = regexp.MustCompile(`#(\d+)`)

// ValidatePhase checks whether the given phase is marked "complete" in state.json.
func (v *DependencyValidator) ValidatePhase(phase string) error {
	path := filepath.Join(v.JobDir, "state.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("phase %q: cannot read state.json: %w", phase, err)
	}

	var state stateJSON
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("phase %q: cannot parse state.json: %w", phase, err)
	}

	entry, ok := state.Phases[phase]
	if !ok {
		return fmt.Errorf("phase %q: not found in state.json", phase)
	}
	if entry.Status != "complete" {
		return fmt.Errorf("phase %q: status is %q, not complete", phase, entry.Status)
	}
	return nil
}

// ValidateArtifact checks whether the artifact file exists under JobDir.
func (v *DependencyValidator) ValidateArtifact(dep model.ArtifactDep) error {
	path := filepath.Join(v.JobDir, dep.Path)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("artifact %q: file not found: %w", dep.Path, err)
	}
	return nil
}

// ValidateAgent checks whether an agent's checklist items are all completed ([o]).
// It looks for a file matching checklists/*_{agentName}.md under JobDir.
// If dep.Items is set, only those specific items (by #id) are checked.
func (v *DependencyValidator) ValidateAgent(dep model.AgentDep) error {
	checklistDir := filepath.Join(v.JobDir, "checklists")
	pattern := filepath.Join(checklistDir, fmt.Sprintf("*_%s.md", dep.Name))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("agent %q: glob error: %w", dep.Name, err)
	}
	if len(matches) == 0 {
		return fmt.Errorf("agent %q: no checklist file found matching *_%s.md", dep.Name, dep.Name)
	}

	filePath := matches[0]
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("agent %q: cannot open checklist: %w", dep.Name, err)
	}
	defer file.Close()

	// If specific items are requested, check only those
	if len(dep.Items) > 0 {
		return v.validateSpecificItems(file, dep.Name, dep.Items)
	}

	// Otherwise, check that ALL items have [o] status
	return v.validateAllItems(file, dep.Name)
}

// validateAllItems ensures every checklist item in the file has [o] status.
func (v *DependencyValidator) validateAllItems(file *os.File, agentName string) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := statusRegex.FindStringSubmatch(line)
		if match == nil {
			continue // not a checklist line
		}
		status := match[1]
		if status != "o" {
			return fmt.Errorf("agent %q: item not complete (status [%s]): %s",
				agentName, status, strings.TrimSpace(line))
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("agent %q: error reading checklist: %w", agentName, err)
	}
	return nil
}

// validateSpecificItems checks that the specified #id items all have [o] status.
func (v *DependencyValidator) validateSpecificItems(file *os.File, agentName string, items []string) error {
	// Build a set of required item IDs (strip # prefix)
	required := make(map[string]bool, len(items))
	for _, item := range items {
		id := strings.TrimPrefix(item, "#")
		required[id] = false // false = not yet found
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line has a status marker
		statusMatch := statusRegex.FindStringSubmatch(line)
		if statusMatch == nil {
			continue
		}

		// Check if this line has an item ID we care about
		idMatch := itemIDRegex.FindStringSubmatch(line)
		if idMatch == nil {
			continue
		}

		id := idMatch[1]
		if _, ok := required[id]; !ok {
			continue // not a required item
		}

		status := statusMatch[1]
		if status != "o" {
			return fmt.Errorf("agent %q: item #%s not complete (status [%s])",
				agentName, id, status)
		}
		required[id] = true // mark as found and complete
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("agent %q: error reading checklist: %w", agentName, err)
	}

	// Check that all required items were found
	for id, found := range required {
		if !found {
			return fmt.Errorf("agent %q: item #%s not found in checklist", agentName, id)
		}
	}
	return nil
}

// ValidateEnv checks that the given environment variable is set and non-empty.
func (v *DependencyValidator) ValidateEnv(varName string) error {
	val := os.Getenv(varName)
	if val == "" {
		return fmt.Errorf("env %q: not set or empty", varName)
	}
	return nil
}

// ValidateService checks that a Docker Compose service is running (and
// optionally healthy). It shells out to `docker compose ps --format json`.
func (v *DependencyValidator) ValidateService(dep model.ServiceDep) error {
	cmd := exec.Command("docker", "compose", "ps", "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("service %q: docker compose ps failed: %w", dep.Name, err)
	}

	// docker outputs one JSON object per line, NOT a JSON array
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry dockerPSEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue // skip malformed lines
		}

		if entry.Service != dep.Name && entry.Name != dep.Name {
			continue
		}

		// Found the service -- check state
		if entry.State != "running" {
			return fmt.Errorf("service %q: state is %q, not running", dep.Name, entry.State)
		}

		if dep.Healthcheck && entry.Health != "healthy" {
			return fmt.Errorf("service %q: health is %q, not healthy", dep.Name, entry.Health)
		}

		return nil // service found, running, and healthy (if required)
	}

	return fmt.Errorf("service %q: not found in docker compose ps output", dep.Name)
}

// ValidateChecklistItem checks that a specific item (by #id) in the main
// checklist.md file has [o] status.
func (v *DependencyValidator) ValidateChecklistItem(itemID string) error {
	path := filepath.Join(v.JobDir, "checklist.md")
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("checklist item %q: cannot open checklist.md: %w", itemID, err)
	}
	defer file.Close()

	id := strings.TrimPrefix(itemID, "#")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line contains the target item ID
		idMatch := itemIDRegex.FindStringSubmatch(line)
		if idMatch == nil || idMatch[1] != id {
			continue
		}

		// Found the item -- check status
		statusMatch := statusRegex.FindStringSubmatch(line)
		if statusMatch == nil {
			return fmt.Errorf("checklist item #%s: found but no status marker", id)
		}

		if statusMatch[1] != "o" {
			return fmt.Errorf("checklist item #%s: status is [%s], not [o]", id, statusMatch[1])
		}
		return nil
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("checklist item #%s: error reading checklist.md: %w", id, err)
	}

	return fmt.Errorf("checklist item #%s: not found in checklist.md", id)
}

// ValidateAll checks all dependencies in the given DependsOn struct and
// returns an aggregated ValidationResult. Required deps that fail go into
// Blocked; soft deps (ArtifactDep with Required=false) that fail go into
// Warnings. OK is true when len(Blocked) == 0.
func (v *DependencyValidator) ValidateAll(deps *model.DependsOn) ValidationResult {
	if deps == nil {
		return ValidationResult{OK: true}
	}

	var result ValidationResult

	// Phases -- all required
	for _, phase := range deps.Phases {
		if err := v.ValidatePhase(phase); err != nil {
			result.Blocked = append(result.Blocked, err.Error())
		}
	}

	// Artifacts -- required flag determines blocked vs warning
	for _, dep := range deps.Artifacts {
		if err := v.ValidateArtifact(dep); err != nil {
			if dep.Required {
				result.Blocked = append(result.Blocked, err.Error())
			} else {
				result.Warnings = append(result.Warnings, err.Error())
			}
		}
	}

	// Agents -- all required
	for _, dep := range deps.Agents {
		if err := v.ValidateAgent(dep); err != nil {
			result.Blocked = append(result.Blocked, err.Error())
		}
	}

	// Env -- all required
	for _, varName := range deps.Env {
		if err := v.ValidateEnv(varName); err != nil {
			result.Blocked = append(result.Blocked, err.Error())
		}
	}

	// Services -- all required
	for _, dep := range deps.Services {
		if err := v.ValidateService(dep); err != nil {
			result.Blocked = append(result.Blocked, err.Error())
		}
	}

	// Checklist items -- all required
	for _, itemID := range deps.ChecklistItems {
		if err := v.ValidateChecklistItem(itemID); err != nil {
			result.Blocked = append(result.Blocked, err.Error())
		}
	}

	result.OK = len(result.Blocked) == 0
	return result
}
