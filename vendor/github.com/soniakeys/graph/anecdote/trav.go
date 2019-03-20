// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func DFSmall() (string, string) {
	b := bits.New(chungLuSmall.Order())
	alt.DepthFirst(chungLuSmall.AdjacencyList,
		chungLuSmallCCRep, alt.Visited(&b))
	return "DepthFirst", chungLuSmallCCTag
}

func DFLarge() (string, string) {
	b := bits.New(chungLuLarge.Order())
	alt.DepthFirst(chungLuLarge.AdjacencyList,
		chungLuLargeCCRep, alt.Visited(&b))
	return "DepthFirst", chungLuLargeCCTag
}

func BFSmall() (string, string) {
	chungLuSmall.AdjacencyList.BreadthFirst(chungLuSmallCCRep,
		func(graph.NI) {})
	return "BreadthFirst", chungLuSmallCCTag
}

func BFLarge() (string, string) {
	chungLuLarge.AdjacencyList.BreadthFirst(chungLuLargeCCRep,
		func(graph.NI) {})
	return "BreadthFirst", chungLuLargeCCTag
}

func AltBFSmall() (string, string) {
	g := chungLuSmall.AdjacencyList
	alt.BreadthFirst2(g, g, chungLuSmallCCma, chungLuSmallCCRep, nil,
		func(graph.NI) bool { return true })
	return "DO BreadthFirst", chungLuSmallCCTag
}

func AltBFLarge() (string, string) {
	g := chungLuLarge.AdjacencyList
	alt.BreadthFirst2(g, g, chungLuLargeCCma, chungLuLargeCCRep, nil,
		func(graph.NI) bool { return true })
	return "DO BreadthFirst", chungLuLargeCCTag
}
