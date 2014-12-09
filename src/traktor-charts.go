package main

import (
	"fmt"
)

func main() {
	fmt.Println("NI directories:", traktorDir(""))

	historyPaths, _ := traktorHistoryPaths(traktorDir(""))
	archiveFiles, _ := traktorArchiveFiles(historyPaths)

	db, err := initializeDB("tracks.db")
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

	db.Close()
}
