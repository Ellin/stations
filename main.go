package main

import (
	"fmt"
	"os"
	"pathinder/internal/parsers"
	"pathinder/internal/paths"
	"time"
)

const (
	reset = "\033[0m"
	green = "\033[32m"
)

func main() {
	// call ParseArgs return arg (type ArgsInfo struct) to make the data more packed rather then returning
	// individual arguments "ps. struct is in model file"
	arg, argumentsErr := parsers.ParseArgs()
	close(argumentsErr)

	if arg == nil {
		return
	}

	// extract the text from map file using arg.MapFile where the file name is stored
	fileText, err := ReadFile(arg.MapFile)
	close(err)

	networkData, networkErr := parsers.ParseNetworkMap(fileText)
	close(networkErr)

	stationErr := parsers.ValidateStartAndEndStation(&networkData, arg.StartStation, arg.EndStation)
	close(stationErr)

	networkData.Start, networkData.End = arg.StartStation, arg.EndStation

	graph := paths.Graph{}
	graph.CreateVertexMaps(networkData)
	graph.CreateEKGraph(networkData)
	// graph.PrintGraph(networkData)

	// Start timer for pathfinding and scheduling
	start := time.Now()
	_, schedulerErr := graph.RunScheduler(arg.StartStation, arg.EndStation, arg.Algo, arg.TrainCount)
	close(schedulerErr)

	// Print algorithm results and performance
	fmt.Println("Algorithm used:", arg.Algo)
	fmt.Printf(green+"Execution Time: %v\n"+reset, time.Since(start))
}

// close exits program with exit code 1 if argument err is not nil
func close(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
