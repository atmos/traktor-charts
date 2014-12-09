package main

import (
	"fmt"
	"time"
)

func lineDelimiter() string {
	return `+---------------------------------------------------------------------------------------------+`
}

func main() {
	fmt.Println("NI directories:", traktorDir(""))

	historyPaths, _ := traktorHistoryPaths(traktorDir(""))
	archiveFiles, _ := traktorArchiveFiles(historyPaths)

	db, err := initializeDB("tracktor-charts.db")
	if err != true {
		fmt.Println("Error initializing db", err)
	}

	fileCount := 0
	for _, fileName := range archiveFiles {
		entries, _ := traktorParseFile(fileName)
		for _, entry := range entries.EntryList {
			insertEntry(db, entries, entry)
		}
		fileCount++
	}
	fmt.Println("Found", fileCount, "NML files")

	totalPlays := countForTable(db, "plays")
	totalTracks := countForTable(db, "tracks")
	fmt.Println("Found", totalTracks, "unique tracks.")
	fmt.Println("For a total of", totalPlays, "plays.")

	for month := 1; month < 13; month++ {
		monthlyEntries := findChartEntriesByMonth(db, month)
		fmt.Println("Status for:", time.Month(month))
		fmt.Println(lineDelimiter())
		for _, chartEntry := range monthlyEntries {
			output := fmt.Sprintf("| %-28s | %-55s | %-02d |", chartEntry.Artist, chartEntry.Title, chartEntry.Count)
			fmt.Println(output)
		}
		fmt.Println(lineDelimiter())
	}

	db.Close()
}
