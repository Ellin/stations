package model

// used for argument data handling instead of being seprate global variables
type ArgsInfo struct {
	MapFile      string
	StartStation string
	EndStation   string
	TrainCount   int
	CliCom       string
}
