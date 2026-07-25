// Scheduler using paths found using Edmonds-Karp

package paths

import (
	"fmt"
	"pathinder/internal/parsers"
	"strconv"
	"strings"
)

type Train struct {
	Name                string              // e.g. "T1"
	PathID              int                 // Index of the path within a pathSet the train is scheduled to use
	CurrentStationIndex int                 // e.g. path: ["waterloo", "victoria", "st_pancras"] ; "waterloo" -> index 0
	CurrentStationName  parsers.StationName // e.g. "waterloo"
}

func (g *Graph) RunScheduler(start, end, alg parsers.StationName, numTrains int) error {

	var pathSet [][]parsers.StationName
	var err error

	switch alg {
	case "EdmondsKarp":
		if pathSet, err = g.FindPathSet(start, end, numTrains); err != nil {
			return err
		}
	case "Dinic":
		if pathSet, err = g.DinicAlg(start, end, numTrains); err != nil {
			return err
		}
	}
	pathAssignments := divideTrains(numTrains, pathSet)
	turnSchedule := runTrains(end, pathSet, pathAssignments)
	printSchedule(turnSchedule)
	fmt.Println("Trains scheduled successfully!")
	fmt.Printf("%d turns to move %d trains from %s to %s using the path set of %d non-overlapping paths:\n%v\n\n", len(turnSchedule), numTrains, start, end, len(pathSet), pathSet)
	return nil
}

// ! NOTE: This should be further optimized to take into account cost given some # of trains, not just max flow
// For example, the current algorithm for choosing a set of paths for 2 trains prefers a set of two very long non-overlapping paths
// vs one very short path that blocks multiple non-overlapping paths.
// But that short blocking path may still be the better path for certain number of trains.
func (g *Graph) FindPathSet(start, end parsers.StationName, numTrains int) ([][]parsers.StationName, error) {

	maxFlow, pathSets := g.EdmondsKarp(start, end, numTrains)

	if maxFlow == 0 {
		return nil, fmt.Errorf("Error: No Path From Start To End Stations")
	}

	var bestSetIndex int
	minTurns := calcAvgTurns(numTrains, pathSets[0])

	for i, pathSet := range pathSets {
		// find avg # of turns per set
		avgTurns := calcAvgTurns(numTrains, pathSet)

		if avgTurns < minTurns {
			minTurns = avgTurns
			bestSetIndex = i
		}
	}

	pathSet := pathSets[bestSetIndex]

	// fmt.Printf("\nCHOSEN PATH SET (%d paths):\n%v\n", len(pathSet), pathSet)

	return pathSet, nil
}

// calcAvgTurns calculate the average number of turns per path in a given pathSet
func calcAvgTurns(numTrains int, pathSet [][]parsers.StationName) int {
	var totalHops int
	for _, path := range pathSet {
		totalHops += len(path) - 1
	}

	// average # of turns per path = ((total # of hops across all paths + total # of trains) / # of paths) - 1 turn [since first set of trains doesn't need to wait]
	// rounding up a/b trick: a + b - 1 / b
	avgTurns := ((totalHops + numTrains + len(pathSet) - 1) / len(pathSet)) - 1

	return avgTurns
}

// Decide which train gets assigned to which path
// The index of each train group assigned to a path in pathAssignments = index of that path in pathSet
func divideTrains(numTrains int, pathSet [][]parsers.StationName) (pathAssignments [][]Train) {
	pathAssignments = make([][]Train, len(pathSet))

	for i := 0; i < numTrains; i++ {

		shortestTurnNum := getNumTurns(0, pathSet[0], pathAssignments)
		shortestPathIndex := 0

		// Find path with shortest # turns
		for pathID, path := range pathSet {
			numTurns := getNumTurns(pathID, path, pathAssignments)
			if numTurns < shortestTurnNum {
				shortestTurnNum = numTurns
				shortestPathIndex = pathID
			}
		}

		// Schedule train
		pathAssignments[shortestPathIndex] = append(pathAssignments[shortestPathIndex], Train{
			Name:                "T" + strconv.Itoa(i+1),
			PathID:              shortestPathIndex,
			CurrentStationIndex: 0,
			CurrentStationName:  pathSet[shortestPathIndex][0],
		})
	}

	return pathAssignments
}

// func divideTrainsImproves(numTrains int, pathSet [][]parsers.StationName) map[int][]Train {
// 	pathMap := make(map[int][]Train)
// 	return pathMap
// }

// getNumTurns returns the number of turns it would take to use a certain path, taking into account the wait time for that path (# of trains already scheduled for that path)
func getNumTurns(pathID int, path []parsers.StationName, pathAssignments [][]Train) int {
	numHops := len(path) - 1 // hops between stations

	if len(pathAssignments[pathID]) == 0 {
		return numHops
	}

	numWait := len(pathAssignments[pathID]) - 1 // -1 because the first train doesn't need a turn of waiting

	return numHops + numWait
}

func runTrains(end parsers.StationName, pathSet [][]parsers.StationName, pathAssignments [][]Train) [][]string {
	var turnSchedule [][]string // e.g. [["T1-a", "T2-b"]["T1-c", "T2-d"]["T1-end","T2-end"]]

	// Path index 0 = shortest path in a pathSet
	// The shortest path in a pathSet will always have the most number of trains scheduled for it

	// While there are trains left to move
	for len(pathAssignments[0]) > 0 {
		var turnGroup []string

		// MOVE ALL TRAINS IN LINE FOR ALL PATHS TO THE NEXT STATION IF NOT MOVING AT SAME STATION TO THE ONE AHEAD
		// IF STATION IS REACHED, DEQUEUE THAT TRAIN
		for pathID, path := range pathSet {

			var canMove bool

			for i := 0; i < len(pathAssignments[pathID]); i++ {
				train := &pathAssignments[pathID][i]

				if train.CurrentStationName == end {
					pathAssignments[pathID] = pathAssignments[pathID][1:] // dequeue train
					i--
					continue
				}

				// ALL TRAINS FIRST IN LINE CAN MOVE ONE STATION FORWARD
				if i == 0 {
					canMove = true

				}

				//	IF NOT FIRST IN LINE, CHECK FIRST IF IT CAN MOVE FORWARD
				if i > 0 {
					prevStation := &pathAssignments[pathID][i-1]
					if path[train.CurrentStationIndex+1] != prevStation.CurrentStationName { // check if train in queue would be moving to the same station as the train ahead
						canMove = true
					} else {
						canMove = false
					}
				}

				if canMove {
					train.CurrentStationIndex++
					nextStation := path[train.CurrentStationIndex]
					train.CurrentStationName = nextStation
					turnGroup = append(turnGroup, train.Name+"-"+nextStation)
				} else {
					break
				}
			}

		}

		turnSchedule = append(turnSchedule, turnGroup)
	}

	turnSchedule = turnSchedule[:len(turnSchedule)-1] // Remove last empty turnGroup

	return turnSchedule
}

func printSchedule(turnSchedule [][]string) {
	var builder strings.Builder

	for _, turn := range turnSchedule {
		for _, train := range turn {
			builder.WriteString(train)
			builder.WriteString(" ")
		}
		builder.WriteString("\n")
	}

	scheduleStr := builder.String()

	fmt.Printf("\nTURN SCHEDULE (%v turns):\n%s\n", len(turnSchedule), scheduleStr)
}
