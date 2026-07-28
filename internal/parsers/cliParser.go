package parsers

import (
	"errors"
	"flag"
	"fmt"
	"pathfinder/model"
	"strings"
)

// using flag here in case we want to add more flag later on
// this function returns ArgsInfo struct  and error if any of the argument are not valid
func ParseArgs() (*model.ArgsInfo, error) {
	var argList model.ArgsInfo // this stores the argument from the cli

	algo := flag.String("alg", "EdmondsKarp", "Use flag -alg to use dinic Algorithm.")
	flag.Usage = func() {
		fmt.Println("Usage: go run . [path to file containing network map] [start station] [end station] [number of trains]")
		fmt.Println()
		fmt.Println("Flags: ")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println(`  go run . network/map2.map tempere helsinki 10`)
		fmt.Println(`  go run . -alg Dinic network/map2.map tempere helsinki 10`)
		fmt.Println(`  go run . network/map2.map hanko kokkola 20`)
	}

	flag.Parse()

	args := flag.Args()

	if len(args) != 4 {
		return nil, errors.New("Error: Incorrect Length Of Arguments. To See The Usage run: \n\n\t\t'go run . -h'\n")
	}

	TrainNumber, isValid := validateTrainNum(args[3])
	startStation, endStation := strings.ToLower(args[1]), strings.ToLower(args[2])

	if err := validateMapFile(args[0]); err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	if !validateStationName(startStation) {
		return nil, fmt.Errorf("Error: Invalid Start Station Name %s.\n", args[1])
	}
	if !validateStationName(endStation) {
		return nil, fmt.Errorf("Error: Invalid End Station Name %s.n\n", args[2])
	}
	if !isValid {
		return nil, fmt.Errorf("Error: Invalid Train Number %s. Must be a positive integer.\n", args[3])
	}

	if TrainNumber > 150000 {
		return nil, fmt.Errorf("Error: Invalid Train Numbers %s The Train Number Exceeded The Maximum Number Of Permitted Trains 150000.\n", args[3])
	}
	if TrainNumber == 0 {
		return nil, nil
	}
	argList.MapFile = args[0]
	argList.StartStation = startStation
	argList.EndStation = endStation
	argList.Algo = *algo
	argList.TrainCount = TrainNumber

	return &argList, nil
}
