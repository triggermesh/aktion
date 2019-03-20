// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt_test

import (
	"reflect"
	"testing"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func TestUndirectedEulerianCycle(t *testing.T) {
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 2)
	c, err := alt.EulerianCycle(g)
	if err != nil {
		t.Fatal(err)
	}
	// reconstruct from node list c
	var r graph.Undirected
	n1 := c[0]
	for _, n2 := range c[1:] {
		r.AddEdge(n1, n2)
		n1 = n2
	}
	// compare
	g.SortArcLists()
	r.SortArcLists()
	if !reflect.DeepEqual(r, g) {
		t.Fatal()
	}
}
