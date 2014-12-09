package main

import (
	"database/sql"
	"fmt"
	"time"
)

func displayOutput(db *sql.DB) {
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
}

func lineDelimiter() string {
	return `+---------------------------------------------------------------------------------------------+`
}
