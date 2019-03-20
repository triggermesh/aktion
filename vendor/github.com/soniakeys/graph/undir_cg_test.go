// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// undir_cg_test.go -- tests for code in undir_cg.go.
//
// These are tests on the labeled versions of methods.
//
// See also undir_ro_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

func ExampleLabeledUndirected_Bipartite() {
	// 0 1 2  5  6
	//  \|/|     |
	//   3 4     7
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{6, 7}, 0)
	b, _, ok := g.Bipartite()
	fmt.Println("ok   ", ok)
	fmt.Println("(bit) 76543210")
	fmt.Println("Color", b.Color)
	fmt.Println("N0   ", b.N0)
	// Output:
	// ok    true
	// (bit) 76543210
	// Color 10011000
	// N0    5
}

func ExampleLabeledUndirected_BipartiteComponent() {
	// 0 1 2
	//  \|/|
	//   3 4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	c1 := bits.New(g.Order())
	c2 := bits.New(g.Order())
	b, n1, n2, _ := g.BipartiteComponent(0, c1, c2)
	if b {
		fmt.Println("n:  43210")
		fmt.Println("c1:", c1, " n1:", n1)
		fmt.Println("c2:", c2, " n2:", n2)
	}
	// Output:
	// n:  43210
	// c1: 00111  n1: 3
	// c2: 11000  n2: 2
}

func ExampleLabeledUndirected_BipartiteComponent_oddCycle() {
	// 0 1  2
	//  \|/ |
	//   3--4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	c1 := bits.New(g.Order())
	c2 := bits.New(g.Order())
	b, _, _, oc := g.BipartiteComponent(0, c1, c2)
	if !b {
		fmt.Println("odd cycle:", oc)
	}
	// Output:
	// odd cycle: [3 4 2]
}

func ExampleLabeledUndirected_BronKerbosch1() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch1(func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleLabeledUndirected_BKPivotMaxDegree() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleLabeledUndirected_BKPivotMinP() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch2(g.BKPivotMinP, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleLabeledUndirected_BronKerbosch2() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleLabeledUndirected_BronKerbosch3() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch3(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [1 2 5]
	// [4 5]
	// [2 3]
	// [3 4]
	// [0 4]
}

func ExampleLabeledUndirected_ConnectedComponentBits() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	f := g.ConnectedComponentBits()
	fmt.Println("n  ma  543210")
	fmt.Println("-  --  ------")
	for {
		n, ma, b := f()
		if n == 0 {
			break
		}
		fmt.Printf("%d  %2d  %s\n", n, ma, b)
	}
	// Output:
	// n  ma  543210
	// -  --  ------
	// 3   6  011001
	// 2   2  100010
	// 1   0  000100
}

func ExampleLabeledUndirected_ConnectedComponentInts() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	ci, nc := g.ConnectedComponentInts()
	fmt.Println(nc, "components.")
	fmt.Println("node  component int")
	for n, i := range ci {
		fmt.Println(n, "   ", i)
	}
	// Output:
	// 3 components.
	// node  component int
	// 0     1
	// 1     2
	// 2     3
	// 3     1
	// 4     1
	// 5     2
}

func ExampleLabeledUndirected_ConnectedComponentLists() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	f := g.ConnectedComponentLists()
	for {
		l, ma := f()
		if l == nil {
			break
		}
		fmt.Println(l, ma)
	}
	// Output:
	// [0 3 4] 6
	// [1 5] 2
	// [2] 0
}

func ExampleLabeledUndirected_ConnectedComponentReps() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	reps, orders, arcSizes := g.ConnectedComponentReps()
	fmt.Println("reps:    ", reps)
	fmt.Println("orders:  ", orders)
	fmt.Println("arcSizes:", arcSizes)
	// Output:
	// reps:     [0 1 2]
	// orders:   [3 2 1]
	// arcSizes: [6 2 0]
}

func ExampleLabeledUndirected_Degeneracy() {
	//   1   ----5
	//  / \ /   / \
	// 0---2---4  |
	//      \   \ /
	//   3   ----6
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{2, 5}, 0)
	g.AddEdge(graph.Edge{2, 6}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 6}, 0)
	g.AddEdge(graph.Edge{5, 6}, 0)
	fmt.Println(g.Degeneracy())
	// Output:
	// 3
}

