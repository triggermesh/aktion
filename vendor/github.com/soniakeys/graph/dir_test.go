// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleDirected_Cycles() {
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
	g.Cycles(func(c []graph.NI) bool {
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

func TestCycles(t *testing.T) {
	// Tushar Roy https://www.youtube.com/watch?v=johyrWospv0
	//   /---->9
	//  8<----/^
	//  ^      |
	//  |      |
	//  1----->2--->7
	//  ^     ^|^
	//   \   / | \
	//    \ / /   5
	//     3<-    ^
	//     |\     |
	//     | ---->4
	//     v      ^
	//     6-----/
	g := graph.Directed{graph.AdjacencyList{
		8: {9},
		9: {8},
		1: {8, 2, 5},
		2: {9, 7, 3},
		3: {1, 2, 4, 6},
		6: {4},
		4: {5},
		5: {2},
	}}
	want := [][]graph.NI{
		{1, 2, 3},
		{1, 5, 2, 3},
		{2, 3},
		{2, 3, 4, 5},
		{2, 3, 6, 4, 5},
		{8, 9},
	}
	i := 0
	g.Cycles(func(c []graph.NI) bool {
		if !reflect.DeepEqual(c, want[i]) {
			t.Fatalf("cycle %d, got %d, want %d", i, c, want[i])
		}
		i++
		return true
	})
	if i != len(want) {
		t.Fatalf("only %d cycles.  want %d.", i, len(want))
	}
}

func ExampleDirected_DAGMaxLenPath() {
	// arcs directed right:
	//      /---\
	//  3--4  1--0--2
	//   \------/
	g := graph.Directed{graph.AdjacencyList{
		3: {0, 4},
		4: {0},
		1: {0},
		0: {2},
	}}
	o, _ := g.Topological()
	fmt.Println(o)
	fmt.Println(g.DAGMaxLenPath(o))
	// Output:
	// [3 4 1 0 2]
	// [3 4 0 2]
}

func ExampleDirected_FromList() {
	//    4   3
	//   / \
	//  2   1
	//       \
	//        0
	g := graph.Directed{graph.AdjacencyList{
		4: {2, 1},
		1: {0},
	}}
	f, sf := g.FromList()
	fmt.Println("simple forest:", sf)
	fmt.Println("N  From")
	for n, e := range f.Paths {
		fmt.Printf("%d %4d\n", n, e.From)
	}
	// Output:
	// simple forest: true
	// N  From
	// 0    1
	// 1    4
	// 2    4
	// 3   -1
	// 4   -1
}

func ExampleDirected_SpanTree() {
	//    0   5
	//   / \
	//  1   2
	//     / \
	//    3-->4
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		2: {3, 4},
		3: {4},
		5: {},
	}}
	var f graph.FromList
	ns, simple := g.SpanTree(0, &f)
	fmt.Println("nodes spanned:", ns)
	fmt.Println("simple tree:", simple)
	fmt.Println("N  From  Len")
	for n, e := range f.Paths {
		fmt.Printf("%d %4d %4d\n", n, e.From, e.Len)
	}
	// Output:
	// nodes spanned: 5
	// simple tree: false
	// N  From  Len
	// 0   -1    1
	// 1    0    2
	// 2    0    2
	// 3    2    3
	// 4    2    3
	// 5   -1    0
}

func ExampleDirected_FromList_nonTree() {
	//    0
	//   / \
	//  1   2
	//   \ /
	//    3
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		1: {3},
		2: {3},
		3: {},
	}}
	f, sf := g.FromList()
	fmt.Println("simple forest:", sf)
	fmt.Println("N  From")
	for n, e := range f.Paths {
		fmt.Printf("%d %4d\n", n, e.From)
	}
	// Output:
	// simple forest: false
	// N  From
	// 0   -1
	// 1    0
	// 2    0
	// 3    1
}

