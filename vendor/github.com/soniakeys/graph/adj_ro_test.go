// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// adj_ro_test.go -- tests on adj_RO.go
//
// These are tests on the code-generated unlabeled versions of methods.
//
// If testing prompts changes in the tested method, be sure to edit
// adj_cg.go, go generate to generate the adj_RO.go, and then retest.
// Do not edit adj_RO.go
//
// See also adj_cg_test.go when editing this file.  Try to keep the tests
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

func ExampleAdjacencyList_AnyLoop_loop() {
	g := graph.AdjacencyList{
		2: {2},
	}
	fmt.Println(g.AnyLoop())
	// Output:
	// true 2
}

func ExampleAdjacencyList_AnyLoop_noLoop() {
	g := graph.AdjacencyList{
		1: {0},
	}
	lp, _ := g.AnyLoop()
	fmt.Println("has loop:", lp)
	// Output:
	// has loop: false
}

func ExampleAdjacencyList_ArcDensity() {
	// 0-->1
	// |
	// v
	// 2-->3
	g := graph.AdjacencyList{
		0: {1, 2},
		2: {3},
		3: {},
	}
	fmt.Println(g.ArcDensity())
	// Output:
	// 0.25
}

func ExampleAdjacencyList_ArcSize() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.Directed{graph.AdjacencyList{ // simple graph
		2: {0, 1},
	}}
	fmt.Println(g.ArcSize())

	// with reciprocals now.
	u := g.Undirected()
	// common term "size" for undirected graph is number of undirected edges.
	// size, m = ArcSize() / 2 here, but only because there are no loops.
	fmt.Println(u.ArcSize())

	g.AdjacencyList[1] = []graph.NI{1} // add a loop
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

func ExampleAdjacencyList_BoundsOk() {
	var g graph.AdjacencyList
	ok, _, _ := g.BoundsOk() // zero value adjacency list is valid
	fmt.Println(ok)
	g = graph.AdjacencyList{
		0: {9},
	}
	fmt.Println(g.BoundsOk()) // arc 0 to 9 invalid with only one node
	// Output:
	// true
	// false 0 9
}

func ExampleAdjacencyList_BreadthFirst() {
	//   <-0->
	//  /  |  \
	// v   v   v
	// 1-->2   4
	// ^   |   ^
	// |   v   |
	// \---3   5
	g := graph.AdjacencyList{
		0: {1, 2, 4},
		1: {2},
		2: {3},
		3: {1},
		5: {4},
	}
	g.BreadthFirst(0, func(n graph.NI) { fmt.Println(n) })
	// Output:
	// 0
	// 1
	// 2
	// 4
	// 3
}

