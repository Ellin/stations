package main

import (
	"strings"
	"unicode"
	"strconv"
	"fmt"
)

// parseNetworkMap extracts the station and connections data from the network map, creating a Station struct for each station and a connections map linking each connected station to each other
func parseNetworkMap(text string) {
	lines := strings.Split(text, "\n")

	var parsingSection string
	var stationsBuffer []string // contains lines within station section to be parsed
	var connectionsBuffer []string // contains lines within connections section to be parsed

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "stations:" {
			parsingSection = "stations"
			continue
		}

		if line == "connections:" {
			parsingSection = "connections"
			continue
		}

		if parsingSection == "stations" {
			stationsBuffer = append(stationsBuffer, line)
		}

		if parsingSection == "connections" {
			connectionsBuffer = append(connectionsBuffer, line)
		}
	}

	parseStations(stationsBuffer)
	parseConnections(connectionsBuffer)
}

type Station struct {
	name string
	x int
	y int
}

// parseStations parses the stations section of the network map and returns a []Station
func parseStations(lines []string) []Station {
	var stations []Station

	for _, line := range lines {
		line,_,_ = strings.Cut(line, "#") // Remove comments

		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := trimSpaceSlice(strings.Split(line, ","))

		if len(parts) != 3 { // Skip malformed station data
			fmt.Printf("Invalid station line: %v\n", line)
			continue
		}

		name, x, y := parts[0], parts[1], parts[2]

		if !validateStationName(name) {
			fmt.Printf("Invalid station name: %v\n", name)
			continue
		}

		xInt, isXValid := validateCoordinate(x)
		yInt, isYValid := validateCoordinate(y)
	
		if !isXValid || !isYValid {
			fmt.Printf("Invalid coordinates: %v, %v\n", x, y)
			continue
		}

		stations = append(stations, Station{name, xInt, yInt})
	}

	fmt.Printf("Valid stations:\n%v\n", stations)
	return stations
}

// parseConnections parses the connections section of the network map and creates a connections map linking all connected stations
func parseConnections(lines []string) map[string]map[string]struct{} {
	// The connections map is a map where each station (name) is a key and the value is a set of all its connecting stations
	connectionsMap := make(map[string]map[string]struct{})

	for _, line := range lines {
		line,_,_ := strings.Cut(line, "#") // Remove comments

		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := trimSpaceSlice(strings.Split(line, "-"))

		if len(parts) != 2 {
			fmt.Printf("Invalid connection line: %v\n", line)
			continue
		}

		start, end := parts[0], parts[1]

		if !validateStationName(start) || !validateStationName(end) {
			fmt.Printf("Malformed station names: %v\n", line)
			continue
		}

		// Add each start/end station as a connection to each other in the connectionsMap
		if connectionsMap[start] == nil {
			connectionsMap[start] = make(map[string]struct{})
		}
		if connectionsMap[end] == nil {
			connectionsMap[end] = make(map[string]struct{})
		}
		connectionsMap[start][end] = struct{}{}
		connectionsMap[end][start] = struct{}{}
	}

	fmt.Printf("\nConnections map:\n%v", connectionsMap)
	return connectionsMap
}

// trimSpaceSlice applies strings.TrimSpace on each string in []string
func trimSpaceSlice(s []string) []string {
	trimmed := make([]string, 0, len(s))
	for _, v := range s {
		v = strings.TrimSpace(v)
		trimmed = append(trimmed, v)	
	}
	return trimmed
}

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

func main() {
	text := `stations:
		# south stations
		waterloo  , 3 , 1
		victoria,6,7

		# north stations
		euston,11,23
		st_pancras,5,15 # international
		inValidStation, 1, 2

		connections:
		waterloo -victoria
		waterloo- euston  
		st_pancras-euston
		victoria-st_pancras`

	parseNetworkMap(text)
}