func ExampleDirected_FromList_multigraphTree() {
	//    0
	//   / \\
	//  1   2
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2, 2},
		2: {},
	}}
	f, sf := g.FromList()
	fmt.Println("simple forest:", sf)
	fmt.Println("N  From")
	for n, e := range f.Paths {
		fmt.Printf("%d %4d\n", n, e.From)
	}
	// Output:
	// simple forest: false
	// N  From
	// 0   -1
	// 1    0
	// 2    0
}

func ExampleDirected_FromList_rootLoop() {
	//     /-\
	//    0--/
	//   / \
	//  1   2
	g := graph.Directed{graph.AdjacencyList{
		0: {0, 1, 2},
		2: {},
	}}
	f, sf := g.FromList()
	fmt.Println("simple forest:", sf)
	fmt.Println("N  From")
	for n, e := range f.Paths {
		fmt.Printf("%d %4d\n", n, e.From)
	}
	// Output:
	// simple forest: false
	// N  From
	// 0   -1
	// 1    0
	// 2    0
}

func ExampleDirected_Transpose() {
	g := graph.Directed{graph.AdjacencyList{
		2: {0, 1},
	}}
	t, m := g.Transpose()
	for n, nbs := range t.AdjacencyList {
		fmt.Printf("%d: %v\n", n, nbs)
	}
	fmt.Println(m)
	// Output:
	// 0: [2]
	// 1: [2]
	// 2: []
	// 2
}

func ExampleDirected_Undirected() {
	// arcs directed down:
	//    0
	//   / \
	//  1   2
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		1: {},
		2: {},
	}}
	u := g.Undirected()
	for fr, to := range u.AdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [1 2]
	// 1 [0]
	// 2 [0]
}

func ExampleDirected_Undirected_loopMultigraph() {
	//  0--\   /->1--\
	//  |  |   |  ^  |
	//  \--/   |  |  |
	//         \--2<-/
	g := graph.Directed{graph.AdjacencyList{
		0: {0},
		1: {2},
		2: {1, 1},
	}}
	u := g.Undirected()
	for fr, to := range u.AdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [0]
	// 1 [2 2]
	// 2 [1 1]
}

func ExampleDominanceFrontiers_Closure() {
	//     0
	//     |
	//     1
	//     |
	// --->2
	// |  / \
	// --3   4
	//      / \
	//     5   6
	//      \ /
	//       7
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2},
		2: {3, 4},
		3: {2},
		4: {5, 6},
		5: {7},
		6: {7},
		7: {},
	}}
	f := g.Dominators(0).Frontiers()
	type ns map[graph.NI]struct{}
	fmt.Println(f.Closure(ns{
		0: struct{}{},
		1: struct{}{},
		3: struct{}{},
	}))
	// Output:
	// map[2:{}]
}

func ExampleDominanceFrontiers_Frontier() {
	//     0
	//     |
	//     1
	//     |
	// --->2
	// |  / \
	// --3   4
	//      / \
	//     5   6
	//      \ /
	//       7
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2},
		2: {3, 4},
		3: {2},
		4: {5, 6},
		5: {7},
		6: {7},
		7: {},
	}}
	f := g.Dominators(0).Frontiers()
	type ns map[graph.NI]struct{}
	fmt.Println(f.Frontier(ns{
		0: struct{}{},
		1: struct{}{},
		3: struct{}{},
	}))
	// Output:
	// map[2:{}]
}

func ExampleDominanceFrontiers_Frontier_labeled() {
	//     0
	//     |
	//     1
	//     |
	// --->2
	// |  / \
	// --3   4
	//      / \
	//     5   6
	//      \ /
	//       7
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}},
		1: {{To: 2}},
		2: {{To: 3}, {To: 4}},
		3: {{To: 2}},
		4: {{To: 5}, {To: 6}},
		5: {{To: 7}},
		6: {{To: 7}},
		7: {},
	}}
	f := g.Dominators(0).Frontiers()
	type ns map[graph.NI]struct{}
	fmt.Println(f.Frontier(ns{
		0: struct{}{},
		1: struct{}{},
		3: struct{}{},
	}))
	// Output:
	// map[2:{}]
}

