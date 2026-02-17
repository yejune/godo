package validator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yejune/godo/internal/model"
)

// --- ValidatePhase ---

func TestValidatePhase_Complete(t *testing.T) {
	dir := t.TempDir()
	writeStateJSON(t, dir, map[string]string{"analysis": "complete"})

	v := &DependencyValidator{JobDir: dir}
	if err := v.ValidatePhase("analysis"); err != nil {
		t.Fatalf("expected no error for complete phase, got: %v", err)
	}
}

func TestValidatePhase_Incomplete(t *testing.T) {
	dir := t.TempDir()
	writeStateJSON(t, dir, map[string]string{"analysis": "in_progress"})

	v := &DependencyValidator{JobDir: dir}
	err := v.ValidatePhase("analysis")
	if err == nil {
		t.Fatal("expected error for incomplete phase")
	}
	if !strings.Contains(err.Error(), "not complete") {
		t.Errorf("expected 'not complete' in error, got: %v", err)
	}
}

func TestValidatePhase_NotFound(t *testing.T) {
	dir := t.TempDir()
	writeStateJSON(t, dir, map[string]string{"analysis": "complete"})

	v := &DependencyValidator{JobDir: dir}
	err := v.ValidatePhase("architecture")
	if err == nil {
		t.Fatal("expected error for missing phase")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestValidatePhase_MissingStateJSON(t *testing.T) {
	dir := t.TempDir()

	v := &DependencyValidator{JobDir: dir}
	err := v.ValidatePhase("analysis")
	if err == nil {
		t.Fatal("expected error when state.json is missing")
	}
	if !strings.Contains(err.Error(), "cannot read state.json") {
		t.Errorf("expected 'cannot read state.json' in error, got: %v", err)
	}
}

func TestValidatePhase_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "state.json"), []byte("{invalid"), 0644); err != nil {
		t.Fatal(err)
	}

	v := &DependencyValidator{JobDir: dir}
	err := v.ValidatePhase("analysis")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "cannot parse state.json") {
		t.Errorf("expected 'cannot parse state.json' in error, got: %v", err)
	}
}

// --- ValidateArtifact ---

func TestValidateArtifact_Exists(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "plan.md"), "# Plan")

	v := &DependencyValidator{JobDir: dir}
	dep := model.ArtifactDep{Path: "plan.md", Required: true}
	if err := v.ValidateArtifact(dep); err != nil {
		t.Fatalf("expected no error for existing artifact, got: %v", err)
	}
}

func TestValidateArtifact_Missing(t *testing.T) {
	dir := t.TempDir()

	v := &DependencyValidator{JobDir: dir}
	dep := model.ArtifactDep{Path: "plan.md", Required: true}
	err := v.ValidateArtifact(dep)
	if err == nil {
		t.Fatal("expected error for missing artifact")
	}
	if !strings.Contains(err.Error(), "file not found") {
		t.Errorf("expected 'file not found' in error, got: %v", err)
	}
}

func TestValidateArtifact_NestedPath(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "checklists")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(subDir, "01_backend.md"), "# Checklist")

	v := &DependencyValidator{JobDir: dir}
	dep := model.ArtifactDep{Path: "checklists/01_backend.md", Required: true}
	if err := v.ValidateArtifact(dep); err != nil {
		t.Fatalf("expected no error for nested artifact, got: %v", err)
	}
}

// --- ValidateAgent ---

func TestValidateAgent_AllComplete(t *testing.T) {
	dir := t.TempDir()
	checklistDir := filepath.Join(dir, "checklists")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# expert-backend\n- [o] #1 Create API\n- [o] #2 Add tests\n"
	writeFile(t, filepath.Join(checklistDir, "01_expert-backend.md"), content)

	v := &DependencyValidator{JobDir: dir}
	dep := model.AgentDep{Name: "expert-backend"}
	if err := v.ValidateAgent(dep); err != nil {
		t.Fatalf("expected no error when all items complete, got: %v", err)
	}
}

