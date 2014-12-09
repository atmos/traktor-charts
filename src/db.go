package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"regexp"
	"strconv"
)

func createTableStatement() string {
	return `
CREATE TABLE 'tracks'(
  id INTEGER PRIMARY KEY,
  artist   STRING,
  name     STRING,
  audio_id STRING UNIQUE
);
CREATE TABLE 'plays'(
  id INTEGER PRIMARY KEY,
  track_id INTEGER,
  year     INTEGER,
  month    INTEGER,
  day      INTEGER,
  hour     INTEGER,
  minute   INTEGER
);
`
}

type ChartEntry struct {
	Title  string
	Count  int
	Artist string
}

func insertTrackStatment() string {
	return `
INSERT INTO tracks (artist,name,audio_id) values(?,?,?)
`
}
func insertPlayStatment() string {
	return `
INSERT INTO plays (track_id, year, month, day, hour, minute) values(?,?,?,?,?,?)
`
}

func playsByMonthAndYearStatement(month int, year int) string {
	return `
SELECT tracks.artist, tracks.name, count(plays.track_id) AS total
FROM plays,tracks
WHERE
  month = ` + strconv.Itoa(month) +
		` AND year = ` + strconv.Itoa(year) +
		` AND plays.track_id = tracks.id
GROUP BY plays.track_id
ORDER by total DESC, tracks.artist ASC
LIMIT 10;
`
}

func countForTable(db *sql.DB, tableName string) int {
	rows, err := db.Query("SELECT COUNT(*) FROM " + tableName)
	if err != nil {
		fmt.Println("Unable to count:", tableName, err, "\n")
	}
	defer rows.Close()

	if rows.Next() {
		var total int
		if err := rows.Scan(&total); err != nil {
			fmt.Println("Unable to find:", tableName, err, "\n")
		}
		return total
	} else {
		return -1
	}
}

func findChartEntriesByMonthAndYear(db *sql.DB, month int, year int) []ChartEntry {
	var entries []ChartEntry
	rows, err := db.Query(playsByMonthAndYearStatement(month, year))
	if err != nil {
		fmt.Println("Unable to query plays by month", err, "\n")
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		var total int
		var artist string

		if err := rows.Scan(&artist, &title, &total); err != nil {
			fmt.Println("Unable to find this entry", err)
		}
		entries = append(entries, ChartEntry{Artist: artist, Title: title, Count: total})
	}
	return entries
}

func findTrackByAudioId(db *sql.DB, id string) int {
	statement := `SELECT id from tracks where audio_id = ?`

	rows, err := db.Query(statement, id)
	if err != nil {
		fmt.Println("Unable to find:\n", id, err)
	}
	defer rows.Close()

	if rows.Next() {
		var trackId int
		if err := rows.Scan(&trackId); err != nil {
			fmt.Println("Unable to find:\n", id, err)
		}
		return trackId
	} else {
		return -1
	}
}

func insertPlay(db *sql.DB, ec EntryCollection, e Entry, id int) {
	_, err := db.Exec(insertPlayStatment(), id, ec.Year, ec.Month, ec.Day, ec.Hour, ec.Minute)
	if err != nil {
		fmt.Println("Error:\n", err)
	}
}

func insertEntry(db *sql.DB, ec EntryCollection, e Entry) {
	_, err := db.Exec(insertTrackStatment(), e.Artist, e.Title, e.AudioId)
	if err != nil {
		matched, _ := regexp.MatchString("UNIQUE constraint", err.Error())
		if !matched {
			fmt.Println("Error:\n", err)
		}
	}
	trackId := findTrackByAudioId(db, e.AudioId)
	insertPlay(db, ec, e, trackId)
}

func initializeDB(s string) (*sql.DB, bool) {
	fullPath := os.ExpandEnv("${HOME}/." + s)
	os.Remove(fullPath)

	db, err := sql.Open("sqlite3", fullPath)
	if err != nil {
		fmt.Println("Error opening db file:", err)
		return db, false
	}

	db.Exec(createTableStatement())

	return db, true
}
