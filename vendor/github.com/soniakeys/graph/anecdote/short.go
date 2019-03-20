// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

func FloydEuc() (string, string) {
	eucSmall.DistanceMatrix(eucSmallWtFunc).FloydWarshall()
	return "Floyd-Warshall", eucSmallTag
}

func FloydGeo() (string, string) {
	geoSmall.DistanceMatrix(geoSmallWtFunc).FloydWarshall()
	return "Floyd-Warshall", geoSmallTag
}

func BellmanSmall() (string, string) {
	eucSmall.BellmanFord(eucSmallWtFunc, eucSmallSCCRep)
	return "Bellman-Ford", eucSmallSCCTag
}

func DijkstraAllSmall() (string, string) {
	geoSmall.Dijkstra(0, -1, geoSmallWtFunc)
	return "Dijkstra all paths", geoSmallTag
}

func DijkstraAllLarge() (string, string) {
	geoLarge.Dijkstra(0, -1, geoLargeWtFunc)
	return "Dijkstra all paths", geoLargeTag
}

func Dijkstra1Small() (string, string) {
	geoSmall.Dijkstra(0, geoSmallEnd, geoSmallWtFunc)
	return "Dijkstra single path", geoSmallTag
}

func Dijkstra1Large() (string, string) {
	geoLarge.Dijkstra(0, geoLargeEnd, geoLargeWtFunc)
	return "Dijkstra single path", geoLargeTag
}

func AStarASmall() (string, string) {
	geoSmall.AStarA(geoSmallWtFunc, 0, geoSmallEnd, geoSmallHeuristic)
	return "AStarA", geoSmallTag
}

func AStarALarge() (string, string) {
	geoLarge.AStarA(geoLargeWtFunc, 0, geoLargeEnd, geoLargeHeuristic)
	return "AStarA", geoLargeTag
}

func AStarMSmall() (string, string) {
	geoSmall.AStarM(geoSmallWtFunc, 0, geoSmallEnd, geoSmallHeuristic)
	return "AStarM", geoSmallTag
}

func AStarMLarge() (string, string) {
	geoLarge.AStarM(geoLargeWtFunc, 0, geoLargeEnd, geoLargeHeuristic)
	return "AStarM", geoLargeTag
}
