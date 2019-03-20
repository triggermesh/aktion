// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleArcDensity() {
	fmt.Println(graph.ArcDensity(4, 3))
	// Output:
	// 0.25
}

func ExampleDensity() {
	fmt.Println(graph.Density(4, 3))
	// Output:
	// 0.5
}

func ExampleUndirected_AddEdge_justOne() {
	//    0
	//   /
	//  1
	var g graph.Undirected // no need to pre-allocate
	g.AddEdge(0, 1)        // AddEdge expands graph as needed
	for fr, to := range g.AdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [1]
	// 1 [0]
}

func ExampleUndirected_AddEdge_more() {
	//    0
	//   / \\
	//  1---2--\
	//       \-/
	//
	// preallocate graph with make for efficiency
	g := graph.Undirected{make(graph.AdjacencyList, 3)} // 3 nodes altogether
	// alternatively, use a literal
	g = graph.Undirected{graph.AdjacencyList{2: nil}} // 2 is last node
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 0) // parallel
	g.AddEdge(2, 2) // loop
	for fr, to := range g.AdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [1 2 2]
	// 1 [0 2]
	// 2 [1 0 0 2]
}

func ExampleUndirected_Edges() {
	//    0
	//   / \\
	//  1---2--\
	//       \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 0) // parallel
	g.AddEdge(2, 2) // loop
	g.Edges(func(e graph.Edge) {
		fmt.Println(e)
	})
	// Output:
	// {1 0}
	// {2 1}
	// {2 0}
	// {2 0}
	// {2 2}
}

func ExampleUndirected_HasEdge() {
	var g graph.Undirected
	g.AddEdge(7, 8)
	g.AddEdge(8, 8)
	g.AddEdge(5, 7)
	g.AddEdge(5, 8)
	g.AddEdge(5, 9)
	fmt.Println(g.HasEdge(4, 7))
	var n1, n2 graph.NI = 5, 8
	has, x1, x2 := g.HasEdge(5, 8)
	fmt.Println(has, x1, x2)
	a := g.AdjacencyList
	fmt.Println(a[n1][x1] == n2, a[n2][x2] == n1)
	// Output:
	// false -1 -1
	// true 1 2
	// true true
}

func ExampleUndirected_RemoveEdge() {
	//    0
	//   / \\
	//  1---2--\
	//       \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 0) // parallel
	g.AddEdge(2, 2) // loop

	fmt.Println(g.RemoveEdge(2, 0)) // remove one of the parallel edges
	fmt.Println(g.RemoveEdge(2, 2)) // remove the loop
	fmt.Println(g.RemoveEdge(2, 2)) // false: there was only one loop

	for fr, to := range g.AdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// true
	// true
	// false
	// 0 [1 2]
	// 1 [0 2]
	// 2 [1 0]
}

func ExampleUndirected_SimpleEdges() {
	//    0
	//   / \\
	//  1---2--\
	//       \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 0) // parallel
	g.AddEdge(2, 2) // loop
	g.SimpleEdges(func(e graph.Edge) {
		fmt.Println(e)
	})
	// Output:
	// {0 1}
	// {0 2}
	// {1 2}
}

func ExampleUndirected_FromList() {
	//    0   5
	//   / \   \
	//  1   2   6
	//     / \
	//    3---4
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(2, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	g.AddEdge(5, 6)
	f, r, s := g.FromList()
	fmt.Println("simple forest:", s)
	fmt.Println("roots:", r)
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %8d\n", n, e.From, e.Len)
	}
	fmt.Println("MaxLen:       ", f.MaxLen)
	// Output:
	// simple forest: false
	// roots: [0 5]
	// n  from  path len
	// 0   -1         1
	// 1    0         2
	// 2    0         2
	// 3    2         3
	// 4    2         3
	// 5   -1         1
	// 6    5         2
	// MaxLen:        3
}

func ExampleUndirected_SpanTree() {
	//    4   3
	//   / \
	//  2   1
	//       \
	//        0
	var g graph.Undirected
	g.AddEdge(2, 4)
	g.AddEdge(4, 1)
	g.AddEdge(1, 0)
	var f graph.FromList
	n, st := g.SpanTree(4, &f)
	fmt.Println("Nodes spanned:", n)
	fmt.Println("Simple tree:", st)
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %8d\n", n, e.From, e.Len)
	}
	fmt.Println("MaxLen:       ", f.MaxLen)
	// Output:
	// Nodes spanned: 4
	// Simple tree: true
	// n  from  path len
	// 0    1         3
	// 1    4         2
	// 2    4         2
	// 3   -1         0
	// 4   -1         1
	// MaxLen:        3
}