func TestValidateAgent_PartialComplete(t *testing.T) {
	dir := t.TempDir()
	checklistDir := filepath.Join(dir, "checklists")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# expert-backend\n- [o] #1 Create API\n- [~] #2 Add tests\n"
	writeFile(t, filepath.Join(checklistDir, "01_expert-backend.md"), content)

	v := &DependencyValidator{JobDir: dir}
	dep := model.AgentDep{Name: "expert-backend"}
	err := v.ValidateAgent(dep)
	if err == nil {
		t.Fatal("expected error when items are incomplete")
	}
	if !strings.Contains(err.Error(), "not complete") {
		t.Errorf("expected 'not complete' in error, got: %v", err)
	}
}

func TestValidateAgent_SpecificItems_AllComplete(t *testing.T) {
	dir := t.TempDir()
	checklistDir := filepath.Join(dir, "checklists")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# expert-backend\n- [o] #1 Create API\n- [ ] #2 Add tests\n- [o] #3 Deploy\n"
	writeFile(t, filepath.Join(checklistDir, "01_expert-backend.md"), content)

	v := &DependencyValidator{JobDir: dir}
	dep := model.AgentDep{Name: "expert-backend", Items: []string{"#1", "#3"}}
	if err := v.ValidateAgent(dep); err != nil {
		t.Fatalf("expected no error when specific items are complete, got: %v", err)
	}
}

func TestValidateAgent_SpecificItems_Incomplete(t *testing.T) {
	dir := t.TempDir()
	checklistDir := filepath.Join(dir, "checklists")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# expert-backend\n- [o] #1 Create API\n- [ ] #2 Add tests\n"
	writeFile(t, filepath.Join(checklistDir, "01_expert-backend.md"), content)

	v := &DependencyValidator{JobDir: dir}
	dep := model.AgentDep{Name: "expert-backend", Items: []string{"#1", "#2"}}
	err := v.ValidateAgent(dep)
	if err == nil {
		t.Fatal("expected error when specific item is incomplete")
	}
	if !strings.Contains(err.Error(), "#2 not complete") {
		t.Errorf("expected '#2 not complete' in error, got: %v", err)
	}
}

func TestValidateAgent_SpecificItems_NotFound(t *testing.T) {
	dir := t.TempDir()
	checklistDir := filepath.Join(dir, "checklists")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# expert-backend\n- [o] #1 Create API\n"
	writeFile(t, filepath.Join(checklistDir, "01_expert-backend.md"), content)

	v := &DependencyValidator{JobDir: dir}
	dep := model.AgentDep{Name: "expert-backend", Items: []string{"#1", "#99"}}
	err := v.ValidateAgent(dep)
	if err == nil {
		t.Fatal("expected error when item not found in checklist")
	}
	if !strings.Contains(err.Error(), "#99 not found") {
		t.Errorf("expected '#99 not found' in error, got: %v", err)
	}
}

func TestValidateAgent_NoChecklistFile(t *testing.T) {
	dir := t.TempDir()
	checklistDir := filepath.Join(dir, "checklists")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatal(err)
	}

	v := &DependencyValidator{JobDir: dir}
	dep := model.AgentDep{Name: "expert-backend"}
	err := v.ValidateAgent(dep)
	if err == nil {
		t.Fatal("expected error when checklist file is missing")
	}
	if !strings.Contains(err.Error(), "no checklist file found") {
		t.Errorf("expected 'no checklist file found' in error, got: %v", err)
	}
}

func TestValidateAgent_AllStatuses(t *testing.T) {
	// Test that all non-[o] statuses are rejected
	statuses := []struct {
		symbol string
		label  string
	}{
		{" ", "pending"},
		{"~", "in_progress"},
		{"*", "testing"},
		{"!", "blocked"},
		{"x", "failed"},
	}

	for _, s := range statuses {
		t.Run(s.label, func(t *testing.T) {
			dir := t.TempDir()
			checklistDir := filepath.Join(dir, "checklists")
			if err := os.MkdirAll(checklistDir, 0755); err != nil {
				t.Fatal(err)
			}
			content := "- [" + s.symbol + "] #1 Some task\n"
			writeFile(t, filepath.Join(checklistDir, "01_test-agent.md"), content)

			v := &DependencyValidator{JobDir: dir}
			dep := model.AgentDep{Name: "test-agent"}
			err := v.ValidateAgent(dep)
			if err == nil {
				t.Fatalf("expected error for status [%s], got nil", s.symbol)
			}
		})
	}
}

