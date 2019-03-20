// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// dir_cg_test.go -- tests for code in dir_cg.go.
//
// These are tests on the labeled versions of methods.
//
// See also dir_ro_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"text/template"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

func ExampleLabeledDirected_Balanced() {
	// 2
	// |
	// v
	// 0----->1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{To: 0}},
		0: {{To: 1}},
	}}
	fmt.Println(g.Balanced())

	// 0<--\
	// |    \
	// v     \
	// 1----->2
	g.LabeledAdjacencyList[1] = []graph.Half{{To: 2}}
	fmt.Println(g.Balanced())
	// Output:
	// false
	// true
}

func ExampleLabeledDirected_Cyclic() {
	//   0
	//  / \
	// 1-->2-->3
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}},
		1: {{To: 2}},
		2: {{To: 3}},
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
	g.LabeledAdjacencyList[3] = []graph.Half{{To: 1}}
	fmt.Println(g.Cyclic())

	// Output:
	// false
	// true 3 {1 0}
}

func ExampleLabeledDirected_DegreeCentralization() {
	// 0<-   ->1
	//    \ /
	//     2
	//    / \
	// 4<-   ->5
	star := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{0, 0}, {1, 0}, {3, 0}, {4, 0}},
		4: {},
	}}
	t, _ := star.Transpose()
	fmt.Println(star.DegreeCentralization(), t.DegreeCentralization())
	//           ->3
	//          /
	// 0<--1<--2
	//          \
	//           ->4
	y := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{1, 0}, {3, 0}, {4, 0}},
		1: {{0, 0}},
		4: {},
	}}
	t, _ = y.Transpose()
	fmt.Println(y.DegreeCentralization(), t.DegreeCentralization())
	//   ->1-->2
	//  /      |
	// 0       |
	// ^       v
	//  \--3<--4
	circle := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 0}},
		1: {{2, 0}},
		2: {{4, 0}},
		4: {{3, 0}},
		3: {{0, 0}},
	}}
	fmt.Println(circle.DegreeCentralization())
	// Output:
	// 1 0.0625
	// 0.6875 0.0625
	// 0
}

func ExampleLabeledDirected_Dominators() {
	//   0   6
	//   |   |
	//   1   |
	//  / \  |
	// 2   3 |
	//  \ / \|
	//   4   5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}},
		1: {{To: 2}, {To: 3}},
		2: {{To: 4}},
		3: {{To: 4}, {To: 5}},
		6: {{To: 5}},
	}}
	d := g.Dominators(0)
	fmt.Println(d.Immediate)
	// Output:
	// [0 0 1 1 1 3 -1]
}

func ExampleLabeledDirected_Doms() {
	//   0   6
	//   |   |
	//   1   |
	//  / \  |
	// 2   3 |
	//  \ / \|
	//   4   5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}},
		1: {{To: 2}, {To: 3}},
		2: {{To: 4}},
		3: {{To: 4}, {To: 5}},
		6: {{To: 5}},
	}}
	// compute postorder with depth-first traversal
	var post []graph.NI
	var vis big.Int
	var f func(graph.NI)
	f = func(n graph.NI) {
		vis.SetBit(&vis, int(n), 1)
		for _, to := range g.LabeledAdjacencyList[n] {
			if vis.Bit(int(to.To)) == 0 {
				f(to.To)
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

func ExampleLabeledDirected_Eulerian() {
	//   /<--------\
	//  /   /<---\  \
	// 0-->1-->\ /  /
	//      \-->2--/
	//         / \
	//        /<--\
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}},
		1: {{To: 2}, {To: 2}},
		2: {{To: 0}, {To: 1}, {To: 2}},
	}}
	fmt.Println(g.Eulerian())
	// Output:
	// -1 -1 <nil>
}

