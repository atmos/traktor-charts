package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

// Exported data types for formatting
type TraktorData struct {
	Plays  int
	Tracks int
	ByYear []TraktorDataByYear
}

type TraktorDataByYear struct {
	Year    int
	Charts  []ChartEntry
	ByMonth []TraktorByYearAndMonth
}

type TraktorByYearAndMonth struct {
	Year  int
	Month int

	Charts []ChartEntry
}

// XML data types for archive parsing
type TraktorXMLEntryInfo struct {
	Key    string `xml:"KEY,attr"`
	Genre  string `xml:"GENRE,attr"`
	Length int    `xml:"PLAYTIME,attr"`
}

type TraktorXMLEntryTempo struct {
	Bpm int `xml:"BPM,attr"`
}

type TraktorXMLEntry struct {
	Info    TraktorXMLEntryInfo  `xml:"INFO"`
	Tempo   TraktorXMLEntryTempo `xml:"TEMPO"`
	Title   string               `xml:"TITLE,attr"`
	Artist  string               `xml:"ARTIST,attr"`
	AudioId string               `xml:"AUDIO_ID,attr"`
}

func (e TraktorXMLEntry) Key() string {
	return e.Info.Key
}

func (e TraktorXMLEntry) Genre() string {
	return e.Info.Genre
}

func (e TraktorXMLEntry) Length() int {
	return e.Info.Length
}

func (e TraktorXMLEntry) Bpm() int {
	return e.Tempo.Bpm
}

func (e TraktorXMLEntry) String() string {
	return fmt.Sprintf("%s - %s - %s", e.Artist, e.Title, e.Info.Genre)
}

type TraktorXMLEntryCollection struct {
	XMLName             xml.Name          `xml:"NML"`
	TraktorXMLEntryList []TraktorXMLEntry `xml:"COLLECTION>ENTRY"`
	Year                int
	Month               int
	Day                 int
	Hour                int
	Minute              int
}

func traktorDir(s string) string {
	return os.ExpandEnv("${HOME}/Documents/Native Instruments") + "/" + s
}

func traktorFilenameExtract(s string, e *TraktorXMLEntryCollection) {
	pattern := `.*/history_(\d+)y(\d+)m(\d+)d_(\d+)h(\d+)m.*`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 0 {
		e.Year, _ = strconv.Atoi(matches[1])
		e.Month, _ = strconv.Atoi(matches[2])
		e.Day, _ = strconv.Atoi(matches[3])
		e.Hour, _ = strconv.Atoi(matches[4])
		e.Minute, _ = strconv.Atoi(matches[5])
	}
}

func traktorParseFile(s string) (TraktorXMLEntryCollection, bool) {
	var entries TraktorXMLEntryCollection

	traktorFilenameExtract(s, &entries)

	xmlFile, err := os.Open(s)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return entries, false
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(b, &entries)
	return entries, false
}

func traktorHistoryPaths(baseDir string) ([]string, bool) {
	var dirs []string

	paths, err := ioutil.ReadDir(baseDir)
	if err != nil {
		fmt.Println("Error reading your NI directories:", err)
		return dirs, false
	}
	for _, dirName := range paths {
		traktorVersionDir := baseDir + "/" + dirName.Name()
		matched, _ := regexp.MatchString("Traktor \\d\\.\\d\\.\\d", traktorVersionDir)
		if matched {
			dirs = append(dirs, traktorVersionDir)
		}
	}
	return dirs, true
}

func traktorArchiveFiles(dirs []string) (map[string]string, bool) {
	files := make(map[string]string)

	for version, dirName := range dirs {
		nmlFiles, err := ioutil.ReadDir(dirName + "/History")
		if err != nil {
			fmt.Println("Error reading your NI", version, "history directories:", err)
			return files, false
		}

		for _, nmlFileName := range nmlFiles {
			fsName := nmlFileName.Name()
			nmlFullFileName := dirName + "/History/" + fsName
			matched, _ := regexp.MatchString("\\.nml$", nmlFullFileName)
			if matched {
				files[fsName] = nmlFullFileName
			}
		}
	}
	return files, true
}
