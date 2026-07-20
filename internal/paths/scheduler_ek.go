// Scheduler using paths found using Edmonds-Karp

package paths

import (
	"pathinder/internal/parsers"
	"strconv" 
	"fmt"
)

type Train struct {
	Name string // e.g. "T1"
	PathID int // Index of the path within a pathSet the train is scheduled to use
	CurrentStationIndex int // e.g. path: ["waterloo", "victoria", "st_pancras"] ; "waterloo" -> index 0
	CurrentStationName parsers.StationName // e.g. "waterloo"
}

func (g *Graph) RunScheduler(start, end parsers.StationName, numTrains int) {
	pathSet := g.FindPathSet(start, end, numTrains)
	pathMap := divideTrains(numTrains, pathSet)
	turnSchedule := runTrains(end, pathSet, pathMap)
	printSchedule(turnSchedule)
}


// ! NOTE: This should be further optimized to take into account cost given some # of trains, not just max flow
// For example, the current algorithm for choosing a set of paths for 2 trains prefers a set of two very long non-overlapping paths
// vs one very short path that blocks multiple non-overlapping paths.
// But that short blocking path may still be the better path for certain number of trains.
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
			CurrentStationIndex: 0,
			CurrentStationName:  pathSet[shortestPathIndex][0],
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

func runTrains(end parsers.StationName, pathSet [][]parsers.StationName, pathMap map[int][]Train) ([][]string) {
	var turnSchedule [][]string // e.g. [["T1-a", "T2-b"]["T1-c", "T2-d"]["T1-end","T2-end"]]

	// Path index 0 = shortest path in a pathSet
	// The shortest path in a pathSet will always have the most number of trains scheduled for it

	// While there are trains left to move
	for len(pathMap[0]) > 0 {
		var turnGroup []string

		// MOVE ALL TRAINS IN LINE FOR ALL PATHS TO THE NEXT STATION IF NOT MOVING AT SAME STATION TO THE ONE AHEAD
		// IF STATION IS REACHED, DEQUEUE THAT TRAIN
		for pathID, path := range pathSet {

			var canMove bool

			for i := 0; i < len(pathMap[pathID]); i++ {
				train := &pathMap[pathID][i]

				if train.CurrentStationName == end {
					pathMap[pathID] = pathMap[pathID][1:] // dequeue train
					i--
					continue
				}

				// ALL TRAINS FIRST IN LINE CAN MOVE ONE STATION FORWARD
				if i == 0 {
					canMove = true

				}

				//	IF NOT FIRST IN LINE, CHECK FIRST IF IT CAN MOVE FORWARD
				if i > 0 {
					prevStation := &pathMap[pathID][i-1]
					if path[train.CurrentStationIndex + 1] != prevStation.CurrentStationName { // check if train in queue would be moving to the same station as the train ahead
						canMove = true
					} else {
						canMove = false
					}
				}

				if canMove {
					train.CurrentStationIndex++
					nextStation := path[train.CurrentStationIndex]
					train.CurrentStationName = nextStation
					turnGroup = append(turnGroup, train.Name + "-" + nextStation)
				} else {
					break
				}
			}
			
		}
		turnSchedule = append(turnSchedule, turnGroup)
	}
	fmt.Println("TURN SCHEDULE:")
	fmt.Println(turnSchedule)
	return turnSchedule
}

func printSchedule(turnSchedule [][]string) {
	var scheduleStr string
	for _, turn := range turnSchedule {
		var line string
		for _, train := range turn {
			line += train + " "
		}
		scheduleStr += line + "\n"
	}

	fmt.Println(scheduleStr)
}
