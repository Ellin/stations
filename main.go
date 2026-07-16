package main

import (
	"fmt"
	"os"
	"pathinder/internal/parsers"
	"pathinder/internal/pathfinding_alg"
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

	fmt.Println(pathfinding_alg.DinicAlg(&networkData, arg.StartStation, arg.EndStation))
	// start := time.Now()
	// bfs := pathfinding_alg.BFSAlg(&networkData, arg.StartStation, arg.EndStation)
	// fmt.Printf("\nprocess Time %dv \nmaybe> bfs %v\n\n", time.Since(start), bfs)

	// start = time.Now()
	// dfs := pathfinding_alg.DFSAlg(&networkData, arg.StartStation, arg.EndStation)
	// fmt.Printf("process Time %v \nmaybe> dfs %v\n", time.Since(start), dfs)

}

// exit progremm with 1 if argument err is not nil
func close(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// func timer(alg func()) {
// 	start := time.Now()
// 	route := alg()
// 	fmt.Println("process time: %v", time.Since(start))
// 	fmt.Println("maybe>", route)
// }
