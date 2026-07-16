package pathfinding_alg

import (
	"fmt"
	"math"
	"pathinder/internal/parsers"
)

type Graph struct {
	Level []int
}

func DinicAlg(network *parsers.NetworkData, start, end string) ([]string, error) {
	networkmap := network.NetworkMap
	Stationmap := network.StationMap
	// g := &Graph{
	// 	Level: make([]int, len(Stationmap)),
	// }

	pathList := []parsers.StationName{}
	max_flow := 0

	setToInf(network, start, end)
	fmt.Println(max_flow)
	fmt.Println(networkmap)
	fmt.Println(Stationmap)

	// for _, v := range Stationmap {
	for range 3 {
		path := BFSAlg2(network, start, end)
		if path == nil {
			break
		}
		for range path {

		}
		max_flow += 1
	}
	// }
	return pathList, nil
}

func setToInf(network *parsers.NetworkData, start, end string) {
	nodeStart := network.StationMap[start]
	nodeStart.Cap = math.MaxInt
	network.StationMap[start] = nodeStart

	nodeEnd := network.StationMap[end]
	nodeEnd.Cap = math.MaxInt
	network.StationMap[end] = nodeEnd
}
