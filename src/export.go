package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

func getExportData(db *sql.DB) []byte {
	tracks := getAllTracks(db)

	data, err := json.MarshalIndent(tracks, "", "  ")
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

	fmt.Println("JSON files for v2 are in ~/.traktor-charts.json.")

	fp.Sync()
	return data
}
