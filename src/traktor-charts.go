package main

import (
	"fmt"
)

func main() {
	fmt.Println("NI directories:", traktorDir(""))

	historyPaths, _ := traktorHistoryPaths(traktorDir(""))
	archiveFiles, _ := traktorArchiveFiles(historyPaths)

	db, err := initializeDB("traktor-charts.db")
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
	fmt.Println("Found", fileCount, "archive files")

	writeOutputFrom(db)

	fmt.Println("Your charts are in ~/.traktor-charts.md.")
	fmt.Println("You should share them on https://gist.github.com")

	fmt.Println("Run 'cat ~/.traktor-charts.md | pbcopy' in your terminal and paste into a new gist.")

	db.Close()
}