// --- ValidateEnv ---

func TestValidateEnv_Set(t *testing.T) {
	t.Setenv("TEST_VALIDATE_ENV_VAR", "some_value")

	v := &DependencyValidator{JobDir: t.TempDir()}
	if err := v.ValidateEnv("TEST_VALIDATE_ENV_VAR"); err != nil {
		t.Fatalf("expected no error for set env var, got: %v", err)
	}
}

func TestValidateEnv_Unset(t *testing.T) {
	v := &DependencyValidator{JobDir: t.TempDir()}
	err := v.ValidateEnv("TEST_VALIDATE_ENV_DEFINITELY_UNSET_12345")
	if err == nil {
		t.Fatal("expected error for unset env var")
	}
	if !strings.Contains(err.Error(), "not set or empty") {
		t.Errorf("expected 'not set or empty' in error, got: %v", err)
	}
}

func TestValidateEnv_EmptyValue(t *testing.T) {
	t.Setenv("TEST_VALIDATE_ENV_EMPTY", "")

	v := &DependencyValidator{JobDir: t.TempDir()}
	err := v.ValidateEnv("TEST_VALIDATE_ENV_EMPTY")
	if err == nil {
		t.Fatal("expected error for empty env var")
	}
	if !strings.Contains(err.Error(), "not set or empty") {
		t.Errorf("expected 'not set or empty' in error, got: %v", err)
	}
}

// --- ValidateService ---
// ValidateService shells out to `docker compose ps`. We test the parsing logic
// indirectly and verify graceful error handling when docker is unavailable.

func TestValidateService_DockerUnavailable(t *testing.T) {
	v := &DependencyValidator{JobDir: t.TempDir()}
	dep := model.ServiceDep{Name: "postgres", Healthcheck: false}

	err := v.ValidateService(dep)
	if err == nil {
		// Docker is available but service not found -- also valid
		return
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "docker compose ps failed") &&
		!strings.Contains(errStr, "not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestDockerPSEntry_Unmarshal tests the JSON parsing of docker compose ps output lines.
func TestDockerPSEntry_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonLine    string
		wantService string
		wantState   string
		wantHealth  string
	}{
		{
			name:        "running_healthy",
			jsonLine:    `{"Name":"project-postgres-1","Service":"postgres","State":"running","Health":"healthy"}`,
			wantService: "postgres",
			wantState:   "running",
			wantHealth:  "healthy",
		},
		{
			name:        "running_no_health",
			jsonLine:    `{"Name":"project-redis-1","Service":"redis","State":"running","Health":""}`,
			wantService: "redis",
			wantState:   "running",
			wantHealth:  "",
		},
		{
			name:        "exited",
			jsonLine:    `{"Name":"project-worker-1","Service":"worker","State":"exited","Health":""}`,
			wantService: "worker",
			wantState:   "exited",
			wantHealth:  "",
		},
		{
			name:        "unhealthy",
			jsonLine:    `{"Name":"project-db-1","Service":"db","State":"running","Health":"unhealthy"}`,
			wantService: "db",
			wantState:   "running",
			wantHealth:  "unhealthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entry dockerPSEntry
			if err := json.Unmarshal([]byte(tt.jsonLine), &entry); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if entry.Service != tt.wantService {
				t.Errorf("Service = %q, want %q", entry.Service, tt.wantService)
			}
			if entry.State != tt.wantState {
				t.Errorf("State = %q, want %q", entry.State, tt.wantState)
			}
			if entry.Health != tt.wantHealth {
				t.Errorf("Health = %q, want %q", entry.Health, tt.wantHealth)
			}
		})
	}
}

// --- ValidateChecklistItem ---

func TestValidateChecklistItem_Complete(t *testing.T) {
	dir := t.TempDir()
	content := "## Tasks\n- [o] #1 Create API endpoint\n- [ ] #2 Add tests\n"
	writeFile(t, filepath.Join(dir, "checklist.md"), content)

	v := &DependencyValidator{JobDir: dir}
	if err := v.ValidateChecklistItem("#1"); err != nil {
		t.Fatalf("expected no error for complete item, got: %v", err)
	}
}

