package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Entry struct {
	Title   string `xml:"TITLE,attr"`
	Artist  string `xml:"ARTIST,attr"`
	AudioId string `xml:"AUDIO_ID,attr"`
}

func (e Entry) String() string {
	return fmt.Sprintf("%s - %s", e.Artist, e.Title)
}

type EntryCollection struct {
	XMLName   xml.Name `xml:"NML"`
	EntryList []Entry  `xml:"COLLECTION>ENTRY"`
}

func traktorDir(s string) string {
	return os.ExpandEnv("${HOME}/Documents/Native Instruments") + "/" + s
}

func traktorParseFile(s string) (EntryCollection, bool) {
	var entries EntryCollection

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