func ExampleLabeledDirected_EulerianCycle() {
	//   /<----------d---\
	//  /      /<--e---\  \
	// 0--a-->1--b-->\ /  /
	//         \--c-->2--/
	//               / \
	//              /   \
	//             /<-f--\
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 'a'}},
		1: {{2, 'b'}, {2, 'c'}},
		2: {{0, 'd'}, {1, 'e'}, {2, 'f'}},
	}}
	c, err := g.EulerianCycle()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	// prettier,
	fmt.Print("\n", c[0].To)
	for _, to := range c[1:] {
		fmt.Printf(" --%c-- %d", to.Label, to.To)
	}
	fmt.Println()
	// Output:
	// [{0 -1} {1 97} {2 98} {1 101} {2 99} {2 102} {0 100}]
	//
	// 0 --a-- 1 --b-- 2 --e-- 1 --c-- 2 --f-- 2 --d-- 0
}

func ExampleLabeledDirected_EulerianCycleD() {
	//   /<----------d---\
	//  /      /<--e---\  \
	// 0--a-->1--b-->\ /  /
	//         \--c-->2--/
	//               / \
	//              /   \
	//             /<-f--\
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 'a'}},
		1: {{2, 'b'}, {2, 'c'}},
		2: {{0, 'd'}, {1, 'e'}, {2, 'f'}},
	}}
	c, err := g.EulerianCycleD(6)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	// prettier,
	fmt.Print("\n", c[0].To)
	for _, to := range c[1:] {
		fmt.Printf(" --%c-- %d", to.Label, to.To)
	}
	fmt.Println()
	// Output:
	// [{0 -1} {1 97} {2 98} {1 101} {2 99} {2 102} {0 100}]
	//
	// 0 --a-- 1 --b-- 2 --e-- 1 --c-- 2 --f-- 2 --d-- 0
}

func ExampleLabeledDirected_EulerianPath() {
	//         /<--e---\
	// 3--a-->1--b-->\ /
	//         \--c-->2--d-->0
	//               / \
	//              /   \
	//             /<-f--\
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		3: {{1, 'a'}},
		1: {{2, 'b'}, {2, 'c'}},
		2: {{0, 'd'}, {1, 'e'}, {2, 'f'}},
	}}
	c, err := g.EulerianPath()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	// prettier,
	fmt.Print("\n", c[0].To)
	for _, to := range c[1:] {
		fmt.Printf(" --%c-- %d", to.Label, to.To)
	}
	fmt.Println()
	// Output:
	// [{3 -1} {1 97} {2 98} {1 101} {2 99} {2 102} {0 100}]
	//
	// 3 --a-- 1 --b-- 2 --e-- 1 --c-- 2 --f-- 2 --d-- 0
}

func ExampleLabeledDirected_EulerianPathD() {
	//         /<--e---\
	// 3--a-->1--b-->\ /
	//         \--c-->2--d-->0
	//               / \
	//              /   \
	//             /<-f--\
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		3: {{1, 'a'}},
		1: {{2, 'b'}, {2, 'c'}},
		2: {{0, 'd'}, {1, 'e'}, {2, 'f'}},
	}}
	c, err := g.EulerianPathD(6, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	// prettier,
	fmt.Print("\n", c[0].To)
	for _, to := range c[1:] {
		fmt.Printf(" --%c-- %d", to.Label, to.To)
	}
	fmt.Println()
	// Output:
	// [{3 -1} {1 97} {2 98} {1 101} {2 99} {2 102} {0 100}]
	//
	// 3 --a-- 1 --b-- 2 --e-- 1 --c-- 2 --f-- 2 --d-- 0
}

func ExampleLabeledDirected_EulerianStart() {
	//      /<---\
	// 3-->1-->\ /
	//      \-->2-->0
	//         / \
	//        /<--\
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		3: {{To: 1}},
		1: {{To: 2}, {To: 2}},
		2: {{To: 0}, {To: 1}, {To: 2}},
	}}
	fmt.Println(g.EulerianStart())
	// Output:
	// 3 <nil>
}

func ExampleLabeledDirected_InDegree() {
	// arcs directed down:
	//  0     2
	//  |
	//  1
	//  |\
	//  | \
	//  3  4<-\
	//     \--/
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}},
		1: {{To: 3}, {To: 4}},
		4: {{To: 4}},
	}}
	fmt.Println("node:    0 1 2 3 4")
	fmt.Println("in-deg:", g.InDegree())
	// Output:
	// node:    0 1 2 3 4
	// in-deg: [0 1 0 1 2]
}

