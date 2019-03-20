// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

var chungLuSmallCCRep graph.NI
var chungLuSmallCCma int
var chungLuSmallCCTag string

func CCSmall() (string, string) {
	reps, orders, ma := chungLuSmall.ConnectedComponentReps()
	max := 0
	for i, o := range orders {
		if o > max {
			max = o
			chungLuSmallCCRep = reps[i]
			chungLuSmallCCma = ma[i]
		}
	}
	chungLuSmallCCTag = "ChungLu giant component " + h(max) + " nds"
	return "Connected Components", chungLuSmallTag
}

var chungLuLargeCCRep graph.NI
var chungLuLargeCCma int
var chungLuLargeCCTag string

func CCLarge() (string, string) {
	reps, orders, ma := chungLuLarge.ConnectedComponentReps()
	max := 0
	for i, o := range orders {
		if o > max {
			max = o
			chungLuLargeCCRep = reps[i]
			chungLuLargeCCma = ma[i]
		}
	}
	chungLuLargeCCTag = "ChungLu giant component " + h(max) + " nds"
	return "Connected Components", chungLuLargeTag
}

var eucSmallSCCRep graph.NI
var eucSmallSCCTag string

func SCCEucSmall() (string, string) {
	max := 0
	eucSmall.StronglyConnectedComponents(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			eucSmallSCCRep = c[0]
		}
		return true
	})
	eucSmallSCCTag = "Euclidean giant component " + h(max) + " nds"
	return "SCC (Pearce)", eucSmallTag
}

var eucLargeSCCRep graph.NI
var eucLargeSCCTag string

func SCCEucLarge() (string, string) {
	max := 0
	eucLarge.StronglyConnectedComponents(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			eucLargeSCCRep = c[0]
		}
		return true
	})
	eucLargeSCCTag = "Euclidean giant component " + h(max) + " nds"
	return "SCC (Pearce)", eucLargeTag
}

var kronDSmallSCCRep graph.NI
var kronDSmallSCCTag string

func SCCPearceSmall() (string, string) {
	max := 0
	kronDSmall.StronglyConnectedComponents(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			kronDSmallSCCRep = c[0]
		}
		return true
	})
	kronDSmallSCCTag = "Kronecker giant component " + h(max) + " nds"
	return "SCC (Pearce)", kronDSmallTag
}

func SCCPearceLarge() (string, string) {
	max := 0
	kronDLarge.StronglyConnectedComponents(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			kronDLargeSCCRep = c[0]
		}
		return true
	})
	kronDLargeSCCTag = "Kronecker giant component " + h(max) + " nds"
	return "SCC (Pearce)", kronDLargeTag
}

func SCCPathSmall() (string, string) {
	max := 0
	alt.SCCPathBased(kronDSmall, func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			kronDSmallSCCRep = c[0]
		}
		return true
	})
	kronDSmallSCCTag = "Kronecker giant component " + h(max) + " nds"
	return "SCC path based", kronDSmallTag
}

var kronDLargeSCCRep graph.NI
var kronDLargeSCCTag string

func SCCPathLarge() (string, string) {
	max := 0
	alt.SCCPathBased(kronDLarge, func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			kronDLargeSCCRep = c[0]
		}
		return true
	})
	kronDLargeSCCTag = "Kronecker giant component " + h(max) + " nds"
	return "SCC path based", kronDLargeTag
}

func SCCTarjanSmall() (string, string) {
	max := 0
	alt.SCCTarjan(kronDSmall, func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			kronDSmallSCCRep = c[0]
		}
		return true
	})
	kronDSmallSCCTag = "Kronecker giant component " + h(max) + " nds"
	return "SCC Tarjan", kronDSmallTag
}

func SCCTarjanLarge() (string, string) {
	max := 0
	alt.SCCTarjan(kronDLarge, func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			kronDLargeSCCRep = c[0]
		}
		return true
	})
	kronDLargeSCCTag = "Kronecker giant component " + h(max) + " nds"
	return "SCC Tarjan", kronDLargeTag
}
func SCCKosarajuSmall() (string, string) {
	max := 0
	alt.SCCKosaraju(kronDSmall, func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			kronDSmallSCCRep = c[0]
		}
		return true
	})
	kronDSmallSCCTag = "Kronecker giant component " + h(max) + " nds"
	return "SCC Kosaraju", kronDSmallTag
}

func SCCKosarajuLarge() (string, string) {
	max := 0
	alt.SCCKosaraju(kronDLarge, func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			kronDLargeSCCRep = c[0]
		}
		return true
	})
	kronDLargeSCCTag = "Kronecker giant component " + h(max) + " nds"
	return "SCC Kosaraju", kronDLargeTag
}
