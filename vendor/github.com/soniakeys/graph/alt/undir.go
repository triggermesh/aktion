// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt

import (
	"errors"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

// EulerianCycle finds an Eulerian cycle in a directed multigraph.
//
// * If g has no nodes, result is nil, nil.
//
// * If g is Eulerian, result is an Eulerian cycle with err = nil.
// The cycle result is a list of nodes, where the first and last
// nodes are the same.
//
// * Otherwise, result is nil, error
//
// The fast Eulerian cycle algorithm of the main graph library is destructive
// on its input graph.  A non-destructive variant simply makes a copy first.
// This non-destructive variant for undirected graphs makes its copy not as
// the usual slice-based adjacency list, but as a custom map-based one.
// The idea was that the destructive removal of reciprocal edges requires
// an additional scan of a slice-based neigbor list, potentially increasing
// the asymptotic time complexity.  A map-based adjacency list keeps node
// retrieval constant time, preserving the time complexity of the original
// algorithm.  In practice, it doesn't seem worth the overhead of maps.
// Anecdote shows it significantly slower.
func EulerianCycle(g graph.Undirected) ([]graph.NI, error) {
	if g.Order() == 0 {
		return nil, nil
	}
	u := newUlerian(g)
	for u.s >= 0 {
		v := u.top()
		u.push()
		if u.top() != v {
			return nil, errors.New("not balanced")
		}
		u.keep()
	}
	if !u.uv.AllZeros() {
		return nil, errors.New("not strongly connected")
	}
	return u.p, nil
}

// undirected variant of similar class in dir.go
type ulerian struct {
	g  []map[graph.NI]int // map (multiset) based copy of graph
	m  int                // number of arcs in g, updated as g is consumed
	uv bits.Bits          // unvisited
	// low end of p is stack of unfinished nodes
	// high end is finished path
	p []graph.NI // stack + path
	s int        // stack pointer
}

func (u *ulerian) top() graph.NI {
	return u.p[u.s]
}

// starting with the node on top of the stack, move nodes with no arcs.
func (u *ulerian) keep() {
	for u.s >= 0 {
		n := u.top()
		if len(u.g[n]) > 0 {
			break
		}
		u.p[u.m] = n
		u.s--
		u.m--
	}
}

func (e *ulerian) push() {
	for u := e.top(); ; {
		e.uv.SetBit(int(u), 0)
		arcs := e.g[u]
		if len(arcs) == 0 {
			return
		}
		// pick an arc
		var w graph.NI
		{
			var c int
			for w, c = range arcs {
				break
			}
			// consume arc
			if c > 1 {
				arcs[w]--
			} else {
				delete(arcs, w)
			}
			// consume reciprocal arc as well
			r := e.g[w]
			if r[u] > 1 {
				r[u]--
			} else {
				delete(r, u)
			}
		}
		e.s++
		e.p[e.s] = w
		u = w
	}
}

func newUlerian(g graph.Undirected) *ulerian {
	a := g.AdjacencyList
	u := &ulerian{
		g:  make([]map[graph.NI]int, len(a)),
		uv: bits.New(len(a)),
	}
	// convert representation, this maintains time complexity for
	// the undirected case.
	m2 := 0
	for n, to := range a {
		m2 += len(to)
		s := map[graph.NI]int{} // a multiset for each node
		for _, to := range to {
			s[to]++
			if to == graph.NI(n) {
				m2++
			}
		}
		u.g[n] = s
	}
	u.m = m2 / 2
	u.p = make([]graph.NI, u.m+1)
	u.uv.SetAll()
	return u
}
