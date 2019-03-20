// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/soniakeys/graph"
)

func ExampleEuclidean() {
	r := rand.New(rand.NewSource(7))
	g, pos, err := graph.Euclidean(4, 6, 1, 1, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(g.Order(), "nodes")
	fmt.Println("n  position")
	for n, p := range pos {
		fmt.Printf("%d  (%.2f, %.2f)\n", n, p.X, p.Y)
	}
	fmt.Println(g.ArcSize(), "arcs:")
	for n, to := range g.AdjacencyList {
		fmt.Println(n, "->", to)
	}
	// Random output:
	// 4 nodes
	// n  position
	// 0  (0.92, 0.23)
	// 1  (0.24, 0.91)
	// 2  (0.70, 0.15)
	// 3  (0.35, 0.34)
	// 6 arcs:
	// 0 -> [2 1]
	// 1 -> [0]
	// 2 -> [3 0]
	// 3 -> [1]
}

func ExampleGeometric() {
	r := rand.New(rand.NewSource(7))
	g, pos, m := graph.Geometric(4, .6, r)
	fmt.Println(g.Order(), "nodes")
	fmt.Println("n  position")
	for n, p := range pos {
		fmt.Printf("%d  (%.2f, %.2f)\n", n, p.X, p.Y)
	}
	fmt.Println(m, "edges:")
	for n, to := range g.AdjacencyList {
		for _, to := range to {
			if graph.NI(n) < to {
				fmt.Println(n, "-", to)
			}
		}
	}
	// Random output:
	// 4 nodes
	// n  position
	// 0  (0.92, 0.23)
	// 1  (0.24, 0.91)
	// 2  (0.70, 0.15)
	// 3  (0.35, 0.34)
	// 4 edges:
	// 0 - 2
	// 0 - 3
	// 1 - 3
	// 2 - 3
}

func ExampleKroneckerDirected() {
	r := rand.New(rand.NewSource(7))
	g, ma := graph.KroneckerDirected(2, 2, r)
	a := g.AdjacencyList
	fmt.Println(len(a), "nodes")
	fmt.Println(ma, "arcs:")
	for fr, to := range a {
		fmt.Println(fr, "->", to)
	}
	// Random output:
	// 4 nodes
	// 5 arcs:
	// 0 -> [2]
	// 1 -> [2]
	// 2 -> [1]
	// 3 -> [2 1]
}

func ExampleKroneckerUndirected() {
	r := rand.New(rand.NewSource(7))
	g, m := graph.KroneckerUndirected(2, 2, r)
	a := g.AdjacencyList
	fmt.Println(len(a), "nodes")
	fmt.Println(m, "edges:")
	for fr, to := range a {
		for _, to := range to {
			if graph.NI(fr) < to {
				fmt.Println(fr, "-", to)
			}
		}
	}
	// Random output:
	// 4 nodes
	// 4 edges:
	// 0 - 2
	// 1 - 2
	// 1 - 3
	// 2 - 3
}

func ExampleLabeledGeometric() {
	r := rand.New(rand.NewSource(7))
	g, pos, wt := graph.LabeledGeometric(4, .6, r)
	fmt.Println(g.Order(), "nodes")
	fmt.Println("n  position")
	for n, p := range pos {
		fmt.Printf("%d  (%.2f, %.2f)\n", n, p.X, p.Y)
	}
	fmt.Println(len(wt), "edges:")
	for n, to := range g.LabeledAdjacencyList {
		for _, to := range to {
			if graph.NI(n) < to.To {
				fmt.Println(n, "-", to.To)
			}
		}
	}
	fmt.Println(g.ArcSize(), "arcs:")
	fmt.Println("arc  label  weight")
	for n, to := range g.LabeledAdjacencyList {
		for _, to := range to {
			fmt.Printf("%d->%d   %d    %.2f\n",
				n, to.To, to.Label, wt[to.Label])
		}
	}
	// Random output:
	// 4 nodes
	// n  position
	// 0  (0.92, 0.23)
	// 1  (0.24, 0.91)
	// 2  (0.70, 0.15)
	// 3  (0.35, 0.34)
	// 4 edges:
	// 0 - 2
	// 0 - 3
	// 1 - 3
	// 2 - 3
	// 8 arcs:
	// arc  label  weight
	// 0->2   0    0.24
	// 0->3   1    0.58
	// 1->3   2    0.58
	// 2->0   0    0.24
	// 2->3   3    0.40
	// 3->0   1    0.58
	// 3->1   2    0.58
	// 3->2   3    0.40
}

func ExampleLabeledEuclidean() {
	r := rand.New(rand.NewSource(7))
	g, pos, wt, err := graph.LabeledEuclidean(4, 6, 1, 1, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(g.Order(), "nodes")
	fmt.Println("n  position")
	for n, p := range pos {
		fmt.Printf("%d  (%.2f, %.2f)\n", n, p.X, p.Y)
	}
	fmt.Println(g.ArcSize(), "arcs:")
	for n, to := range g.LabeledAdjacencyList {
		fmt.Println(n, "->", to)
	}
	fmt.Println("arc  label  weight")
	for n, to := range g.LabeledAdjacencyList {
		for _, to := range to {
			fmt.Printf("%d->%d   %d    %.2f\n",
				n, to.To, to.Label, wt[to.Label])
		}
	}
	// Random output:
	// 4 nodes
	// n  position
	// 0  (0.92, 0.23)
	// 1  (0.24, 0.91)
	// 2  (0.70, 0.15)
	// 3  (0.35, 0.34)
	// 6 arcs:
	// 0 -> [{2 1} {1 5}]
	// 1 -> [{0 4}]
	// 2 -> [{3 0} {0 2}]
	// 3 -> [{1 3}]
	// arc  label  weight
	// 0->2   1    0.24
	// 0->1   5    0.96
	// 1->0   4    0.96
	// 2->3   0    0.40
	// 2->0   2    0.24
	// 3->1   3    0.58
}

func TestEuclidean(t *testing.T) {
	var g graph.Directed
	var err error
	for {
		if g, _, err = graph.Euclidean(10, 30, 2, 10, nil); err == nil {
			break
		}
	}
	if s, n := g.IsSimple(); !s {
		t.Fatalf("Euclidean returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
}

func TestKroneckerDir(t *testing.T) {
	g, _ := graph.KroneckerDirected(10, 10, nil)
	if s, n := g.IsSimple(); !s {
		t.Fatalf("KroneckerDir returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
}

func TestKroneckerUndir(t *testing.T) {
	g, _ := graph.KroneckerUndirected(10, 10, nil)
	if s, n := g.IsSimple(); !s {
		t.Fatalf("KroneckerUndir returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
	if u, from, to := g.IsUndirected(); !u {
		t.Fatalf("KroneckerUndir returned directed graph.  "+
			"Arc %d->%d has no reciprocal.", from, to)
	}
}

func TestGnmUndirected(t *testing.T) {
	u := graph.GnmUndirected(15, 21, nil)
	if ok, _, _ := u.IsUndirected(); !ok {
		t.Fatal("GnmUndirected returned directed graph")
	}
	if ok, _ := u.IsSimple(); !ok {
		t.Fatal("GnmUndirected returned non-simple graph")
	}

	u = graph.GnmUndirected(15, 84, nil)
	/*
		rand.New(rand.NewSource(time.Now().UnixNano())))
		for fr, to := range u.AdjacencyList {
			t.Log(fr, to)
		}
		t.Log("order, size: ", u.Order(), u.Size())
		t.Log("density: ", u.Density())
	*/
	if ok, _, _ := u.IsUndirected(); !ok {
		t.Fatal("GnmUndirected returned directed graph")
	}
	if ok, _ := u.IsSimple(); !ok {
		t.Fatal("GnmUndirected returned non-simple graph")
	}
}

func TestGnmDirected(t *testing.T) {
	d := graph.GnmDirected(15, 189, nil)
	if ok, _ := d.IsSimple(); !ok {
		t.Fatal("GnmDirected returned non-simple graph")
	}
	d = graph.GnmDirected(15, 21, nil)
	if ok, _ := d.IsSimple(); !ok {
		t.Fatal("GnmDirected returned non-simple graph")
	}
}

func TestGnm3Directed(t *testing.T) {
	d := graph.Gnm3Directed(15, 42, nil)
	if ok, _ := d.IsSimple(); !ok {
		t.Fatal("Gnm3Directed returned non-simple graph")
	}
}

func TestGnm3Undirected(t *testing.T) {
	u := graph.Gnm3Undirected(15, 21, nil)
	if ok, _, _ := u.IsUndirected(); !ok {
		t.Fatal("Gnm3Undirected returned directed graph")
	}
	if ok, _ := u.IsSimple(); !ok {
		t.Fatal("Gnm3Undirected returned non-simple graph")
	}
}

func TestGnpUndirected(t *testing.T) {
	u, _ := graph.GnpUndirected(15, .4, nil)
	if ok, _, _ := u.IsUndirected(); !ok {
		t.Fatal("GnpUndirected returned directed graph")
	}
	if ok, _ := u.IsSimple(); !ok {
		t.Fatal("GnpUndirected returned non-simple graph")
	}
}

func TestGnpDirected(t *testing.T) {
	u, _ := graph.GnpDirected(15, .4, nil)
	if ok, _ := u.IsSimple(); !ok {
		t.Fatal("GnpDirected returned non-simple graph")
	}
}

func TestChungLu(t *testing.T) {
	w := make([]float64, 15)
	for i := range w {
		w[i] = (15 - float64(i)) * .8
	}
	u, _ := graph.ChungLu(w, rand.New(rand.NewSource(time.Now().UnixNano())))
	if ok, _ := u.IsSimple(); !ok {
		t.Fatal("ChungLu returned non-simple graph")
	}
}
