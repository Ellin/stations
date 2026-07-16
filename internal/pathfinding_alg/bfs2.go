package pathfinding_alg

import (
	"container/list"
	"pathinder/internal/parsers"
	"slices"
)

func BFSAlg2(network *parsers.NetworkData, start, end string) []parsers.StationName {
	networkmap := network.NetworkMap
	Stationmap := network.StationMap

	visited := map[parsers.StationName]bool{start: true}
	backtrack := map[parsers.StationName]parsers.StationName{}

	queue := list.New()
	queue.PushBack(start)

	found := false
	for queue.Len() > 0 {
		firstNode := queue.Front()
		queue.Remove(firstNode)
		current := firstNode.Value.(parsers.StationName)
		// fmt.Println(current)

		if current == end {
			found = true
			break
		}
		for neighbor := range networkmap[current] {
			if !visited[neighbor] && Stationmap[neighbor].Cap > 0 {
				visited[neighbor] = true
				queue.PushBack(neighbor)
				backtrack[neighbor] = current
			}

		}
	}
	if !found {
		return nil
	}

	track := []parsers.StationName{}

	for from := end; from != start; from = backtrack[from] {
		track = append(track, from)
	}
	slices.Reverse(track)
	return track
}
