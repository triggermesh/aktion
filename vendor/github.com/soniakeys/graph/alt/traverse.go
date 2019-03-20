// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt

import (
	"math/rand"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

type config struct {
	start         graph.NI
	arcVisitor    func(n graph.NI, x int)
	iterateFrom   func(n graph.NI)
	nodeVisitor   func(n graph.NI)
	okArcVisitor  func(n graph.NI, x int) bool
	okNodeVisitor func(n graph.NI) bool
	rand          *rand.Rand
	visBits       *bits.Bits
	pathBits      *bits.Bits
	fromList      *graph.FromList
}

// A TraverseOption specifies an option for a breadth first or depth first
// traversal.
//
// Values of this type are returned by various TraverseOption constructor
// functions.  These constructors take optional values to be used
// in a traversal and wrap them in the TraverseOption type.  This type
// is actually a function.  The BreadthFirst and DepthFirst traversal
// methods call these functions in order, to initialize state that controls
// the traversal.
type TraverseOption func(*config)

// ArcVisitor specifies a visitor function to call at each arc.
//
// See also OkArcVisitor.
func ArcVisitor(v func(n graph.NI, x int)) TraverseOption {
	return func(c *config) {
		c.arcVisitor = v
	}
}

// From specifies a graph.FromList to populate.
func From(f *graph.FromList) TraverseOption {
	return func(c *config) {
		c.fromList = f
	}
}

// NodeVisitor specifies a visitor function to call at each node.
//
// The node visitor function is called before any arc visitor functions.
//
// See also OkNodeVisitor.
func NodeVisitor(v func(graph.NI)) TraverseOption {
	return func(c *config) {
		c.nodeVisitor = v
	}
}

// OkArcVisitor specifies a visitor function to perform some test at each arc
// and return a boolean result.
//
// As long as v return a result of true, the traverse progresses to traverse all
// arcs.
//
// If v returns false, the traverse terminates immediately.
//
// See also ArcVisitor.
func OkArcVisitor(v func(n graph.NI, x int) bool) TraverseOption {
	return func(c *config) {
		c.okArcVisitor = v
	}
}

// OkNodeVisitor specifies a visitor function to perform some test at each node
// and return a boolean result.
//
// As long as v returns a result of true, the traverse progresses to traverse
// all nodes.  If v returns false, the traverse terminates immediately.
//
// The node visitor function is called before any arc visitor functions.
//
// See also NodeVisitor.
func OkNodeVisitor(v func(graph.NI) bool) TraverseOption {
	return func(c *config) {
		c.okNodeVisitor = v
	}
}

// PathBits specifies a bits.Bits value for nodes of the path to the
// currently visited node.
//
// A use for PathBits is identifying back arcs in a traverse.
//
// Unlike Visited, PathBits are zeroed at the start of a traverse.
func PathBits(b *bits.Bits) TraverseOption {
	return func(c *config) { c.pathBits = b }
}

// Rand specifies to traverse edges from each visited node in random order.
func Rand(r *rand.Rand) TraverseOption {
	return func(c *config) { c.rand = r }
}

// Visited specifies a bits.Bits value to record visited nodes.
//
// For each node visited, the corresponding bit is set to 1.  Other bits
// are not modified.
//
// The traverse algorithm controls the traverse using a bits.Bits.  If this
// function is used, argument b will be used as the controlling value.
//
// Bits are not zeroed at the start of a traverse, so the initial Bits value
// passed in should generally be zero.  Non-zero bits will limit the traverse.
func Visited(b *bits.Bits) TraverseOption {
	return func(c *config) { c.visBits = b }
}

// BreadthFirst traverses a directed or undirected graph in breadth first order.
//
// Argument start is the start node for the traversal.  Argument opt can be
// any number of values returned by a supported TraverseOption function.
//
// Supported:
//
//   From
//   NodeVisitor
//   OkNodeVisitor
//   Rand
//
// Unsupported:
//
//   ArcVisitor
//   OkArcVisitor
//   Visited
//   PathBits
//
// See also alt.BreadthFirst2, a direction optimizing breadth first algorithm.
func BreadthFirst(g graph.AdjacencyList, start graph.NI, opt ...TraverseOption) {
	cf := &config{start: start}
	for _, o := range opt {
		o(cf)
	}
	f := cf.fromList
	switch {
	case f == nil:
		e := graph.NewFromList(len(g))
		f = &e
	case f.Paths == nil:
		*f = graph.NewFromList(len(g))
	}
	rp := f.Paths
	// the frontier consists of nodes all at the same level
	frontier := []graph.NI{cf.start}
	level := 1
	// assign path when node is put on frontier
	rp[cf.start] = graph.PathEnd{Len: level, From: -1}
	for {
		f.MaxLen = level
		level++
		var next []graph.NI
		if cf.rand == nil {
			for _, n := range frontier {
				// visit nodes as they come off frontier
				if cf.nodeVisitor != nil {
					cf.nodeVisitor(n)
				}
				if cf.okNodeVisitor != nil {
					if !cf.okNodeVisitor(n) {
						return
					}
				}
				for _, nb := range g[n] {
					if rp[nb].Len == 0 {
						next = append(next, nb)
						rp[nb] = graph.PathEnd{From: n, Len: level}
					}
				}
			}
		} else { // take nodes off frontier at random
			for _, i := range cf.rand.Perm(len(frontier)) {
				n := frontier[i]
				// remainder of block same as above
				if cf.nodeVisitor != nil {
					cf.nodeVisitor(n)
				}
				if cf.okNodeVisitor != nil {
					if !cf.okNodeVisitor(n) {
						return
					}
				}
				for _, nb := range g[n] {
					if rp[nb].Len == 0 {
						next = append(next, nb)
						rp[nb] = graph.PathEnd{From: n, Len: level}
					}
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
	}
}

// DepthFirst traverses a directed or undirected graph in depth first order.
//
// Argument start is the start node for the traversal.  Argument opt can be
// any number of values returned by a supported TraverseOption function.
//
// Supported:
//
//   NodeVisitor
//   OkNodeVisitor
//   ArcVisitor
//   OkArcVisitor
//   Visited
//   PathBits
//   Rand
//
// Unsupported:
//
//   From
func DepthFirst(g graph.AdjacencyList, start graph.NI, options ...TraverseOption) {
	cf := &config{start: start}
	for _, o := range options {
		o(cf)
	}
	b := cf.visBits
	if b == nil {
		n := bits.New(len(g))
		b = &n
	} else if b.Bit(int(cf.start)) != 0 {
		return
	}
	if cf.pathBits != nil {
		cf.pathBits.ClearAll()
	}
	var df func(graph.NI) bool
	df = func(n graph.NI) bool {
		b.SetBit(int(n), 1)
		if cf.pathBits != nil {
			cf.pathBits.SetBit(int(n), 1)
		}

		if cf.nodeVisitor != nil {
			cf.nodeVisitor(n)
		}
		if cf.okNodeVisitor != nil {
			if !cf.okNodeVisitor(n) {
				return false
			}
		}

		if cf.rand == nil {
			for x, to := range g[n] {
				if cf.arcVisitor != nil {
					cf.arcVisitor(n, x)
				}
				if cf.okArcVisitor != nil {
					if !cf.okArcVisitor(n, x) {
						return false
					}
				}
				if b.Bit(int(to)) != 0 {
					continue
				}
				if !df(to) {
					return false
				}
			}
		} else {
			to := g[n]
			for _, x := range cf.rand.Perm(len(to)) {
				if cf.arcVisitor != nil {
					cf.arcVisitor(n, x)
				}
				if cf.okArcVisitor != nil {
					if !cf.okArcVisitor(n, x) {
						return false
					}
				}
				if b.Bit(int(to[x])) != 0 {
					continue
				}
				if !df(to[x]) {
					return false
				}
			}
		}
		if cf.pathBits != nil {
			cf.pathBits.SetBit(int(n), 0)
		}
		return true
	}
	df(cf.start)
}
