// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"os"
	"text/template"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_construction() {
	_ = make(graph.AdjacencyList, 10) // make an empty graph with 10 nodes
	_ = graph.AdjacencyList{{1}, {0}} // a graph with 2 nodes and 2 arcs
	// the same graph, with "keyed elements".  See the Go spec!
	_ = graph.AdjacencyList{
		0: {1},
		1: {0},
	}
	_ = graph.AdjacencyList{9: nil} // empty graph with 10 nodes
}

func ExampleAdjacencyList_AnyParallel_parallelArcs() {
	//   0
	//  / \
	// 1==>2
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2, 2},
		2: {},
	}
	// result true 1 2 means parallel arcs from node 1 to node 2
	fmt.Println(g.AnyParallel())
	// Output:
	// true 1 2
}

func ExampleAdjacencyList_AnyParallel_noParallelArcs() {
	//   0
	//  / \
	// 1-->2
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {},
	}
	// result false -1 -1 means no parallel arc
	fmt.Println(g.AnyParallel())
	// Output:
	// false -1 -1
}

func ExampleAdjacencyList_Complement() {
	//  0            0<-
	//  |\           ^  \
	//  v v  comp=>  |   \
	//  2 1          2<==>1
	g := graph.AdjacencyList{
		0: {1, 2},
		2: {},
	}
	for fr, to := range g.Complement() {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 []
	// 1 [0 2]
	// 2 [0 1]
}

func ExampleAdjacencyList_IsUndirected() {
	// 0<--    2<--\
	//  \  \   |   |
	//   -->1  \---/
	g := graph.AdjacencyList{
		0: {1},
		1: {0},
		2: {2},
	}
	ud, _, _ := g.IsUndirected()
	fmt.Println(ud)
	// 0<--
	//  \  \
	//   -->1<--2
	g = graph.AdjacencyList{
		0: {1},
		1: {0},
		2: {1},
	}
	fmt.Println(g.IsUndirected())
	// Output:
	// true
	// false 2 1
}

// A directed graph with negative arc weights.
// Arc weights are encoded simply as label numbers.
func ExampleLabeledAdjacencyList_DistanceMatrix() {
	//   (-1)   (4)
	//  0---->2---->1
	//  ^     |     |
	//  |(2)  |(3)  |(-2)
	//  |     v     |
	//  ------3<-----
	g := graph.LabeledAdjacencyList{
		0: {{To: 2, Label: -1}},
		1: {{To: 3, Label: -2}},
		2: {{To: 1, Label: 4}, {To: 3, Label: 3}},
		3: {{To: 0, Label: 2}},
	}
	d := g.DistanceMatrix(func(l graph.LI) float64 { return float64(l) })
	for _, di := range d {
		fmt.Printf("%5.0f\n", di)
	}
	// Output:
	// [    0  +Inf    -1  +Inf]
	// [ +Inf     0  +Inf    -2]
	// [ +Inf     4     0     3]
	// [    2  +Inf  +Inf     0]
}

func ExampleLabeledAdjacencyList_HasArcLabel() {
	//    /--\
	//   2<--/
	//  ||  b
	// a||c
	//  ||
	//  vv
	//   0
	g := graph.LabeledAdjacencyList{
		2: {{0, 'a'}, {2, 'b'}, {0, 'c'}},
	}
	fmt.Println(g.HasArcLabel(2, 0, 'c'))
	fmt.Println(g.HasArcLabel(1, 0, 'c'))
	fmt.Println(g.HasArcLabel(2, 0, 'z'))
	fmt.Println(g.HasArcLabel(2, 2, 'b')) // test for loop
	// Output:
	// true 2
	// false -1
	// false -1
	// true 1
}

func ExampleLabeledAdjacencyList_AnyParallel_parallelArcs() {
	//    0
	//   / \\
	//  /  a\\b
	// 1----->2
	g := graph.LabeledAdjacencyList{
		0: {{1, 0}, {2, 'a'}, {2, 'b'}}, // different labels, still parallel
		1: {{2, 0}},
		2: {},
	}
	// result true 0 2 means parallel arc from node 0 to node 2
	fmt.Println(g.AnyParallel())
	// Output:
	// true 0 2
}

func ExampleLabeledAdjacencyList_AnyParallel_noParallelArcs() {
	g := graph.LabeledAdjacencyList{
		1: {{0, 0}},
	}
	// result false -1 -1 means no parallel arc
	fmt.Println(g.AnyParallel())
	// Output:
	// false -1 -1
}

func ExampleLabeledAdjacencyList_AnyParallelLabel() {
	//    0
	//   / \\
	//  /  a\\a
	// 1----->2
	g := graph.LabeledAdjacencyList{
		0: {{1, 0}, {2, 'a'}, {2, 'a'}}, // same labels
		1: {{2, 0}},
		2: {},
	}
	fmt.Println(g.AnyParallelLabel())
	// Output:
	// true 0 {2 97}
}

func ExampleLabeledAdjacencyList_AnyParallelLabel_none() {
	//    0
	//   / \\
	//  /  a\\b
	// 1----->2
	g := graph.LabeledAdjacencyList{
		0: {{1, 0}, {2, 'a'}, {2, 'b'}}, // different labels
		1: {{2, 0}},
		2: {},
	}
	fmt.Println(g.AnyParallelLabel())
	// Output:
	// false -1 {0 0}
}

func ExampleLabeledAdjacencyList_IsUndirected() {
	//             0<--
	// (label 'A')  \  \ (matching label 'A' on reciprocal)
	//               -->1
	// 2<--\
	// |   | (label 'B' on loop)
	// \---/
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 'A'}},
		1: {{0, 'A'}},
		2: {{2, 'B'}},
	}
	ok, _, _ := g.IsUndirected()
	fmt.Println(ok)
	// Output:
	// true
}

