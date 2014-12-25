package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
		for _, entry := range entries.TraktorXMLEntryList {
			insertEntry(db, entries, entry)
		}
		fileCount++
	}
	fmt.Println("Found", fileCount, "archive files")

	writeMarkdownFile(getTraktorData(db))

	jsonBytes := getExportData(db)
	httpPostResults(jsonBytes)

	fmt.Println("Your charts are in ~/.traktor-charts.md.")
	fmt.Println("You should share them on https://gist.github.com")

	fmt.Println("Run 'cat ~/.traktor-charts.md | pbcopy' in your terminal and paste into a new gist.")

	db.Close()
}

func httpPostResults(traktorBody []byte) {
	url := "https://djcharts.io/api/import"
	fmt.Println("URL:>", url)

	token, _ := ioutil.ReadFile(os.ExpandEnv("${HOME}/.traktor-charts"))
	basicAuthToken := strings.TrimSuffix(string(token), "\n")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(traktorBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("X", basicAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(body))
}
