// Scheduler using paths found using Edmonds-Karp

package paths

import (
	"fmt"
	"pathinder/internal/parsers"
	"strconv"
	"strings"
	"log"
)

type Train struct {
	Name                string              // e.g. "T1"
	PathID              int                 // Index of the path within a pathSet the train is scheduled to use
	CurrentStationIndex int                 // e.g. path: ["waterloo", "victoria", "st_pancras"] ; "waterloo" -> index 0
	CurrentStationName  parsers.StationName // e.g. "waterloo"
}

func (g *Graph) RunScheduler(start, end parsers.StationName, alg string, numTrains int) error {

	var pathSet [][]parsers.StationName
	var err error
	var minTurns int

	switch alg {
	case "EdmondsKarp":
		if pathSet, minTurns, err = g.FindPathSet(start, end, numTrains); err != nil {
			return err
		}
	case "Dinic":
		if pathSet, err = g.DinicAlg(start, end, numTrains); err != nil {
			return err
		}
		// fmt.Println("PATH SET FOUND")
		// fmt.Println(pathSet)
	}
	pathAssignments := divideTrains(numTrains, pathSet, minTurns)
	turnSchedule := runTrains(end, pathSet, pathAssignments)
	printSchedule(turnSchedule)
	return nil
}

// ! NOTE: This should be further optimized to take into account cost given some # of trains, not just max flow
// For example, the current algorithm for choosing a set of paths for 2 trains prefers a set of two very long non-overlapping paths
// vs one very short path that blocks multiple non-overlapping paths.
// But that short blocking path may still be the better path for certain number of trains.
func (g *Graph) FindPathSet(start, end parsers.StationName, numTrains int) ([][]parsers.StationName, int, error) {

	maxFlow, pathSets := g.EdmondsKarp(start, end, numTrains)

	if maxFlow == 0 {
		return nil, 0, fmt.Errorf("Error: No Path From Start To End Stations")
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

	fmt.Printf("\nCHOSEN PATH SET (%d paths):\n%v\n", len(pathSet), pathSet)

	return pathSet, minTurns, nil
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
func divideTrains(numTrains int, pathSet [][]parsers.StationName, avgTurns int) (pathAssignments [][]Train) {
	pathAssignments = make([][]Train, len(pathSet))
	var lastTrainID int
	var totalTrainsAssigned int

	for pathID, path := range pathSet {
		numHops := len(path) - 1

		if numHops <= avgTurns {
			// The number of trains to add to a path's queue to reach the avg number of turns
			// uses the same formula for calculating the # of turns for a path, but solving for the number of trains:
			// # turns for a path = # hops + # trains in queue - 1 wait turn [first train in queue doesn't increase # of turns]
			// ==> # trains = # turns - # hops + 1

			numTrainsToAdd := avgTurns - numHops + 1

			// Add trains
			for i := 0; i < numTrainsToAdd && totalTrainsAssigned < numTrains; i++ {
				totalTrainsAssigned++

				pathAssignments[pathID] = append(pathAssignments[pathID], Train{
					Name:                "T" + strconv.Itoa(i + lastTrainID + 1),
					PathID:              pathID,
					CurrentStationIndex: 0,
					CurrentStationName:  path[0],
				})
			}
		}
		currentPathTrains := pathAssignments[pathID]
		lastTrainID += len(currentPathTrains)
	}

	if totalTrainsAssigned != numTrains {
		log.Fatalf("Wrong num trains assigned: %d", totalTrainsAssigned)
	}
	return pathAssignments
}

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
