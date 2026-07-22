package pathfinding_alg

import (
	"pathinder/model"
)

func (g *Graph) DinicAlg(start, end string) ([][]string, error) {

	pathList := [][]model.StationName{} // contains the path list from dfs search
	max_flow := 0

	// checks if there is actually a path available if so it keep looping untill all node cap are full
	for g.BFSAlg(start, end) {

		for {
			foundFlow, _ := g.DFSAlg(false, start, end)
			if foundFlow == 0 {
				break
			}
			max_flow += foundFlow
		}
	}

	for range max_flow {
		if !g.BFSFlowCorrection(start, end) { // rebuild depth graph for the flow correction for dfs again, based on the reverse cap
			break
		}
		foundFlow, path := g.DFSAlg(true, start, end)
		if foundFlow == 0 {
			break
		}
		pathList = append(pathList, path)
	}
	// fmt.Println("max Flow at a time :", max_flow) // max flow is the same lenght as the path we found
	return pathList, nil
}
