// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/soniakeys/graph"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var chungLuSmall graph.Undirected
var chungLuSmallTag string

func ChungLuSmall() (string, int, int) {
	const n = 1e4
	w := make([]float64, n)
	for i := range w {
		w[i] = 5 + 10*float64(n-i)/float64(n)
	}
	chungLuSmall, m = graph.ChungLu(w, r)
	chungLuSmallTag = "ChungLu " + h(n) + " nds"
	return "Chung Lu (undirected)", n, m
}

var chungLuLarge graph.Undirected
var chungLuLargeTag string

func ChungLuLarge() (string, int, int) {
	const n = 2e5
	w := make([]float64, n)
	for i := range w {
		w[i] = 2 + 50*n/float64(i+1)
	}
	chungLuLarge, m = graph.ChungLu(w, r)
	chungLuLargeTag = "ChungLu " + h(n) + " nds"
	return "Chung Lu (undirected)", n, m
}

var eucSmall graph.LabeledDirected
var eucSmallTag string
var eucSmallWt []float64
var eucSmallWtFunc = func(n graph.LI) float64 { return eucSmallWt[n] }

func EucSmall() (string, int, int) {
	const n = 1024
	const ma = 5e3
	var err error
	eucSmall, _, eucSmallWt, err = graph.LabeledEuclidean(n, ma, 1, 1, r)
	if err != nil {
		return "nope", n, ma
	}
	eucSmallTag = "Euclidean " + h(n) + " nds"
	return "Euclidean (directed)", n, ma
}

var eucLarge graph.Directed
var eucLargeTag string

func EucLarge() (string, int, int) {
	const n = 1048576
	const ma = 5e6
	var err error
	eucLarge, _, err = graph.Euclidean(n, ma, 1, 1, r)
	if err != nil {
		return "nope", n, ma
	}
	eucLargeTag = "Euclidean " + h(n) + " nds"
	return "Euclidean (directed)", n, ma
}

var geoSmall graph.LabeledUndirected
var geoSmallPos []struct{ X, Y float64 }
var geoSmallWt []float64
var geoSmallWtFunc = func(n graph.LI) float64 { return geoSmallWt[n] }
var geoSmallTag string

func GeoSmall() (string, int, int) {
	const n = 1000
	geoSmall, geoSmallPos, geoSmallWt = graph.LabeledGeometric(n, .1, r)
	geoSmallTag = "Geometric " + h(n) + " nds"
	return "Geometric (undirected)", n, len(geoSmallWt)
}

var geoLarge graph.LabeledUndirected
var geoLargePos []struct{ X, Y float64 }
var geoLargeWt []float64
var geoLargeWtFunc = func(n graph.LI) float64 { return geoLargeWt[n] }
var geoLargeTag string

func GeoLarge() (string, int, int) {
	const n = 3e4
	geoLarge, geoLargePos, geoLargeWt = graph.LabeledGeometric(n, .01, r)
	geoLargeTag = "Geometric " + h(n) + " nds"
	return "Geometric (undirected)", n, len(geoLargeWt)
}

var gnpUSmall graph.Undirected
var gnpUSmallTag string

func GnpUSmall() (string, int, int) {
	const n = 1000
	gnpUSmallTag = fmt.Sprint("Gnp ", n, " nds")
	gnpUSmall, m = graph.GnpUndirected(n, .2, r)
	return "Gnp undirected", n, m
}

var gnpULarge graph.Undirected
var gnpULargeTag string

func GnpULarge() (string, int, int) {
	const n = 2e4
	gnpULargeTag = fmt.Sprint("Gnp ", n, " nds")
	gnpULarge, m = graph.GnpUndirected(n, .105, r)
	return "Gnp undirected", n, m
}

var gnmUSmall graph.Undirected
var gnmUSmallTag string

func GnmUSmall() (string, int, int) {
	const n = 1000
	const m = 100e3
	gnmUSmallTag = fmt.Sprint("Gnm ", n, " nds")
	gnmUSmall = graph.GnmUndirected(n, m, r)
	return "Gnm undirected", n, m
}

var gnmULarge graph.Undirected
var gnmULargeTag string

func GnmULarge() (string, int, int) {
	const n = 20e3
	const m = 20e6
	gnmULargeTag = fmt.Sprint("Gnm ", n, " nds")
	gnmULarge = graph.GnmUndirected(n, m, r)
	return "Gnm undirected", n, m
}

func Gnm3USmall() (string, int, int) {
	const n = 1000
	const m = 100e3
	graph.Gnm3Undirected(n, m, r)
	return "Gnm3 undirected", n, m
}

func Gnm3ULarge() (string, int, int) {
	const n = 20e3
	const m = 20e6
	graph.Gnm3Undirected(n, m, r)
	return "Gnm3 undirected", n, m
}

var gnpDSmall graph.Directed
var gnpDSmallTag string

func GnpDSmall() (string, int, int) {
	const n = 1000
	gnpDSmallTag = fmt.Sprint("Gnp ", n, " nds")
	gnpDSmall, ma = graph.GnpDirected(n, .101, r)
	return "Gnp directed", n, ma
}

var gnpDLarge graph.Directed
var gnpDLargeTag string

func GnpDLarge() (string, int, int) {
	const n = 2e4
	gnpDLargeTag = fmt.Sprint("Gnp ", n, " nds")
	gnpDLarge, ma = graph.GnpDirected(n, .05, r)
	return "Gnp directed", n, ma
}

var gnmDSmall graph.Directed
var gnmDSmallTag string

func GnmDSmall() (string, int, int) {
	const n = 1000
	const ma = 100e3
	gnmDSmallTag = fmt.Sprint("Gnm ", n, " nds")
	gnmDSmall = graph.GnmDirected(n, ma, r)
	return "Gnm directed", n, ma
}

var gnmDLarge graph.Directed
var gnmDLargeTag string

func GnmDLarge() (string, int, int) {
	const n = 20e3
	const ma = 20e6
	gnmDLargeTag = fmt.Sprint("Gnm ", n, " nds")
	gnmDLarge = graph.GnmDirected(n, ma, r)
	return "Gnm directed", n, ma
}

var kronDSmall graph.Directed
var kronDSmallTag string
var m, ma int

func KronDSmall() (string, int, int) {
	kronDSmall, ma = graph.KroneckerDirected(11, 7, r)
	kronDSmallTag = "Kronecker " + h(kronDSmall.Order()) + "nds"
	return "Kronecker directed", kronDSmall.Order(), ma
}

var kronDLarge graph.Directed
var kronDLargeTag string

func KronDLarge() (string, int, int) {
	kronDLarge, ma = graph.KroneckerDirected(17, 21, r)
	kronDLargeTag = "Kronecker " + h(kronDLarge.Order()) + "nds"
	return "Kronecker directed", kronDLarge.Order(), ma
}

var kronUSmall graph.Undirected
var kronUSmallTag string

func KronUSmall() (string, int, int) {
	kronUSmall, m = graph.KroneckerUndirected(11, 7, r)
	kronUSmallTag = "Kronecker " + h(kronUSmall.Order()) + "nds"
	return "Kronecker undirected", kronUSmall.Order(), m
}

var kronULarge graph.Undirected
var kronULargeTag string

func KronULarge() (string, int, int) {
	kronULarge, m = graph.KroneckerUndirected(17, 21, r)
	kronULargeTag = "Kronecker " + h(kronULarge.Order()) + " nds"
	return "Kronecker undirected", kronULarge.Order(), m
}
