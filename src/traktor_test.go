package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestTraktorDirFromLib(t *testing.T) {
	traktorHomeDir := os.ExpandEnv("${HOME}/Documents/Native Instruments")

	if dirName := traktorDir(""); !reflect.DeepEqual(dirName, traktorHomeDir+"/") {
		t.Errorf("traktorDir() = %+v, want %+v", dirName, traktorHomeDir)
	}
}

func TestNMLFileParsing(t *testing.T) {
	testFile := "../test/fixtures/Traktor 2.7.1/History/history_2014y04m23d_08h45m58s.nml"

	entries, _ := traktorParseFile(testFile)
	assert.Equal(t, 2, len(entries.EntryList), "Should have 2 entries")

	entry := entries.EntryList[1]
	assert.Equal(t, entry.Title, "Night Drive (Original Mix)", "Title does not match.")
	assert.Equal(t, entry.Artist, "Tony Rohr", "Artist does not match.")
	assert.Equal(t, entry.Genre(), "Techno", "Genre does not match.")
	assert.Equal(t, entry.Bpm(), 126, "BPM does not match.")
	assert.Equal(t, entry.Length(), 488, "Length does not match.")

	for _, entry := range entries.EntryList {
		if assert.NotNil(t, entry) {
			assert.NotNil(t, entry.Title, "No title provided.")
			assert.NotNil(t, entry.Artist, "No artist provided.")
			assert.NotNil(t, entry.AudioId, "No audio id provided.")
		}
	}
}

func TestTraktorHistoryPaths(t *testing.T) {
	historyPaths, _ := traktorHistoryPaths("../test/fixtures")
	assert.Equal(t, 2, len(historyPaths), "Should have 2 paths")
}

func TestTraktorArchiveFiles(t *testing.T) {
	historyPaths, _ := traktorHistoryPaths("../test/fixtures")
	archiveFiles, _ := traktorArchiveFiles(historyPaths)
	assert.Equal(t, 2, len(archiveFiles), "Should only have 2 archive files.")
}