func ExampleLabeledDirected_InduceBits() {
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
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{0, 'a'}, {3, 'b'}, {2, 'c'}, {2, 'd'}},
		0: {{3, 'e'}},
		2: {{3, 'f'}},
		3: {},
	}}
	s := g.InduceBits(bits.NewGivens(2, 1, 3))
	fmt.Println("Subgraph:")
	for fr, to := range s.LabeledDirected.LabeledAdjacencyList {
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

func ExampleLabeledDirected_InduceList() {
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
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{0, 'a'}, {3, 'b'}, {2, 'c'}, {2, 'd'}},
		0: {{3, 'e'}},
		2: {{3, 'f'}},
		3: {},
	}}
	s := g.InduceList([]graph.NI{2, 1, 2, 3})
	fmt.Println("Subgraph:")
	for fr, to := range s.LabeledDirected.LabeledAdjacencyList {
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

func ExampleLabeledDirected_IsTree() {
	// Example graph
	// Arcs point down unless otherwise indicated
	//           1
	//          / \
	//         0   5
	//        /   / \
	//       2   3-->4
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{To: 0}, {To: 5}},
		0: {{To: 2}},
		5: {{To: 3}, {To: 4}},
		3: {{To: 4}},
	}}
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	// Output:
	// true false
	// false false
}

func ExampleLabeledDirected_MaximalNonBranchingPaths() {
	//   a    b     c
	// 0--->1---->2---->3
	//             \ d
	//    -->6      ---->4
	// e /  / f
	//  5<--
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 'a'}},
		1: {{2, 'b'}},
		2: {{3, 'c'}, {4, 'd'}},
		5: {{6, 'e'}},
		6: {{5, 'f'}},
	}}
	g.MaximalNonBranchingPaths(func(p []graph.Half) bool {
		fmt.Println(p)
		// prettier,
		fmt.Print(p[0].To)
		for _, to := range p[1:] {
			fmt.Printf(" --%c-- %d", to.Label, to.To)
		}
		fmt.Print("\n\n")
		return true
	})
	// Output:
	// [{0 -1} {1 97} {2 98}]
	// 0 --a-- 1 --b-- 2
	//
	// [{2 -1} {3 99}]
	// 2 --c-- 3
	//
	// [{2 -1} {4 100}]
	// 2 --d-- 4
	//
	// [{5 -1} {6 101} {5 102}]
	// 5 --e-- 6 --f-- 5
}

func ExampleLabeledDirected_PostDominators() {
	// Example graph here is transpose of that in the Dominators example
	// to show result is the same.
	//   4   5
	//  / \ /|
	// 2   3 |
	//  \ /  |
	//   1   |
	//   |   |
	//   0   6
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		4: {{To: 2}, {To: 3}},
		5: {{To: 3}, {To: 6}},
		2: {{To: 1}},
		3: {{To: 1}},
		1: {{To: 0}},
		6: {},
	}}
	d := g.PostDominators(0)
	fmt.Println(d.Immediate)
	// Output:
	// [0 0 1 1 1 3 -1]
}

func ExampleLabeledDirected_PageRank() {
	//     0<-\
	//    / \ |
	//   /   \|
	//  1---->2<---3
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}},
		1: {{To: 2}},
		2: {{To: 0}},
		3: {{To: 2}},
	}}
	fmt.Printf("%.2f\n", g.PageRank(.85, 20))
	// Output:
	// [1.49 0.78 1.58 0.15]
}

func ExampleLabeledDirected_StronglyConnectedComponents() {
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
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 0}, {To: 5}, {To: 7}},
		5: {{To: 4}, {To: 6}},
		4: {{To: 5}, {To: 2}, {To: 3}},
		7: {{To: 6}},
		6: {{To: 7}, {To: 3}},
		3: {{To: 1}},
		1: {{To: 2}},
		2: {{To: 3}},
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

