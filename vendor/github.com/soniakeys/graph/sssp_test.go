// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleLabeledAdjacencyList_AStarAPath() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label graph.LI) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	h := func(from graph.NI) float64 { return h4[from] }
	p, d := g.AStarAPath(0, 4, h, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path distance:", d)
	// Output:
	// Shortest path: {0 [{2 9} {3 11} {4 6}]}
	// Path distance: 26
}

func ExampleLabeledAdjacencyList_AStarMPath() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label graph.LI) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	h := func(from graph.NI) float64 { return h4[from] }
	p, d := g.AStarMPath(0, 4, h, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path distance:", d)
	// Output:
	// Shortest path: {0 [{2 9} {3 11} {4 6}]}
	// Path distance: 26
}

func ExampleHeuristic_Admissible() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label graph.LI) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	var h graph.Heuristic = func(from graph.NI) float64 { return h4[from] }
	fmt.Println(h.Admissible(g, w, 4))
	// Output:
	// true
}

func ExampleHeuristic_Monotonic() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label graph.LI) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	var h graph.Heuristic = func(from graph.NI) float64 { return h4[from] }
	fmt.Println(h.Monotonic(g, w))
	// Output:
	// true
}

func ExampleLabeledDirected_BellmanFord() {
	//              /--------3        4<-------9
	//              |        ^        |   (6)  ^
	//              |(1)     |        |        |
	//              |        |(-2)    |(3)     |
	//    (wt: 10)  v   (2)  |        v        |
	//  1---------->2------->6<-------5--------/
	//  |           ^        ^   (-1)    (-10)
	//  |(8)        |        |
	//  |           |(-4)    |(-1)
	//  v     (1)   |        |
	//  8---------->7--------/
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{2, 10}, {8, 8}},
		2: {{6, 2}},
		3: {{2, 1}},
		4: {{5, 3}},
		5: {{6, -1}, {9, -10}},
		6: {{3, -2}},
		7: {{6, -1}, {2, -4}},
		8: {{7, 1}},
		9: {{4, 6}},
	}}
	w := func(label graph.LI) float64 { return float64(label) }
	// graph contains negative cycle somewhere
	fmt.Println("negative cycle:", g.HasNegativeCycle(w))

	// but negative cycle not reached starting at node 1
	start := graph.NI(1)
	fmt.Println("start:", start)
	f, labels, dist, end := g.BellmanFord(w, start)
	if end >= 0 {
		fmt.Println("negative cycle")
		return
	} else {
		fmt.Println("no negative cycle reachable from", start)
	}
	fmt.Println("end       path  path")
	fmt.Println("node  LI  len   dist   path")
	p := make([]graph.NI, f.MaxLen)
	for n, e := range f.Paths {
		fmt.Printf("%d    %3d    %d   %4.0f   %d\n",
			n, labels[n], e.Len, dist[n], f.PathTo(graph.NI(n), p))
	}
	// Output:
	// negative cycle: true
	// start: 1
	// no negative cycle reachable from 1
	// end       path  path
	// node  LI  len   dist   path
	// 0      0    0   +Inf   []
	// 1      0    1      0   [1]
	// 2     -4    4      5   [1 8 7 2]
	// 3     -2    6      5   [1 8 7 2 6 3]
	// 4      0    0   +Inf   []
	// 5      0    0   +Inf   []
	// 6      2    5      7   [1 8 7 2 6]
	// 7      1    3      9   [1 8 7]
	// 8      8    2      8   [1 8]
	// 9      0    0   +Inf   []
}

