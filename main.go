package main

import (
	"fmt"
	"os"
	"pathinder/internal/parsers"
)

func main() {
	//call ParseArgs return arg (type ArgsInfo struct) to make the data more packed rather then returning
	// individual arguments "ps. struct is in model file"
	arg, err := parsers.ParseArgs()
	close(err)

	// fmt.Println("returned str arg:", arg, err)
	// extract the text from map file using arg.MapFile where the file name is stored
	fileText, err := ReadFile(arg.MapFile)

	close(err)
	parsers.ParseNetworkMap(fileText)
}

// exit progremm with 1 if argument err is not nil
func close(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
