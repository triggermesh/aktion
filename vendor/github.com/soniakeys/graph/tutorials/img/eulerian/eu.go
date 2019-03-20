// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// +build ignore

// go generate in this directory generates images for the tutorials.

//go:generate go run eu.go

package main

import (
	"log"
	"os/exec"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func main() {
	for _, f := range []func() error{cycle, cycleOrder} {
		if err := f(); err != nil {
			log.Fatal(err)
		}
	}
}

func cycle() error {
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 'a'}},
		1: {{3, 'b'}, {2, 'e'}},
		2: {{1, 'd'}, {3, 'f'}},
		3: {{2, 'c'}, {0, 'g'}},
	}}
	//	g := graph.Directed{graph.AdjacencyList{
	//		0: {1},
	//		1: {3, 2},
	//		2: {1, 3},
	//		3: {2, 0},
	//	}}
	c := exec.Command("dot", "-Tsvg", "-o", "cycle.svg")
	w, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Start()
	log.Print(dot.String(g,
		dot.GraphAttr("rankdir", "LR"),
		dot.EdgeLabel(func(l graph.LI) string { return "" }),
		dot.EdgeAttr(func(l graph.LI) []dot.AttrVal {
			return []dot.AttrVal{{"arrowhead", "none"}}
		})))
	dot.Write(g, w,
		dot.GraphAttr("rankdir", "LR"),
		dot.EdgeLabel(func(l graph.LI) string { return "" }),
		dot.EdgeAttr(func(l graph.LI) []dot.AttrVal {
			return []dot.AttrVal{{"arrowhead", "none"}}
		}))
	w.Close()
	return c.Wait()
}

func cycleOrder() error {
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{1, 'a'}},
		1: {{3, 'b'}, {2, 'e'}},
		2: {{1, 'd'}, {3, 'f'}},
		3: {{2, 'c'}, {0, 'g'}},
	}}
	c := exec.Command("dot", "-Tsvg", "-o", "cycleOrder.svg")
	w, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Start()
	dot.Write(g, w,
		dot.GraphAttr("rankdir", "LR"),
		dot.GraphAttr("start", "2"),
		dot.EdgeLabel(func(l graph.LI) string { return string(l) }))
	w.Close()
	return c.Wait()
}