func ExampleUndirected_SpanTree_cycle() {
	//    0
	//   / \
	//  1   2
	//     / \
	//    3---4
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(2, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	var f graph.FromList
	n, st := g.SpanTree(0, &f)
	fmt.Println("Nodes spanned:", n)
	fmt.Println("Simple tree:", st)
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %8d\n", n, e.From, e.Len)
	}
	fmt.Println("MaxLen:       ", f.MaxLen)
	// Output:
	// Nodes spanned: 5
	// Simple tree: false
	// n  from  path len
	// 0   -1         1
	// 1    0         2
	// 2    0         2
	// 3    2         3
	// 4    2         3
	// MaxLen:        3
}

func ExampleUndirected_SpanTree_loop() {
	//    0
	//   / \ /-\
	//  1   2--/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(2, 2)
	var f graph.FromList
	n, st := g.SpanTree(0, &f)
	fmt.Println("Nodes spanned:", n)
	fmt.Println("Simple tree:", st)
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %8d\n", n, e.From, e.Len)
	}
	fmt.Println("MaxLen:       ", f.MaxLen)
	// Output:
	// Nodes spanned: 3
	// Simple tree: false
	// n  from  path len
	// 0   -1         1
	// 1    0         2
	// 2    0         2
	// MaxLen:        2
}

func ExampleUndirected_SpanTree_loopDisconnected() {
	//    0
	//   /   /-\
	//  1   2--/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(2, 2)
	var f graph.FromList
	n, st := g.SpanTree(0, &f)
	fmt.Println("Nodes spanned:", n)
	fmt.Println("Simple tree:", st)
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %8d\n", n, e.From, e.Len)
	}
	fmt.Println("MaxLen:       ", f.MaxLen)
	// Output:
	// Nodes spanned: 2
	// Simple tree: true
	// n  from  path len
	// 0   -1         1
	// 1    0         2
	// 2   -1         0
	// MaxLen:        2
}

func ExampleUndirected_SpanTree_multigraph() {
	//    0
	//   / \
	//  1   2==3
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(2, 3)
	g.AddEdge(2, 3)
	var f graph.FromList
	n, st := g.SpanTree(0, &f)
	fmt.Println("Nodes spanned:", n)
	fmt.Println("Simple tree:", st)
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %8d\n", n, e.From, e.Len)
	}
	fmt.Println("MaxLen:       ", f.MaxLen)
	// Output:
	// Nodes spanned: 4
	// Simple tree: false
	// n  from  path len
	// 0   -1         1
	// 1    0         2
	// 2    0         2
	// 3    2         3
	// MaxLen:        3
}

func ExampleUndirected_TarjanBiconnectedComponents() {
	// undirected edges:
	// 3---2---1---7---9
	//  \ / \ / \   \ /
	//   4   5---6   8
	var g graph.Undirected
	g.AddEdge(3, 4)
	g.AddEdge(3, 2)
	g.AddEdge(2, 4)
	g.AddEdge(2, 5)
	g.AddEdge(2, 1)
	g.AddEdge(5, 1)
	g.AddEdge(6, 1)
	g.AddEdge(6, 5)
	g.AddEdge(7, 1)
	g.AddEdge(7, 9)
	g.AddEdge(7, 8)
	g.AddEdge(9, 8)
	g.TarjanBiconnectedComponents(func(bcc []graph.Edge) bool {
		fmt.Println("Edges:")
		for _, e := range bcc {
			fmt.Println(e)
		}
		return true
	})
	// Output:
	// Edges:
	// {4 2}
	// {3 4}
	// {2 3}
	// Edges:
	// {6 1}
	// {5 6}
	// {5 1}
	// {2 5}
	// {1 2}
	// Edges:
	// {8 7}
	// {9 8}
	// {7 9}
	// Edges:
	// {1 7}
}

