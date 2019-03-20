// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func ExampleTarjanCycles() {
	// 0-->1--->2-\
	// ^  ^^\  ^|  v
	// | / || / |  3
	// |/  \v/  v /
	// 4<---5<--6<
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 5},
		2: {3, 6},
		3: {6},
		4: {0, 1},
		5: {1, 2, 4},
		6: {5},
	}}
	alt.TarjanCycles(g, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 1 2 3 6 5 4]
	// [0 1 2 6 5 4]
	// [0 1 5 4]
	// [1 2 3 6 5]
	// [1 2 3 6 5 4]
	// [1 2 6 5]
	// [1 2 6 5 4]
	// [1 5]
	// [1 5 4]
	// [2 3 6 5]
	// [2 6 5]
}

func ExampleTarjanCyclesLabeled() {
	//     -1      2
	//  0----->1------>2--\
	//  ^     ^^\     ^|   \2
	// 3|  -4/ ||1   / |-1  v
	//  |   /  ||   /  |    3
	//  |  /   ||  /   |   /
	//  | /  -2|| /-2  |  /-1
	//  |/     \v/     v /
	//  4<------5<-----6<
	//      -1      1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1, Label: -1}},
		1: {{To: 2, Label: 2}, {To: 5, Label: 1}},
		2: {{To: 3, Label: 2}, {To: 6, Label: -1}},
		3: {{To: 6, Label: -1}},
		4: {{To: 0, Label: 3}, {To: 1, Label: -4}},
		5: {{To: 1, Label: -2}, {To: 2, Label: -2}, {To: 4, Label: -1}},
		6: {{To: 5, Label: 1}},
	}}
	alt.TarjanCyclesLabeled(g, func(c []graph.Half) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [{1 -1} {2 2} {3 2} {6 -1} {5 1} {4 -1} {0 3}]
	// [{1 -1} {2 2} {6 -1} {5 1} {4 -1} {0 3}]
	// [{1 -1} {5 1} {4 -1} {0 3}]
	// [{2 2} {3 2} {6 -1} {5 1} {1 -2}]
	// [{2 2} {3 2} {6 -1} {5 1} {4 -1} {1 -4}]
	// [{2 2} {6 -1} {5 1} {1 -2}]
	// [{2 2} {6 -1} {5 1} {4 -1} {1 -4}]
	// [{5 1} {1 -2}]
	// [{5 1} {4 -1} {1 -4}]
	// [{3 2} {6 -1} {5 1} {2 -2}]
	// [{6 -1} {5 1} {2 -2}]
}

func ExampleSCCKosaraju() {
	// /---0---\
	// |   |\--/
	// |   v
	// |   5<=>4---\
	// |   |   |   |
	// v   v   |   |
	// 7<=>6   |   |
	//     |   v   v
	//     \-->3<--2
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.Directed{graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
	}}
	alt.SCCKosaraju(g, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0]
	// [5 4]
	// [6 7]
	// [2 1 3]
}

func ExampleSCCPathBased() {
	// /---0---\
	// |   |\--/
	// |   v
	// |   5<=>4---\
	// |   |   |   |
	// v   v   |   |
	// 7<=>6   |   |
	//     |   v   v
	//     \-->3<--2
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.Directed{graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
	}}
	alt.SCCPathBased(g, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [1 3 2]
	// [7 6]
	// [4 5]
	// [0]
}

func ExampleSCCTarjan() {
	// /---0---\
	// |   |\--/
	// |   v
	// |   5<=>4---\
	// |   |   |   |
	// v   v   |   |
	// 7<=>6   |   |
	//     |   v   v
	//     \-->3<--2
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.Directed{graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
	}}
	alt.SCCTarjan(g, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [1 3 2]
	// [7 6]
	// [4 5]
	// [0]
}
