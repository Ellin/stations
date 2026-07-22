package pathfinding_alg

import (
	"container/list"
)

func (g *Graph) BFSAlg(start, end string) bool {

	//reset level node to -1 to mark it as not visited
	g.ResetLevel()

	startID := g.VertexNameMap[start]
	endID := g.VertexNameMap[end]

	queue := list.New()
	queue.PushBack(startID)
	g.LeveL[startID] = 0 // we set the start level count to 0 since its the sourse and its 0 distance away from the start
	for queue.Len() > 0 {
		firstNode := queue.Front()
		queue.Remove(firstNode)
		current := firstNode.Value.(int)

		// v := g.VertexNameMap[current]
		for _, e := range g.EKGraph[current] { // go throught the neighbor edges of the current node
			if g.LeveL[e.To] < 0 && e.Cap > 0 { // check if level node is count is empty = -1 and cap is not used in another path
				g.LeveL[e.To] = g.LeveL[current] + 1 // increment the neighbor level with parent level count
				queue.PushBack(e.To)
			}

		}
	}

	return g.LeveL[endID] >= 0 // return false if a level was not set for the end path
	//  meaning there is no path between start and end
}

func (g *Graph) BFSFlowCorrection(start, end string) bool {

	//reset level node to -1 to mark it as not visited
	g.ResetLevel()

	startID := g.VertexNameMap[start]
	endID := g.VertexNameMap[end]

	queue := list.New()
	queue.PushBack(startID)
	g.LeveL[startID] = 0 // we set the start level count to 0 since its the sourse and its 0 distance away from the start
	for queue.Len() > 0 {
		firstNode := queue.Front()
		queue.Remove(firstNode)
		current := firstNode.Value.(int)

		// v := g.VertexNameMap[current]
		for _, e := range g.EKGraph[current] { // go throught the neighbor edges of the current node

			if e.Reverse {
				continue
			}
			rev := &g.EKGraph[e.To][e.Rev]
			if g.LeveL[e.To] < 0 && rev.Cap > 0 { // check if level node is count is empty = -1 and cap is not used in another path
				g.LeveL[e.To] = g.LeveL[current] + 1 // increment the neighbor level with parent level count
				queue.PushBack(e.To)
			}

		}
	}

	return g.LeveL[endID] >= 0 // return false if a level was not set for the end path
	//  meaning there is no path between start and end
}

// graph level reseter func  reset each node to -1 meaning empty in this case
func (g *Graph) ResetLevel() {
	for i := range g.LeveL {
		g.LeveL[i] = -1
		g.DeadEnd[i] = false
	}
}
