package main

import (
	"fmt"
	"os"
)

// readFile using profided filepath and return error or string data
func ReadFile(filepath string) (string, error) {
	mapdata, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("Error : Could Not Read Map File %s Error:%v", filepath, err)
	}
	return string(mapdata), nil
}