func ExampleDominators_Frontiers() {
	//   0
	//   |
	//   1
	//  / \
	// 2   3
	//  \ / \
	//   4   5   6
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 3},
		2: {4},
		3: {4, 5},
		6: {},
	}}
	for n, f := range g.Dominators(0).Frontiers() {
		fmt.Print(n, ":")
		if f == nil {
			fmt.Println(" nil")
			continue
		}
		for n := range f {
			fmt.Print(" ", n)
		}
		fmt.Println()
	}
	// Output:
	// 0:
	// 1:
	// 2: 4
	// 3: 4
	// 4:
	// 5:
	// 6: nil
}

func ExampleDominators_Set() {
	//   0
	//   |
	//   1
	//  / \
	// 2   3
	//  \ / \
	//   4   5   6
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 3},
		2: {4},
		3: {4, 5},
		6: {},
	}}
	d := g.Dominators(0)
	for n := range g.AdjacencyList {
		fmt.Println(n, d.Set(graph.NI(n)))
	}
	// Output:
	// 0 [0]
	// 1 [1 0]
	// 2 [2 1 0]
	// 3 [3 1 0]
	// 4 [4 1 0]
	// 5 [5 3 1 0]
	// 6 []
}

// ------- Labeled examples -------

func ExampleLabeledDirected_NegativeCycles() {
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
	w := func(l graph.LI) float64 { return float64(l) }
	fmt.Println("dist  cycle")
	g.NegativeCycles(w, func(c []graph.Half) bool {
		d := 0.
		for _, h := range c {
			d += w(h.Label)
		}
		fmt.Printf("%3.0f   %d\n", d, c)
		return true
	})
	// Output:
	// dist  cycle
	//  -2   [{2 -2} {6 -1} {5 1}]
	//  -4   [{1 -4} {5 1} {4 -1}]
	//  -1   [{1 -2} {5 1}]
	//  -3   [{1 -4} {2 2} {6 -1} {5 1} {4 -1}]
	//  -1   [{1 -4} {2 2} {3 2} {6 -1} {5 1} {4 -1}]
}

