// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package treevis_test

import (
	"os"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/treevis"
)

func ExampleWrite() {
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2, 3},
		1: {4, 5},
		5: {6},
		6: {},
	}}
	treevis.Write(g, 0, os.Stdout)
	// Output:
	// ┐0
	// ├─┐1
	// │ ├─╴4
	// │ └─┐5
	// │   └─╴6
	// ├─╴2
	// └─╴3
}

func ExampleGlyphs() {
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2, 3},
		1: {4, 5},
		5: {6},
		6: {},
	}}
	treevis.Write(g, 0, os.Stdout, treevis.Glyphs(treevis.G{
		Leaf:      "-",
		NonLeaf:   "-",
		Child:     " |",
		Vertical:  " |",
		LastChild: " `",
		Indent:    "  ",
	}))
	// Output:
	// -0
	//  |-1
	//  | |-4
	//  | `-5
	//  |   `-6
	//  |-2
	//  `-3
}

func ExampleNodeLabel() {
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2, 3},
		1: {4, 5},
		5: {6},
		6: {},
	}}
	labels := []string{
		0: "root",
		1: "a",
		2: "b",
		3: "c",
		4: "d",
		5: "e",
		6: "f",
	}
	treevis.Write(g, 0, os.Stdout,
		treevis.NodeLabel(func(n graph.NI) string {
			return " " + labels[n]
		}))
	// Output:
	// ┐ root
	// ├─┐ a
	// │ ├─╴ d
	// │ └─┐ e
	// │   └─╴ f
	// ├─╴ b
	// └─╴ c
}
