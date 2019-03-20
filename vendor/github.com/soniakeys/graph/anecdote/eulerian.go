// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Anecdotal timings
package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

type dirEuGraph struct {
	g   graph.Directed
	ma  int
	tag string
}

// generate directed graph with n nodes, 10*n arcs
func dirEu(n int) dirEuGraph {
	unvis := make([]int, n-1)
	for i := range unvis {
		unvis[i] = i + 1
	}
	ma := n * 10
	a := make(graph.AdjacencyList, n)
	var fr, to int
	for i := 1; i < ma; i++ {
		// draw either from unvis or all nodes, p of drawing from unvis must
		// be 1 when (ma-i) == len(unvis).  p of drawing from unvis is 0 at
		// first.
		if len(unvis) == 0 || i < rand.Intn(ma-len(unvis)) {
			to = rand.Intn(n)
		} else {
			x := rand.Intn(len(unvis))
			to = unvis[x]
			last := len(unvis) - 1
			unvis[x] = unvis[last]
			unvis = unvis[:last]
		}
		a[fr] = append(a[fr], graph.NI(to))
		fr = to
	}
	a[fr] = append(a[fr], 0)
	return dirEuGraph{
		graph.Directed{a},
		ma,
		"directed, " + h(ma) + " arcs",
	}
}

type uEuGraph struct {
	g   graph.Undirected
	m   int
	tag string
}

// generate undirected graph with n nodes, 10*n edges
func uEu(n int) uEuGraph {
	unvis := make([]int, n-1)
	for i := range unvis {
		unvis[i] = i + 1
	}
	m := n * 10
	var u graph.Undirected
	var n1, n2 int
	for i := 1; i < m; i++ {
		// draw either from unvis or all nodes, p of drawing from unvis must
		// be 1 when (ma-i) == len(unvis).  p of drawing from unvis is 0 at
		// first.
		if len(unvis) == 0 || i < rand.Intn(m-len(unvis)) {
			n2 = rand.Intn(n)
		} else {
			x := rand.Intn(len(unvis))
			n2 = unvis[x]
			last := len(unvis) - 1
			unvis[x] = unvis[last]
			unvis = unvis[:last]
		}
		u.AddEdge(graph.NI(n1), graph.NI(n2))
		n1 = n2
	}
	u.AddEdge(graph.NI(n1), 0)
	return uEuGraph{u, m, "undirected, " + h(m) + " edges"}
}

type euResult struct {
	method string
	tag    string
	d      time.Duration
}

func dirEuTest(g dirEuGraph) euResult {
	t := time.Now()
	_, err := g.g.EulerianCycle()
	d := time.Now().Sub(t)
	if err != nil {
		log.Fatal(err)
	}
	return euResult{"EulerianCycle", g.tag, d}
}

func dirEuDTest(g dirEuGraph) euResult {
	t := time.Now()
	_, err := g.g.EulerianCycleD(g.ma)
	d := time.Now().Sub(t)
	if err != nil {
		log.Fatal(err)
	}
	return euResult{"EulerianCycleD", g.tag, d}
}

func uEuTest(g uEuGraph) euResult {
	t := time.Now()
	_, err := alt.EulerianCycle(g.g)
	d := time.Now().Sub(t)
	if err != nil {
		log.Fatal(err)
	}
	return euResult{"Alt EulerianCycle", g.tag, d}
}

func uEuDTest(g uEuGraph) euResult {
	t := time.Now()
	_, err := g.g.EulerianCycleD(g.m)
	d := time.Now().Sub(t)
	if err != nil {
		log.Fatal(err)
	}
	return euResult{"EulerianCycleD", g.tag, d}
}
