// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// +build ignore

//go:generate go run Euclidean.go

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
	g, pos, err := graph.Euclidean(n, n*3, 12, 100, r)
	if err != nil {
		log.Fatal(err)
	}
	c := exec.Command("neato", "-Tsvg", "-o", "Euclidean.svg")
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
