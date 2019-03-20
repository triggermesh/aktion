// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// +build ignore

// go generate in this directory generates images for the tutorials.

//go:generate go run al.go
//go:generate dia -t svg -e almem.svg almem.dia

package main

import (
	"log"
	"os/exec"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func main() {
	for _, f := range []func() error{al, al0, alpair, ald} {
		if err := f(); err != nil {
			log.Fatal(err)
		}
	}
}

func al() error {
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	c := exec.Command("dot", "-Tsvg", "-o", "al.svg")
	w, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Start()
	dot.Write(g, w, dot.GraphAttr("rankdir", "LR"))
	w.Close()
	return c.Wait()
}

func al0() error {
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	// al0 output file name
	c := exec.Command("dot", "-Tsvg", "-o", "al0.svg")
	w, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Start()
	// otherwise, only difference from al() is dot.Isolated
	dot.Write(g, w, dot.GraphAttr("rankdir", "LR"), dot.Isolated(true))
	w.Close()
	return c.Wait()
}

func alpair() error {
	g := graph.AdjacencyList{
		0: {1},
		1: {0},
	}
	c := exec.Command("dot", "-Tsvg", "-o", "alpair.svg")
	w, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Start()
	dot.Write(g, w, dot.GraphAttr("rankdir", "LR"))
	w.Close()
	return c.Wait()
}

func ald() error {
	g := graph.LabeledAdjacencyList{
		1: {{To: 2, Label: 7}, {To: 3, Label: 9}, {To: 6, Label: 11}},
		2: {{To: 3, Label: 10}, {To: 4, Label: 15}},
		3: {{To: 4, Label: 11}, {To: 6, Label: 2}},
		4: {{To: 5, Label: 7}},
		6: {{To: 5, Label: 9}},
	}
	c := exec.Command("dot", "-Tsvg", "-o", "ald.svg")
	w, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Start()
	dot.Write(g, w, dot.GraphAttr("rankdir", "LR"))
	w.Close()
	return c.Wait()
}