func ExampleFromList_BellmanFordCycle() {
	//              /--------3        4<-------9<------10
	//              |        ^        |   (6)  ^   (7)
	//              |(1)     |        |        |
	//              |        |(-2)    |(3)     |
	//    (wt: 10)  v   (2)  |        v        |
	//  1---------->2------->6<-------5--------/
	//  |           ^        ^   (-1)    (-10)
	//  |(8)        |        |
	//  |           |(-4)    |(-1)
	//  v     (1)   |        |
	//  8---------->7--------/
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1:  {{2, 10}, {8, 8}},
		2:  {{6, 2}},
		3:  {{2, 1}},
		4:  {{5, 3}},
		5:  {{6, -1}, {9, -10}},
		6:  {{3, -2}},
		7:  {{6, -1}, {2, -4}},
		8:  {{7, 1}},
		9:  {{4, 6}},
		10: {{9, 7}},
	}}
	w := func(label graph.LI) float64 { return float64(label) }
	start := graph.NI(10)
	fmt.Println("start:", start)
	f, _, _, end := g.BellmanFord(w, start)
	fmt.Println("end of path with negative cycle:", end)
	fmt.Println("negative cycle:", f.BellmanFordCycle(end))
	// Output:
	// start: 10
	// end of path with negative cycle: 3
	// negative cycle: [9 4 5]
}

func ExampleLabeledDirected_NegativeCycle() {
	//              /--------3        4<-------9
	//              |        ^        |   (6)  ^
	//              |(1)     |        |        |
	//              |        |(-2)    |(3)     |
	//    (wt: 10)  v   (2)  |        v        |
	//  1---------->2------->6<-------5--------/
	//  |           ^        ^   (-1)    (-10)
	//  |(8)        |        |
	//  |           |(-4)    |(-1)
	//  v     (1)   |        |
	//  8---------->7--------/
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{2, 10}, {8, 8}},
		2: {{6, 2}},
		3: {{2, 1}},
		4: {{5, 3}},
		5: {{6, -1}, {9, -10}},
		6: {{3, -2}},
		7: {{6, -1}, {2, -4}},
		8: {{7, 1}},
		9: {{4, 6}},
	}}
	w := func(label graph.LI) float64 { return float64(label) }
	fmt.Println(g.NegativeCycle(w))
	// Output:
	// [{9 -10} {4 6} {5 3}]
}

func ExampleLabeledDirected_DAGOptimalPaths_allShortestPaths() {
	// arcs are directed right:
	//   (11)
	// 0------2         4
	//                   \(11)
	//      (11)    (10)  \   (30)   (10)
	//    1-------3--------5-------7-------9
	//                      \     /
	//                   (10)\   /(20)
	//                        \ /
	//                         6------8
	//                           (10)
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 2, Label: 11}},
		1: {{3, 11}},
		3: {{5, 10}},
		4: {{5, 11}},
		5: {{6, 10}, {7, 30}},
		6: {{7, 20}, {8, 10}},
		7: {{9, 10}},
		9: {},
	}}
	o := []graph.NI{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	w := func(l graph.LI) float64 { return float64(l) }
	// find all shortest paths from 3
	var start, end graph.NI = 3, -1
	f, labels, dist, reached := g.DAGOptimalPaths(start, end, o, w, false)
	fmt.Println("node  LI  path dist  path len  leaf")
	for n, pd := range dist {
		fmt.Printf("%d  %5d %7.0f  %9d %7d\n",
			n, labels[n], pd, f.Paths[n].Len, f.Leaves.Bit(n))
	}
	fmt.Println()
	fmt.Println("Nodes reached:       ", reached)
	fmt.Println("Max path len:        ", f.MaxLen)
	// Output:
	// node  LI  path dist  path len  leaf
	// 0      0       0          0       0
	// 1      0       0          0       0
	// 2      0       0          0       0
	// 3      0       0          1       0
	// 4      0       0          0       0
	// 5     10      10          2       0
	// 6     10      20          3       0
	// 7     30      40          3       0
	// 8     10      30          4       1
	// 9     10      50          4       1
	//
	// Nodes reached:        6
	// Max path len:         4
}