func ExampleLabeledUndirected_DegeneracyOrdering() {
	//   1   ----5
	//  / \ /   / \
	// 0---2---4  |
	//      \   \ /
	//   3   ----6
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{2, 5}, 0)
	g.AddEdge(graph.Edge{2, 6}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 6}, 0)
	g.AddEdge(graph.Edge{5, 6}, 0)
	ord, breaks := g.DegeneracyOrdering()
	fmt.Println("Degeneracy:", len(breaks)-1)
	fmt.Println("k-breaks:", breaks)
	fmt.Println("Ordering:", ord)
	for k, x := range breaks {
		fmt.Printf("nodes of %d-core(s): %d\n", k, ord[:x])
	}
	// Output:
	// Degeneracy: 3
	// k-breaks: [7 6 6 4]
	// Ordering: [4 5 6 2 0 1 3]
	// nodes of 0-core(s): [4 5 6 2 0 1 3]
	// nodes of 1-core(s): [4 5 6 2 0 1]
	// nodes of 2-core(s): [4 5 6 2 0 1]
	// nodes of 3-core(s): [4 5 6 2]
}

func ExampleLabeledUndirected_Degree() {
	// 0---1--\
	//      \-/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{1, 1}, 0)
	fmt.Println(g.Degree(0))
	fmt.Println(g.Degree(1))
	// Output:
	// 1
	// 3
}

func ExampleLabeledUndirected_DegreeCentralization() {
	// 0   1
	//  \ /
	//   2
	//  / \
	// 3   4
	var star graph.LabeledUndirected
	star.AddEdge(graph.Edge{2, 0}, 0)
	star.AddEdge(graph.Edge{2, 1}, 0)
	star.AddEdge(graph.Edge{2, 3}, 0)
	star.AddEdge(graph.Edge{2, 4}, 0)
	fmt.Println(star.DegreeCentralization())
	//           3
	//          /
	// 0---1---2
	//          \
	//           4
	var y graph.LabeledUndirected
	y.AddEdge(graph.Edge{0, 1}, 0)
	y.AddEdge(graph.Edge{1, 2}, 0)
	y.AddEdge(graph.Edge{2, 3}, 0)
	y.AddEdge(graph.Edge{2, 4}, 0)
	fmt.Println(y.DegreeCentralization())
	//
	// 0---1---2---3---4
	var line graph.LabeledUndirected
	line.AddEdge(graph.Edge{0, 1}, 0)
	line.AddEdge(graph.Edge{1, 2}, 0)
	line.AddEdge(graph.Edge{2, 3}, 0)
	line.AddEdge(graph.Edge{3, 4}, 0)
	fmt.Println(line.DegreeCentralization())
	//   1---2
	//  /    |
	// 0     |
	//  \    |
	//   3---4
	var circle graph.LabeledUndirected
	circle.AddEdge(graph.Edge{0, 1}, 0)
	circle.AddEdge(graph.Edge{0, 3}, 0)
	circle.AddEdge(graph.Edge{1, 2}, 0)
	circle.AddEdge(graph.Edge{3, 4}, 0)
	circle.AddEdge(graph.Edge{2, 4}, 0)
	fmt.Println(circle.DegreeCentralization())
	// Output:
	// 1
	// 0.5833333333333334
	// 0.16666666666666666
	// 0
}

func ExampleLabeledUndirected_Density() {
	// 0---1
	// |
	// 2---3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	fmt.Println(g.Density())
	// Output:
	// 0.5
}

func ExampleLabeledUndirected_Eulerian() {
	//   0--
	//  /   \
	//  \    \
	//   1----3
	//  / \  /
	//  \ / /
	//   2--
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	fmt.Println(g.Eulerian())
	// Output:
	// 2 3 <nil>
}

