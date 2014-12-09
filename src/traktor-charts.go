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
	for key, fileName := range archiveFiles {
		fmt.Println("Found an NML File:" + key)
		entries, _ := traktorParseFile(fileName)
		for _, entry := range entries.EntryList {
			fmt.Println("Inserting", entry.Title)
			insertEntry(db, entries, entry)
		}
		fmt.Printf("Found %d entries\n", len(entries.EntryList))
		fileCount++
	}
	fmt.Println("Found", fileCount, "NML files")

	db.Close()
}