func ExampleLabeledDirected_DAGMinDistPath() {
	// arcs are directed right:
	//             4
	//        (-3)/ \(-2)
	//           /   \
	//    (10)  /     \   (5)    (10)
	// 3-------1-------0-------6-------2
	//           (-10)  \     /
	//                   \   /
	//                 (2)\ /(3)
	//                     5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 5, Label: 2}, {6, 5}},
		1: {{0, -10}, {4, -3}},
		3: {{1, 10}},
		4: {{0, -2}},
		5: {{6, 3}},
		6: {{2, 10}},
	}}
	var start, end graph.NI = 3, 2
	w := func(l graph.LI) float64 { return float64(l) }
	p, dist, err := g.DAGMinDistPath(start, end, w)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Path:", p)
	fmt.Println("Distance:", dist)
	// Output:
	// Path: {3 [{1 10} {0 -10} {6 5} {2 10}]}
	// Distance: 15
}

func ExampleLabeledDirected_DAGMaxDistPath() {
	// arcs are directed right:
	//             4
	//        (-3)/ \(-2)
	//           /   \
	//    (10)  /     \   (5)    (10)
	// 3-------1-------0-------6-------2
	//           (-10)  \     /
	//                   \   /
	//                 (2)\ /(3)
	//                     5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 5, Label: 2}, {6, 5}},
		1: {{0, -10}, {4, -3}},
		3: {{1, 10}},
		4: {{0, -2}},
		5: {{6, 3}},
		6: {{2, 10}},
	}}
	var start, end graph.NI = 3, 2
	w := func(l graph.LI) float64 { return float64(l) }
	p, dist, err := g.DAGMaxDistPath(start, end, w)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Path:", p)
	fmt.Println("Distance:", dist)
	// Output:
	// Path: {3 [{1 10} {4 -3} {0 -2} {5 2} {6 3} {2 10}]}
	// Distance: 20
}

func ExampleLabeledAdjacencyList_DijkstraPath() {
	// arcs are directed right:
	//          (wt: 11)
	//       --------------6----
	//      /             /     \
	//     /             /(2)    \(9)
	//    /     (9)     /         \
	//   1-------------3----       5
	//    \           /     \     /
	//     \     (10)/   (11)\   /(7)
	//   (7)\       /         \ /
	//       ------2-----------4
	//                 (15)
	g := graph.LabeledAdjacencyList{
		1: {{To: 2, Label: 7}, {To: 3, Label: 9}, {To: 6, Label: 11}},
		2: {{To: 3, Label: 10}, {To: 4, Label: 15}},
		3: {{To: 4, Label: 11}, {To: 6, Label: 2}},
		4: {{To: 5, Label: 7}},
		6: {{To: 5, Label: 9}},
	}
	w := func(label graph.LI) float64 { return float64(label) }
	p, d := g.DijkstraPath(1, 5, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path distance:", d)
	// Output:
	// Shortest path: {1 [{6 11} {5 9}]}
	// Path distance: 20
}

func ExampleLabeledAdjacencyList_Dijkstra_allPaths() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label graph.LI) float64 { return float64(label) }
	start := graph.NI(2)
	f, labels, dist, n := g.Dijkstra(start, -1, w)
	fmt.Println(n, "paths found.")
	fmt.Println("node:  path      len  dist   LI")
	p := make([]graph.NI, f.MaxLen)
	for nd := range g {
		r := &f.Paths[nd]
		path := f.PathTo(graph.NI(nd), p)
		if r.Len > 0 {
			fmt.Printf("%d:     %-11s %d    %2.0f  %3d\n",
				nd, fmt.Sprint(path), r.Len, dist[nd], labels[nd])
		}
	}

	// Output:
	// 4 paths found.
	// node:  path      len  dist   LI
	// 2:     [2]         1     0    0
	// 3:     [2 3]       2    11   11
	// 4:     [2 3 4]     3    17    6
	// 5:     [2 5]       2     2    2
}