func ExampleLabeledUndirected_EulerianCycle() {
	//    0---
	// a /    \ b
	//   \  c  \
	//    1-----3
	// d / \  f/ \ g
	//   \ e\ /  /
	//    ---2---
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'a')
	g.AddEdge(graph.Edge{0, 3}, 'b')
	g.AddEdge(graph.Edge{1, 2}, 'd')
	g.AddEdge(graph.Edge{1, 2}, 'e')
	g.AddEdge(graph.Edge{1, 3}, 'c')
	g.AddEdge(graph.Edge{2, 3}, 'f')
	g.AddEdge(graph.Edge{2, 3}, 'g')
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
	// [{0 -1} {1 97} {3 99} {2 103} {1 100} {2 101} {3 102} {0 98}]
	//
	// 0 --a-- 1 --c-- 3 --g-- 2 --d-- 1 --e-- 2 --f-- 3 --b-- 0
}

func ExampleLabeledUndirected_EulerianCycleD() {
	//    a
	// 0-------1
	//  \    c/|\
	//   \   /d| \
	//  b \  \ | /e
	//     \  \|/
	//      ---2--\
	//          \-/f
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'a')
	g.AddEdge(graph.Edge{0, 2}, 'b')
	g.AddEdge(graph.Edge{1, 2}, 'c')
	g.AddEdge(graph.Edge{1, 2}, 'd')
	g.AddEdge(graph.Edge{1, 2}, 'e')
	g.AddEdge(graph.Edge{2, 2}, 'f') // 6 edges total
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
	// [{0 -1} {1 97} {2 101} {1 99} {2 100} {2 102} {0 98}]
	//
	// 0 --a-- 1 --e-- 2 --c-- 1 --d-- 2 --f-- 2 --b-- 0
}

func ExampleLabeledUndirected_EulerianPath() {
	//    0
	//  a/|\
	//  /b| \
	//  \ | /c
	//   \|/
	//    1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'a')
	g.AddEdge(graph.Edge{0, 1}, 'b')
	g.AddEdge(graph.Edge{0, 1}, 'c')
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
	// [{0 -1} {1 97} {0 99} {1 98}]
	//
	// 0 --a-- 1 --c-- 0 --b-- 1
}

func ExampleLabeledUndirected_EulerianPathD() {
	//    0
	//  a/|\
	//  /b| \
	//  \ | /c
	//   \|/
	//    1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'a')
	g.AddEdge(graph.Edge{0, 1}, 'b')
	g.AddEdge(graph.Edge{0, 1}, 'c')
	c, err := g.EulerianPathD(3, 0)
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
	// [{0 -1} {1 97} {0 99} {1 98}]
	//
	// 0 --a-- 1 --c-- 0 --b-- 1
}

func ExampleLabeledUndirected_EulerianStart() {
	//   0--
	//  /   \
	//  \    \
	//   1----3
	//  / \  /
	//  \ / /
	//   2--
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	fmt.Println(g.EulerianStart())
	// Output:
	// 2
}

func ExampleLabeledUndirected_InduceBits() {
	// undirected graph:
	//     1
	//    /|\\
	//  a/ | \\
	//  / b| c\\d
	// 0   |    2
	//  \  |   /
	//  e\ |  /f
	//    \| /
	//     3-
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{1, 0}, 'a')
	g.AddEdge(graph.Edge{1, 3}, 'b')
	g.AddEdge(graph.Edge{1, 2}, 'c')
	g.AddEdge(graph.Edge{1, 2}, 'd')
	g.AddEdge(graph.Edge{0, 3}, 'e')
	g.AddEdge(graph.Edge{2, 3}, 'f')
	s := g.InduceBits(bits.NewGivens(2, 1, 3))
	fmt.Println("Subgraph:")
	for fr, to := range s.LabeledUndirected.LabeledAdjacencyList {
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
	// 1: [{0, c} {0, d} {2, f} ]
	// 2: [{0, b} {1, f} ]
	// Sub NI -> Super NI
	//   0         1
	//   1         2
	//   2         3
	// Super NI -> Sub NI
	//   1         0
	//   2         1
	//   3         2
}

