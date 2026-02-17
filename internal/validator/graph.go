// Package validator provides dependency graph construction and cycle detection
// for agent dependency validation.
package validator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yejune/godo/internal/model"
)

// DependencyGraph represents a directed acyclic graph of agent dependencies.
// Edges map from an agent to the list of agents it depends on.
type DependencyGraph struct {
	// Adjacency list: agent -> list of agents it depends on
	Edges map[string][]string
	// All known agent names
	Nodes map[string]bool
}

// GraphBuilder accumulates agent definitions and builds a DependencyGraph.
type GraphBuilder struct {
	agents map[string]*model.DependsOn
	nodes  map[string]bool
}

// NewGraphBuilder creates a new GraphBuilder for accumulating agent definitions.
func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		agents: make(map[string]*model.DependsOn),
		nodes:  make(map[string]bool),
	}
}

// AddAgent registers an agent and its dependencies in the builder.
// Only agents-type dependencies form graph edges; phases, artifacts, env,
// and services are runtime checks and do not contribute to the graph structure.
// A nil deps is handled gracefully (agent with no dependencies).
func (b *GraphBuilder) AddAgent(name string, deps *model.DependsOn) {
	b.nodes[name] = true
	b.agents[name] = deps

	// Forward references: add dependency target nodes even if they haven't
	// been explicitly added via AddAgent yet.
	if deps != nil {
		for _, agentDep := range deps.Agents {
			b.nodes[agentDep.Name] = true
		}
	}
}

// Build finalizes the builder and returns a DependencyGraph.
func (b *GraphBuilder) Build() *DependencyGraph {
	g := &DependencyGraph{
		Edges: make(map[string][]string),
		Nodes: make(map[string]bool),
	}

	// Copy all nodes.
	for name := range b.nodes {
		g.Nodes[name] = true
	}

	// Build adjacency list from agent dependencies.
	for name, deps := range b.agents {
		if deps == nil {
			continue
		}
		for _, agentDep := range deps.Agents {
			g.Edges[name] = append(g.Edges[name], agentDep.Name)
		}
	}

	// Sort each edge list for deterministic output.
	for name := range g.Edges {
		sort.Strings(g.Edges[name])
	}

	return g
}

// color constants for DFS-based cycle detection (3-color marking).
const (
	white = 0 // unvisited
	gray  = 1 // in current DFS path (visiting)
	black = 2 // fully processed
)

// DetectCycles uses DFS with 3-color marking to find all cycles in the graph.
// Returns a slice of cycle paths (e.g., ["A", "B", "C", "A"]) and an error
// summarizing the cycles found. Returns nil slice and nil error if no cycles.
func (g *DependencyGraph) DetectCycles() ([][]string, error) {
	color := make(map[string]int)    // default white (0)
	parent := make(map[string]string) // tracks DFS path
	var cycles [][]string

	// Get sorted node list for deterministic traversal.
	nodes := g.sortedNodes()

	for _, node := range nodes {
		if color[node] == white {
			g.dfsDetectCycles(node, color, parent, &cycles)
		}
	}

	if len(cycles) == 0 {
		return nil, nil
	}

	// Build error message listing all cycles.
	msgs := make([]string, len(cycles))
	for i, cycle := range cycles {
		msgs[i] = strings.Join(cycle, " -> ")
	}
	return cycles, fmt.Errorf("dependency cycles detected:\n  %s", strings.Join(msgs, "\n  "))
}

// dfsDetectCycles performs DFS from the given node, detecting back edges (cycles).
func (g *DependencyGraph) dfsDetectCycles(
	node string,
	color map[string]int,
	parent map[string]string,
	cycles *[][]string,
) {
	color[node] = gray

	deps := g.Edges[node]
	for _, dep := range deps {
		if color[dep] == gray {
			// Back edge found: reconstruct cycle path.
			cycle := g.reconstructCycle(node, dep, parent)
			*cycles = append(*cycles, cycle)
		} else if color[dep] == white {
			parent[dep] = node
			g.dfsDetectCycles(dep, color, parent, cycles)
		}
		// black nodes are fully processed, skip them.
	}

	color[node] = black
}

// reconstructCycle builds the cycle path from the back edge target back to itself.
// Given that we found a back edge from `from` to `to` (where `to` is gray),
// we walk parent pointers from `from` back to `to` to reconstruct the cycle.
func (g *DependencyGraph) reconstructCycle(from, to string, parent map[string]string) []string {
	// Build path from `to` -> ... -> `from` -> `to`
	var path []string
	current := from
	for current != to {
		path = append(path, current)
		current = parent[current]
	}
	path = append(path, to)

	// Reverse to get: to -> ... -> from
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	// Close the cycle: to -> ... -> from -> to
	path = append(path, to)
	return path
}

// TopologicalSort returns agents in dependency order using Kahn's algorithm
// (BFS, in-degree based). Nodes with the same in-degree are sorted alphabetically
// for deterministic, reproducible output.
// Returns an error if a cycle prevents complete topological ordering.
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	// Compute in-degrees. In this graph, edges go from dependent -> dependency,
	// so for topological sort we need dependencies processed first.
	// We reverse the edge direction: if A depends on B, B should come before A.
	inDegree := make(map[string]int)
	reverseEdges := make(map[string][]string) // dependency -> list of dependents

	for node := range g.Nodes {
		inDegree[node] = 0
	}

	for dependent, deps := range g.Edges {
		for _, dep := range deps {
			reverseEdges[dep] = append(reverseEdges[dep], dependent)
			inDegree[dependent]++
		}
	}

	// Collect nodes with in-degree 0 (no dependencies).
	var queue []string
	for node := range g.Nodes {
		if inDegree[node] == 0 {
			queue = append(queue, node)
		}
	}
	sort.Strings(queue)

	var result []string
	for len(queue) > 0 {
		// Pop the first (alphabetically smallest) node.
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// For each dependent of this node, decrement in-degree.
		dependents := reverseEdges[node]
		sort.Strings(dependents)
		for _, dependent := range dependents {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
				sort.Strings(queue)
			}
		}
	}

	if len(result) != len(g.Nodes) {
		return nil, fmt.Errorf(
			"topological sort incomplete: processed %d of %d nodes (cycle exists)",
			len(result), len(g.Nodes),
		)
	}

	return result, nil
}

// Visualize returns a text-based DAG representation for debugging.
// Each node is listed alphabetically, followed by its dependencies indented below.
//
// Example output:
//
//	expert-backend
//	  <- manager-spec (depends on)
//	expert-frontend
//	  <- expert-backend (depends on)
//	  <- manager-spec (depends on)
func (g *DependencyGraph) Visualize() string {
	nodes := g.sortedNodes()

	var sb strings.Builder
	for i, node := range nodes {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(node)
		sb.WriteString("\n")

		deps := g.Edges[node]
		if len(deps) == 0 {
			sb.WriteString("  (no dependencies)\n")
			continue
		}

		// deps are already sorted from Build()
		for _, dep := range deps {
			sb.WriteString("  <- ")
			sb.WriteString(dep)
			sb.WriteString(" (depends on)\n")
		}
	}

	return sb.String()
}

// sortedNodes returns all node names in alphabetical order.
func (g *DependencyGraph) sortedNodes() []string {
	nodes := make([]string, 0, len(g.Nodes))
	for node := range g.Nodes {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	return nodes
}
