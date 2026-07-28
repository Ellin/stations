package paths

import (
	"fmt"
	"log"
	"pathfinder/internal/parsers"
	"slices"
	"strings"
)

type VertexID = int

type Parent struct {
	ID          VertexID // parent ID
	ToEdgeIndex int      // Index within the parent's []Edge array where the "to" edge can be found (using the EKGraph)
}

func (g *Graph) Bfs(version string, start, end VertexID, usedMap map[VertexID]map[VertexID]struct{}) (parents map[VertexID]Parent, found bool) {
	var queue = []VertexID{start}
	var seenMap = make(map[VertexID]struct{})
	parents = make(map[VertexID]Parent)

	seenMap[start] = struct{}{}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:] // dequeue current

		connections, ok := g.EKGraph[curr]
		if !ok {
			log.Fatal("Non-existent key in EK graph")
		}

		for i, edge := range connections {
			// enqueue all unexplored neighbors that have residual capacity along with parent information
			_, seen := seenMap[edge.To]
			if !seen && canTravel(version, edge, curr, usedMap) {
				queue = append(queue, edge.To)
				parents[edge.To] = Parent{curr, i}
				seenMap[edge.To] = struct{}{}

				if edge.To == end {
					// path found
					return parents, true
				}
			}
		}
	}

	return nil, false
}

// canTravel is a helper for the Bfs function that checks if an edge can be travelled upon
// The edge travelling conditions depend on the Bfs version: looking for augmenting paths vs real paths
func canTravel(version string, edge Edge, from VertexID, usedMap map[VertexID]map[VertexID]struct{}) bool {
	switch version {
	case "aug":
		return edge.Cap > 0

	case "real":
		// Only real edges (no reverse edges) that had flow when finding augmenting paths (reducing its capacity to 0)
		// that haven't already been used in real path finding should be used
		return edge.Cap == 0 && edge.Real && !isUsed(usedMap, from, edge.To)

	default:
		return false
	}
}

// traceAugmentedPath follows the path from the sink up to the source using the parents map
// The residual capacities of the used edges in the augmented path are updated
// Returns a list of nodes in the path from start to sink
func (g *Graph) traceAugmentedPath(parents map[VertexID]Parent, end VertexID) ([]VertexID, bool) {
	path := []VertexID{end}
	isPreviousBlocking := false

	parent, ok := parents[end]
	for ok {
		// Decrease residual cap (.Cap) of forward edges and increase on reverse edges
		forward := &g.EKGraph[parent.ID][parent.ToEdgeIndex]
		reverse := &g.EKGraph[forward.To][forward.Rev]
		forward.Cap--
		reverse.Cap++

		// Check if reverse edge was used -> if yes, previous augmented path found was blocking and real path extraction will result in different path choices
		if !forward.Real {
			isPreviousBlocking = true
		}

		path = append(path, parent.ID)
		parent, ok = parents[parent.ID]
	}

	slices.Reverse(path)

	return path, isPreviousBlocking
}

func (g *Graph) EdmondsKarp(start, end parsers.StationName, numTrains int) (maxFlow int, stationPaths [][][]parsers.StationName) {
	startID, endID := g.VertexNameMap[start], g.VertexNameMap[end]
	parents, found := g.Bfs("aug", startID, endID, nil)

	var augmentingPaths [][]VertexID
	var realPaths [][][]VertexID
	var currentPathSet [][]VertexID // non-overlapping paths

	for found {
		// Update the residual capacities by tracing the augmented path
		// No need to find minimal residual capacity bottleneck for updating the maxFlow as it will always be 1 in our case
		augmentingPath, isPreviousBlocking := g.traceAugmentedPath(parents, endID)
		augmentingPaths = append(augmentingPaths, augmentingPath)
		maxFlow++

		// Run BFS for extracting real paths only if the newest augmenting path found contains a reverse edge (meaning a path in the previous path set blocked flow and would not be found in the new path set with more flow)
		// Otherwise, the augmenting path is a real path that doesn't overlap with any other paths in the current path set.
		if isPreviousBlocking {
			realPaths = append(realPaths, currentPathSet)

			// Find REAL paths without any reverse edges
			// Run BFS maxFlow amount of times as we know maxFlow == # of non-overlapping paths that can be found
			usedMap := make(map[VertexID]map[VertexID]struct{})
			var foundPath []VertexID
			var realPathSet [][]VertexID
			for i := maxFlow; i > 0; i-- {
				parents_real, _ := g.Bfs("real", startID, endID, usedMap)
				foundPath = g.traceRealPath(parents_real, usedMap, startID, endID)
				realPathSet = append(realPathSet, foundPath)
			}

			currentPathSet = realPathSet
		} else {
			currentPathSet = append(currentPathSet, augmentingPath)
		}

		parents, found = g.Bfs("aug", startID, endID, nil)

		if maxFlow >= numTrains {
			break
		}
	}

	if len(currentPathSet) > 0 {
		realPaths = append(realPaths, currentPathSet)
	}

	g.printEdmondsKarpResults(maxFlow, augmentingPaths, realPaths)
	stationPaths = g.prettifyPaths(realPaths)
	return maxFlow, stationPaths
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
	fmt.Println("Path sets found:", realStationPaths)
	return realStationPaths
}

func (g *Graph) printEdmondsKarpResults(maxFlow int, augmentingPaths [][]VertexID, realPaths [][][]VertexID) {
	fmt.Println("Max flow:", maxFlow)
	fmt.Println("!AUGMENTING PATHS:")
	g.printPaths(augmentingPaths)

	fmt.Println("!REAL PATHS:")
	for _, paths := range realPaths {
		fmt.Printf("Flow %v:\n", len(paths))
		g.printPaths(paths)
	}
}

func (g *Graph) printPaths(paths [][]VertexID) {
	var builder strings.Builder

	for i, path := range paths {
		var prevStation string
		fmt.Fprintf(&builder, "%v. ", i+1)

		station := g.VertexIDMap[0]
		builder.WriteString(station)

		for j := 1; j < len(path); j++ {
			vertexID := path[j]
			station := g.VertexIDMap[vertexID]
			if prevStation != station { // Collapse split nodes
				builder.WriteString("->")
				builder.WriteString(station)
				prevStation = station
			}
		}
		builder.WriteString("\n")
	}
	fmt.Println(builder.String())
}

// isUsed checks if an edge has already been used during real path finding
func isUsed(usedMap map[VertexID]map[VertexID]struct{}, fromID int, toID int) bool {
	if _, ok := usedMap[fromID]; ok {
		_, used := usedMap[fromID][toID]
		return used
	}
	return false
}

// traceRealPath follows the path from the sink (end) up to the source (start) using the parents map
// Used edges are recorded in usedMap so that next iterations of Bfs for finding real paths can't use the same edges
// Returns a list of nodes in the path from start to sink
func (g *Graph) traceRealPath(parents map[VertexID]Parent, usedMap map[VertexID]map[VertexID]struct{}, start, end VertexID) []VertexID {
	path := []VertexID{end}

	parent, ok := parents[end]
	for ok {
		toID := g.EKGraph[parent.ID][parent.ToEdgeIndex].To

		if _, ok := usedMap[parent.ID]; !ok {
			usedMap[parent.ID] = make(map[VertexID]struct{})
		}
		usedMap[parent.ID][toID] = struct{}{}

		path = append(path, parent.ID)
		parent, ok = parents[parent.ID]
	}

	slices.Reverse(path)

	return path
}