func ExampleLabeledUndirected_InduceList() {
	// undirected graph:
	//     1
	//    /|\\
	//  a/ | \\
	//  / b| c\\d
	// 0   |    2
	//  \  |   /
	//  e\ |  /f
	//    \| /
	//     3-
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{1, 0}, 'a')
	g.AddEdge(graph.Edge{1, 3}, 'b')
	g.AddEdge(graph.Edge{1, 2}, 'c')
	g.AddEdge(graph.Edge{1, 2}, 'd')
	g.AddEdge(graph.Edge{0, 3}, 'e')
	g.AddEdge(graph.Edge{2, 3}, 'f')
	s := g.InduceList([]graph.NI{2, 1, 2, 3})
	fmt.Println("Subgraph:")
	for fr, to := range s.LabeledUndirected.LabeledAdjacencyList {
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
	// 0: [{1, c} {1, d} {2, f} ]
	// 1: [{2, b} {0, c} {0, d} ]
	// 2: [{1, b} {0, f} ]
	// Sub NI -> Super NI
	//   0         2
	//   1         1
	//   2         3
	// Super NI -> Sub NI
	//   1         1
	//   2         0
	//   3         2
}

func ExampleLabeledUndirected_IsConnected() {
	// undirected graph:
	//   0
	//  / \
	// 1   2
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// connected:  true
}

func ExampleLabeledUndirected_IsConnected_notConnected() {
	// undirected graph:
	//   0   1
	//  / \
	// 2   3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{0, 3}, 0)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// connected:  false
}
func ExampleLabeledUndirected_IsTree() {
	//  0--\
	//  |  |
	//  \--/   1   3
	//        /   / \
	//       2   4---5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 0}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{3, 5}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	fmt.Println(g.IsTree(3))
	// Output:
	// false false
	// true false
	// false false
}

func ExampleLabeledUndirected_Size() {
	//   0--\
	//  / \-/
	// 1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 0}, 0)
	g.AddEdge(graph.Edge{0, 1}, 0)
	fmt.Println("Size:", g.Size())
	fmt.Printf("(Arc size = %d)\n", g.ArcSize())
	// Output:
	// Size: 2
	// (Arc size = 3)
}

func ExampleLabeledUndirectedSubgraph_AddNode() {
	// supergraph:
	//    0
	//   / \
	//  1---2
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, -1)
	g.AddEdge(graph.Edge{0, 2}, -1)
	g.AddEdge(graph.Edge{1, 2}, -1)
	s := g.InduceList(nil)    // construct empty subgraph
	fmt.Println(s.AddNode(2)) // first node added will have NI = 0
	fmt.Println(s.AddNode(1)) // next node added will have NI = 1
	fmt.Println(s.AddNode(1)) // returns existing mapping
	fmt.Println(s.AddNode(2)) // returns existing mapping
	fmt.Println("Subgraph:")  // (it has no arcs)
	for fr, to := range s.LabeledUndirected.LabeledAdjacencyList {
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

func ExampleLabeledUndirectedSubgraph_AddNode_panic() {
	// supergraph:
	//    0
	//   / \
	//  1---2
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, -1)
	g.AddEdge(graph.Edge{0, 2}, -1)
	g.AddEdge(graph.Edge{1, 2}, -1)
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

func ExampleLabeledBipartite_Density() {
	// 0 1 2
	//  \|/|
	//   3 4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	b, _, _ := g.Bipartite()
	fmt.Printf("%.2f\n", b.Density())
	// Output:
	// 0.67
}

func ExampleLabeledBipartite_PermuteBiadjacency() {
	// 3 1 4
	//  \|/|
	//   2 0
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{4, 2}, 0)
	g.AddEdge(graph.Edge{4, 0}, 0)
	b, _, _ := g.Bipartite()
	fmt.Println("Permutation:", b.PermuteBiadjacency())
	fmt.Println("Biadjacency:")
	for fr, to := range g.LabeledAdjacencyList[:b.N0] {
		fmt.Println(fr, to)
	}
	fmt.Println("Full graph:")
	for fr, to := range g.LabeledAdjacencyList {
		fmt.Println(fr, to)
	}
	fmt.Println("(bit) 43210")
	fmt.Println("Color", b.Color)
	fmt.Println("N0   ", b.N0)
	// Output:
	// Permutation: [0 2 1 3 4]
	// Biadjacency:
	// 0 [{4 0}]
	// 1 [{3 0} {2 0} {4 0}]
	// Full graph:
	// 0 [{4 0}]
	// 1 [{3 0} {2 0} {4 0}]
	// 2 [{1 0}]
	// 3 [{1 0}]
	// 4 [{1 0} {0 0}]
	// (bit) 43210
	// Color 11100
	// N0    2
}
