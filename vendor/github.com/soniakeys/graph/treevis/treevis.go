// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Treevis draws trees with text.  Capabilities are currently quite limited.
package treevis

import (
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/soniakeys/graph"
)

type G struct {
	Leaf      string
	NonLeaf   string
	Child     string
	Vertical  string
	LastChild string
	Indent    string
}

type Config struct {
	NodeLabel func(graph.NI) string
	Glyphs    G
}

var Defaults = Config{
	NodeLabel: func(n graph.NI) string { return strconv.Itoa(int(n)) },
	Glyphs: G{
		Leaf:      "╴",
		NonLeaf:   "┐",
		Child:     "├─",
		Vertical:  "│ ",
		LastChild: "└─",
		Indent:    "  ",
	},
}

type Option func(*Config)

func NodeLabel(f func(graph.NI) string) Option {
	return func(c *Config) { c.NodeLabel = f }
}

func Glyphs(g G) Option {
	return func(c *Config) { c.Glyphs = g }
}

func Write(g graph.Directed, root graph.NI, w io.Writer, options ...Option) (err error) {
	cf := Defaults
	for _, o := range options {
		o(&cf)
	}
	var vis big.Int
	var f func(graph.NI, string) bool
	f = func(n graph.NI, pre string) bool {
		if vis.Bit(int(n)) != 0 {
			fmt.Fprintln(w, "%!(NONTREE)")
			err = fmt.Errorf("non-tree")
			return false
		}
		vis.SetBit(&vis, int(n), 1)
		to := g.AdjacencyList[n]
		if len(to) == 0 {
			_, err = fmt.Fprint(w, cf.Glyphs.Leaf, cf.NodeLabel(n), "\n")
			return err == nil
		}
		_, err = fmt.Fprint(w, cf.Glyphs.NonLeaf, cf.NodeLabel(n), "\n")
		if err != nil {
			return false
		}
		last := len(to) - 1
		for _, to := range to[:last] {
			if _, err = fmt.Fprint(w, pre, cf.Glyphs.Child); err != nil {
				return false
			}
			if !f(to, pre+cf.Glyphs.Vertical) {
				return false
			}
		}
		if _, err = fmt.Fprint(w, pre, cf.Glyphs.LastChild); err != nil {
			return false
		}
		return f(to[last], pre+cf.Glyphs.Indent)
	}
	f(root, "")
	return
}
