package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func requiresUpdate(count int) bool {
	shouldUpdate := false
	countFile := os.ExpandEnv("${HOME}/.traktor-charts.count")

	oldCount, fileErr := ioutil.ReadFile(countFile)
	if fileErr != nil {
		shouldUpdate = true
	} else {
		oldFileCount, _ := strconv.Atoi(string(oldCount))
		fmt.Printf("Found %d old entries\n", oldFileCount)
		if count != oldFileCount {
			shouldUpdate = true
		}
	}
	_ = ioutil.WriteFile(countFile, []byte(strconv.Itoa(count)), 0600)
	return shouldUpdate
}

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

	jsonBytes := getExportData(db)
	db.Close()

	if requiresUpdate(fileCount) {
		httpPostResults(jsonBytes)
	} else {
		fmt.Println("No new traktor archive files found")
		os.Exit(3)
	}

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
