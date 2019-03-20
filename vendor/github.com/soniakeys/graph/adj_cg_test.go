// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// adj_cg_test.go -- tests for code in adj_cg.go.
//
// These are tests on the labeled versions of methods.
//
// See also adj_ro_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"testing"
	"text/template"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

func ExampleLabeledAdjacencyList_AnyLoop_loop() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 2}},
	}
	fmt.Println(g.AnyLoop())
	// Output:
	// true 2
}

func ExampleLabeledAdjacencyList_AnyLoop_noLoop() {
	g := graph.LabeledAdjacencyList{
		1: {{To: 0}},
	}
	lp, _ := g.AnyLoop()
	fmt.Println("has loop:", lp)
	// Output:
	// has loop: false
}

func ExampleLabeledAdjacencyList_ArcDensity() {
	// 0-->1
	// |
	// v
	// 2-->3
	g := graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}},
		2: {{To: 3}},
		3: {},
	}
	fmt.Println(g.ArcDensity())
	// Output:
	// 0.25
}

func ExampleLabeledAdjacencyList_ArcSize() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{ // simple graph
		2: {{To: 0}, {To: 1}},
	}}
	fmt.Println(g.ArcSize())

	// with reciprocals now.
	u := g.Undirected()
	// common term "size" for undirected graph is number of undirected edges.
	// size, m = ArcSize() / 2 here, but only because there are no loops.
	fmt.Println(u.ArcSize())

	g.LabeledAdjacencyList[1] = []graph.Half{{To: 1}} // add a loop
	//   2
	//  / \
	// 0   1---\
	//      \--/
	fmt.Println(g.ArcSize())

	// loops have no reciprocals.  ArcSize() / 2 no longer meaningful.
	fmt.Println(g.Undirected().ArcSize())
	// Output:
	// 2
	// 4
	// 3
	// 5
}

func ExampleLabeledAdjacencyList_BoundsOk() {
	var g graph.LabeledAdjacencyList
	ok, _, _ := g.BoundsOk() // zero value adjacency list is valid
	fmt.Println(ok)
	g = graph.LabeledAdjacencyList{
		0: {{To: 9}},
	}
	fmt.Println(g.BoundsOk()) // arc 0 to 9 invalid with only one node
	// Output:
	// true
	// false 0 {9 0}
}

func ExampleLabeledAdjacencyList_BreadthFirst() {
	//   <-0->
	//  /  |  \
	// v   v   v
	// 1-->2   4
	// ^   |   ^
	// |   v   |
	// \---3   5
	g := graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}, {To: 4}},
		1: {{To: 2}},
		2: {{To: 3}},
		3: {{To: 1}},
		5: {{To: 4}},
	}
	g.BreadthFirst(0, func(n graph.NI) { fmt.Println(n) })
	// Output:
	// 0
	// 1
	// 2
	// 4
	// 3
}

func ExampleLabeledAdjacencyList_DepthFirst() {
	//   <-0->
	//  /  |  \
	// v   v   v
	// 1-->2   4
	// ^   |   ^
	// |   v   |
	// \---3   5
	g := graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}, {To: 4}},
		1: {{To: 2}},
		2: {{To: 3}},
		3: {{To: 1}},
		5: {{To: 4}},
	}
	g.DepthFirst(0, func(n graph.NI) { fmt.Println(n) })
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleLabeledAdjacencyList_HasArc() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 2}, {To: 0}, {To: 1}, {To: 1}},
	}
	fmt.Println(g.HasArc(2, 1))
	fmt.Println(g.HasArc(2, 2)) // test for loop
	// Output:
	// true 3
	// true 1
}

func ExampleLabeledAdjacencyList_InduceBits() {
	// arcs directed down:
	//     1
	//    /|\\
	//  a/ | \\
	//  / b| c\\d
	// 0   |    2
	//  \  |   /
	//  e\ |  /f
	//    \| /
	//     3-
	g := graph.LabeledAdjacencyList{
		1: {{0, 'a'}, {3, 'b'}, {2, 'c'}, {2, 'd'}},
		0: {{3, 'e'}},
		2: {{3, 'f'}},
		3: {},
	}
	s := g.InduceBits(bits.NewGivens(2, 1, 3))
	fmt.Println("Subgraph:")
	for fr, to := range s.LabeledAdjacencyList {
		fmt.Print(fr, ": [")
		for _, h := range to {
			fmt.Printf("{%d, %c} ", h.To, h.Label)
		}
		fmt.Println("]")
	}
	fmt.Println("Sub NI -> Super NI")
	for b, p := range s.SuperNI {
		fmt.Printf("  %d         %d\n", b, p)
	}
	fmt.Println("Super NI -> Sub NI")
	template.Must(template.New("").Parse(
		`{{range $k, $v := .}}  {{$k}}         {{$v}}
{{end}}`)).Execute(os.Stdout, s.SubNI)
	// Output:
	// Subgraph:
	// 0: [{2, b} {1, c} {1, d} ]
	// 1: [{2, f} ]
	// 2: []
	// Sub NI -> Super NI
	//   0         1
	//   1         2
	//   2         3
	// Super NI -> Sub NI
	//   1         0
	//   2         1
	//   3         2
}