func ExampleLabeledAdjacencyList_IsUndirected_undirectedMultigraph() {
	// lines shown here are edges (arcs with reciprocals.)
	//               0---
	//  (Label: 'A')  \  \  (Label: 'B')
	//                 ---1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'A')
	g.AddEdge(graph.Edge{0, 1}, 'B')
	ok, _, _ := g.IsUndirected()
	fmt.Println(ok)
	// Output:
	// true
}

func ExampleLabeledAdjacencyList_IsUndirected_labelMismatch() {
	// directed graph, arcs with different labels
	//               0<--
	//  (Label: 'A')  \  \  (Label: 'B')
	//                 -->1
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 'A'}},
		1: {{To: 0, Label: 'B'}},
	}
	ok, fr, to := g.IsUndirected()
	fmt.Printf("%t %d {%d %c}\n", ok, fr, to.To, to.Label)
	// Output:
	// false 0 {1 A}
}

func ExampleLabeledAdjacencyList_ArcLabels() {
	//      0
	//    // \
	//  x//y  \x
	//  //     \
	// 1------->2
	//     z
	g := graph.LabeledAdjacencyList{
		0: {{1, 'x'}, {1, 'y'}, {2, 'x'}},
		1: {{2, 'z'}},
		2: {},
	}
	l := g.ArcLabels()
	fmt.Println(graph.OrderMap(l))
	// prettier,
	template.Must(template.New("").Parse(
		`{{range $k, $v := .}}{{printf "%c" $k}}: {{$v}}
{{end}}`)).Execute(os.Stdout, l)
	// Output:
	// map[120:2 121:1 122:1 ]
	// x: 2
	// y: 1
	// z: 1
}

func ExampleLabeledAdjacencyList_ArcLabels_undirected() {
	//         /----\
	// 0------1-----/
	//   3772   2089
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 3772) // edge has reciprocal arcs
	g.AddEdge(graph.Edge{1, 1}, 2089) // loop has just one arc
	fmt.Println(graph.OrderMap(g.ArcLabels()))
	// Output:
	// map[2089:1 3772:2 ]
}

