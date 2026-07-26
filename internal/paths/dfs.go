package paths

import (
	"pathinder/internal/parsers"
	"slices"
)

func (g *Graph) DFSAlg(dfsFlow bool, start, end string) (int, []parsers.StationName) {

	var route []parsers.StationName
	startID, endID := g.VertexNameMap[start], g.VertexNameMap[end]
	flow := 1 // the flow for each station is alway 1 so we initialize it with 1

	if dfsFlow {
		// after we find the paths we need to remove blocking vertices
		//dfsFlowCorrection is same as the normal dfs functions one diffrence is uses the  residual graph's reverse cap
		if route, flow = g.dfsFlowCorrection(flow, startID, endID); len(route) != 0 {
			slices.Reverse(route)
		}
	} else {
		// call the dfs algorithm with initail flow 1 and source and sink node id
		route, flow = g.dfs(flow, startID, endID)
		slices.Reverse(route) // reverse the paths retuned as they are currently in reverse
	}

	return flow, route
}

func (g *Graph) dfs(flow int, current, end int) ([]parsers.StationName, int) {

	// if the start and end are same the path is found return the path slice and the flow count
	if current == end {
		return []parsers.StationName{g.VertexIDMap[current]}, flow
	}

	// g.EdgeBook is used to book keep the current edges of the parent
	// if we have parent A connected to 1000,
	// if we already check 998 we dont want to go throung the whole neighbors list again and edge book should give you the id to the 999 edge
	for i := g.EdgeBook[current]; i < len(g.EKGraph[current]); i++ { // loop around the neighbor of the current vertex
		g.EdgeBook[current] = i
		e := &g.EKGraph[current][i]
		if g.DeadEnd[e.To] { // skip neighbors that dont lead anywhere
			continue
		}
		if g.LeveL[current] < g.LeveL[e.To] && e.Cap > 0 { // it checks if the neighbor edge is one level deeper(current vertex is alway higher in the level graph) and there is still a cap remaining for the edge
			minCp := flow
			if e.Cap < minCp {
				minCp = e.Cap
			}

			path, flowReturn := g.dfs(int(minCp), e.To, end) // call the dfs again this time the neighbor vertex as the source

			if flowReturn > 0 { // if path found in the previous update the cap
				e.Cap -= flowReturn
				g.EKGraph[e.To][e.Rev].Cap += flowReturn

				if len(path) > 0 && path[len(path)-1] == g.VertexIDMap[current] { // skip for the split node
					return path, flowReturn
				}
				return append(path, g.VertexIDMap[current]), flowReturn
			}
		}

	}
	// if the flow returned is 0 or there is no neighbors there no path from this node
	// so set it to deadEnd so it doesnt maroon the search
	g.DeadEnd[current] = true
	return []parsers.StationName{}, 0
}

// we construct the path with same step as we used in the max flow dfs function
// it does a normal dfs pathfinding algorithm and check if the neighbors split edge has been used to backtrack if so it ignores that path
// since it a blocking edge eg*( a - b. b - c, the we back track using the edge c - b ) if this is cound on edge c then we skip that neighbor
func (g *Graph) dfsFlowCorrection(flow int, current, end int) ([]parsers.StationName, int) {

	// if the start and end are same the path is found return the path slice and the flow count
	if current == end {
		return []parsers.StationName{g.VertexIDMap[current]}, flow
	}

	for i := range g.EKGraph[current] { // loop around the neighbor of the current vertex
		e := &g.EKGraph[current][i]
		rev := &g.EKGraph[e.To][e.Rev]

		if !e.Real {
			continue
		}
		if rev.Cap > 0 { // it checks if the neighbor's reverse edge (b-c reverse edge = c-b) has a remaining cap
			minCp := flow
			if rev.Cap < minCp {
				minCp = rev.Cap
			}
			path, flowReturn := g.dfsFlowCorrection(int(minCp), e.To, end) // call the dfs again this time the neighbor vertex as the source
			if flowReturn > 0 {                                            // if path found in the previous update the cap
				e.Cap += flowReturn                                               // we reset the e.cap backto its original capacity
				g.EKGraph[e.To][e.Rev].Cap -= flowReturn                          // we reset the e.reverse.cap backto its original capacity
				if len(path) > 0 && path[len(path)-1] == g.VertexIDMap[current] { // skip for the split node
					return path, flowReturn
				}
				return append(path, g.VertexIDMap[current]), flowReturn
			}

		}

	}

	return []parsers.StationName{}, 0
}
