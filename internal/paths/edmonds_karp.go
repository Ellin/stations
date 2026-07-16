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
					g.createPath(parents, end)
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

func (g *Graph) EdmondsKarp(start, end parsers.StationName) (maxFlow int, allPaths [][]VertexID){
	startID, endID := g.VertexNameMap[start], g.VertexNameMap[end]	

	parents, found := g.Bfs(startID, endID)
	for found {
		// Update the residual capacities (no need to find minimal residual capacity bottleneck as it will always be 1 in our case)
		allPaths = append(allPaths, g.createPath(parents, endID))
		maxFlow++
		parents, found = g.Bfs(startID, endID)
	}

	fmt.Println(maxFlow, allPaths)
	fmt.Println("!Found path:")
	for _, path := range allPaths{
		for _, id := range path {
			fmt.Print(g.VertexIDMap[id], "->")
		}
		fmt.Println()
	}

	return maxFlow, allPaths
}