func ExampleLabeledAdjacencyList_NegativeArc() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 0}, {To: 1, Label: 1}},
	}
	arcWeights := []float64{0, .5}
	w := func(label graph.LI) float64 { return arcWeights[label] }
	fmt.Println(g.NegativeArc(w))
	g[0] = []graph.Half{{To: 1, Label: graph.LI(len(arcWeights))}}
	arcWeights = append(arcWeights, -2)
	fmt.Println(g.NegativeArc(w))
	// Output:
	// false
	// true
}

func ExampleLabeledAdjacencyList_ParallelArcsLabel() {
	g := graph.LabeledAdjacencyList{
		2: {{0, 10}, {2, 20}, {0, 10}, {0, 30}},
	}
	fmt.Println(g.ParallelArcsLabel(2, 0, 10))
	fmt.Println(g.ParallelArcsLabel(2, 0, 30))
	fmt.Println(g.ParallelArcsLabel(0, 0, 30)) // returns loops on 0
	fmt.Println(g.ParallelArcsLabel(2, 0, 100))
	// Output:
	// [0 2]
	// [3]
	// []
	// []
}

func ExampleLabeledAdjacencyList_Unlabeled() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}
	fmt.Println("Input:")
	for fr, to := range g {
		fmt.Printf("%d, %#v\n", fr, to)
	}
	fmt.Println("\nUnlabeled:")
	for fr, to := range g.Unlabeled() {
		fmt.Printf("%d, %#v\n", fr, to)
	}
	// Output:
	// Input:
	// 0, []graph.Half(nil)
	// 1, []graph.Half(nil)
	// 2, []graph.Half{graph.Half{To:0, Label:7}, graph.Half{To:1, Label:9}}
	//
	// Unlabeled:
	// 0, []graph.NI{}
	// 1, []graph.NI{}
	// 2, []graph.NI{0, 1}
}

func ExampleLabeledAdjacencyList_WeightedInDegree() {
	//  0
	//  | (weight = label: 3)
	//  v
	//  1
	//  | (4)
	//  v
	//  2<-\
	//  \--/ (5)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 3}},
		1: {{2, 4}},
		2: {{2, 5}},
	}
	w := func(l graph.LI) float64 { return float64(l) }
	fmt.Println(g.WeightedInDegree(w))
	// Output:
	// [0 3 9]
}

func ExampleLabeledAdjacencyList_WeightedOutDegree() {
	//  0
	//  | (weight = label: 3)
	//  v
	//  1
	//  | (4)
	//  v
	//  2<-\
	//  \--/ (5)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 3}},
		1: {{2, 4}},
		2: {{2, 5}},
	}
	w := func(l graph.LI) float64 { return float64(l) }
	fmt.Println("node  weighted out degree")
	for n := range g {
		fmt.Println(n, "   ", g.WeightedOutDegree(graph.NI(n), w))
	}
	// Output:
	// node  weighted out degree
	// 0     3
	// 1     4
	// 2     5
}

func ExampleLabeledAdjacencyList_WeightedOutDegree_undirected() {
	//  0
	//  | (weight = label: 3)
	//  |
	//  1
	//  | (4)
	//  |
	//  2--\
	//  \--/ (5)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 3}},
		1: {{0, 3}, {2, 4}},
		2: {{1, 4}, {2, 5}},
	}
	w := func(l graph.LI) float64 { return float64(l) }
	ok, _, _ := g.IsUndirected()
	fmt.Println("undirected:", ok)
	fmt.Println()
	fmt.Println("node  weighted out-degree  weighted in-degree")
	ind := g.WeightedInDegree(w)
	for n := range g {
		fmt.Println(n, "   ", g.WeightedOutDegree(graph.NI(n), w),
			"                  ", ind[n])
	}
	// Output:
	// undirected: true
	//
	// node  weighted out-degree  weighted in-degree
	// 0     3                    3
	// 1     7                    7
	// 2     9                    9
}
