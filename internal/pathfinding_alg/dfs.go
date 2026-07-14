package pathfinding_alg

import (
	"pathinder/internal/parsers"
	"slices"
)

func DFSAlg(network *parsers.NetworkData, start, end string) []parsers.StationName {
	visited := map[parsers.StationName]bool{}
	var route []parsers.StationName

	if route = dfs(network, visited, start, end); len(route) != 0 {
		slices.Reverse(route)
	}
	return route
}
func dfs(network *parsers.NetworkData, visited map[parsers.StationName]bool, start, end string) []parsers.StationName {
	nodemap := network.NetworkMap
	visited[start] = true
	if start == end {
		return []parsers.StationName{start}
	}
	for neighbor := range nodemap[start] {
		if !visited[neighbor] {
			if next := dfs(network, visited, neighbor, end); next != nil {
				return append(next, start)
			}

		}

	}

	return nil
}
