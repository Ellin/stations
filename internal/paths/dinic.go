package paths

import (
	"errors"
	"pathinder/model"
)

func (g *Graph) DinicAlg(start, end string, numTrains int) ([][]model.StationName, error) {

	var optionsList [][][]model.StationName
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
			if foundFlow == 0 {
				break
			}
			pathList = append(pathList, path)
			max_flow += foundFlow

			if max_flow >= numTrains {
				break exit
			}
		}

		// we save the augmenting path as a probable path to the opionslist
		// we can only do this in the first level graph as the graph change there will
		// mostlikly be an overlap
		if layercount == 1 {
			optionsList = append(optionsList, pathList)
		}
	}

	// exit for non existing flow no need to do othe checks
	if max_flow == 0 {
		err = errors.New("Error: No Path From Start To End Stations")
		return [][]model.StationName{}, err
	}

	// if the maxflow is equal to the train number and the max flow was found on the first level graph with only augmenting path
	//  there is no need for flow correction and we can skip this part and use the path found in the dfs
	if layercount > 1 {
		FlowCorrectedPaths := make([][]model.StationName, 0, max_flow)
		for range max_flow {
			foundFlow, path := g.DFSAlg(true, start, end)
			if foundFlow == 0 {
				break
			}
			FlowCorrectedPaths = append(FlowCorrectedPaths, path)

		}
		optionsList = append(optionsList, FlowCorrectedPaths)
	} else if layercount == 1 {
		return pathList, nil
	}

	return bestSet(numTrains, optionsList), err
}

// this func takes a list of pathlist from diffrent augmenting flows and returned the one with lowest avarageturns
func bestSet(numTrains int, optionsList [][][]model.StationName) [][]model.StationName {

	var bestSetIndex int
	bestTurn := calcAvgTurns(numTrains, optionsList[0])
	for i, pathlist := range optionsList {

		avgTurn := calcAvgTurns(numTrains, pathlist)
		if avgTurn < bestTurn {
			bestTurn = avgTurn
			bestSetIndex = i
		}
	}
	bestPathSet := optionsList[bestSetIndex]
	return bestPathSet
}
