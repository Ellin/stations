package paths

import (
	"container/list"
	"pathfinder/internal/parsers"
	"pathfinder/model"
	"slices"
)

func BFSAlg4(network *model.NetworkData, start, end string) []parsers.StationName {
	nodemap := network.NetworkMap

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
		for neighbor := range nodemap[current] {
			if !visited[neighbor] {
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
