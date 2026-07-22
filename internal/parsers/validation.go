package parsers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pathinder/model"
	"strconv"
	"unicode"
)

// validateStationName returns true if the station name meets the following criteria:
// Comprised of lower-case letters, numbers and underscores (_) only. No special characters or other punctuation are allowed.
func validateStationName(name string) bool {
	if name == "" {
		return false
	}

	for _, char := range name {
		if !(unicode.IsLower(char) || unicode.IsNumber(char) || char == '_') {
			return false
		}
	}

	return true
}

// validateCoordinate returns the coordinate as an int and true if valid
// Coordinate must be positive integer
func validateCoordinate(c string) (int, bool) {
	num, err := strconv.Atoi(c)

	if err != nil || num < 1 {
		return -1, false
	}

	return num, true
}

// validateMapFile returns error if file doesn't exist or file extention is not map else it returns nil
// it takes file name type string
func validateMapFile(file string) error {
	ex := filepath.Ext(file)
	if ex != ".map" {
		return fmt.Errorf("Invalid File Extension '%s', Allowed extentions = [.map].", ex)
	}
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return fmt.Errorf("File %s Does Not Exist.\n", file)
	}

	return nil
}

func ValidateStartAndEndStation(network *model.NetworkData, startStation, endStation string) error {
	if startStation == endStation {
		return errors.New("Error: Start Station And End Station Cannot Be The Same.")
	}
	_, startFound := network.StationMap[startStation]
	_, endFound := network.StationMap[endStation]
	_, StartConnectionFound := network.NetworkMap[startStation]
	_, endConnectionFound := network.NetworkMap[endStation]

	if !startFound {
		return errors.New("Error: Start Station Not Found In The Map.")
	}
	if !endFound {
		return errors.New("Error: End Station Not Found In The Map.")
	}
	if !StartConnectionFound {
		return errors.New("Error: Start Station Has No Connections.")
	}
	if !endConnectionFound {
		return errors.New("Error: End Station Has No Connections.")
	}
	return nil
}
