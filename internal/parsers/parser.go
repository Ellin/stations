package parsers

import (
	"errors"
	"fmt"

	// "fmt"
	"strings"
)

type NetworkData struct {
	StationMap map[StationName]Station
	NetworkMap map[StationName]map[StationName]struct{} // The network map links each station to a set of all connecting stations
}

type StationName = string

type Station struct {
	name      string
	x         int
	y         int
	Cap, Flow int
}

// parseNetworkMap extracts the station and connections data from the network map
// It returns a station map, connections map, and an error if data in the network map is malformed or invalid
func ParseNetworkMap(text string) (NetworkData, error) {
	lines := strings.Split(text, "\n")

	var parsingSection string
	var stationsBuffer []string    // contains lines within station section to be parsed
	var connectionsBuffer []string // contains lines within connections section to be parsed
	var stationLineNum int         // the line number in the text at which the "stations:" section starts
	var connectionsLineNum int     // the line number in the text at which the "connections:" section starts

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
		return NetworkData{}, errors.New("Error: Missing \"stations:\" section")
	}

	if connectionsLineNum == 0 {
		return NetworkData{}, errors.New("Error: Missing \"connections:\" section")
	}
	stationMap, err1 := ParseStations(stationsBuffer, stationLineNum)
	if err1 != nil {
		fmt.Println(err1)
		return NetworkData{}, err1
	}
	// fmt.Println("stations map", stationMap)
	networkMap, err2 := parseConnections(connectionsBuffer, connectionsLineNum, stationMap)
	if err2 != nil {
		fmt.Println(err2)
		return NetworkData{}, err2
	}

	networkData := NetworkData{stationMap, networkMap}

	return networkData, nil
}

// parseStations parses the stations section of the network map and returns a station map: map[StationName]Station
func ParseStations(lines []string, lineNum int) (map[StationName]Station, error) {
	stationMap := make(map[StationName]Station)
	var errs []error
	seenCoordinates := make(map[string]struct{})

	for i, line := range lines {
		var lineHasError bool

		line, _, _ = strings.Cut(line, "#") // Remove comments

		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := trimSpaceSlice(strings.Split(line, ","))

		if len(parts) != 3 { // Skip malformed station data
			errs = append(errs, fmt.Errorf("Malformed station format on line #%d: %s\n", lineNum+i+1, line))
			continue
		}

		name, x, y := parts[0], parts[1], parts[2]

		if !validateStationName(name) {
			errs = append(errs, fmt.Errorf("Invalid station name in line #%d: %v\n", lineNum+i+1, name))
			lineHasError = true
		}

		xInt, isXValid := validateCoordinate(x)
		yInt, isYValid := validateCoordinate(y)

		if !isXValid || !isYValid {
			errs = append(errs, fmt.Errorf("Invalid coordinates at station %s in line #%d: %v, %v\n", name, lineNum+i+1, x, y))
			continue
		}

		// Check for duplicate coordinates
		var coordinates string = x + "," + y
		if _, ok := seenCoordinates[coordinates]; ok {
			errs = append(errs, fmt.Errorf("Duplicate coordinates in line #%d: %s\n", lineNum+i+1, coordinates))
			continue
		}
		seenCoordinates[coordinates] = struct{}{}

		if !lineHasError {
			_, ok := stationMap[name]
			if ok {
				errs = append(errs, fmt.Errorf("Duplicate station name %s in line #%d\n", name, lineNum+i+1))
				continue
			}
			stationMap[name] = Station{name, xInt, yInt, 1, 0}
		}
	}

	// fmt.Printf("Valid stations:\n%v\n", stationMap)
	return stationMap, errors.Join(errs...)
}

// parseConnections parses the connections section of the network map and creates a connections map linking all connected stations
func parseConnections(lines []string, lineNum int, stationMap map[StationName]Station) (map[StationName]map[StationName]struct{}, error) {
	// The connections map is a map where each station (name) is a key and the value is a set of all its connecting stations
	connectionsMap := make(map[StationName]map[StationName]struct{})
	var errs []error

	for i, line := range lines {
		line, _, _ := strings.Cut(line, "#") // Remove comments

		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := trimSpaceSlice(strings.Split(line, "-"))

		if len(parts) != 2 {
			errs = append(errs, fmt.Errorf("Malformed connection format in line #%d: %s\n", lineNum+i+1, line))
			continue
		}

		start, end := parts[0], parts[1]

		if !validateStationName(start) || !validateStationName(end) {
			errs = append(errs, fmt.Errorf("Malformed station names in line #%d: %s\n", lineNum+i+1, line))
			continue
		}

		// Check that the stations exist
		_, startExists := stationMap[start]
		_, endExists := stationMap[end]
		if !startExists || !endExists {
			if !startExists {
				errs = append(errs, fmt.Errorf("Non-existent start station in line #%d: %s\n", lineNum+i+1, start))
			}
			if !endExists {
				errs = append(errs, fmt.Errorf("Non-existent end station in line #%d: %s\n", lineNum+i+1, end))
			}
			continue
		}

		if start == end {
			errs = append(errs, fmt.Errorf("Error on line #%d. Start and end connections are the same: %s", lineNum+i+1, line))
			continue
		}

		if connectionsMap[start] == nil {
			connectionsMap[start] = make(map[StationName]struct{})
		}
		if connectionsMap[end] == nil {
			connectionsMap[end] = make(map[StationName]struct{})
		}

		// Check duplicate connections
		_, ok := connectionsMap[start][end]
		if ok {
			errs = append(errs, fmt.Errorf("Error on line #%d. Duplicate connections between %s and %s", lineNum+i+1, start, end))
			continue
		}

		// Add each start/end station as a connection to each other in the connectionsMap
		connectionsMap[start][end] = struct{}{}
		connectionsMap[end][start] = struct{}{}
	}

	// fmt.Printf("\nConnections map:\n%v\n\n", connectionsMap)
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
