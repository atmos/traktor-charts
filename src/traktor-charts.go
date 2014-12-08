package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("NI directories:", traktorDir(""))

	db, err := sql.Open("sqlite3", "tracks.db")
	if err != nil {
		fmt.Println("Error opening db file:", err)
		return
	}
	res, _ := db.Exec("CREATE TABLE 'tracks' id INTEGER, audio_id STRING")
	fmt.Sprintf("%v\n", res)

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
