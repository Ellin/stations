package model

// used for argument data handling instead of being seprate global variables
type ArgsInfo struct {
	MapFile      string
	StartStation string
	EndStation   string
	TrainCount   int
	Algo         string
	CliCom       string
}

type NetworkData struct {
	StationMap map[StationName]Station
	NetworkMap map[StationName]map[StationName]struct{} // The network map links each station to a set of all connecting stations
	Start      StationName
	End        StationName
}

type StationName = string

type Station struct {
	Name string
	X    int
	Y    int
}
