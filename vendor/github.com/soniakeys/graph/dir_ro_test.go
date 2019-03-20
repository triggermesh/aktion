// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// dir_ro_test.go -- tests on dir_RO.go
//
// These are tests on the code-generated unlabeled versions of methods.
//
// If testing prompts changes in the tested method, be sure to edit
// dir_cg.go, go generate to generate the dir_RO.go, and then retest.
// Do not edit dir_RO.go
//
// See also dir_cg_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"text/template"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

func ExampleDirected_Balanced() {
	// 2
	// |
	// v
	// 0----->1
	g := graph.Directed{graph.AdjacencyList{
		2: {0},
		0: {1},
	}}
	fmt.Println(g.Balanced())

	// 0<--\
	// |    \
	// v     \
	// 1----->2
	g.AdjacencyList[1] = []graph.NI{2}
	fmt.Println(g.Balanced())
	// Output:
	// false
	// true
}

func ExampleDirected_Cyclic() {
	//   0
	//  / \
	// 1-->2-->3
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {},
	}}
	cyclic, _, _ := g.Cyclic()
	fmt.Println(cyclic)

	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g.AdjacencyList[3] = []graph.NI{1}
	fmt.Println(g.Cyclic())

	// Output:
	// false
	// true 3 1
}

func ExampleDirected_DegreeCentralization() {
	// 0<-   ->1
	//    \ /
	//     2
	//    / \
	// 4<-   ->5
	star := graph.Directed{graph.AdjacencyList{
		2: {0, 1, 3, 4},
		4: {},
	}}
	t, _ := star.Transpose()
	fmt.Println(star.DegreeCentralization(), t.DegreeCentralization())
	//           ->3
	//          /
	// 0<--1<--2
	//          \
	//           ->4
	y := graph.Directed{graph.AdjacencyList{
		2: {1, 3, 4},
		1: {0},
		4: {},
	}}
	t, _ = y.Transpose()
	fmt.Println(y.DegreeCentralization(), t.DegreeCentralization())
	//   ->1-->2
	//  /      |
	// 0       |
	// ^       v
	//  \--3<--4
	circle := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2},
		2: {4},
		4: {3},
		3: {0},
	}}
	fmt.Println(circle.DegreeCentralization())
	// Output:
	// 1 0.0625
	// 0.6875 0.0625
	// 0
}

func ExampleDirected_Dominators() {
	//   0   6
	//   |   |
	//   1   |
	//  / \  |
	// 2   3 |
	//  \ / \|
	//   4   5
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 3},
		2: {4},
		3: {4, 5},
		6: {5},
	}}
	d := g.Dominators(0)
	fmt.Println(d.Immediate)
	// Output:
	// [0 0 1 1 1 3 -1]
}

func ExampleDirected_Doms() {
	//   0   6
	//   |   |
	//   1   |
	//  / \  |
	// 2   3 |
	//  \ / \|
	//   4   5
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 3},
		2: {4},
		3: {4, 5},
		6: {5},
	}}
	// compute postorder with depth-first traversal
	var post []graph.NI
	var vis big.Int
	var f func(graph.NI)
	f = func(n graph.NI) {
		vis.SetBit(&vis, int(n), 1)
		for _, to := range g.AdjacencyList[n] {
			if vis.Bit(int(to)) == 0 {
				f(to)
			}
		}
		post = append(post, n)
	}
	f(0)
	fmt.Println("post:", post)
	tr, _ := g.Transpose()
	d := g.Doms(tr, post)
	fmt.Println("doms:", d.Immediate)
	// Output:
	// post: [4 2 5 3 1 0]
	// doms: [0 0 1 1 1 3 -1]
}

func ExampleDirected_Eulerian() {
	//   /<--------\
	//  /   /<---\  \
	// 0-->1-->\ /  /
	//      \-->2--/
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println(g.Eulerian())
	// Output:
	// -1 -1 <nil>
}

func ExampleDirected_EulerianCycle() {
	//   /<--------\
	//  /   /<---\  \
	// 0-->1-->\ /  /
	//      \-->2--/
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println(g.EulerianCycle())
	// Output:
	// [0 1 2 1 2 2 0] <nil>
}

