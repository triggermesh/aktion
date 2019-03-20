// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func ExampleBreadthFirst2_allPaths() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	var f graph.FromList
	alt.BreadthFirst2(g, nil, 0, 1, &f, func(n graph.NI) bool {
		return true
	})
	fmt.Println("Max path length:", f.MaxLen)
	p := make([]graph.NI, f.MaxLen)
	for n := range g {
		fmt.Println(n, f.PathTo(graph.NI(n), p))
	}
	// Output:
	// Max path length: 4
	// 0 []
	// 1 [1]
	// 2 []
	// 3 [1 4 3]
	// 4 [1 4]
	// 5 [1 4 3 5]
	// 6 [1 4 6]
}
