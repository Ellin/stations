package main

import (
	"fmt"
	"os"
	"pathinder/internal/parsers"
	"pathinder/internal/pathfinding_alg_test"
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

	res := pathfinding_alg_test.DFSAlg(&networkData, arg.StartStation, arg.EndStation)
	fmt.Println("maybe>", res)
}

// exit progremm with 1 if argument err is not nil
func close(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
