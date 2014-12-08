package main

import (
	"fmt"
)

func main() {
	fmt.Println("NI directories:", traktorDir(""))

	historyPaths, _ := traktorHistoryPaths(traktorDir(""))
	archiveFiles, _ := traktorArchiveFiles(historyPaths)

	fileCount := 0
	for key, fileName := range archiveFiles {
		fmt.Println("Found an NML File:" + key)
		entries, _ := traktorParseFile(fileName)
		for _, entry := range entries.EntryList {
			fmt.Printf("%v\n", entry)
		}
		fmt.Printf("Found %d entries\n", len(entries.EntryList))
		fileCount++
	}
	fmt.Println("Found", fileCount, "NML files")
}