func ExampleUndirected_BlockCut() {
	// undirected edges:
	// 3---2---1---7---9
	//  \ / \ / \   \ /
	//   4   5---6   8
	var g graph.Undirected
	g.AddEdge(3, 4)
	g.AddEdge(3, 2)
	g.AddEdge(2, 4)
	g.AddEdge(2, 5)
	g.AddEdge(2, 1)
	g.AddEdge(5, 1)
	g.AddEdge(6, 1)
	g.AddEdge(6, 5)
	g.AddEdge(7, 1)
	g.AddEdge(7, 9)
	g.AddEdge(7, 8)
	g.AddEdge(9, 8)
	g.BlockCut(
		func(bcc []graph.Edge) bool {
			fmt.Println("Edges:")
			for _, e := range bcc {
				fmt.Println(e)
			}
			return true
		},
		func(ap graph.NI) bool {
			fmt.Println("Articulation Point:", ap)
			return true
		},
		func(ip graph.NI) bool {
			fmt.Println("Isolated Point:", ip)
			return true
		})
	// Output:
	// Isolated Point: 0
	// Articulation Point: 2
	// Edges:
	// {4 2}
	// {3 4}
	// {2 3}
	// Edges:
	// {6 1}
	// {5 6}
	// {5 1}
	// {2 5}
	// {1 2}
	// Articulation Point: 7
	// Edges:
	// {8 7}
	// {9 8}
	// {7 9}
	// Edges:
	// {1 7}
	// Articulation Point: 1
}

func ExampleLabeledUndirected_AddEdge() {
	//       --0--
	//      /     \\6001
	// 5000/   6000\\
	//    /         \\
	//   1-----------2--\8000
	//        5000    \-/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 5000)
	g.AddEdge(graph.Edge{1, 2}, 5000)
	g.AddEdge(graph.Edge{2, 0}, 6000)
	g.AddEdge(graph.Edge{2, 0}, 6001) // parallel
	g.AddEdge(graph.Edge{2, 2}, 8000) // loop
	for fr, to := range g.LabeledAdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [{1 5000} {2 6000} {2 6001}]
	// 1 [{0 5000} {2 5000}]
	// 2 [{1 5000} {0 6000} {0 6001} {2 8000}]
}

func ExampleLabeledUndirected_Edges() {
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 0}, 'A')
	g.AddEdge(graph.Edge{0, 1}, 'B')
	g.AddEdge(graph.Edge{1, 2}, 'C')
	g.AddEdge(graph.Edge{1, 2}, 'D')
	g.Edges(func(e graph.LabeledEdge) {
		fmt.Printf("%d %c\n", e.Edge, e.LI)
	})
	// Output:
	// {0 0} A
	// {1 0} B
	// {2 1} C
	// {2 1} D
}

func ExampleLabeledUndirected_HasEdge() {
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{7, 8}, 'A')
	g.AddEdge(graph.Edge{8, 8}, 'B')
	g.AddEdge(graph.Edge{5, 7}, 'C')
	g.AddEdge(graph.Edge{5, 8}, 'D')
	g.AddEdge(graph.Edge{5, 9}, 'E')
	a := g.LabeledAdjacencyList
	var n1, n2 graph.NI = 5, 8
	has, x1, x2 := g.HasEdge(n1, n2)
	fmt.Printf("%t %d %d %c %c\n",
		has, x1, x2, a[n1][x1].Label, a[n2][x2].Label)
	// Output:
	// true 1 2 D D
}

func ExampleLabeledUndirected_HasEdgeLabel() {
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{7, 8}, 'A')
	g.AddEdge(graph.Edge{8, 8}, 'B')
	g.AddEdge(graph.Edge{5, 7}, 'C')
	g.AddEdge(graph.Edge{5, 8}, 'D')
	g.AddEdge(graph.Edge{5, 9}, 'E')
	fmt.Println(g.HasEdgeLabel(7, 8, 'D'))
	var n1, n2 graph.NI = 5, 8
	var l graph.LI = 'D'
	has, x1, x2 := g.HasEdgeLabel(5, 8, l)
	fmt.Println(has, x1, x2)
	a := g.LabeledAdjacencyList
	fmt.Println(a[n1][x1] == graph.Half{n2, l}, a[n2][x2] == graph.Half{n1, l})
	// Output:
	// false -1 -1
	// true 1 2
	// true true
}

func ExampleLabeledUndirected_FromList() {
	//      0
	// 'A' / \ 'B'
	//    1---2
	//     'C' \ 'D'
	//          3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'A')
	g.AddEdge(graph.Edge{0, 2}, 'B')
	g.AddEdge(graph.Edge{1, 2}, 'C')
	g.AddEdge(graph.Edge{2, 3}, 'D')
	f, l, r, s := g.FromList()
	fmt.Println("simple forest:", s)
	fmt.Println("roots:", r)
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
	// simple forest: false
	// roots: [0]
	// n  from  label
	// 0   -1
	// 1    0     A
	// 2    0     B
	// 3    2     D
}

