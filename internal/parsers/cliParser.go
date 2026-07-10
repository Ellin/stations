package parsers

import (
	"errors"
	"flag"
	"fmt"
	"pathinder/model"
)

// using flag here in case we want to add more flag later on
// this function returns ArgsInfo struct  and error if any of the argument are not valid
func ParseArgs() (*model.ArgsInfo, error) {
	var argList model.ArgsInfo // this stores the argument from the cli

	flag.Parse()
	args := flag.Args()

	if len(args) < 4 {
		return nil, errors.New("Error: Insufficient number of arguments.\n")
	}
	TrainNumber, isValid := validateCoordinate(args[3])

	if err := validateMapFile(args[0]); err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	if !validateStationName(args[1]) {
		return nil, fmt.Errorf("Error: Invalid Start Station Name %s.\n", args[1])
	}
	if !validateStationName(args[2]) {
		return nil, fmt.Errorf("Error: Invalid End Station Name %s.n\n", args[2])
	}
	if !isValid {
		return nil, fmt.Errorf("Error: Invalid Train Number %s.\n", args[3])
	}

	argList.MapFile = args[0]
	argList.StartStation = args[1]
	argList.EndStation = args[2]
	argList.TrainCount = TrainNumber

	return &argList, nil
}