func TestSSSP(t *testing.T) {
	r100 := r(100, 200, 62)
	testSSSP(r100, t)
}

func testSSSP(tc testCase, t *testing.T) {
	w := func(label graph.LI) float64 { return tc.w[label] }
	f, labels, dist, _ := tc.l.LabeledAdjacencyList.Dijkstra(tc.start, tc.end, w)
	pathD := f.PathToLabeled(tc.end, labels, nil)
	distD := dist[tc.end]
	// A*
	pathA, distA := tc.l.AStarAPath(tc.start, tc.end, tc.h, w)
	// test that a* path is same distance and length as dijkstra path
	if len(pathA.Path) != len(pathD.Path) {
		t.Log("pathA:", pathA)
		t.Log("pathD:", pathD)
		t.Fatal(len(tc.w), "A, D len mismatch")
	}
	if distA != distD {
		t.Log("distA:", distA)
		t.Log("distD:", distD)
		t.Log("delta:", math.Abs(distA-distD))
		t.Fatal(len(tc.w), "A, D dist mismatch")
	}
	// test Bellman Ford against Dijkstra all paths
	dr, _, _, _ := tc.l.LabeledAdjacencyList.Dijkstra(tc.start, -1, w)
	br, _, _, _ := tc.l.BellmanFord(w, tc.start)
	// result objects should be identical
	if len(dr.Paths) != len(br.Paths) {
		t.Fatal("len(dr.Paths), len(br.Paths)",
			len(dr.Paths), len(br.Paths))
	}
	/* breadth first, compare to dijkstra with unit weights
	w = func(graph.LI) float64 { return 1 }
	ur, _, _ := tc.l.LabeledAdjacencyList.Dijkstra(tc.start, -1, w)
	var bfsr graph.FromList
	np, _ := tc.g.AdjacencyList.breadthFirst(tc.start, nil, &bfsr,
		func(n graph.NI) bool { return true })
	var ml, npf int
	for i, ue := range ur.Paths {
		bl := bfsr.Paths[i].Len
		if bl != ue.Len {
			t.Fatal("ue.From.Len, bfsr.Paths[i].Len", ue.Len, bl)
		}
		if bl > ml {
			ml = bl
		}
		if bl > 0 {
			npf++
		}
	}
	if ml != bfsr.MaxLen {
		t.Fatal("bfsr.MaxLen, recomputed", bfsr.MaxLen, ml)
	}
	if npf != np {
		t.Fatal("bfs all paths returned", np, "recount:", npf)
	}
	*/
}

type testCase struct {
	l graph.LabeledDirected // generated labeled directed graph
	w []float64             // arc weights for l
	// variants
	g graph.Directed // unlabeled
	t graph.Directed // transpose

	h graph.Heuristic

	start, end graph.NI
	m          int
}

func r(nNodes, nArcs int, seed int64) testCase {
	s := rand.New(rand.NewSource(59))
	s.Seed(seed)
	l, coords, w, err := graph.LabeledEuclidean(nNodes, nArcs, 1, 1, s)
	if err != nil {
		panic(err)
	}
	tc := testCase{
		l:     l,
		w:     w,
		start: graph.NI(s.Intn(nNodes)), // random start
	}
	// end is point at distance nearest target distance
	const target = .3
	nearest := 2.
	c1 := coords[tc.start]
	for i, c2 := range coords {
		d := math.Abs(target - math.Hypot(c2.X-c1.X, c2.Y-c1.Y))
		if d < nearest {
			tc.end = graph.NI(i)
			nearest = d
		}
	}
	// with end chosen, define heuristic
	ce := coords[tc.end]
	tc.h = func(n graph.NI) float64 {
		cn := &coords[n]
		return math.Hypot(ce.X-cn.X, ce.Y-cn.Y)
	}
	// variants
	tc.g = tc.l.Unlabeled()
	tc.t, tc.m = tc.g.Transpose()
	return tc
}
