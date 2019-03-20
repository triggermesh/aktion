// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package dot_test

import (
	"fmt"
	"os"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func ExampleString() {
	// arcs directed down:
	// 0  4
	// | /|
	// |/ |
	// 2  3
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	s, _ := dot.String(g)
	fmt.Println(s)
	// Output:
	// digraph {
	//   0 -> 2
	//   4 -> {2 3}
	// }
}

func ExampleWrite_adjacencyList() {
	// arcs directed down:
	// 0  4
	// | /|
	// |/ |
	// 2  3
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	dot.Write(g, os.Stdout)
	// Output:
	// digraph {
	//   0 -> 2
	//   4 -> {2 3}
	// }
}

func ExampleWrite_adjacencyListParallelArcs() {
	// arcs directed down:
	// 0  4
	// | /|\
	// |/ \|
	// 2   3
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3, 3},
	}
	dot.Write(g, os.Stdout)
	// Output:
	// digraph {
	//   0 -> 2
	//   4 -> {2 3}
	//   4 -> {3}
	// }
}

func ExampleWrite_directed() {
	// arcs directed down:
	// 0  4
	// | /|
	// |/ |
	// 2  3
	g := graph.Directed{graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}}
	dot.Write(g, os.Stdout)
	// Output:
	// digraph {
	//   0 -> 2
	//   4 -> {2 3}
	// }
}

func ExampleWrite_directedLabeled() {
	// arcs directed down:
	//     0      4
	// (30)|     /|
	//     | (20) |
	//     |/     |(10)
	//     2      3
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{2, 30}},
		4: {{2, 20}, {3, 10}},
	}}
	dot.Write(g, os.Stdout)
	// Output:
	// digraph {
	//   0 -> 2 [label = 30]
	//   4 -> 2 [label = 20]
	//   4 -> 3 [label = 10]
	// }
}

func ExampleWrite_fromList() {
	//     0
	//    / \
	//   /   2
	//  /     \
	// 1       3
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
	}}
	f.RecalcLeaves()
	dot.Write(f, os.Stdout)
	// Output:
	// digraph {
	//   rankdir = BT
	//   1 -> 0
	//   2 -> 0
	//   3 -> 2
	//   {rank = same 1 3}
	// }
}

func ExampleWrite_labeledAdjacencyList() {
	// arcs directed down:
	//     0      4
	// (30)|     /|
	//     | (20) |
	//     |/     |(10)
	//     2      3
	g := graph.LabeledAdjacencyList{
		0: {{2, 30}},
		4: {{2, 20}, {3, 10}},
	}
	dot.Write(g, os.Stdout)
	// Output:
	// digraph {
	//   0 -> 2 [label = 30]
	//   4 -> 2 [label = 20]
	//   4 -> 3 [label = 10]
	// }
}

func ExampleWrite_undirected() {
	//   0
	//  / \
	// 1---2
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	dot.Write(g, os.Stdout)
	// Output:
	// graph {
	//   0 -- {1 2}
	//   1 -- 2
	// }
}

func ExampleWrite_undirectedLabeled() {
	//       0
	// (12) / \ (17)
	//     1---2
	//      (64)
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 12)
	g.AddEdge(graph.Edge{0, 2}, 17)
	g.AddEdge(graph.Edge{1, 2}, 64)
	dot.Write(g, os.Stdout)
	// Output:
	// graph {
	//   0 -- 1 [label = 12]
	//   0 -- 2 [label = 17]
	//   1 -- 2 [label = 64]
	// }
}

func ExampleWrite_weightedEdgeList() {
	//              (label 0, wt 1.6)
	//          0----------------------2
	// (label 1 |                     /
	//  wt .33) |  ------------------/
	//          | / (label 2, wt 1.7)
	//          |/
	//          1
	weights := []float64{
		0: 1.6,
		1: .33,
		2: 1.7,
	}
	g := graph.WeightedEdgeList{
		WeightFunc: func(l graph.LI) float64 { return weights[int(l)] },
		Order:      3,
		Edges: []graph.LabeledEdge{
			{graph.Edge{0, 1}, 1},
			{graph.Edge{1, 0}, 1},

			{graph.Edge{0, 2}, 0},
			{graph.Edge{2, 0}, 0},

			{graph.Edge{1, 2}, 2},
			{graph.Edge{2, 1}, 2},
		},
	}
	dot.Write(g, os.Stdout)
	// Output:
	// graph {
	//   0 -- 1 [label = "0.33"]
	//   0 -- 2 [label = "1.6"]
	//   1 -- 2 [label = "1.7"]
	// }
}
