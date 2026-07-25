package paths

func (g *Graph) BFSAlg(start, end string) bool {

	//reset level node to -1 with no level set yet
	g.ResetLevel()

	startID, endID := g.VertexNameMap[start], g.VertexNameMap[end]

	queue := []int{startID}
	g.LeveL[startID] = 0 // we set the start level count to 0 since its the sourse and its 0 distance away from the start

	for len(queue) > 0 {

		firstNode := queue[0]
		queue = queue[1:]
		current := firstNode

		for _, e := range g.EKGraph[current] { // go throught the neighbor edges of the current station
			if g.LeveL[e.To] < 0 && e.Cap > 0 { // check if level for that station is no set = -1 and cap is not used in another path
				g.LeveL[e.To] = g.LeveL[current] + 1 // increment the neighbor level with parent level count
				queue = append(queue, e.To)
			}

		}
	}

	return g.LeveL[endID] >= 0 // returns false if a level was not set for the end path
	//  meaning there is no path between start and end
}

// graph level reseter func  reset each node to -1 meaning empty in this case
func (g *Graph) ResetLevel() {
	for i := range g.LeveL {
		g.LeveL[i] = -1
		g.DeadEnd[i] = false

	}
	clear(g.EdgeBook)
}
