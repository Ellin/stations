package paths

import (
	"pathinder/internal/parsers"
	"fmt"
	"slices"
	"log"
)

type VertexID = int


// Parents map: the []int value holds a pair of values: the parent vertex ID and the INDEX within its []Edge array where the "to" edge can be found 
func (g *Graph) Bfs(start, end VertexID) (parents map[VertexID][]int, found bool) {
	var queue = []VertexID{start}
	var seen = make(map[VertexID]struct{})
	parents = make(map[VertexID][]int) 

	seen[start] = struct{}{}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:] // dequeue current

		connections, ok := g.EKGraph[curr]
		if !ok {
			log.Fatal("Non-existent key in EK graph")
		}

		for i, edge := range connections {
			// enqueue all unexplored neighbors that have residual capacity along with parent information
			if _, ok := seen[edge.To]; !ok && edge.Cap > 0{
				queue = append(queue, edge.To)
				parents[edge.To] = []int{curr, i} 
				seen[edge.To] = struct{}{}

				if edge.To == end {
					// path found
					return parents, true
				}
			}
		}
	}

	return nil, false
}

// traceAugmentedPath follows the path from the sink up to the source using the parents map
// The residual capacities of the used edges in the augmented path are updated
// Returns a list of nodes in the path from start to sink
func (g *Graph) traceAugmentedPath(parents map[VertexID][]int, end VertexID) []VertexID {
	path := []VertexID{end}

	parent, ok := parents[end]
	for ok {
		parentID, toIndex := parent[0], parent[1]
		
		// Decrease residual cap (.Cap) of forward edges and increase on reverse edges
		forward := &g.EKGraph[parentID][toIndex]
		reverse := &g.EKGraph[forward.To][forward.Rev]	
		forward.Cap--
		reverse.Cap++

		path = append(path, parentID)
		parent, ok = parents[parentID]
	}

	slices.Reverse(path)

	return path
}

func (g *Graph) EdmondsKarp(start, end parsers.StationName) (maxFlow int, augmentingPaths [][]VertexID, realPaths [][][]VertexID){
	startID, endID := g.VertexNameMap[start], g.VertexNameMap[end]	
	parents, found := g.Bfs(startID, endID)
	for found {
		// Update the residual capacities by tracing the augmented path
		// No need to find minimal residual capacity bottleneck for updating the maxFlow as it will always be 1 in our case
		augmentingPaths = append(augmentingPaths, g.traceAugmentedPath(parents, endID))
		maxFlow++ 

		// Find REAL paths without any reverse edges
		// Run BFS maxFlow amount of times as we know maxFlow == # of non-overlapping paths that can be found
		usedMap := make(map[VertexID]map[VertexID]struct{})
		var foundPath []VertexID
		var flowPathSet [][]VertexID
		for i := maxFlow; i > 0; i-- {
			parents_real, _ := g.Bfs_flow(startID, endID, usedMap)
			foundPath = g.traceRealPath(parents_real, usedMap, startID, endID)
			flowPathSet = append(flowPathSet, foundPath)
		}
		realPaths = append(realPaths, flowPathSet)

		parents, found = g.Bfs(startID, endID)
	}

	g.prettifyPaths(realPaths)
	g.printEdmondsKarpResults(maxFlow, augmentingPaths, realPaths)
	
	return maxFlow, augmentingPaths, realPaths
}

// prettifyPaths converts paths containing split vertexes to a readable form
// Splits vertexes are rejoined and vertex ids are converted back to station names
func (g *Graph) prettifyPaths(realPaths [][][]VertexID) (realStationPaths [][][]string) {
	for _, pathSet := range realPaths {
		var convertedPathSet [][]string
		for _, path := range pathSet {
			var prev string
			var convertedPath []string
			for _, vertexID := range path {
				stationName := g.VertexIDMap[vertexID]
				if prev != stationName {
					convertedPath = append(convertedPath, stationName)
					prev = stationName
				}
			}
			convertedPathSet = append(convertedPathSet, convertedPath)
		}
		realStationPaths = append(realStationPaths, convertedPathSet)
	}
	fmt.Println("PRETTY", realStationPaths)
	return realStationPaths
}

func (g *Graph) printEdmondsKarpResults(maxFlow int, augmentingPaths [][]VertexID, realPaths [][][]VertexID) {
	fmt.Println("Max flow:", maxFlow)
	fmt.Println("!AUGMENTING PATHS:")
	g.printPaths(augmentingPaths)

	fmt.Println("!REAL PATHS:")
	for i, paths := range realPaths {
		fmt.Printf("Flow %v:\n", i+1)
		g.printPaths(paths)
	}
}

func (g *Graph) printPaths(paths [][]VertexID) {
	for _, path := range paths {
		var prevStation string
		for _, id := range path {
			station := g.VertexIDMap[id]
			if prevStation != station { // Collapse split nodes
				fmt.Print(station, "->")
				prevStation = station
			}
		}
		fmt.Println()
	}
}

func (g *Graph) Bfs_flow(start, end VertexID, usedMap map[VertexID]map[VertexID]struct{}) (parents map[VertexID][]int, found bool) {
	var queue = []VertexID{start}
	var seen = make(map[VertexID]struct{})
	parents = make(map[VertexID][]int) 

	seen[start] = struct{}{}

	for len(queue) > 0 {
		curr := queue[0]		
		queue = queue[1:] // dequeue current

		connections, ok := g.EKGraph[curr]
		if !ok {
			log.Fatal("Non-existent key in EK graph")
		}

		for i, edge := range connections {
			if _, ok := seen[edge.To]; !ok && edge.Cap == 0 && edge.Real && !isUsed(usedMap, curr, edge.To) {
				queue = append(queue, edge.To)
				parents[edge.To] = []int{curr, i} 
				seen[edge.To] = struct{}{}

				if edge.To == end {
					// path found)
					return parents, true
				}
			}
		}
	}

	return nil, false
}

func isUsed(usedMap map[VertexID]map[VertexID]struct{}, fromID int, toID int) bool {
	if _, ok := usedMap[fromID]; ok {
		_, used := usedMap[fromID][toID]
		return used
	}
	return false
}

func (g *Graph) traceRealPath(parents map[VertexID][]int, usedMap map[VertexID]map[VertexID]struct{}, start, end VertexID) []VertexID {
	path := []VertexID{end}

	parent, ok := parents[end]
	for ok {
		parentID, toIndex := parent[0], parent[1]
		toID := g.EKGraph[parentID][toIndex].To

		if _, ok := usedMap[parentID]; !ok {
			usedMap[parentID] = make(map[VertexID]struct{})
		}
		usedMap[parentID][toID] = struct{}{}

		path = append(path, parentID)
		parent, ok = parents[parentID]
	}

	slices.Reverse(path)

	return path
}
