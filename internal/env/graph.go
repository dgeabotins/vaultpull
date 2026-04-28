// Package env provides utilities for working with .env files.
package env

import (
	"fmt"
	"sort"
	"strings"
)

// GraphNode represents a single key in the dependency graph.
type GraphNode struct {
	Key  string
	Deps []string // keys this node references via ${VAR} or $VAR
}

// GraphResult holds the output of a dependency graph analysis.
type GraphResult struct {
	Nodes   []GraphNode
	Cycles  [][]string // each inner slice is one cycle path
	Ordered []string   // topological order (empty if cycles exist)
}

// Summary returns a human-readable summary of the graph result.
func (r GraphResult) Summary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("keys: %d", len(r.Nodes)))
	if len(r.Cycles) > 0 {
		sb.WriteString(fmt.Sprintf(", cycles: %d", len(r.Cycles)))
	} else {
		sb.WriteString(", no cycles detected")
	}
	return sb.String()
}

// Graph analyses interpolation dependencies between keys in m.
// It returns a GraphResult with dependency edges, any detected cycles,
// and a topological ordering when the graph is acyclic.
func Graph(m map[string]string) GraphResult {
	nodes := buildNodes(m)
	cycles := detectCycles(nodes)

	var ordered []string
	if len(cycles) == 0 {
		ordered = topoSort(nodes)
	}

	return GraphResult{
		Nodes:   nodes,
		Cycles:  cycles,
		Ordered: ordered,
	}
}

// buildNodes constructs a GraphNode for every key, parsing $VAR / ${VAR}
// references from values.
func buildNodes(m map[string]string) []GraphNode {
	nodes := make([]GraphNode, 0, len(m))
	for k, v := range m {
		deps := extractRefs(v)
		// only keep deps that actually exist as keys in the map
		filtered := deps[:0]
		for _, d := range deps {
			if _, ok := m[d]; ok {
				filtered = append(filtered, d)
			}
		}
		nodes = append(nodes, GraphNode{Key: k, Deps: filtered})
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Key < nodes[j].Key })
	return nodes
}

// extractRefs returns all variable names referenced in s via ${VAR} or $VAR.
func extractRefs(s string) []string {
	seen := map[string]bool{}
	var refs []string
	i := 0
	for i < len(s) {
		if s[i] != '$' {
			i++
			continue
		}
		i++ // skip '$'
		if i >= len(s) {
			break
		}
		var name string
		if s[i] == '{' {
			i++ // skip '{'
			j := i
			for j < len(s) && s[j] != '}' {
				j++
			}
			name = s[i:j]
			i = j + 1
		} else {
			j := i
			for j < len(s) && isIdentChar(s[j]) {
				j++
			}
			name = s[i:j]
			i = j
		}
		if name != "" && !seen[name] {
			seen[name] = true
			refs = append(refs, name)
		}
	}
	return refs
}

func isIdentChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') || b == '_'
}

// detectCycles uses DFS to find all cycles in the dependency graph.
func detectCycles(nodes []GraphNode) [][]string {
	adj := make(map[string][]string, len(nodes))
	for _, n := range nodes {
		adj[n.Key] = n.Deps
	}

	const (
		unvisited = 0
		visiting  = 1
		visited   = 2
	)
	state := make(map[string]int, len(nodes))
	var cycles [][]string
	path := []string{}

	var dfs func(key string)
	dfs = func(key string) {
		state[key] = visiting
		path = append(path, key)
		for _, dep := range adj[key] {
			switch state[dep] {
			case visiting:
				// found a cycle — extract the loop portion
				start := 0
				for start < len(path) && path[start] != dep {
					start++
				}
				cycle := make([]string, len(path)-start)
				copy(cycle, path[start:])
				cycles = append(cycles, cycle)
			case unvisited:
				dfs(dep)
			}
		}
		path = path[:len(path)-1]
		state[key] = visited
	}

	for _, n := range nodes {
		if state[n.Key] == unvisited {
			dfs(n.Key)
		}
	}
	return cycles
}

// topoSort returns keys in dependency order (dependencies before dependents).
func topoSort(nodes []GraphNode) []string {
	adj := make(map[string][]string, len(nodes))
	for _, n := range nodes {
		adj[n.Key] = n.Deps
	}

	visited := make(map[string]bool, len(nodes))
	var order []string

	var visit func(key string)
	visit = func(key string) {
		if visited[key] {
			return
		}
		visited[key] = true
		for _, dep := range adj[key] {
			visit(dep)
		}
		order = append(order, key)
	}

	for _, n := range nodes {
		visit(n.Key)
	}
	return order
}
