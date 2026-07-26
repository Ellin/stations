package paths

func (g *Graph) BFSAlg(start, end string) bool {

	//reset all node levels to -1(unvisited) with no level set yet
	g.ResetLevel()

	startID, endID := g.VertexNameMap[start], g.VertexNameMap[end]

	queue := []int{startID}
	g.LeveL[startID] = 0 // we set the start level count to 0 since its the source and its 0 distance away from the start

	for len(queue) > 0 {

		firstNode := queue[0]
		queue = queue[1:]
		current := firstNode

		for _, e := range g.EKGraph[current] { // go through the neighbor edges of the current station
			if g.LeveL[e.To] < 0 && e.Cap > 0 { // check if level for that station is no set = -1 and the edge still has a residual capacity
				g.LeveL[e.To] = g.LeveL[current] + 1 // increment the neighbor level with parent level count
				queue = append(queue, e.To)
			}

		}
	}

	return g.LeveL[endID] >= 0 // returns false if a level was not set for the end path
	//  meaning there is no augmenting path between start and end
}

// reset graph level to -1 (unvisited) this is done for every new BFS level graph
// reset node DeadEnd to false since a node that was a deadEnd in previos phase may not be a deadEnd now since residual caps change
// and clear the EdgBook(the residual graph neighbor book keeper)

func (g *Graph) ResetLevel() {
	for i := range g.LeveL {
		g.LeveL[i] = -1
		g.DeadEnd[i] = false

	}
	clear(g.EdgeBook)
}