func ExampleDirected_EulerianCycleD() {
	//   /<--------\
	//  /   /<---\  \
	// 0-->1-->\ /  /
	//      \-->2--/
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println("ma =", g.ArcSize())
	fmt.Println(g.EulerianCycleD(6))
	// Output:
	// ma = 6
	// [0 1 2 1 2 2 0] <nil>
}

func TestEulerianCycle(t *testing.T) {
	same := func(a, b []graph.NI) bool {
		if len(a) != len(b) {
			return false
		}
		for i, x := range a {
			if b[i] != x {
				return false
			}
		}
		return true
	}
	var msg string
	for _, tc := range []struct {
		g     graph.AdjacencyList
		cycle []graph.NI
		ok    bool
	}{
		{nil, nil, true},
		{graph.AdjacencyList{nil}, []graph.NI{0}, true},    // 1 node, 0 arcs
		{graph.AdjacencyList{{0}}, []graph.NI{0, 0}, true}, // loop
		{graph.AdjacencyList{nil, nil}, nil, false},        // not connected
		{graph.AdjacencyList{{1}, nil}, nil, false},        // not balanced
		{graph.AdjacencyList{nil, {0}}, nil, false},        // not balanced
	} {
		got, err := graph.Directed{tc.g}.EulerianCycle()
		switch {
		case err != nil:
			if !tc.ok {
				continue
			}
			msg = "g.EulerianCycle() returned error" + err.Error()
		case !tc.ok:
			msg = fmt.Sprintf("g.EulerianCycle() = %v, want error", got)
		case !same(got, tc.cycle):
			msg = fmt.Sprintf("g.EulerianCycle() = %v, want %v", got, tc.cycle)
		default:
			continue
		}
		t.Log("g:", tc.g)
		t.Fatal(msg)
	}
}

func ExampleDirected_EulerianPath() {
	//      /<---\
	// 3-->1-->\ /
	//      \-->2-->0
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		3: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println(g.EulerianPath())
	// Output:
	// [3 1 2 1 2 2 0] <nil>
}

func ExampleDirected_EulerianPathD() {
	//      /<---\
	// 3-->1-->\ /
	//      \-->2-->0
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		3: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println("ma =", g.ArcSize())
	fmt.Print("start: ")
	fmt.Println(g.EulerianStart())
	fmt.Println(g.EulerianPathD(6, 3))
	// Output:
	// ma = 6
	// start: 3 <nil>
	// [3 1 2 1 2 2 0] <nil>
}

func ExampleDirected_EulerianStart() {
	//      /<---\
	// 3-->1-->\ /
	//      \-->2-->0
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		3: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println(g.EulerianStart())
	// Output:
	// 3 <nil>
}

func TestEulerianPath(t *testing.T) {
	same := func(a, b []graph.NI) bool {
		if len(a) != len(b) {
			return false
		}
		for i, x := range a {
			if b[i] != x {
				return false
			}
		}
		return true
	}
	var msg string
	for _, tc := range []struct {
		g    graph.AdjacencyList
		path []graph.NI
		ok   bool
	}{
		{nil, nil, true},
		{graph.AdjacencyList{nil}, []graph.NI{0}, true},    // 1 node, 0 arcs
		{graph.AdjacencyList{{0}}, []graph.NI{0, 0}, true}, // loop
		{graph.AdjacencyList{{1}, nil}, []graph.NI{0, 1}, true},
		{graph.AdjacencyList{nil, {0}}, []graph.NI{1, 0}, true},
		{graph.AdjacencyList{nil, nil}, nil, false},         // not connected
		{graph.AdjacencyList{{1}, nil, {1}}, nil, false},    // two starts
		{graph.AdjacencyList{nil, nil, {0, 1}}, nil, false}, // two ends
	} {
		got, err := graph.Directed{tc.g}.EulerianPath()
		switch {
		case err != nil:
			if !tc.ok {
				continue
			}
			msg = "g.EulerianPath() returned error" + err.Error()
		case !tc.ok:
			msg = fmt.Sprintf("g.EulerianPath() = %v, want error", got)
		case !same(got, tc.path):
			msg = fmt.Sprintf("g.EulerianPath() = %v, want %v", got, tc.path)
		default:
			continue
		}
		t.Log("g:", tc.g)
		t.Fatal(msg)
	}
}

