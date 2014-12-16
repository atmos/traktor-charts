package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func writeJSONFile(traktorData TraktorData) {
	data, err := json.MarshalIndent(traktorData, "", "  ")
	if err != nil {
		fmt.Println("Unable to marshal shit, yo", err)
	}

	jsonFile := os.ExpandEnv("${HOME}/.traktor-charts.json")
	fp, err := os.Create(jsonFile)
	if err != nil {
		fmt.Println("Unable to create", jsonFile)
	}
	defer fp.Close()

	fp.Write(data)

	fmt.Println("JSON files for your charts are in ~/.traktor-charts.json.")

	fp.Sync()
}