func TestValidateChecklistItem_Incomplete(t *testing.T) {
	dir := t.TempDir()
	content := "## Tasks\n- [o] #1 Create API endpoint\n- [ ] #2 Add tests\n"
	writeFile(t, filepath.Join(dir, "checklist.md"), content)

	v := &DependencyValidator{JobDir: dir}
	err := v.ValidateChecklistItem("#2")
	if err == nil {
		t.Fatal("expected error for incomplete item")
	}
	if !strings.Contains(err.Error(), "status is [ ]") {
		t.Errorf("expected 'status is [ ]' in error, got: %v", err)
	}
}

func TestValidateChecklistItem_NotFound(t *testing.T) {
	dir := t.TempDir()
	content := "## Tasks\n- [o] #1 Create API endpoint\n"
	writeFile(t, filepath.Join(dir, "checklist.md"), content)

	v := &DependencyValidator{JobDir: dir}
	err := v.ValidateChecklistItem("#99")
	if err == nil {
		t.Fatal("expected error for missing item")
	}
	if !strings.Contains(err.Error(), "not found in checklist.md") {
		t.Errorf("expected 'not found in checklist.md' in error, got: %v", err)
	}
}

func TestValidateChecklistItem_NoStatusMarker(t *testing.T) {
	dir := t.TempDir()
	content := "## Tasks\n- #1 Create API endpoint (no status)\n"
	writeFile(t, filepath.Join(dir, "checklist.md"), content)

	v := &DependencyValidator{JobDir: dir}
	err := v.ValidateChecklistItem("#1")
	if err == nil {
		t.Fatal("expected error for item without status marker")
	}
	if !strings.Contains(err.Error(), "no status marker") {
		t.Errorf("expected 'no status marker' in error, got: %v", err)
	}
}

func TestValidateChecklistItem_MissingChecklistFile(t *testing.T) {
	dir := t.TempDir()

	v := &DependencyValidator{JobDir: dir}
	err := v.ValidateChecklistItem("#1")
	if err == nil {
		t.Fatal("expected error when checklist.md is missing")
	}
	if !strings.Contains(err.Error(), "cannot open checklist.md") {
		t.Errorf("expected 'cannot open checklist.md' in error, got: %v", err)
	}
}

func TestValidateChecklistItem_WithoutHashPrefix(t *testing.T) {
	dir := t.TempDir()
	content := "## Tasks\n- [o] #1 Create API endpoint\n"
	writeFile(t, filepath.Join(dir, "checklist.md"), content)

	v := &DependencyValidator{JobDir: dir}
	// ValidateChecklistItem strips # prefix, so "1" should also work
	if err := v.ValidateChecklistItem("1"); err != nil {
		t.Fatalf("expected no error when passing id without # prefix, got: %v", err)
	}
}

func TestValidateChecklistItem_AllStatuses(t *testing.T) {
	statuses := []struct {
		symbol  string
		label   string
		wantErr bool
	}{
		{"o", "complete", false},
		{" ", "pending", true},
		{"~", "in_progress", true},
		{"*", "testing", true},
		{"!", "blocked", true},
		{"x", "failed", true},
	}

	for _, s := range statuses {
		t.Run(s.label, func(t *testing.T) {
			dir := t.TempDir()
			content := "- [" + s.symbol + "] #1 Some task\n"
			writeFile(t, filepath.Join(dir, "checklist.md"), content)

			v := &DependencyValidator{JobDir: dir}
			err := v.ValidateChecklistItem("#1")
			if s.wantErr && err == nil {
				t.Errorf("expected error for status [%s], got nil", s.symbol)
			}
			if !s.wantErr && err != nil {
				t.Errorf("expected no error for status [%s], got: %v", s.symbol, err)
			}
		})
	}
}

// --- ValidateAll ---

func TestValidateAll_NilDeps(t *testing.T) {
	v := &DependencyValidator{JobDir: t.TempDir()}
	result := v.ValidateAll(nil)
	if !result.OK {
		t.Fatal("expected OK for nil deps")
	}
	if len(result.Blocked) != 0 {
		t.Errorf("expected no blocked items, got: %v", result.Blocked)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings, got: %v", result.Warnings)
	}
}

