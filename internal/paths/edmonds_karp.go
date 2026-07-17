package paths

import (
	"pathinder/internal/parsers"
	"fmt"
	"slices"
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

		// dequeue current
		queue = queue[1:]

		connections, ok := g.EKGraph[curr]
		fmt.Printf("CONNECTIONS to %v: %v", g.VertexIDMap[curr], connections)
		if !ok {
			fmt.Println("BFS error: Something went wrong")
		}

		for i, edge := range connections {
			fmt.Println("edge: ", edge)
			fmt.Println("Cap:", edge.Cap)

			// enqueue all unexplored neighbors that have residual capacity along with parent information
			if _, ok := seen[edge.To]; !ok && edge.Cap > 0{
				fmt.Println("Adding unexplored neighbour: %v", g.VertexIDMap[edge.To])
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

func (g *Graph) createPath(parents map[VertexID][]int, end VertexID) []VertexID {
	path := []VertexID{end}

	parent, ok := parents[end]
	for ok {
		parentID, toIndex := parent[0], parent[1]
		path = append(path, parentID)

		// Decrease residual cap (.Cap) of forward edges and increase on reverse edges
		forward := &g.EKGraph[parentID][toIndex]
		reverse := &g.EKGraph[forward.To][forward.Rev]
		
		forward.Cap--
		reverse.Cap++
		fmt.Printf("forward cap from %v to %v: %v\n", g.VertexIDMap[parentID], g.VertexIDMap[forward.To], forward.Cap)
		fmt.Println("reverse cap", reverse.Cap)

		parent, ok = parents[parentID]
	}

	slices.Reverse(path)


	return path
}

func (g *Graph) EdmondsKarp(start, end parsers.StationName) (maxFlow int, allPaths [][]VertexID, realPaths [][][]VertexID){
	startID, endID := g.VertexNameMap[start], g.VertexNameMap[end]	

	

	parents, found := g.Bfs(startID, endID)
	for found {
		// Update the residual capacities (no need to find minimal residual capacity bottleneck as it will always be 1 in our case)
		allPaths = append(allPaths, g.createPath(parents, endID))
		maxFlow++

		// find REAL paths
		usedMap := make(map[VertexID]struct{})
		var foundPath []VertexID
		var flowPathSet [][]VertexID
		fmt.Printf("FOUND PATHS FOR FLOW %v:\n", maxFlow)
		for i := maxFlow; i > 0; i-- {
			parents_real, _ := g.Bfs_flow(startID, endID, usedMap)
			foundPath, usedMap = g.createPath_real(parents_real, startID, endID)
			g.printPaths([][]VertexID{foundPath})
			flowPathSet = append(flowPathSet, foundPath)
		}
		realPaths = append(realPaths, flowPathSet)

		parents, found = g.Bfs(startID, endID)
	}

	fmt.Println(maxFlow, allPaths)
	fmt.Println("!AUGMENTING PATHS:")
	g.printPaths(allPaths)

	fmt.Println("!REAL PATHS:")
	for i, paths := range realPaths {
		fmt.Printf("flow %v:\n", i+1)
		g.printPaths(paths)
	}

	return maxFlow, allPaths, realPaths
}

func (g *Graph) printPaths(paths [][]VertexID) {
	for _, path := range paths {
		for _, id := range path {
			fmt.Print(g.VertexIDMap[id], "->")
		}
		fmt.Println()
	}
}

func (g *Graph) Bfs_flow(start, end VertexID, usedMap map[VertexID]struct{}) (parents map[VertexID][]int, found bool) {
	var queue = []VertexID{start}
	var seen = make(map[VertexID]struct{})
	parents = make(map[VertexID][]int) 

	seen[start] = struct{}{}

	for len(queue) > 0 {
		curr := queue[0]

		// dequeue current
		queue = queue[1:]

		connections, ok := g.EKGraph[curr]
		fmt.Printf("CONNECTIONS to %v: %v", g.VertexIDMap[curr], connections)
		if !ok {
			fmt.Println("BFS error: Something went wrong")
		}

		for i, edge := range connections {
			fmt.Println("edge: ", edge)
			fmt.Println("Cap:", edge.Cap)

			// enqueue all unexplored neighbors that have residual capacity along with parent information

			_, used := usedMap[edge.To]

			if _, ok := seen[edge.To]; !ok && edge.Cap == 0 && edge.Real && !used {
				fmt.Println("Adding unexplored neighbour: %v", g.VertexIDMap[edge.To])
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

func (g *Graph) createPath_real(parents map[VertexID][]int, start, end VertexID) (path []VertexID, usedMap map[VertexID]struct{}) {
	usedMap = make(map[VertexID]struct{})
	path = append(path, end)

	parent, ok := parents[end]
	for ok {
		parentID, _ := parent[0], parent[1]
		path = append(path, parentID)
		usedMap[parentID] = struct{}{}

		// fmt.Printf("forward cap from %v to %v: %v\n", g.VertexIDMap[parentID], g.VertexIDMap[forward.To], forward.Cap)
		// fmt.Println("reverse cap", reverse.Cap)

		parent, ok = parents[parentID]

	}

	delete(usedMap, start)
	delete(usedMap, end)

	slices.Reverse(path)


	return path, usedMap
}
