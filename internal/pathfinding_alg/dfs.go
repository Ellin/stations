package pathfinding_alg

import (
	"math"
	"pathinder/internal/parsers"
	"slices"
)

func (g *Graph) DFSAlg(start, end string) (int, []parsers.StationName) {

	var route []parsers.StationName
	flow := 1
	startID := g.VertexNameMap[start]
	endID := g.VertexNameMap[end]

	// call the dfs algorithm with initail flow 1 and source and sink node id
	if route, flow = g.dfs(flow, startID, endID); len(route) != 0 {
		slices.Reverse(route)
	}
	return flow, route
}
func (g *Graph) dfs(flow int, start, end int) ([]parsers.StationName, int) {

	// of the start and end are same the path is found return the path slice and the flow count
	if start == end {
		return []parsers.StationName{g.VertexIDMap[start]}, flow
	}

	for i := range g.EKGraph[start] { // loop around the neighbor of the current vertex
		e := &g.EKGraph[start][i]
		if g.LeveL[start] < g.LeveL[e.To] && e.Cap > 0 { // it checks if the neighbor edge is one level deeper(current vertex is alway higher in the level graph) and there is still a cap remaining for the edge
			minCp := math.Min(float64(flow), float64(e.Cap)) // in this implementaion the flow and cap are always the same if cap is greater than 0 meaning its 1 then flow is also one but to be true to algorithm we are doing this part as well

			path, flowReturn := g.dfs(int(minCp), e.To, end) // call the dfs again this time the neighbor vertex as the source

			if flowReturn > 0 { // if path found in the previous update the cap
				e.Cap -= flowReturn
				g.EKGraph[e.To][e.Rev].Cap += flowReturn

				return append(path, g.VertexIDMap[start]), flowReturn
			}

		}

	}

	return []parsers.StationName{}, 0
}
