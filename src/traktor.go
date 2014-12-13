package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type Entry struct {
	Info    EntryInfo `xml:"INFO"`
	Title   string    `xml:"TITLE,attr"`
	Artist  string    `xml:"ARTIST,attr"`
	AudioId string    `xml:"AUDIO_ID,attr"`
}

type EntryInfo struct {
	Genre string `xml:"GENRE,attr"`
}

func (e Entry) String() string {
	return fmt.Sprintf("%s - %s - %s", e.Artist, e.Title, e.Info.Genre)
}

type EntryCollection struct {
	XMLName   xml.Name `xml:"NML"`
	EntryList []Entry  `xml:"COLLECTION>ENTRY"`
	Year      int
	Month     int
	Day       int
	Hour      int
	Minute    int
}

func traktorDir(s string) string {
	return os.ExpandEnv("${HOME}/Documents/Native Instruments") + "/" + s
}

func traktorFilenameExtract(s string, e *EntryCollection) {
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

func traktorParseFile(s string) (EntryCollection, bool) {
	var entries EntryCollection

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