func ExampleDirected_InDegree() {
	// arcs directed down:
	//  0     2
	//  |
	//  1
	//  |\
	//  | \
	//  3  4<-\
	//     \--/
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {3, 4},
		4: {4},
	}}
	fmt.Println("node:    0 1 2 3 4")
	fmt.Println("in-deg:", g.InDegree())
	// Output:
	// node:    0 1 2 3 4
	// in-deg: [0 1 0 1 2]
}

func ExampleDirected_InduceBits() {
	// arcs directed down:
	//   1
	//  /|\\
	// 0 |  2
	//  \| /
	//   3-
	g := graph.Directed{graph.AdjacencyList{
		1: {0, 3, 2, 2},
		0: {3},
		2: {3},
		3: {},
	}}
	s := g.InduceBits(bits.NewGivens(2, 1, 3))
	fmt.Println("Subgraph:")
	for fr, to := range s.Directed.AdjacencyList {
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

func ExampleDirected_InduceList() {
	// arcs directed down:
	//   1
	//  /|\\
	// 0 |  2
	//  \| /
	//   3-
	g := graph.Directed{graph.AdjacencyList{
		1: {0, 3, 2, 2},
		0: {3},
		2: {3},
		3: {},
	}}
	s := g.InduceList([]graph.NI{2, 1, 2, 3})
	fmt.Println("Subgraph:")
	for fr, to := range s.Directed.AdjacencyList {
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

func ExampleDirected_IsTree() {
	// Example graph
	// Arcs point down unless otherwise indicated
	//           1
	//          / \
	//         0   5
	//        /   / \
	//       2   3-->4
	g := graph.Directed{graph.AdjacencyList{
		1: {0, 5},
		0: {2},
		5: {3, 4},
		3: {4},
	}}
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	// Output:
	// true false
	// false false
}

func ExampleDirected_MaximalNonBranchingPaths() {
	// 0-->1-->2-->3
	//          \
	//   -->6    ->4
	//  /  /
	// 5<--
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2},
		2: {3, 4},
		5: {6},
		6: {5},
	}}
	g.MaximalNonBranchingPaths(func(p []graph.NI) bool {
		fmt.Println(p)
		return true
	})
	// Output:
	// [0 1 2]
	// [2 3]
	// [2 4]
	// [5 6 5]
}

func ExampleDirected_PostDominators() {
	// Example graph here is transpose of that in the Dominators example
	// to show result is the same.
	//   4   5
	//  / \ /|
	// 2   3 |
	//  \ /  |
	//   1   |
	//   |   |
	//   0   6
	g := graph.Directed{graph.AdjacencyList{
		4: {2, 3},
		5: {3, 6},
		2: {1},
		3: {1},
		1: {0},
		6: {},
	}}
	d := g.PostDominators(0)
	fmt.Println(d.Immediate)
	// Output:
	// [0 0 1 1 1 3 -1]
}

func ExampleDirected_PageRank() {
	//     0<-\
	//    / \ |
	//   /   \|
	//  1---->2<---3
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {0},
		3: {2},
	}}
	fmt.Printf("%.2f\n", g.PageRank(.85, 20))
	// Output:
	// [1.49 0.78 1.58 0.15]
}

