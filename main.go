package main

import (
	"fmt"
	"os"
	"pathinder/internal/parsers"
	"pathinder/internal/pathfinding_alg"
	"pathinder/internal/paths"
	"time"
)

func main() {
	//call ParseArgs return arg (type ArgsInfo struct) to make the data more packed rather then returning
	// individual arguments "ps. struct is in model file"
	arg, argumentsErr := parsers.ParseArgs()
	close(argumentsErr)

	// fmt.Println("returned str arg:", arg, err)
	// extract the text from map file using arg.MapFile where the file name is stored
	fileText, err := ReadFile(arg.MapFile)
	close(err)

	networkData, networkErr := parsers.ParseNetworkMap(fileText)
	close(networkErr)
	// fmt.Println(networkData.NetworkMap)
	// fmt.Println("con")
	// fmt.Println(networkData.StationMap)
	StationErr := parsers.ValidateStartAndEndStation(&networkData, arg.StartStation, arg.EndStation)
	close(StationErr)

	networkData.Start, networkData.End = arg.StartStation, arg.EndStation
	paths.Bfs(networkData.NetworkMap, arg.StartStation, arg.EndStation)

	graph := pathfinding_alg.Graph{}
	graph.CreateVertexMaps(networkData)
	// fmt.Printf("\nVertexIDMap: %#v\n", graph.VertexIDMap)
	graph.CreateEKGraph(networkData)
	// graph.PrintGraph(networkData)
	// fmt.Println("\n\n", graph)

	paths, err := graph.DinicAlg(arg.StartStation, arg.EndStation)
	close(err)

	for _, p := range paths {
		fmt.Println("\n\n> path: ", p)
	}

}

// exit progremm with 1 if argument err is not nil
func close(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func timer(alg func()) {
	start := time.Now()
	alg()
	fmt.Printf("\033[32m Execution Time: %v \033[0m\n", time.Since(start))
}
