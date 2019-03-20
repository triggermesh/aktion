// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func ExampleAnyParallelMap_parallelArcs() {
	g := graph.AdjacencyList{
		1: {0, 0},
	}
	// result true 1 0 means parallel arcs from node 1 to node 0
	fmt.Println(alt.AnyParallelMap(g))
	// Output:
	// true 1 0
}

func ExampleAnyParallelMap_noParallelArcs() {
	g := graph.AdjacencyList{
		1: {0},
	}
	fmt.Println(alt.AnyParallelMap(g)) // result false -1 -1 means no parallel arc
	// Output:
	// false -1 -1
}