func ExampleLabeledUndirected_SpanTree() {
	//      0
	// 'A' / \ 'B'
	//    1---2
	//     'C' \ 'D'
	//          3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'A')
	g.AddEdge(graph.Edge{0, 2}, 'B')
	g.AddEdge(graph.Edge{1, 2}, 'C')
	g.AddEdge(graph.Edge{2, 3}, 'D')
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
	// nodes spanned: 4
	// simple tree: false
	// n  from  label
	// 0    2     B
	// 1    2     C
	// 2   -1
	// 3    2     D
}

func ExampleLabeledUndirected_RemoveEdge() {
	//       --0--
	//      /     \\6001
	// 5000/   6000\\
	//    /         \\
	//   1-----------2--\8000
	//        5000    \-/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 5000)
	g.AddEdge(graph.Edge{1, 2}, 5000)
	g.AddEdge(graph.Edge{2, 0}, 6000)
	g.AddEdge(graph.Edge{2, 0}, 6001) // parallel
	g.AddEdge(graph.Edge{2, 2}, 8000) // loop

	fmt.Println(g.RemoveEdge(0, 2)) // remove one of the parallel arcs
	fmt.Println(g.RemoveEdge(0, 0)) // false: no loop on 0

	for fr, to := range g.LabeledAdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// true 6000
	// false 0
	// 0 [{1 5000} {2 6001}]
	// 1 [{0 5000} {2 5000}]
	// 2 [{1 5000} {2 8000} {0 6001}]
}

func ExampleLabeledUndirected_RemoveEdgeLabel() {
	//       --0--
	//      /     \\6001
	// 5000/   6000\\
	//    /         \\
	//   1-----------2--\8000
	//        5000    \-/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 5000)
	g.AddEdge(graph.Edge{1, 2}, 5000)
	g.AddEdge(graph.Edge{2, 0}, 6000)
	g.AddEdge(graph.Edge{2, 0}, 6001) // parallel
	g.AddEdge(graph.Edge{2, 2}, 8000) // loop

	fmt.Println(g.RemoveEdgeLabel(2, 0, 6001))
	fmt.Println(g.RemoveEdgeLabel(1, 2, 1000))

	for fr, to := range g.LabeledAdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// true
	// false
	// 0 [{1 5000} {2 6000}]
	// 1 [{0 5000} {2 5000}]
	// 2 [{1 5000} {0 6000} {2 8000}]
}

func ExampleLabeledUndirected_TarjanBiconnectedComponents() {
	// undirected edges:
	// 3---2---1---7---9
	//  \ / \ / \   \ /
	//   4   5---6   8
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{2, 5}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{6, 1}, 0)
	g.AddEdge(graph.Edge{6, 5}, 0)
	g.AddEdge(graph.Edge{7, 1}, 0)
	g.AddEdge(graph.Edge{7, 9}, 0)
	g.AddEdge(graph.Edge{7, 8}, 0)
	g.AddEdge(graph.Edge{9, 8}, 0)
	g.TarjanBiconnectedComponents(func(bcc []graph.LabeledEdge) bool {
		fmt.Println("Edges:")
		for _, e := range bcc {
			fmt.Println(e.Edge)
		}
		return true
	})
	// Output:
	// Edges:
	// {4 2}
	// {3 4}
	// {2 3}
	// Edges:
	// {6 1}
	// {5 6}
	// {5 1}
	// {2 5}
	// {1 2}
	// Edges:
	// {8 7}
	// {9 8}
	// {7 9}
	// Edges:
	// {1 7}
}

/* shelved
func ExampleBiconnectedComponents_Find() {
	g := graph.AdjacencyList{
		0:  {1, 7},
		1:  {2, 4, 0},
		2:  {3, 1},
		3:  {2, 4},
		4:  {3, 1},
		5:  {6, 12},
		6:  {5, 12, 8},
		7:  {8, 0},
		8:  {6, 7, 9, 10},
		9:  {8},
		10: {8, 13, 11},
		11: {10},
		12: {5, 6, 13},
		13: {12, 10},
	}
	b := graph.NewBiconnectedComponents(g)
	b.Find(0)
	fmt.Println("n: cut from")
	for n, f := range b.From {
		fmt.Printf("%d: %d %d\n",
			n, b.Cuts.Bit(n), f)
	}
	fmt.Println("Leaves:", b.Leaves)
	// Output:
	// n: cut from
	// 0: 1 -1
	// 1: 1 0
	// 2: 0 1
	// 3: 0 2
	// 4: 0 3
	// 5: 0 6
	// 6: 0 8
	// 7: 1 0
	// 8: 1 7
	// 9: 0 8
	// 10: 1 13
	// 11: 0 10
	// 12: 0 5
	// 13: 0 12
	// Leaves: [4 11 9]
}
*/
