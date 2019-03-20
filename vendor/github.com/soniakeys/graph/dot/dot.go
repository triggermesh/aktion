// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// Package dot writes graphs from package graph in the Graphviz dot format.
//
// This package provides a minimal capability to output graphs simply and
// efficiently.
//
// There is no goal to provide a rich API to the many capabilities of the
// dot format.  Someday, maybe, another package.  Not currently.
//
// The scheme
//
// The dot package is a separate package from graph.  It imports graph;
// graph knows nothing of dot.  This keeps the graph package uncluttered by
// file format specific code.  Dot functions are functions then, not methods
// of graph representations.
//
// The function Write() takes any type of graph, an io.Writer, and optional
// arguments that control the output.  For convenience, there is also a
// String function that does not require an io.Writer and simply returns the
// dot format as a string.
//
// Optional arguments are variadic and constructed by calls to configuration
// functions defined in this package.  Not all configuration functions are
// meaningful for all graph types.  When a Write or String function is called
// it (1) initializes a Config struct from the package variable Defaults,
// then (2) in some cases initializes some members according to the graph type,
// then (3) calls the option functions in order.  Each option function can
// modify the Config struct.  After processing options, the funcion generates
// a dot file using the options specified in the Config struct.
package dot

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

// String generates a dot format string for a graph.
//
// g may be any of:
//
//   AdjacencyList
//   Directed
//   Undirected
//   LabeledAdjacencyList
//   LabeledDirected
//   LabeledUndirected
//   FromList
//   WeightedEdgeList
//
// or a pointer to any of these types.
//
// See also Write().
func String(g interface{}, options ...Option) (string, error) {
	var b bytes.Buffer
	if err := Write(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

// Write writes dot format text for a graph to an io.Writer.
//
// g may be any of:
//
//   AdjacencyList
//   Directed
//   Undirected
//   LabeledAdjacencyList
//   LabeledDirected
//   LabeledUndirected
//   FromList
//   WeightedEdgeList
//
// or a pointer to any of these types.
//
// When g is an undirected graph type, Config.Directed is initialized to
// false.
//
// The WeightedEdgeList, as used by the Kruskal methods, is a bit strange
// in that Kruskal interprets it as an undirected graph, but does not require
// that reciprocal edges be present.  Depending on how you construct a
// WeightedEdgeList, you may or may not have reciprocal edges.  If you do
// have reciprocal edges, the Directed(false) option is appropriate for
// collapsing reciprocals as usual and writing an undirected dot file.
// If, for Kruskal, for example, you constructed a WeightedEdgeList without
// reciprocals, then the UndirectArcs(true) is appropriate for writing an
// undirected dot file.  Specifying neither option and using the default of
// Directed(true) will produce a directed dot file.
//
// See also String().
func Write(g interface{}, w io.Writer, options ...Option) error {
	switch t := g.(type) {
	case graph.AdjacencyList:
		return writeAdjacencyList(t, w, options)
	case *graph.AdjacencyList:
		return writeAdjacencyList(*t, w, options)
	case graph.Directed:
		return writeAdjacencyList(t.AdjacencyList, w, options)
	case *graph.Directed:
		return writeAdjacencyList(t.AdjacencyList, w, options)
	case graph.Undirected:
		return writeUndirected(t.AdjacencyList, w, options)
	case *graph.Undirected:
		return writeUndirected(t.AdjacencyList, w, options)
	case graph.LabeledAdjacencyList:
		return writeLabeledAdjacencyList(t, w, options)
	case *graph.LabeledAdjacencyList:
		return writeLabeledAdjacencyList(*t, w, options)
	case graph.LabeledDirected:
		return writeLabeledAdjacencyList(t.LabeledAdjacencyList, w, options)
	case *graph.LabeledDirected:
		return writeLabeledAdjacencyList(t.LabeledAdjacencyList, w, options)
	case graph.LabeledUndirected:
		return writeLabeledUndirected(t.LabeledAdjacencyList, w, options)
	case *graph.LabeledUndirected:
		return writeLabeledUndirected(t.LabeledAdjacencyList, w, options)
	case graph.FromList:
		return writeFromList(t, w, options)
	case *graph.FromList:
		return writeFromList(*t, w, options)
	case graph.WeightedEdgeList:
		return writeWeightedEdgeList(t, w, options)
	case *graph.WeightedEdgeList:
		return writeWeightedEdgeList(*t, w, options)
	}
	return fmt.Errorf("dot: unknown graph type")
}

func writeAdjacencyList(g graph.AdjacencyList, w io.Writer, options []Option) error {
	cf := Defaults
	for _, o := range options {
		o(&cf)
	}
	return writeAL(g, w, &cf)
}

func writeUndirected(g graph.AdjacencyList, w io.Writer, options []Option) error {
	cf := Defaults
	cf.Directed = false
	for _, o := range options {
		o(&cf)
	}
	return writeAL(g, w, &cf)
}

func writeAL(g graph.AdjacencyList, w io.Writer, cf *Config) (err error) {
	b := bufio.NewWriter(w)
	if err = writeHead(cf, b); err != nil {
		return
	}
	if cf.NodePos != nil {
		_, err = fmt.Fprint(b, cf.Indent, "node [shape=point]\n")
		for n := range g {
			_, err = fmt.Fprintf(b, "%s%d [pos=\"%s!\"]\n",
				cf.Indent, n, cf.NodePos(graph.NI(n)))
			if err != nil {
				return
			}
		}
	}
	var iso bits.Bits
	if cf.Isolated {
		iso = g.IsolatedNodes()
		if iso.AllZeros() {
			cf.Isolated = false // optimization. turn off checking
		}
	}
	wf := writeALUndirected
	if cf.Directed {
		wf = writeALDirected
	}
	if err := wf(g, cf, iso, b); err != nil {
		return err
	}
	return writeTail(b)
}

func writeHead(cf *Config, b *bufio.Writer) error {
	t := "graph"
	if cf.Directed {
		t = "digraph"
	}
	if _, err := fmt.Fprintf(b, "%s {\n", t); err != nil {
		return err
	}
	for _, av := range cf.GraphAttr {
		_, err := fmt.Fprintf(b, "%s%s = %s\n", cf.Indent, av.Attr, av.Val)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeTail(b *bufio.Writer) error {
	if err := b.WriteByte('}'); err != nil {
		return err
	}
	return b.Flush()
}

func writeALDirected(g graph.AdjacencyList, cf *Config, iso bits.Bits, b *bufio.Writer) error {
	for fr, to := range g {
		err := writeALEdgeStmt(graph.NI(fr), to, "->", cf, iso, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeALEdgeStmt(fr graph.NI, to []graph.NI, op string, cf *Config, iso bits.Bits, b *bufio.Writer) (err error) {
	attr := ""
	if cf.EdgeAttr != nil {
		attr = " " + fmtAttr(cf.EdgeAttr(0))
	}
	if len(to) == 0 { // fast path
		if cf.Isolated && iso.Bit(int(fr)) == 1 {
			_, err = fmt.Fprintf(b, "%s%s\n", cf.Indent, cf.NodeID(fr))
		}
		return
	}
	if len(to) == 1 { // fast path
		_, err = fmt.Fprintf(b, "%s%s %s %s%s\n",
			cf.Indent, cf.NodeID(fr), op, cf.NodeID(to[0]), attr)
		return
	}
	// otherwise it's complicated.  we like to use a subgraph rhs to keep
	// output compact, but graphviz (some version) won't separate parallel
	// arcs in a subgraph, so in that case we write multiple edge statements.
	_, err = fmt.Fprintf(b, "%s%s %s ",
		cf.Indent, cf.NodeID(fr), op)
	if err != nil {
		return
	}
	s1 := map[graph.NI]bool{}
	m := map[graph.NI]int{} // multiset of defered duplicates
	c := "{"
	// first pass is over the to-list, the slice
	for _, to := range to {
		if !s1[to] {
			if _, err = b.WriteString(c + cf.NodeID(to)); err != nil {
				return
			}
			c = " "
			s1[to] = true
		} else {
			m[to]++
		}
	}
	if _, err = b.WriteString("}" + attr + "\n"); err != nil {
		return
	}
	// make additional passes over the map until it's fully consumed
	for len(m) > 0 {
		_, err = fmt.Fprintf(b, "%s%s %s ",
			cf.Indent, cf.NodeID(graph.NI(fr)), op)
		if err != nil {
			return
		}
		c1 := "{"
		for n, c := range m {
			if _, err = b.WriteString(c1 + cf.NodeID(n)); err != nil {
				return
			}
			if c == 1 {
				delete(m, n)
			} else {
				m[n]--
			}
			c1 = " "
		}
		if _, err = b.WriteString("}" + attr + "\n"); err != nil {
			return
		}
	}
	return
}

func writeALUndirected(g graph.AdjacencyList, cf *Config, iso bits.Bits, b *bufio.Writer) error {
	// Similar code in undir.go at IsUndirected
	unpaired := make(graph.AdjacencyList, len(g))
	for fr, to := range g {
		// first collect unpaired subset of to
		var uto []graph.NI
	arc: // for each arc in g
		for _, to := range to {
			if to == graph.NI(fr) {
				uto = append(uto, to) // loop
				continue
			}
			// search unpaired arcs
			ut := unpaired[to]
			for i, u := range ut {
				if u == graph.NI(fr) { // found reciprocal
					last := len(ut) - 1
					ut[i] = ut[last]
					unpaired[to] = ut[:last]
					continue arc
				}
			}
			// reciprocal not found
			uto = append(uto, to)
			unpaired[fr] = append(unpaired[fr], to)
		}
		err := writeALEdgeStmt(graph.NI(fr), uto, "--", cf, iso, b)
		if err != nil {
			return err
		}
	}
	for _, to := range unpaired {
		if len(to) > 0 {
			return fmt.Errorf("directed graph")
		}
	}
	return nil
}

func writeLabeledAdjacencyList(g graph.LabeledAdjacencyList, w io.Writer, options []Option) error {
	cf := Defaults
	for _, o := range options {
		o(&cf)
	}
	return writeLAL(g, w, &cf)
}

func writeLabeledUndirected(g graph.LabeledAdjacencyList, w io.Writer, options []Option) error {
	cf := Defaults
	cf.Directed = false
	for _, o := range options {
		o(&cf)
	}
	return writeLAL(g, w, &cf)
}

func writeLAL(g graph.LabeledAdjacencyList, w io.Writer, cf *Config) (err error) {
	b := bufio.NewWriter(w)
	if err = writeHead(cf, b); err != nil {
		return
	}
	if cf.NodePos != nil {
		_, err = fmt.Fprint(b, cf.Indent, "node [shape=point]\n")
		for n := range g {
			_, err = fmt.Fprintf(b, "%s%d [pos=\"%s!\"]\n",
				cf.Indent, n, cf.NodePos(graph.NI(n)))
			if err != nil {
				return
			}
		}
	}
	var iso bits.Bits
	if cf.Isolated {
		iso = g.IsolatedNodes()
		if iso.AllZeros() {
			cf.Isolated = false // optimization. turn off checking
		}
	}
	wf := writeLALUndirected
	if cf.Directed {
		wf = writeLALDirected
	}
	if err = wf(g, cf, iso, b); err != nil {
		return
	}
	return writeTail(b)
}

func writeLALDirected(g graph.LabeledAdjacencyList, cf *Config, iso bits.Bits, b *bufio.Writer) error {
	for fr, to := range g {
		err := writeLALEdgeStmt(graph.NI(fr), to, "->", cf, iso, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeLALEdgeStmt(fr graph.NI, to []graph.Half, op string, cf *Config, iso bits.Bits, b *bufio.Writer) (err error) {
	if len(to) == 0 {
		if cf.Isolated && iso.Bit(int(fr)) == 1 {
			_, err = fmt.Fprintf(b, "%s%s\n",
				cf.Indent, cf.NodeID(fr))
		}
		return
	}
	for _, to := range to {
		var attr []AttrVal
		if cf.EdgeAttr != nil {
			attr = cf.EdgeAttr(to.Label)
		}
		if el := cf.EdgeLabel(to.Label); el > "" {
			attr = append(attr, AttrVal{"label", cf.EdgeLabel(to.Label)})
		}
		_, err = fmt.Fprintf(b, "%s%s %s %s %s\n",
			cf.Indent, cf.NodeID(fr), op, cf.NodeID(to.To),
			fmtAttr(attr))
		if err != nil {
			return
		}
	}
	return
}

func fmtAttr(attr []AttrVal) string {
	if len(attr) == 0 {
		return ""
	}
	f := fmt.Sprintf("[%s = %s", attr[0].Attr, attr[0].Val)
	for _, a := range attr[1:] {
		f = fmt.Sprintf("%s, %s = %s", f, a.Attr, a.Val)
	}
	return f + "]"
}

func writeLALUndirected(g graph.LabeledAdjacencyList, cf *Config, iso bits.Bits, b *bufio.Writer) error {
	// Similar code in undir.go at IsUndirected
	unpaired := make(graph.LabeledAdjacencyList, len(g))
	for fr, to := range g {
		// first collect unpaired subset of to
		var uto []graph.Half
	arc: // for each arc in g
		for _, to := range to {
			if to.To == graph.NI(fr) {
				uto = append(uto, to) // loop
				continue
			}
			// search unpaired arcs
			ut := unpaired[to.To]
			for i, u := range ut {
				if u.To == graph.NI(fr) && u.Label == to.Label { // found reciprocal
					last := len(ut) - 1
					ut[i] = ut[last]
					unpaired[to.To] = ut[:last]
					continue arc
				}
			}
			// reciprocal not found
			uto = append(uto, to)
			unpaired[fr] = append(unpaired[fr], to)
		}
		err := writeLALEdgeStmt(graph.NI(fr), uto, "--", cf, iso, b)
		if err != nil {
			return err
		}
	}
	for _, to := range unpaired {
		if len(to) > 0 {
			return fmt.Errorf("directed graph")
		}
	}
	return nil
}

func writeFromList(f graph.FromList, w io.Writer, options []Option) error {
	cf := Defaults
	GraphAttr("rankdir", "BT")(&cf)
	for _, o := range options {
		o(&cf)
	}
	b := bufio.NewWriter(w)
	if err := writeHead(&cf, b); err != nil {
		return err
	}
	//var iso bits.Bits
	//if cf.Isolated {
	iso := f.IsolatedNodes()
	//}
	for i, e := range f.Paths {
		n := graph.NI(i)
		fr := e.From
		if fr < 0 {
			if cf.Isolated && iso.Bit(int(n)) != 0 {
				_, err := fmt.Fprintln(b, cf.Indent+cf.NodeID(graph.NI(fr)))
				if err != nil {
					return err
				}
			}
			continue
		}
		_, err := fmt.Fprintf(b, "%s%s -> %s\n",
			cf.Indent, cf.NodeID(n), cf.NodeID(fr))
		if err != nil {
			return err
		}
	}
	// repurpose iso for ranked same leaves.
	// leaves are ranked same if they not isolated nodes and there are
	// at least two of them.
	iso.AndNot(f.Leaves, iso)
	if !iso.AllZeros() && !iso.Single() { // rank:
		if _, err := b.WriteString(cf.Indent + "{rank = same"); err != nil {
			return err
		}
		for n := iso.OneFrom(0); n >= 0; n = iso.OneFrom(n + 1) {
			if _, err := b.WriteString(" "); err != nil {
				return err
			}
			if _, err := b.WriteString(cf.NodeID(graph.NI(n))); err != nil {
				return err
			}
		}
		if _, err := b.WriteString("}\n"); err != nil {
			return err
		}
	}
	return writeTail(b)
}

func writeWeightedEdgeList(g graph.WeightedEdgeList, w io.Writer, options []Option) error {
	cf := Defaults
	cf.Directed = false
	cf.EdgeLabel = func(l graph.LI) string {
		return fmt.Sprintf(`"%g"`, g.WeightFunc(l))
	}
	for _, o := range options {
		o(&cf)
	}
	if cf.UndirectArcs {
		cf.Directed = false
	}
	b := bufio.NewWriter(w)
	if err := writeHead(&cf, b); err != nil {
		return err
	}
	wf := writeWELNoRecip
	if cf.UndirectArcs || cf.Directed {
		wf = writeWELAllArcs
	}
	if err := wf(g, &cf, b); err != nil {
		return err
	}
	return writeTail(b)
}

func writeWELNoRecip(g graph.WeightedEdgeList, cf *Config, b *bufio.Writer) error {
	unpaired := make(graph.LabeledAdjacencyList, g.Order)
edge:
	for _, e := range g.Edges {
		// search unpaired arcs
		u2 := unpaired[e.N2]
		for i, u := range u2 {
			if u.To == e.N1 && u.Label == e.LI { // found reciprocal
				// write the edge
				_, err := fmt.Fprintf(b, "%s%s -- %s [label = %s]\n",
					cf.Indent, cf.NodeID(e.N2), cf.NodeID(e.N1),
					cf.EdgeLabel(e.LI))
				if err != nil {
					return err
				}
				// delete reciprocal
				last := len(u2) - 1
				u2[i] = u2[last]
				unpaired[e.N2] = u2[:last]
				continue edge
			}
		}
		// reciprocal not found
		unpaired[e.N1] = append(unpaired[e.N1], graph.Half{e.N2, e.LI})
	}
	for _, to := range unpaired {
		if len(to) > 0 {
			return fmt.Errorf("directed graph")
		}
	}
	return nil
}

func writeWELAllArcs(g graph.WeightedEdgeList, cf *Config, b *bufio.Writer) error {
	op := "--"
	if cf.Directed {
		op = "->"
	}
	for _, e := range g.Edges {
		_, err := fmt.Fprintf(b, "%s%s %s %s [label = %s]\n",
			cf.Indent, cf.NodeID(e.N1), op, cf.NodeID(e.N2),
			cf.EdgeLabel(e.LI))
		if err != nil {
			return err
		}
	}
	return nil
}
