package validator

import (
	"strings"
	"testing"

	"github.com/yejune/godo/internal/model"
)

// --- DetectCycles ---

func TestDetectCycles_NoCycle(t *testing.T) {
	// Linear chain: A -> B -> C (A depends on B, B depends on C)
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "B"}},
	})
	b.AddAgent("B", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "C"}},
	})
	b.AddAgent("C", nil)

	g := b.Build()
	cycles, err := g.DetectCycles()
	if err != nil {
		t.Fatalf("expected no cycles, got error: %v", err)
	}
	if cycles != nil {
		t.Fatalf("expected nil cycles, got: %v", cycles)
	}
}

func TestDetectCycles_SimpleCycle(t *testing.T) {
	// A -> B -> A
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "B"}},
	})
	b.AddAgent("B", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "A"}},
	})

	g := b.Build()
	cycles, err := g.DetectCycles()
	if err == nil {
		t.Fatal("expected error for cycle, got nil")
	}
	if len(cycles) == 0 {
		t.Fatal("expected at least one cycle")
	}

	// Verify the cycle contains both A and B
	cycleStr := strings.Join(cycles[0], " -> ")
	if !strings.Contains(cycleStr, "A") || !strings.Contains(cycleStr, "B") {
		t.Errorf("expected cycle to contain A and B, got: %s", cycleStr)
	}

	// Verify the cycle is closed (first == last)
	cycle := cycles[0]
	if cycle[0] != cycle[len(cycle)-1] {
		t.Errorf("expected cycle to be closed (first == last), got: %v", cycle)
	}
}

func TestDetectCycles_ComplexCycle(t *testing.T) {
	// A -> B -> C -> A
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "B"}},
	})
	b.AddAgent("B", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "C"}},
	})
	b.AddAgent("C", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "A"}},
	})

	g := b.Build()
	cycles, err := g.DetectCycles()
	if err == nil {
		t.Fatal("expected error for cycle, got nil")
	}
	if len(cycles) == 0 {
		t.Fatal("expected at least one cycle")
	}

	// Verify the cycle contains A, B, and C
	cycleStr := strings.Join(cycles[0], " -> ")
	if !strings.Contains(cycleStr, "A") ||
		!strings.Contains(cycleStr, "B") ||
		!strings.Contains(cycleStr, "C") {
		t.Errorf("expected cycle to contain A, B, and C, got: %s", cycleStr)
	}

	// Verify the cycle is closed
	cycle := cycles[0]
	if cycle[0] != cycle[len(cycle)-1] {
		t.Errorf("expected cycle to be closed, got: %v", cycle)
	}

	// Verify the error message contains "dependency cycles detected"
	if !strings.Contains(err.Error(), "dependency cycles detected") {
		t.Errorf("expected 'dependency cycles detected' in error, got: %v", err)
	}
}

func TestDetectCycles_SelfCycle(t *testing.T) {
	// A -> A
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "A"}},
	})

	g := b.Build()
	cycles, err := g.DetectCycles()
	if err == nil {
		t.Fatal("expected error for self-cycle, got nil")
	}
	if len(cycles) == 0 {
		t.Fatal("expected at least one cycle")
	}
}

func TestDetectCycles_DiamondNoCycle(t *testing.T) {
	// Diamond shape: A -> B, A -> C, B -> D, C -> D (no cycle)
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "B"}, {Name: "C"}},
	})
	b.AddAgent("B", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "D"}},
	})
	b.AddAgent("C", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "D"}},
	})
	b.AddAgent("D", nil)

	g := b.Build()
	cycles, err := g.DetectCycles()
	if err != nil {
		t.Fatalf("expected no cycles in diamond graph, got: %v", err)
	}
	if cycles != nil {
		t.Fatalf("expected nil cycles, got: %v", cycles)
	}
}

// --- TopologicalSort ---

func TestTopologicalSort_LinearChain(t *testing.T) {
	// A -> B -> C: execution order should be C, B, A
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "B"}},
	})
	b.AddAgent("B", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "C"}},
	})
	b.AddAgent("C", nil)

	g := b.Build()
	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// C must come before B, B must come before A
	posC := indexOf(order, "C")
	posB := indexOf(order, "B")
	posA := indexOf(order, "A")

	if posC == -1 || posB == -1 || posA == -1 {
		t.Fatalf("expected all nodes in result, got: %v", order)
	}
	if posC >= posB {
		t.Errorf("expected C before B, got order: %v", order)
	}
	if posB >= posA {
		t.Errorf("expected B before A, got order: %v", order)
	}
}

func TestTopologicalSort_Parallel(t *testing.T) {
	// Independent nodes: A, B, C (no dependencies)
	// Should be sorted alphabetically
	b := NewGraphBuilder()
	b.AddAgent("C", nil)
	b.AddAgent("A", nil)
	b.AddAgent("B", nil)

	g := b.Build()
	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"A", "B", "C"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d nodes, got %d: %v", len(expected), len(order), order)
	}
	for i, name := range expected {
		if order[i] != name {
			t.Errorf("order[%d] = %q, want %q (full order: %v)", i, order[i], name, order)
		}
	}
}

