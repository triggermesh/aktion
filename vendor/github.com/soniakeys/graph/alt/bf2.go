// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt

import (
	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

// BreadthFirst2 traverses a graph breadth first using a direction
// optimizing algorithm pioneered by Scott Beamer.
//
// The algorithm is supposed to be faster than the conventional breadth first
// algorithm but I haven't seen it yet.
func BreadthFirst2(g, tr graph.AdjacencyList, ma int, start graph.NI, f *graph.FromList, v func(graph.NI) bool) int {
	if tr == nil {
		var d graph.Directed
		d, ma = graph.Directed{g}.Transpose()
		tr = d.AdjacencyList
	}
	switch {
	case f == nil:
		e := graph.NewFromList(len(g))
		f = &e
	case f.Paths == nil:
		*f = graph.NewFromList(len(g))
	}
	if ma <= 0 {
		ma = g.ArcSize()
	}
	rp := f.Paths
	level := 1
	rp[start] = graph.PathEnd{Len: level, From: -1}
	if !v(start) {
		f.MaxLen = level
		return -1
	}
	nReached := 1 // accumulated for a return value
	// the frontier consists of nodes all at the same level
	frontier := []graph.NI{start}
	mf := len(g[start])     // number of arcs leading out from frontier
	ctb := ma / 10          // threshold change from top-down to bottom-up
	k14 := 14 * ma / len(g) // 14 * mean degree
	cbt := len(g) / k14     // threshold change from bottom-up to top-down
	// well shoot.  part of the speed problem is the bitmap implementation.
	// big.Ints are slow.  this custom bits package is faster.  faster still
	// is just a slice of bools. :(
	fBits := bits.New(len(g))
	nextb := bits.New(len(g))
	for {
		// top down step
		level++
		var next []graph.NI
		for _, n := range frontier {
			for _, nb := range g[n] {
				if rp[nb].Len == 0 {
					rp[nb] = graph.PathEnd{From: n, Len: level}
					if !v(nb) {
						f.MaxLen = level
						return -1
					}
					next = append(next, nb)
					nReached++
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
		if mf > ctb {
			// switch to bottom up!
		} else {
			// stick with top down
			continue
		}
		// convert frontier representation
		nf := 0 // number of vertices on the frontier
		for _, n := range frontier {
			fBits.SetBit(int(n), 1)
			nf++
		}
	bottomUpLoop:
		level++
		nNext := 0
		for n := range tr {
			if rp[n].Len == 0 {
				for _, nb := range tr[n] {
					if fBits.Bit(int(nb)) == 1 {
						rp[n] = graph.PathEnd{From: nb, Len: level}
						if !v(nb) {
							f.MaxLen = level
							return -1
						}
						nextb.SetBit(n, 1)
						nReached++
						nNext++
						break
					}
				}
			}
		}
		if nNext == 0 {
			break
		}
		fBits, nextb = nextb, fBits
		nextb.ClearAll()
		nf = nNext
		if nf < cbt {
			// switch back to top down!
		} else {
			// stick with bottom up
			goto bottomUpLoop
		}
		// convert frontier representation
		mf = 0
		frontier = frontier[:0]
		// TODO use exported bits representation here.  quick search for a
		// non zero word, then convert all bits in the word, then clear the
		// word.
		for n := range g {
			if fBits.Bit(n) == 1 {
				frontier = append(frontier, graph.NI(n))
				mf += len(g[n])
				// fBits.SetBit(n, 0) // alternative to clear all.
			}
		}
		// or clear individual bits, but in bottom up step we expect a
		// relatively dense frontier, so ClearAll seems better.
		fBits.ClearAll()
	}
	f.MaxLen = level - 1
	return nReached
}