func TestValidateAll_AllPass(t *testing.T) {
	dir := t.TempDir()

	// Set up phase
	writeStateJSON(t, dir, map[string]string{"analysis": "complete"})

	// Set up artifact
	writeFile(t, filepath.Join(dir, "plan.md"), "# Plan")

	// Set up agent checklist
	checklistDir := filepath.Join(dir, "checklists")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(checklistDir, "01_expert-backend.md"), "- [o] #1 Done\n")

	// Set up env
	t.Setenv("TEST_VALIDATE_ALL_VAR", "value")

	// Set up checklist item
	writeFile(t, filepath.Join(dir, "checklist.md"), "- [o] #1 Main task done\n")

	deps := &model.DependsOn{
		Phases:         []string{"analysis"},
		Artifacts:      []model.ArtifactDep{{Path: "plan.md", Required: true}},
		Agents:         []model.AgentDep{{Name: "expert-backend"}},
		Env:            []string{"TEST_VALIDATE_ALL_VAR"},
		ChecklistItems: []string{"#1"},
		// Services omitted -- requires docker
	}

	v := &DependencyValidator{JobDir: dir}
	result := v.ValidateAll(deps)
	if !result.OK {
		t.Fatalf("expected OK, got blocked: %v", result.Blocked)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings, got: %v", result.Warnings)
	}
}

func TestValidateAll_Blocked(t *testing.T) {
	dir := t.TempDir()

	// Phase incomplete
	writeStateJSON(t, dir, map[string]string{"analysis": "in_progress"})

	deps := &model.DependsOn{
		Phases:    []string{"analysis"},
		Artifacts: []model.ArtifactDep{{Path: "plan.md", Required: true}},
	}

	v := &DependencyValidator{JobDir: dir}
	result := v.ValidateAll(deps)
	if result.OK {
		t.Fatal("expected NOT OK when deps are blocked")
	}
	if len(result.Blocked) != 2 {
		t.Errorf("expected 2 blocked items, got %d: %v", len(result.Blocked), result.Blocked)
	}
}

func TestValidateAll_Warnings(t *testing.T) {
	dir := t.TempDir()

	// Artifact missing but NOT required -> warning, not blocked
	deps := &model.DependsOn{
		Artifacts: []model.ArtifactDep{{Path: "optional.md", Required: false}},
	}

	v := &DependencyValidator{JobDir: dir}
	result := v.ValidateAll(deps)
	if !result.OK {
		t.Fatal("expected OK when only optional deps fail")
	}
	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d: %v", len(result.Warnings), result.Warnings)
	}
	if len(result.Blocked) != 0 {
		t.Errorf("expected 0 blocked, got %d: %v", len(result.Blocked), result.Blocked)
	}
}

func TestValidateAll_MixedBlockedAndWarnings(t *testing.T) {
	dir := t.TempDir()

	// Phase incomplete (blocked)
	writeStateJSON(t, dir, map[string]string{"analysis": "in_progress"})

	deps := &model.DependsOn{
		Phases:    []string{"analysis"},
		Artifacts: []model.ArtifactDep{{Path: "optional.md", Required: false}},
		Env:       []string{"TEST_VALIDATE_ALL_MISSING_VAR_XYZ"},
	}

	v := &DependencyValidator{JobDir: dir}
	result := v.ValidateAll(deps)
	if result.OK {
		t.Fatal("expected NOT OK with blocked deps")
	}
	if len(result.Blocked) != 2 {
		t.Errorf("expected 2 blocked (phase + env), got %d: %v", len(result.Blocked), result.Blocked)
	}
	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning (optional artifact), got %d: %v", len(result.Warnings), result.Warnings)
	}
}

// --- Test Helpers ---

// writeStateJSON creates a state.json file in dir with given phase statuses.
func writeStateJSON(t *testing.T, dir string, phases map[string]string) {
	t.Helper()
	state := stateJSON{
		Phases: make(map[string]phaseEntry, len(phases)),
	}
	for name, status := range phases {
		state.Phases[name] = phaseEntry{Status: status}
	}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "state.json"), string(data))
}

// writeFile creates a file with the given content.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
