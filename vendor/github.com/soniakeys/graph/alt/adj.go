// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt

import "github.com/soniakeys/graph"

// AnyParallelMap identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the method returns true and
// results fr and to represent an example where there are parallel arcs
// from node `fr` to node `to`.
//
// If there are no parallel arcs, the method returns false, -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Map" in the method name indicates that a Go map is used to detect parallel
// arcs.  Compared to method AnyParallelSort, this gives better asymptotic
// performance for large dense graphs but may have increased overhead for
// small or sparse graphs.
func AnyParallelMap(g graph.AdjacencyList) (has bool, fr, to graph.NI) {
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		m := map[graph.NI]struct{}{}
		for _, to := range to {
			if _, ok := m[to]; ok {
				return true, graph.NI(n), to
			}
			m[to] = struct{}{}
		}
	}
	return false, -1, -1
}
