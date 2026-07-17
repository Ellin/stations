package pathfinding_alg

import (
	"fmt"
	"pathinder/internal/parsers"
)

func (g *Graph) DinicAlg(start, end string) ([][]string, error) {

	pathList := [][]parsers.StationName{} // contains the path list from dfs search
	max_flow := 0

	// checks if there is actually a path available if so it keep looping untill all node cap are full
	for g.BFSAlg2(start, end) {

		for {
			foundFlow, path := g.DFSAlg(start, end)
			if foundFlow == 0 {
				break
			}
			max_flow += foundFlow
			pathList = append(pathList, path)
		}
	}
	fmt.Println("max Flow at a time :", max_flow) // max flow is the same lenght as the path we found
	return pathList, nil
}
