package main

import (
	"strings"
	"unicode"
	"strconv"
	"errors"
	"fmt"
)

// parseNetworkMap extracts the station and connections data from the network map, creating a Station struct for each station and a connections map linking each connected station to each other
func parseNetworkMap(text string) error {
	lines := strings.Split(text, "\n")

	var parsingSection string
	var stationsBuffer []string // contains lines within station section to be parsed
	var connectionsBuffer []string // contains lines within connections section to be parsed
	var stationLineNum int // the line number in the text at which the "stations:" section starts
	var connectionsLineNum int // the line number in the text at which the "connections:" section starts
	
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if line == "stations:" {
			parsingSection = "stations"
			stationLineNum = i + 1
			continue
		}

		if line == "connections:" {
			parsingSection = "connections"
			connectionsLineNum = i + 1
			continue
		}

		if parsingSection == "stations" {
			stationsBuffer = append(stationsBuffer, line)
		}

		if parsingSection == "connections" {
			connectionsBuffer = append(connectionsBuffer, line)
		}
	}

	if stationLineNum == 0 {
		return errors.New("Error: Missing \"stations:\" section")
	}

	if connectionsLineNum == 0 {
		return errors.New("Error: Missing \"connections:\" section")
	}

	_, err1 := parseStations(stationsBuffer, stationLineNum)
	_, err2 := parseConnections(connectionsBuffer, connectionsLineNum)
	fmt.Println(err1)
	fmt.Println(err2)

	return errors.Join(err1, err2)
}

type Station struct {
	name string
	x int
	y int
}

// parseStations parses the stations section of the network map and returns a []Station
func parseStations(lines []string, lineNum int) ([]Station, error) {
	var stations []Station
	var errs []error

	for i, line := range lines {
		var lineHasError bool

		line,_,_ = strings.Cut(line, "#") // Remove comments

		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := trimSpaceSlice(strings.Split(line, ","))

		if len(parts) != 3 { // Skip malformed station data
			errs = append(errs, fmt.Errorf("Malformed station format on line #%d: %s\n", lineNum + i + 1, line))
			continue
		}

		name, x, y := parts[0], parts[1], parts[2]

		if !validateStationName(name) {
			errs = append(errs, fmt.Errorf("Invalid station name in line #%d: %v\n", lineNum + i + 1, name))
			lineHasError = true
		}

		xInt, isXValid := validateCoordinate(x)
		yInt, isYValid := validateCoordinate(y)
	
		if !isXValid || !isYValid {
			errs = append(errs, fmt.Errorf("Invalid coordinates at station %s in line #%d: %v, %v\n", name, lineNum + i + 1, x, y))
			continue
		}

		if !lineHasError {
			stations = append(stations, Station{name, xInt, yInt})
		}
	}

	fmt.Printf("Valid stations:\n%v\n", stations)
	return stations, errors.Join(errs...)
}

// parseConnections parses the connections section of the network map and creates a connections map linking all connected stations
func parseConnections(lines []string, lineNum int) (map[string]map[string]struct{}, error) {
	// The connections map is a map where each station (name) is a key and the value is a set of all its connecting stations
	connectionsMap := make(map[string]map[string]struct{})
	var errs []error

	for i, line := range lines {
		line,_,_ := strings.Cut(line, "#") // Remove comments

		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := trimSpaceSlice(strings.Split(line, "-"))

		if len(parts) != 2 {
			errs = append(errs, fmt.Errorf("Malformed connection format in line #%d: %s\n", lineNum + i + 1, line))
			continue
		}

		start, end := parts[0], parts[1]

		if !validateStationName(start) || !validateStationName(end) {
			errs = append(errs, fmt.Errorf("Malformed station names in line #%d: %s\n", lineNum + i + 1, line))
			continue
		}

		if start == end {
			errs = append(errs, fmt.Errorf("Error on line #%d. Start and end connections are the same: %s", lineNum + i + 1, line))
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

	fmt.Printf("\nConnections map:\n%v\n\n", connectionsMap)
	return connectionsMap, errors.Join(errs...)
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
		euston,11,-23
		st_pancras,5,15 # international
		inValidStation, 1, 2

		connections:
		waterloo -victoria
		waterloo- euston  
		st_pancras-euston
		victoria-st_pancras
		toronto--waterloo
		toronto-WATERLOO
		new_york - new_york
		`

	parseNetworkMap(text)
}
