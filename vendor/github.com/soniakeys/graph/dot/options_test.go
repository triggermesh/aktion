// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package dot_test

import (
	"fmt"
	"os"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func ExampleDirected() {
	// arcs directed down:
	// 0  2
	// | /|
	// |/ |
	// 3  4
	g := graph.AdjacencyList{
		0: {3},
		2: {3, 4},
		4: {},
	}
	// default for AdjacencyList is directed
	dot.Write(g, os.Stdout)
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout)

	// Directed(false) generates error witout reciprocal arcs
	err := dot.Write(g, os.Stdout, dot.Directed(false))
	fmt.Fprintln(os.Stdout, "Error:", err)
	fmt.Fprintln(os.Stdout)

	// undirected
	u := graph.Directed{g}.Undirected()
	dot.Write(u.AdjacencyList, os.Stdout, dot.Directed(false))

	// Output:
	// digraph {
	//   0 -> 3
	//   2 -> {3 4}
	// }
	//
	// Error: directed graph
	//
	// graph {
	//   0 -- 3
	//   2 -- {3 4}
	// }
}

func ExampleEdgeLabel() {
	// arcs directed down:
	//      0       4
	// (.33)|      /|
	//      | (1.7) |
	//      |/      |(2e117)
	//      2       3
	weights := map[int]float64{
		30: .33,
		20: 1.7,
		10: 2e117,
	}
	lf := func(l graph.LI) string {
		return fmt.Sprintf(`"%g"`, weights[int(l)])
	}
	g := graph.LabeledAdjacencyList{
		0: {{2, 30}},
		4: {{2, 20}, {3, 10}},
	}
	dot.Write(g, os.Stdout, dot.EdgeLabel(lf))
	// Output:
	// digraph {
	//   0 -> 2 [label = "0.33"]
	//   4 -> 2 [label = "1.7"]
	//   4 -> 3 [label = "2e+117"]
	// }
}

func ExampleGraphAttr() {
	// arcs directed right:
	// 0---2
	//  \ / \
	//   1---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {},
	}
	dot.Write(g, os.Stdout, dot.GraphAttr("rankdir", "LR"))
	// Output:
	// digraph {
	//   rankdir = LR
	//   0 -> {1 2}
	//   1 -> {2 3}
	//   2 -> 3
	// }
}

func ExampleIndent() {
	// arcs directed down:
	// 0  4
	// | /|
	// |/ |
	// 2  3
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	dot.Write(g, os.Stdout, dot.Indent("")) // (default indent is 2 spaces)
	// Output:
	// digraph {
	// 0 -> 2
	// 4 -> {2 3}
	// }
}

func ExampleIsolated() {
	// 0  1-->2
	g := graph.AdjacencyList{
		1: {2},
		2: {},
	}
	dot.Write(g, os.Stdout, dot.Isolated(true))
	// Output:
	// digraph {
	//   0
	//   1 -> 2
	// }
}

func ExampleNodePos() {
	// 0--1
	// |\
	// | \
	// 2  3
	f := graph.AdjacencyList{
		0: {1, 2, 3},
		3: {},
	}
	pos := []struct{ x, y float64 }{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}
	dot.Write(f, os.Stdout, dot.NodePos(func(n graph.NI) string {
		return fmt.Sprintf("%.0f,%.0f", 4*pos[n].x, 4*pos[n].y)
	}))
	// Output:
	// digraph {
	//   node [shape=point]
	//   0 [pos="0,0!"]
	//   1 [pos="0,4!"]
	//   2 [pos="4,0!"]
	//   3 [pos="4,4!"]
	//   0 -> {1 2 3}
	// }
}

func ExampleNodeID() {
	// arcs directed down:
	// A  D
	// | /|
	// |/ |
	// B  C
	labels := []string{
		0: "A",
		4: "D",
		2: "B",
		3: "C",
	}
	lf := func(n graph.NI) string { return labels[n] }
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	dot.Write(g, os.Stdout, dot.NodeID(lf))
	// Output:
	// digraph {
	//   A -> B
	//   D -> {B C}
	// }
}

func ExampleNodeID_construction() {
	// arcs directed down:
	// A  D
	// | /|
	// |/ |
	// B  C
	var g graph.AdjacencyList

	// example graph construction mechanism
	labels := []string{}
	nodes := map[string]graph.NI{}
	node := func(l string) graph.NI {
		if n, ok := nodes[l]; ok {
			return n
		}
		n := graph.NI(len(labels))
		labels = append(labels, l)
		g = append(g, nil)
		nodes[l] = n
		return n
	}
	addArc := func(fr, to string) {
		f := node(fr)
		g[f] = append(g[f], node(to))
	}

	// construct graph
	addArc("A", "B")
	addArc("D", "B")
	addArc("D", "C")

	// generate dot
	lf := func(n graph.NI) string { return labels[n] }
	dot.Write(g, os.Stdout, dot.NodeID(lf))

	// Output:
	// digraph {
	//   A -> B
	//   D -> {B C}
	// }
}

func ExampleUndirectArcs() {
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
			{graph.Edge{0, 2}, 0},
			{graph.Edge{1, 2}, 2},
		},
	}
	dot.Write(g, os.Stdout, dot.UndirectArcs(true))
	// Output:
	// graph {
	//   0 -- 1 [label = "0.33"]
	//   0 -- 2 [label = "1.6"]
	//   1 -- 2 [label = "1.7"]
	// }
}