func ExampleDirected_StronglyConnectedComponents() {
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
	g.StronglyConnectedComponents(func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [3 1 2]
	// [7 6]
	// [4 5]
	// [0]
}

func ExampleDirected_Condensation() {
	// input:          condensation:
	// /---0---\      <->  /---3
	// |   |\--/           |   |
	// |   v               |   v
	// |   5<=>4---\  <->  |   2--\
	// |   |   |   |       |   |  |
	// v   v   |   |       |   v  |
	// 7<=>6   |   |  <->  \-->1  |
	//     |   v   v           |  v
	//     \-->3<--2  <->      \->0
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
	scc, cd := g.Condensation()
	fmt.Println(len(scc), "components:")
	for cn, c := range scc {
		fmt.Println(cn, c)
	}
	fmt.Println("condensation:")
	for cn, to := range cd {
		fmt.Println(cn, to)
	}
	// Output:
	// 4 components:
	// 0 [3 1 2]
	// 1 [7 6]
	// 2 [4 5]
	// 3 [0]
	// condensation:
	// 0 []
	// 1 [0]
	// 2 [0 1]
	// 3 [2 1]
}

func ExampleDirected_Topological() {
	g := graph.Directed{graph.AdjacencyList{
		1: {2},
		3: {1, 2},
		4: {3, 2},
	}}
	fmt.Println(g.Topological())
	g.AdjacencyList[2] = []graph.NI{3}
	fmt.Println(g.Topological())
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleDirected_TopologicalKahn() {
	g := graph.Directed{graph.AdjacencyList{
		1: {2},
		3: {1, 2},
		4: {3, 2},
	}}
	tr, _ := g.Transpose()
	fmt.Println(g.TopologicalKahn(tr))

	g.AdjacencyList[2] = []graph.NI{3}
	tr, _ = g.Transpose()
	fmt.Println(g.TopologicalKahn(tr))
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleDirected_TopologicalSubgraph() {
	// arcs directected down unless otherwise indicated
	// 0       1<-\
	//  \     / \ /
	//   2   3   4
	//    \ / \
	//     5   6
	g := graph.Directed{graph.AdjacencyList{
		0: {2},
		1: {3, 4},
		2: {5},
		3: {5, 6},
		4: {1},
		6: {},
	}}
	fmt.Println(g.TopologicalSubgraph([]graph.NI{0, 3}))
	// Output:
	// [3 6 0 2 5] []
}

func ExampleDirected_TransitiveClosure() {
	//     1-->2----
	//     ^   |   |
	//  0  |   v   v
	//     ----3-->4-->5<=>7
	//             |   ^
	//             |   |
	//             --->6<=>8
	g := graph.Directed{graph.AdjacencyList{
		1: {2},
		2: {3, 4},
		3: {1, 4},
		4: {5, 6},
		5: {7},
		6: {5, 8},
		7: {5},
		8: {6},
	}}
	t := g.TransitiveClosure()
	fmt.Println(".  0 1 2 3 4 5 6 7 8")
	fmt.Println("   -----------------")
	for fr, tn := range t {
		fmt.Print(fr, ":")
		for to := range t {
			fmt.Print(" ", tn.Bit(to))
		}
		fmt.Println()
	}
	// Output:
	// .  0 1 2 3 4 5 6 7 8
	//    -----------------
	// 0: 0 0 0 0 0 0 0 0 0
	// 1: 0 1 1 1 1 1 1 1 1
	// 2: 0 1 1 1 1 1 1 1 1
	// 3: 0 1 1 1 1 1 1 1 1
	// 4: 0 0 0 0 0 1 1 1 1
	// 5: 0 0 0 0 0 1 0 1 0
	// 6: 0 0 0 0 0 1 1 1 1
	// 7: 0 0 0 0 0 1 0 1 0
	// 8: 0 0 0 0 0 1 1 1 1
}

func ExampleDirectedSubgraph_AddNode() {
	// supergraph:
	//    0
	//   / \
	//  1-->2
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {},
	}}
	s := g.InduceList(nil)    // construct empty subgraph
	fmt.Println(s.AddNode(2)) // first node added will have NI = 0
	fmt.Println(s.AddNode(1)) // next node added will have NI = 1
	fmt.Println(s.AddNode(1)) // returns existing mapping
	fmt.Println(s.AddNode(2)) // returns existing mapping
	fmt.Println("Subgraph:")  // (it has no arcs)
	for fr, to := range s.Directed.AdjacencyList {
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

func ExampleDirectedSubgraph_AddNode_panic() {
	// supergraph:
	//    0
	//   / \
	//  1-->2
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {},
	}}
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

func ExampleDirectedSubgraph_AddArc() {
	// supergraph:
	//    0
	//   / \\
	//  1    2
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2, 2},
		2: {},
	}}
	s := g.InduceList(nil)      // construct empty subgraph
	fmt.Println(s.AddArc(0, 2)) // okay
	fmt.Println(s.AddArc(0, 2)) // adding one parallel arc okay
	fmt.Println(s.AddArc(0, 2)) // adding another not okay
	fmt.Println(s.AddArc(1, 2)) // arc not in supergraph at all
	fmt.Println("Subgraph:")
	for fr, to := range s.Directed.AdjacencyList {
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

func ExampleDirectedSubgraph_AddArc_panic() {
	// supergraph:
	//    0
	//   / \\
	//  1    2
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2, 2},
		2: {},
	}}
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