func ExampleAdjacencyList_DepthFirst() {
	//   <-0->
	//  /  |  \
	// v   v   v
	// 1-->2   4
	// ^   |   ^
	// |   v   |
	// \---3   5
	g := graph.AdjacencyList{
		0: {1, 2, 4},
		1: {2},
		2: {3},
		3: {1},
		5: {4},
	}
	g.DepthFirst(0, func(n graph.NI) { fmt.Println(n) })
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleAdjacencyList_HasArc() {
	g := graph.AdjacencyList{
		2: {0, 2, 0, 1, 1},
	}
	fmt.Println(g.HasArc(2, 1))
	fmt.Println(g.HasArc(2, 2)) // test for loop
	// Output:
	// true 3
	// true 1
}

func ExampleAdjacencyList_InduceBits() {
	// arcs directed down:
	//   1
	//  /|\\
	// 0 |  2
	//  \| /
	//   3-
	g := graph.AdjacencyList{
		1: {0, 3, 2, 2},
		0: {3},
		2: {3},
		3: {},
	}
	s := g.InduceBits(bits.NewGivens(2, 1, 3))
	fmt.Println("Subgraph:")
	for fr, to := range s.AdjacencyList {
		fmt.Printf("%d: %d\n", fr, to)
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
	// 0: [2 1 1]
	// 1: [2]
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

func ExampleAdjacencyList_InduceList() {
	// arcs directed down:
	//   1
	//  /|\\
	// 0 |  2
	//  \| /
	//   3-
	g := graph.AdjacencyList{
		1: {0, 3, 2, 2},
		0: {3},
		2: {3},
		3: {},
	}
	s := g.InduceList([]graph.NI{2, 1, 2, 3})
	fmt.Println("Subgraph:")
	for fr, to := range s.AdjacencyList {
		fmt.Printf("%d: %d\n", fr, to)
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
	// 0: [2]
	// 1: [2 0 0]
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
func ExampleAdjacencyList_IsolatedNodes() {
	//   0  1
	//  / \
	// 2   3  4
	g := graph.AdjacencyList{
		0: {2, 3},
		4: {},
	}
	fmt.Println(g.IsolatedNodes().Slice())
	// Output:
	// [1 4]
}

func ExampleAdjacencyList_IsSimple() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.AdjacencyList{
		2: {0, 1},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// true -1
}

func ExampleAdjacencyList_IsSimple_loop() {
	// arcs directed down
	//   2
	//  / \
	// 0   1---\
	//      \--/
	g := graph.AdjacencyList{
		2: {0, 1},
		1: {1}, // loop
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 1
}

func ExampleAdjacencyList_IsSimple_parallelArc() {
	// arcs directed down
	//   2
	//  /|\
	//  |/ \
	//  0   1
	g := graph.AdjacencyList{
		2: {0, 1, 0},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 2
}

func ExampleAdjacencyList_Order() {
	g := graph.AdjacencyList{ // maybe you think of node 0 as "unused",
		1: {2},
		2: {1},
	}
	fmt.Println(g.Order())             // but graph still has 3 nodes.
	g = make(graph.AdjacencyList, 300) // empty graph,
	fmt.Println(g.Order())             // with 300 nodes.
	fmt.Println(len(g))                // equivalent.

	var u graph.Undirected
	u.AddEdge(0, 99)
	fmt.Println(len(u.AdjacencyList)) // explicit.
	fmt.Println(u.Order())            // a little more concise.
	// Output:
	// 3
	// 300
	// 300
	// 100
	// 100
}

func ExampleAdjacencyList_ParallelArcs() {
	g := graph.AdjacencyList{
		2: {0, 2, 0, 1, 1},
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

func ExampleAdjacencyList_Permute() {
	//    0                             2
	//   / \   Permute([2 0 1]) gives  / \
	//  1-->2                         0-->1
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {},
	}
	g.Permute([]int{2, 0, 1})
	for fr, to := range g {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [1]
	// 1 []
	// 2 [0 1]
}

// not much of a test.  doesn't actually test that shuffle did anything,
// does a bunch of unrelated stuff.  It at least tests that shuffle doesn't
// corrupt the graph.
func TestShuffleArcLists(t *testing.T) {
	testCase := func(p float64, r *rand.Rand) {
		g, _ := graph.GnpUndirected(10, p, r)
		c, _ := g.AdjacencyList.Copy()
		c.ShuffleArcLists(r)
		for fr, to := range g.AdjacencyList {
			sort.Slice(to, func(i, j int) bool { return to[i] < to[j] })
			sh := c[fr]
			sort.Slice(sh, func(i, j int) bool { return sh[i] < sh[j] })
			if (len(to) > 0 || len(sh) > 0) && !reflect.DeepEqual(to, sh) {
				t.Fatal(p, r, fr, to, sh)
			}
		}
	}
	testCase(.1, nil)
	testCase(.9, rand.New(rand.NewSource(3)))
}

func ExampleSubgraph_AddNode() {
	// supergraph:
	//    0
	//   / \
	//  1-->2
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {},
	}
	s := g.InduceList(nil)    // construct empty subgraph
	fmt.Println(s.AddNode(2)) // first node added will have NI = 0
	fmt.Println(s.AddNode(1)) // next node added will have NI = 1
	fmt.Println(s.AddNode(1)) // returns existing mapping
	fmt.Println(s.AddNode(2)) // returns existing mapping
	fmt.Println("Subgraph:")  // (it has no arcs)
	for fr, to := range s.AdjacencyList {
		fmt.Printf("%d: %d\n", fr, to)
	}
	fmt.Println("Mappings:")
	// mapping from subgraph NIs to supergraph NIs
	fmt.Println(s.SuperNI)
	// mapping from supergraph NIs to subgraph NIs
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

func ExampleSubgraph_AddNode_panic() {
	// supergraph:
	//    0
	//   / \
	//  1-->2
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
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

func ExampleSubgraph_AddArc() {
	// supergraph:
	//    0
	//   / \\
	//  1    2
	g := graph.AdjacencyList{
		0: {1, 2, 2},
		2: {},
	}
	s := g.InduceList(nil)      // construct empty subgraph
	fmt.Println(s.AddArc(0, 2)) // okay
	fmt.Println(s.AddArc(0, 2)) // adding one parallel arc okay
	fmt.Println(s.AddArc(0, 2)) // adding another not okay
	fmt.Println(s.AddArc(1, 2)) // arc not in supergraph at all
	fmt.Println("Subgraph:")
	for fr, to := range s.AdjacencyList {
		fmt.Printf("%d: %d\n", fr, to)
	}
	fmt.Println("Mappings:")
	// mapping from subgraph NIs to supergraph NIs
	fmt.Println(s.SuperNI)
	// mapping from supergraph NIs to subgraph NIs
	fmt.Println(graph.OrderMap(s.SubNI))
	// Output:
	// <nil>
	// <nil>
	// arc not available in supergraph
	// arc not available in supergraph
	// Subgraph:
	// 0: [1 1]
	// 1: []
	// Mappings:
	// [0 2]
	// map[0:0 2:1 ]
}

func ExampleSubgraph_AddArc_panic() {
	// supergraph:
	//    0
	//   / \\
	//  1    2
	g := graph.AdjacencyList{
		0: {1, 2, 2},
		2: {},
	}
	s := g.InduceList(nil)
	func() {
		defer func() { fmt.Println(recover()) }()
		fmt.Println(s.AddArc(0, -1))
	}()
	func() {
		defer func() { fmt.Println(recover()) }()
		fmt.Println(s.AddArc(3, 0))
	}()
	// Output:
	// AddArc: NI -1 not in supergraph
	// AddArc: NI 3 not in supergraph
}
