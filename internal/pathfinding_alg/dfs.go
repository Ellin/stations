package pathfinding_alg

import (
	"math"
	"pathinder/internal/parsers"
	"slices"
)

func (g *Graph) DFSAlg(dfsFlow bool, start, end string) (int, []parsers.StationName) {

	var route []parsers.StationName
	flow := 1
	startID := g.VertexNameMap[start]
	endID := g.VertexNameMap[end]

	// call the dfs algorithm with initail flow 1 and source and sink node id
	if dfsFlow {
		if route, flow = g.dfsFlowCorrection(flow, startID, endID); len(route) != 0 {
			slices.Reverse(route)
		}
	} else {
		flow = g.dfs(flow, startID, endID)
	}

	return flow, route
}
func (g *Graph) dfs(flow int, current, end int) int {
	// fmt.Println("in dfs", g.EKGraph[current])
	// of the start and end are same the path is found return the path slice and the flow count

	if current == end {
		return flow
	}

	for i := range g.EKGraph[current] { // loop around the neighbor of the current vertex
		e := &g.EKGraph[current][i]
		if g.LeveL[current] < g.LeveL[e.To] && e.Cap > 0 { // it checks if the neighbor edge is one level deeper(current vertex is alway higher in the level graph) and there is still a cap remaining for the edge
			minCp := math.Min(float64(flow), float64(e.Cap)) // in this implementaion the flow and cap are always the same if cap is greater than 0 meaning its 1 then flow is also one but to be true to algorithm we are doing this part as well

			flowReturn := g.dfs(int(minCp), e.To, end) // call the dfs again this time the neighbor vertex as the source

			if flowReturn > 0 { // if path found in the previous update the cap
				e.Cap -= flowReturn
				g.EKGraph[e.To][e.Rev].Cap += flowReturn

				return flowReturn
			}

		}

	}

	return 0
}

// we construct the path with same step as we used in the max flow dfs function
// it does a normal dfs pathfinding algorithm and check if the neighbors split edge has been used to backtrack if so it ignores that path
// since it a blocking edge eg*( a - b. b - c, the we back track using the edge c - b ) if this is cound on edge c then we skip that neighbor
func (g *Graph) dfsFlowCorrection(flow int, current, end int) ([]parsers.StationName, int) {
	// fmt.Println("in dfs", g.EKGraph[current])
	// of the start and end are same the path is found return the path slice and the flow count

	if current == end {
		return []parsers.StationName{g.VertexIDMap[current]}, flow
	}

	for i := range g.EKGraph[current] { // loop around the neighbor of the current vertex
		e := &g.EKGraph[current][i]
		rev := &g.EKGraph[e.To][e.Rev]

		if rev.Cap > 0 { // it checks if the neighbor's reverse edge (b-c reverse edge = c-b) has a remaining cap
			minCp := math.Min(float64(flow), float64(rev.Cap)) // in this implementaion the flow and reverse cap are always the same if reverse cap is greater than 0 meaning its 1 then flow is also one but to be true to algorithm we are doing this part as well

			path, flowReturn := g.dfsFlowCorrection(int(minCp), e.To, end) // call the dfs again this time the neighbor vertex as the source

			if flowReturn > 0 { // if path found in the previous update the cap
				e.Cap += flowReturn                      // we reset the e.cap backto its original capacity
				g.EKGraph[e.To][e.Rev].Cap -= flowReturn // we reset the e.reverse.cap backto its original capacity

				return append(path, g.VertexIDMap[current]), flowReturn
			}

		}

	}

	return []parsers.StationName{}, 0
}
