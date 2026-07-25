package paths

import (
	"errors"
	"fmt"
	"pathinder/model"
)

func (g *Graph) DinicAlg(start, end string, numTrains int) ([][]string, error) {

	pathList := make([][]model.StationName, 0, numTrains) // contains the path list from dfs search
	max_flow := 0
	var err error
	layercount := 0

	// checks if there is actually a path available if so it keep looping untill all node cap are full

exit:
	for g.BFSAlg(start, end) {
		layercount++
		for {
			foundFlow, path := g.DFSAlg(false, start, end)
			fmt.Println(">>", path)
			if foundFlow == 0 {
				break
			}
			pathList = append(pathList, path)
			max_flow += foundFlow
			if max_flow >= numTrains {
				break exit
			}
		}
	}

	// if the maxflow is equal to the train number and the max flow was found on the first layer
	//  there is no need for flow correction and we can skip this part and use the path found in the dfs
	if layercount > 1 {
		pathList = pathList[:0]
		for range max_flow {
			foundFlow, path := g.DFSAlg(true, start, end)
			if foundFlow == 0 {
				break
			}
			pathList = append(pathList, path)

		}
	}
	if max_flow == 0 {
		err = errors.New("Error: No Path From Start To End Stations")
	}

	return pathList, err
}