func ExampleLabeledDirected_Cycles() {
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
	g.Cycles(func(c []graph.Half) bool {
		fmt.Println(c)
		/*
			d := 0.
			for _, h := range c {
				d += float64(h.Label)
			}
			if d < 0 {
				fmt.Println("^^^ neg:", d)
			}*/
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

func TestNegativeCycles(t *testing.T) {
	// example from the paper, by table 1.  (arrows on fig 3 seem wrong)
	arcs := []struct {
		fr, to graph.NI
		wt     float64
	}{
		1:  {1, 2, 53},
		2:  {1, 6, 180},
		3:  {1, 11, 353},
		4:  {3, 2, -10},
		5:  {4, 2, 70},
		6:  {2, 5, -30},
		7:  {6, 2, -20},
		8:  {3, 4, -104},
		9:  {8, 3, 183},
		10: {5, 4, 10},
		11: {7, 4, -40},
		12: {4, 8, 172},
		13: {5, 6, 98},
		14: {5, 7, -20},
		15: {6, 7, 133},
		16: {6, 11, 175},
		17: {6, 12, 190},
		18: {6, 20, 162},
		19: {8, 7, 10},
		20: {9, 7, 159},
		21: {7, 10, 30},
		22: {7, 20, 120},
		23: {9, 8, -60},
		24: {8, 16, 338},
		25: {10, 9, 94},
		26: {15, 9, -90},
		27: {9, 16, 201},
		28: {12, 10, 134},
		29: {13, 10, 150},
		30: {10, 14, -50},
		31: {10, 15, 169},
		32: {20, 10, 40},
		33: {11, 12, 60},
		34: {11, 19, 234},
		35: {12, 13, 104},
		36: {12, 20, 104},
		37: {14, 13, 47},
		38: {18, 19, 86},
		39: {14, 15, 30},
		40: {19, 14, 55},
		41: {15, 16, 115},
		42: {15, 18, 92},
		43: {16, 17, 125},
		44: {18, 16, 137},
		45: {17, 18, 98},
		46: {19, 18, 100},
	}
	a := make(graph.LabeledAdjacencyList, 21)
	for e := 1; e <= 46; e++ {
		arc := arcs[e]
		a[arc.fr] = append(a[arc.fr], graph.Half{arc.to, graph.LI(e)})
	}
	c, _ := a.Copy()
	w := func(l graph.LI) float64 { return arcs[l].wt }
	// results from paragraph of section 4, p. 284.
	want := [][]graph.LI{
		{21, 30, 39, 26, 23, 19},             // C2
		{6, 14, 11, 5},                       // C1
		{21, 30, 39, 26, 23, 9, 8, 5, 6, 14}, // C3
		{21, 30, 39, 26, 23, 9, 4, 6, 14},    // C4 (corrected)
	}
	i := 0
	graph.LabeledDirected{a}.NegativeCycles(w, func(c []graph.Half) bool {
		got := make([]graph.LI, len(c))
		for i, e := range c {
			got[i] = e.Label
		}
		if !reflect.DeepEqual(got, want[i]) {
			t.Fatal(got, want[i])
		}
		i++
		return true
	})
	if !reflect.DeepEqual(a[1:], c[1:]) {
		t.Fatal("graph altered")
	}
}

func TestC6(t *testing.T) {
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 1}, {2, 1}, {3, 1}, {4, -1}, {5, -1}},
		1: {{0, 1}, {2, 1}, {3, 1}, {4, 1}, {5, -1}},
		2: {{0, 1}, {1, 1}, {3, 1}, {4, 1}, {5, 1}},
		3: {{0, 1}, {1, 1}, {2, 1}, {4, 1}, {5, 1}},
		4: {{0, -1}, {1, 1}, {2, 1}, {3, 1}, {5, 1}},
		5: {{0, -1}, {1, -1}, {2, 1}, {3, 1}, {4, 1}},
	}}
	w := func(l graph.LI) float64 { return float64(l) }
	want := [][]graph.Half{
		{{4, -1}, {0, -1}},
		{{5, -1}, {0, -1}},
		{{1, -1}, {5, -1}},
		{{1, 1}, {5, -1}, {0, -1}},
		{{5, -1}, {1, -1}, {4, 1}, {0, -1}},
		{{5, -1}, {4, 1}, {0, -1}},
		{{5, -1}, {1, -1}, {0, 1}},
		{{5, -1}, {1, -1}, {2, 1}, {4, 1}, {0, -1}},
		{{5, -1}, {1, -1}, {3, 1}, {4, 1}, {0, -1}},
		{{4, -1}, {5, 1}, {0, -1}},
		{{4, -1}, {1, 1}, {5, -1}, {0, -1}},
		{{4, -1}, {3, 1}, {1, 1}, {5, -1}, {0, -1}},
		{{4, -1}, {2, 1}, {1, 1}, {5, -1}, {0, -1}},
	}
	c, _ := g.Copy()
	i := 0
	g.NegativeCycles(w, func(c []graph.Half) bool {
		if !reflect.DeepEqual(c, want[i]) {
			t.Fatal("i: ", i, " got: ", c, " want: ", want[i])
		}
		i++
		return true
	})
	if i != len(want) {
		t.Fatal("only ", i)
	}
	if !reflect.DeepEqual(g, c) {
		for fr, to := range g.LabeledAdjacencyList {
			log.Println(fr, to)
		}
		t.Fatal("graph altered")
	}

	// Test that graph is unaltered even after early return
	for i := range want {
		j := 0
		g.NegativeCycles(w, func(c []graph.Half) bool {
			j++
			return j <= i
		})
		if !reflect.DeepEqual(g, c) {
			t.Fatal("graph altered. i = ", i)
		}
	}
}