func ExampleLabeledAdjacencyList_InduceList() {
	// arcs directed down:
	//     1
	//    /|\\
	//  a/ | \\
	//  / b| c\\d
	// 0   |    2
	//  \  |   /
	//  e\ |  /f
	//    \| /
	//     3-
	g := graph.LabeledAdjacencyList{
		1: {{0, 'a'}, {3, 'b'}, {2, 'c'}, {2, 'd'}},
		0: {{3, 'e'}},
		2: {{3, 'f'}},
		3: {},
	}
	s := g.InduceList([]graph.NI{2, 1, 2, 3})
	fmt.Println("Subgraph:")
	for fr, to := range s.LabeledAdjacencyList {
		fmt.Print(fr, ": [")
		for _, h := range to {
			fmt.Printf("{%d, %c} ", h.To, h.Label)
		}
		fmt.Println("]")
	}
	fmt.Println("Sub NI -> Super NI")
	for b, p := range s.SuperNI {
		fmt.Printf("  %d         %d\n", b, p)
	}
	fmt.Println("Super NI -> Sub NI")
	template.Must(template.New("").Parse(
		`{{range $k, $v := .}}  {{$k}}         {{$v}}
{{end}}`)).Execute(os.Stdout, s.SubNI)
	// Output:
	// Subgraph:
	// 0: [{2, f} ]
	// 1: [{2, b} {0, c} {0, d} ]
	// 2: []
	// Sub NI -> Super NI
	//   0         2
	//   1         1
	//   2         3
	// Super NI -> Sub NI
	//   1         1
	//   2         0
	//   3         2
}

func ExampleLabeledAdjacencyList_IsolatedNodes() {
	//   0  1
	//  / \
	// 2   3  4
	g := graph.LabeledAdjacencyList{
		0: {{To: 2}, {To: 3}},
		4: {},
	}
	fmt.Println(g.IsolatedNodes().Slice())
	// Output:
	// [1 4]
}

func ExampleLabeledAdjacencyList_IsSimple() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 1}},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// true -1
}

func ExampleLabeledAdjacencyList_IsSimple_loop() {
	// arcs directed down
	//   2
	//  / \
	// 0   1---\
	//      \--/
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 1}},
		1: {{To: 1}}, // loop
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 1
}

func ExampleLabeledAdjacencyList_IsSimple_parallelArc() {
	// arcs directed down
	//   2
	//  /|\
	//  |/ \
	//  0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 1}, {To: 0}},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 2
}

func ExampleLabeledAdjacencyList_Order() {
	g := graph.LabeledAdjacencyList{ // maybe you think of node 0 as "unused",
		1: {{To: 2}},
		2: {{To: 1}},
	}
	fmt.Println(g.Order())                    // but graph still has 3 nodes.
	g = make(graph.LabeledAdjacencyList, 300) // empty graph,
	fmt.Println(g.Order())                    // with 300 nodes.
	fmt.Println(len(g))                       // equivalent.

	var u graph.LabeledUndirected
	u.AddEdge(graph.Edge{0, 99}, 0)
	fmt.Println(len(u.LabeledAdjacencyList)) // explicit.
	fmt.Println(u.Order())                   // a little more concise.
	// Output:
	// 3
	// 300
	// 300
	// 100
	// 100
}

func ExampleLabeledAdjacencyList_ParallelArcs() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 2}, {To: 0}, {To: 1}, {To: 1}},
	}
	fmt.Println(g.ParallelArcs(0, 2))
	fmt.Println(g.ParallelArcs(2, 0))
	fmt.Println(g.ParallelArcs(2, 1))
	fmt.Println(g.ParallelArcs(2, 2)) // returns loops on 2
	// Output:
	// []
	// [0 2]
	// [3 4]
	// [1]
}

func ExampleLabeledAdjacencyList_Permute() {
	//    0                              2
	//  x/ \y  Permute([2 0 1]) gives  x/ \y
	//  1-->2                          0-->1
	//    z                              z
	g := graph.LabeledAdjacencyList{
		0: {{1, 'x'}, {2, 'y'}},
		1: {{2, 'z'}},
		2: {},
	}
	g.Permute([]int{2, 0, 1})
	for fr, to := range g {
		fmt.Print(fr, ":")
		for _, to := range to {
			fmt.Printf(" (%d %c)", to.To, to.Label)
		}
		fmt.Println()
	}
	// Output:
	// 0: (1 z)
	// 1:
	// 2: (0 x) (1 y)
}