func ExampleLabeledDirected_Condensation() {
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
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 0}, {To: 5}, {To: 7}},
		5: {{To: 4}, {To: 6}},
		4: {{To: 5}, {To: 2}, {To: 3}},
		7: {{To: 6}},
		6: {{To: 7}, {To: 3}},
		3: {{To: 1}},
		1: {{To: 2}},
		2: {{To: 3}},
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

func ExampleLabeledDirected_Topological() {
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{To: 2}},
		3: {{To: 1}, {To: 2}},
		4: {{To: 3}, {To: 2}},
	}}
	fmt.Println(g.Topological())
	g.LabeledAdjacencyList[2] = []graph.Half{{To: 3}}
	fmt.Println(g.Topological())
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleLabeledDirected_TopologicalKahn() {
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{To: 2}},
		3: {{To: 1}, {To: 2}},
		4: {{To: 3}, {To: 2}},
	}}
	tr, _ := g.UnlabeledTranspose()
	fmt.Println(g.TopologicalKahn(tr))

	g.LabeledAdjacencyList[2] = []graph.Half{{To: 3}}
	tr, _ = g.UnlabeledTranspose()
	fmt.Println(g.TopologicalKahn(tr))
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleLabeledDirected_TopologicalSubgraph() {
	// arcs directected down unless otherwise indicated
	// 0       1<-\
	//  \     / \ /
	//   2   3   4
	//    \ / \
	//     5   6
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 2}},
		1: {{To: 3}, {To: 4}},
		2: {{To: 5}},
		3: {{To: 5}, {To: 6}},
		4: {{To: 1}},
		6: {},
	}}
	fmt.Println(g.TopologicalSubgraph([]graph.NI{0, 3}))
	// Output:
	// [3 6 0 2 5] []
}

func ExampleLabeledDirected_TransitiveClosure() {
	//     1-->2----
	//     ^   |   |
	//  0  |   v   v
	//     ----3-->4-->5<=>7
	//             |   ^
	//             |   |
	//             --->6<=>8
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{To: 2}},
		2: {{To: 3}, {To: 4}},
		3: {{To: 1}, {To: 4}},
		4: {{To: 5}, {To: 6}},
		5: {{To: 7}},
		6: {{To: 5}, {To: 8}},
		7: {{To: 5}},
		8: {{To: 6}},
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

func ExampleLabeledDirectedSubgraph_AddNode() {
	// supergraph:
	//    0
	//   / \
	//  1-->2
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, -1}, {2, -1}},
		1: {{2, -1}},
		2: {},
	}}
	s := g.InduceList(nil)    // construct empty subgraph
	fmt.Println(s.AddNode(2)) // first node added will have NI = 0
	fmt.Println(s.AddNode(1)) // next node added will have NI = 1
	fmt.Println(s.AddNode(1)) // returns existing mapping
	fmt.Println(s.AddNode(2)) // returns existing mapping
	fmt.Println("Subgraph:")  // (it has no arcs)
	for fr, to := range s.LabeledDirected.LabeledAdjacencyList {
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

func ExampleLabeledDirectedSubgraph_AddNode_panic() {
	// supergraph:
	//    0
	//   / \
	//  1-->2
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, -1}, {2, -1}},
		1: {{2, -1}},
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

func ExampleLabeledDirectedSubgraph_AddArc() {
	// supergraph:
	//     0
	//    / \\
	//  x/  y\\z
	//  1      2
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 'x'}, {2, 'y'}, {2, 'z'}},
		2: {},
	}}
	s := g.InduceList(nil)                       // construct empty subgraph
	fmt.Println(s.AddArc(1, graph.Half{0, 'x'})) // no, not that direction
	fmt.Println(s.AddArc(0, graph.Half{1, 'y'})) // no, not with that label
	fmt.Println(s.AddArc(0, graph.Half{2, 'y'})) // okay
	fmt.Println("Subgraph:")
	for fr, to := range s.LabeledDirected.LabeledAdjacencyList {
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

func ExampleLabeledDirectedSubgraph_AddArc_panic() {
	// supergraph:
	//    0
	//   / \\
	//  1    2
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, -1}, {2, -1}, {2, -1}},
		2: {},
	}}
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