func ExampleLabeledDirected_DAGMaxLenPath() {
	// arcs directed right:
	//            (M)
	//    (W)  /---------\
	//  3-----4   1-------0-----2
	//         \    (S)  /  (P)
	//          \       /
	//           \-----/ (Q)
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		3: {{To: 0, Label: 'Q'}, {4, 'W'}},
		4: {{0, 'M'}},
		1: {{0, 'S'}},
		0: {{2, 'P'}},
	}}
	o, _ := g.Topological()
	fmt.Println("ordering:", o)
	n, p := g.DAGMaxLenPath(o)
	fmt.Printf("path from %d:", n)
	for _, e := range p {
		fmt.Printf(" {%d, '%c'}", e.To, e.Label)
	}
	fmt.Println()
	fmt.Print("label path: ")
	for _, h := range p {
		fmt.Print(string(h.Label))
	}
	fmt.Println()
	// Output:
	// ordering: [3 4 1 0 2]
	// path from 3: {4, 'W'} {0, 'M'} {2, 'P'}
	// label path: WMP
}

func ExampleLabeledDirected_FromList() {
	//      0
	// 'A' / \ 'B'
	//    1   2
	//         \ 'C'
	//          3
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 'A'}, {2, 'B'}},
		2: {{3, 'C'}},
		3: {},
	}}
	f, l, s := g.FromList()
	fmt.Println("simple forest:", s)
	fmt.Println("n  from  label")
	for n, e := range f.Paths {
		fmt.Printf("%d   %2d", n, e.From)
		if e.From < 0 {
			fmt.Println()
		} else {
			fmt.Printf("     %c\n", l[n])
		}
	}
	// Output:
	// simple forest: true
	// n  from  label
	// 0   -1
	// 1    0     A
	// 2    0     B
	// 3    2     C
}

func ExampleLabeledDirected_SpanTree() {
	//      0       4
	// 'A' / \ 'B'   \ 'D'
	//    1   2       5
	//         \ 'C'
	//          3
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 'A'}, {2, 'B'}},
		2: {{3, 'C'}},
		4: {{5, 'D'}},
		5: {},
	}}
	var f graph.FromList
	l := make([]graph.LI, g.Order())
	ns, simple := g.SpanTree(2, &f, l)
	fmt.Println("nodes spanned:", ns)
	fmt.Println("simple tree:", simple)
	fmt.Println("n  from  label")
	for n, e := range f.Paths {
		fmt.Printf("%d   %2d", n, e.From)
		if e.From < 0 {
			fmt.Println()
		} else {
			fmt.Printf("     %c\n", l[n])
		}
	}
	// Output:
	// nodes spanned: 2
	// simple tree: true
	// n  from  label
	// 0   -1
	// 1   -1
	// 2   -1
	// 3    2     C
	// 4   -1
	// 5   -1
}

func ExampleLabeledDirected_Transpose() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}}
	tr, m := g.Transpose()
	for fr, to := range tr.LabeledAdjacencyList {
		fmt.Printf("%d %#v\n", fr, to)
	}
	fmt.Println(m, "arcs")
	// Output:
	// 0 []graph.Half{graph.Half{To:2, Label:7}}
	// 1 []graph.Half{graph.Half{To:2, Label:9}}
	// 2 []graph.Half(nil)
	// 2 arcs
}

func ExampleLabeledDirected_Undirected() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}}
	for fr, to := range g.Undirected().LabeledAdjacencyList {
		fmt.Printf("%d %#v\n", fr, to)
	}
	// Output:
	// 0 []graph.Half{graph.Half{To:2, Label:7}}
	// 1 []graph.Half{graph.Half{To:2, Label:9}}
	// 2 []graph.Half{graph.Half{To:0, Label:7}, graph.Half{To:1, Label:9}}
}

func ExampleLabeledDirected_UnlabeledTranspose() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}}

	fmt.Println("two steps:")
	ut, m := g.Unlabeled().Transpose()
	for fr, to := range ut.AdjacencyList {
		fmt.Println(fr, to)
	}
	fmt.Println(m, "arcs")

	fmt.Println("direct:")
	ut, m = g.UnlabeledTranspose()
	for fr, to := range ut.AdjacencyList {
		fmt.Println(fr, to)
	}
	fmt.Println(m, "arcs")
	// Output:
	// two steps:
	// 0 [2]
	// 1 [2]
	// 2 []
	// 2 arcs
	// direct:
	// 0 [2]
	// 1 [2]
	// 2 []
	// 2 arcs
}