// not much of a test.  doesn't actually test that shuffle did anything,
// does a bunch of unrelated stuff.  It at least tests that shuffle doesn't
// corrupt the graph.
func TestShuffleArcListsLabeled(t *testing.T) {
	testCase := func(rad float64, r *rand.Rand) {
		g, _, _ := graph.LabeledGeometric(10, rad, r)
		c, _ := g.LabeledAdjacencyList.Copy()
		c.ShuffleArcLists(r)
		for fr, to := range g.LabeledAdjacencyList {
			sort.Slice(to, func(i, j int) bool {
				return to[i].To < to[j].To ||
					to[i].To == to[j].To && to[i].Label < to[j].Label
			})
			sh := c[fr]
			sort.Slice(sh, func(i, j int) bool {
				return sh[i].To < sh[j].To ||
					sh[i].To == sh[j].To && sh[i].Label < sh[j].Label
			})
			if (len(to) > 0 || len(sh) > 0) && !reflect.DeepEqual(to, sh) {
				t.Fatal(rad, r, fr, to, sh)
			}
		}
	}
	testCase(.2, nil)
	testCase(.6, rand.New(rand.NewSource(3)))
}

func ExampleLabeledSubgraph_AddNode() {
	// supergraph:
	//    0
	//   / \
	//  1-->2
	g := graph.LabeledAdjacencyList{
		0: {{1, -1}, {2, -1}},
		1: {{2, -1}},
		2: {},
	}
	s := g.InduceList(nil)    // construct empty subgraph
	fmt.Println(s.AddNode(2)) // first node added will have NI = 0
	fmt.Println(s.AddNode(1)) // next node added will have NI = 1
	fmt.Println(s.AddNode(1)) // returns existing mapping
	fmt.Println(s.AddNode(2)) // returns existing mapping
	fmt.Println("Subgraph:")  // (it has no arcs)
	for fr, to := range s.LabeledAdjacencyList {
		fmt.Printf("%d: %d\n", fr, to)
	}
	fmt.Println("Mappings:")
	fmt.Println(s.SuperNI)
	fmt.Println(graph.OrderMap(s.SubNI))
	// Output:
	// 0
	// 1
	// 1
	// 0
	// Subgraph:
	// 0: []
	// 1: []
	// Mappings:
	// [2 1]
	// map[1:1 2:0 ]
}

func ExampleLabeledSubgraph_AddNode_panic() {
	// supergraph:
	//    0
	//   / \
	//  1-->2
	g := graph.LabeledAdjacencyList{
		0: {{1, -1}, {2, -1}},
		1: {{2, -1}},
		2: {},
	}
	s := g.InduceList(nil)
	func() {
		defer func() { fmt.Println(recover()) }()
		fmt.Println(s.AddNode(-1))
	}()
	s.AddNode(0) // ok
	s.AddNode(2) // ok
	func() {
		defer func() { fmt.Println(recover()) }()
		fmt.Println(s.AddNode(3))
	}()
	// Output:
	// AddNode: NI -1 not in supergraph
	// AddNode: NI 3 not in supergraph
}

func ExampleLabeledSubgraph_AddArc() {
	// supergraph:
	//     0
	//    / \\
	//  x/  y\\z
	//  1      2
	g := graph.LabeledAdjacencyList{
		0: {{1, 'x'}, {2, 'y'}, {2, 'z'}},
		2: {},
	}
	s := g.InduceList(nil)                       // construct empty subgraph
	fmt.Println(s.AddArc(1, graph.Half{0, 'x'})) // no, not that direction
	fmt.Println(s.AddArc(0, graph.Half{1, 'y'})) // no, not with that label
	fmt.Println(s.AddArc(0, graph.Half{2, 'y'})) // okay
	fmt.Println("Subgraph:")
	for fr, to := range s.LabeledAdjacencyList {
		fmt.Print(fr, ":")
		for _, h := range to {
			fmt.Printf(" {%d, %c}", h.To, h.Label)
		}
		fmt.Println()
	}
	fmt.Println("Mappings:")
	// mapping from subgraph NIs to supergraph NIs
	fmt.Println(s.SuperNI)
	// mapping from supergraph NIs to subgraph NIs
	fmt.Println(graph.OrderMap(s.SubNI))
	// Output:
	// arc not available in supergraph
	// arc not available in supergraph
	// <nil>
	// Subgraph:
	// 0: {1, y}
	// 1:
	// Mappings:
	// [0 2]
	// map[0:0 2:1 ]
}

func ExampleLabeledSubgraph_AddArc_panic() {
	// supergraph:
	//    0
	//   / \\
	//  1    2
	g := graph.LabeledAdjacencyList{
		0: {{1, -1}, {2, -1}, {2, -1}},
		2: {},
	}
	s := g.InduceList(nil)
	func() {
		defer func() { fmt.Println(recover()) }()
		fmt.Println(s.AddArc(0, graph.Half{-1, -1}))
	}()
	func() {
		defer func() { fmt.Println(recover()) }()
		fmt.Println(s.AddArc(3, graph.Half{0, -1}))
	}()
	// Output:
	// AddArc: NI -1 not in supergraph
	// AddArc: NI 3 not in supergraph
}
