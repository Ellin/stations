// Scheduler using paths found using Edmonds-Karp

package paths

import (
	"pathinder/internal/parsers"
	"strconv" 
	"fmt"
)

type Train struct {
	Name string // e.g. "T1"
	PathID int // Index of the path the train is scheduled to use
	NextStation parsers.StationName // e.g. "waterloo"
}

func (g *Graph) RunScheduler(start, end parsers.StationName, numTrains int) {
	pathSet := g.FindPathSet(start, end, numTrains)
	divideTrains(numTrains, pathSet)
}


func (g *Graph) FindPathSet(start, end parsers.StationName, numTrains int) (pathSet [][]parsers.StationName) {
	flow, realPaths := g.EdmondsKarp(start, end, numTrains)
	pathSet = realPaths[flow-1]
	fmt.Println("PATH SET FOUND")
	fmt.Println(pathSet)

	return pathSet
}

// Decide which train gets assigned to which path
func divideTrains(numTrains int, pathSet [][]parsers.StationName) map[int][]Train {
	pathMap := make(map[int][]Train)
	for i, _ := range pathSet {
		pathMap[i] = []Train{}
	}

	for i:=0; i < numTrains; i++ {

		shortestTurnNum := getNumTurns(0, pathSet, pathMap)
		shortestPathIndex := 0

		// Find path with shortest # turns
		for pathID, _ := range pathMap {
			numTurns := getNumTurns(pathID, pathSet, pathMap)
			if numTurns < shortestTurnNum {
				shortestTurnNum = numTurns
				shortestPathIndex = pathID
			}
		}

		// Schedule train
		pathMap[shortestPathIndex] = append(pathMap[shortestPathIndex], Train{
			Name: "T" + strconv.Itoa(i+1),
			PathID: shortestPathIndex,
			NextStation: pathSet[shortestPathIndex][0],
		})
	}

	fmt.Println(pathMap)

	return pathMap
}

// getNumTurns returns the number of turns it would take to use a certain path, taking into account the wait time for that path (# of trains already scheduled for that path)
func getNumTurns(pathID int, pathSet [][]parsers.StationName, pathMap map[int][]Train) int {
	numHops := len(pathSet[pathID]) - 1 // hops between stations
	numWait := len(pathMap[pathID]) - 1 // -1 because the first train doesn't need a turn of waiting
	return numHops + numWait
}