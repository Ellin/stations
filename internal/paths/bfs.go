package paths

import (
	"fmt"
	"slices"
)

func Bfs(networkMap map[StationName]map[StationName]struct{}, start string, end string) []StationName {
	var queue = []StationName{start}
	var seen = make(map[StationName]struct{})
	var parents = make(map[StationName]StationName)

	seen[start] = struct{}{}

	for len(queue) > 0 {
		curr := queue[0]

		// dequeue current
		queue = queue[1:]

		connections, ok := networkMap[curr]
		if !ok {
			fmt.Println("BFS error: Something went wrong")
		}

		for connection := range connections {
			// enqueue all unexplored direct connections along with parent information
			if _, ok := seen[connection]; !ok {
				queue = append(queue, connection)
				parents[connection] = curr
				seen[connection] = struct{}{}
			}

			if connection == end {
				// path found -> retrace steps
				return createPath(parents, end)
			}
		}
	}

	return nil
}

func createPath(parents map[StationName]StationName, end StationName) []StationName {
	path := []StationName{end}

	parent, ok := parents[end]
	for ok {
		path = append(path, parent)
		parent, ok = parents[parent]
	}

	slices.Reverse(path)

	fmt.Println("Found path:", path)
	return path
}
