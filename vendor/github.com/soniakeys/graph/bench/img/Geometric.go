// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// +build ignore

//go:generate go run Geometric.go

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	n := 100
	g, pos, _ := graph.Geometric(n, .2, r)
	c := exec.Command("neato", "-Tsvg", "-o", "Geometric.svg")
	w, err := c.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	dot.Write(g, w, dot.NodePos(func(n graph.NI) string {
		return fmt.Sprintf("%.3f,%.3f", 6*pos[n].X, 6*pos[n].Y)
	}))
	if err = w.Close(); err != nil {
		log.Fatal(err)
	}
	if err = c.Wait(); err != nil {
		log.Fatal(err)
	}
}