func TestTopologicalSort_Diamond(t *testing.T) {
	// A -> B, A -> C, B -> D, C -> D
	// Valid order: D, B, C, A or D, C, B, A (B and C are interchangeable)
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "B"}, {Name: "C"}},
	})
	b.AddAgent("B", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "D"}},
	})
	b.AddAgent("C", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "D"}},
	})
	b.AddAgent("D", nil)

	g := b.Build()
	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order) != 4 {
		t.Fatalf("expected 4 nodes, got %d: %v", len(order), order)
	}

	// D must be first (no deps), A must be last (depends on everything)
	if order[0] != "D" {
		t.Errorf("expected D first, got: %v", order)
	}
	if order[3] != "A" {
		t.Errorf("expected A last, got: %v", order)
	}

	// B and C should be alphabetically ordered between D and A
	if order[1] != "B" || order[2] != "C" {
		t.Errorf("expected B then C in middle, got: %v", order)
	}
}

func TestTopologicalSort_WithCycle(t *testing.T) {
	// A -> B -> A: topological sort should fail
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "B"}},
	})
	b.AddAgent("B", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "A"}},
	})

	g := b.Build()
	_, err := g.TopologicalSort()
	if err == nil {
		t.Fatal("expected error for cyclic graph")
	}
	if !strings.Contains(err.Error(), "cycle exists") {
		t.Errorf("expected 'cycle exists' in error, got: %v", err)
	}
}

func TestTopologicalSort_SingleNode(t *testing.T) {
	b := NewGraphBuilder()
	b.AddAgent("A", nil)

	g := b.Build()
	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 1 || order[0] != "A" {
		t.Errorf("expected [A], got: %v", order)
	}
}

// --- Visualize ---

func TestVisualize_WithDeps(t *testing.T) {
	b := NewGraphBuilder()
	b.AddAgent("expert-frontend", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "expert-backend"}, {Name: "manager-spec"}},
	})
	b.AddAgent("expert-backend", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "manager-spec"}},
	})
	b.AddAgent("manager-spec", nil)

	g := b.Build()
	output := g.Visualize()

	// Verify all nodes appear in the output
	if !strings.Contains(output, "expert-frontend") {
		t.Error("expected 'expert-frontend' in output")
	}
	if !strings.Contains(output, "expert-backend") {
		t.Error("expected 'expert-backend' in output")
	}
	if !strings.Contains(output, "manager-spec") {
		t.Error("expected 'manager-spec' in output")
	}

	// Verify dependency arrows appear
	if !strings.Contains(output, "<- expert-backend (depends on)") {
		t.Error("expected '<- expert-backend (depends on)' in output")
	}
	if !strings.Contains(output, "<- manager-spec (depends on)") {
		t.Error("expected '<- manager-spec (depends on)' in output")
	}

	// Verify node without deps shows "(no dependencies)"
	if !strings.Contains(output, "(no dependencies)") {
		t.Error("expected '(no dependencies)' in output")
	}
}

func TestVisualize_EmptyGraph(t *testing.T) {
	g := &DependencyGraph{
		Edges: make(map[string][]string),
		Nodes: make(map[string]bool),
	}
	output := g.Visualize()
	if output != "" {
		t.Errorf("expected empty output for empty graph, got: %q", output)
	}
}

// --- GraphBuilder ---

func TestBuildGraph_NilDeps(t *testing.T) {
	b := NewGraphBuilder()
	b.AddAgent("A", nil)
	b.AddAgent("B", nil)

	g := b.Build()
	if len(g.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(g.Nodes))
	}
	if len(g.Edges) != 0 {
		t.Errorf("expected no edges for nil deps, got: %v", g.Edges)
	}
}

func TestBuildGraph_ForwardReference(t *testing.T) {
	// A depends on B, but B is added after A
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "B"}},
	})
	// B not yet added via AddAgent, but should be in nodes as forward reference

	g := b.Build()
	if !g.Nodes["B"] {
		t.Error("expected forward-referenced node B to be in graph")
	}
	if len(g.Edges["A"]) != 1 || g.Edges["A"][0] != "B" {
		t.Errorf("expected A -> B edge, got: %v", g.Edges["A"])
	}
}

func TestBuildGraph_NonAgentDepsIgnored(t *testing.T) {
	// Phase, artifact, env, service deps should NOT create graph edges
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Phases:         []string{"analysis"},
		Artifacts:      []model.ArtifactDep{{Path: "plan.md", Required: true}},
		Env:            []string{"MY_VAR"},
		Services:       []model.ServiceDep{{Name: "postgres"}},
		ChecklistItems: []string{"#1"},
	})

	g := b.Build()
	if len(g.Edges["A"]) != 0 {
		t.Errorf("expected no edges for non-agent deps, got: %v", g.Edges["A"])
	}
}

func TestBuildGraph_EdgesAreSorted(t *testing.T) {
	b := NewGraphBuilder()
	b.AddAgent("A", &model.DependsOn{
		Agents: []model.AgentDep{{Name: "C"}, {Name: "B"}, {Name: "D"}},
	})
	b.AddAgent("B", nil)
	b.AddAgent("C", nil)
	b.AddAgent("D", nil)

	g := b.Build()
	edges := g.Edges["A"]
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(edges))
	}
	if edges[0] != "B" || edges[1] != "C" || edges[2] != "D" {
		t.Errorf("expected edges sorted [B, C, D], got: %v", edges)
	}
}

// --- Helpers ---

// indexOf returns the position of name in slice, or -1 if not found.
func indexOf(slice []string, name string) int {
	for i, v := range slice {
		if v == name {
			return i
		}
	}
	return -1
}
