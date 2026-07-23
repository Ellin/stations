package paths

import (
	"errors"
	"fmt"
	"pathinder/model"
	"slices"
)

func (g *Graph) DinicAlg(start, end string, numTrains int) ([][]string, error) {

	pathList := make([][]model.StationName, 0, numTrains) // contains the path list from dfs search
	max_flow := 0
	var err error

	// checks if there is actually a path available if so it keep looping untill all node cap are full
exit:
	for g.BFSAlg(start, end) {
		for {
			foundFlow, path := g.DFSAlg(false, start, end)
			if foundFlow == 0 {
				break
			}
			pathList = append(pathList, slices.Compact(path))
			max_flow += foundFlow
			if max_flow >= numTrains {
				break exit
			}
		}
	}

	fmt.Println("max ", max_flow)
	if max_flow <= numTrains {
		clear(pathList)
		pathList = pathList[:0]
		for range max_flow {
			// if !g.BFSFlowCorrection(start, end) { // rebuild depth graph for the flow correction for dfs again, based on the reverse cap
			// 	break
			// }
			foundFlow, path := g.DFSAlg(true, start, end)
			if foundFlow == 0 {
				break
			}
			pathList = append(pathList, slices.Compact(path))

		}
	}
	if max_flow == 0 {
		err = errors.New("Error: No Path From Start To End Stations")
	}

	// fmt.Println("max Flow at a time :", max_flow, "\n paths: ", pathList) // max flow is the same lenght as the path we found
	return pathList, err
}
