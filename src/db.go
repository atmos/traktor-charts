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
  bpm      INTEGER,
  key      STRING,
  name     STRING,
  genre    STRING,
  artist   STRING,
  length   INTEGER,
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
	Bpm    int
	Count  int
	Key    string
	Title  string
	Genre  string
	Length int
	Artist string
}

func (ce ChartEntry) StringLength() string {
	return fmt.Sprintf("%dm:%02ds", ce.Length/60, ce.Length%60)
}

func insertTrackStatment() string {
	return `
INSERT INTO tracks (artist,name,genre,bpm,key,length,audio_id) values(?,?,?,?,?,?,?)
`
}
func updateTrackStatment() string {
	return `
UPDATE tracks
SET genre = ?, bpm = ?, key = ?, length = ?
WHERE id = ?
`
}

func insertPlayStatment() string {
	return `
INSERT INTO plays (track_id, year, month, day, hour, minute) values(?,?,?,?,?,?)
`
}

func playsByMonthAndYearStatement(month int, year int) string {
	return `
SELECT tracks.artist, tracks.name, tracks.genre, tracks.bpm, tracks.key, tracks.length, count(plays.track_id) AS total
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

func playsByYearStatement(year int) string {
	return `
SELECT tracks.artist, tracks.name, tracks.genre, tracks.bpm, tracks.key, tracks.length, count(plays.track_id) AS total
FROM plays,tracks
WHERE year = ` + strconv.Itoa(year) +
		` AND plays.track_id = tracks.id
GROUP BY plays.track_id
ORDER by total DESC, tracks.artist ASC
LIMIT 15;
`
}

func countForTable(db *sql.DB, tableName string) int {
	rows, err := db.Query("SELECT COUNT(*) FROM " + tableName)
	if err != nil {
		fmt.Println("Unable to count:", tableName, err)
	}
	defer rows.Close()

	if rows.Next() {
		var total int
		if err := rows.Scan(&total); err != nil {
			fmt.Println("Unable to find:", tableName, err)
		}
		return total
	} else {
		return -1
	}
}

func chartEntryFindBySql(db *sql.DB, query string) []ChartEntry {
	var entries []ChartEntry
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Unable to execute query", err)
		return entries
	}
	defer rows.Close()

	for rows.Next() {
		entries = append(entries, chartEntryForRows(rows))
	}
	return entries
}

func chartEntryForRows(rows *sql.Rows) ChartEntry {
	var bpm int
	var total int
	var key string
	var title string
	var artist string
	var length int
	var genre string

	if err := rows.Scan(&artist, &title, &genre, &bpm, &key, &length, &total); err != nil {
		fmt.Println("Unable to find this entry", err)
	}
	return ChartEntry{Artist: artist, Title: title, Bpm: bpm, Key: key, Genre: genre, Length: length, Count: total}
}

func findChartEntriesByYear(db *sql.DB, year int) []ChartEntry {
	return chartEntryFindBySql(db, playsByYearStatement(year))
}

func findChartEntriesByMonthAndYear(db *sql.DB, month int, year int) []ChartEntry {
	return chartEntryFindBySql(db, playsByMonthAndYearStatement(month, year))
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
	_, err := db.Exec(insertTrackStatment(), e.Artist, e.Title, e.Genre(), e.Bpm(), e.Key(), e.Length(), e.AudioId)
	if err != nil {
		matched, _ := regexp.MatchString("UNIQUE constraint", err.Error())
		if !matched {
			fmt.Println("Error:\n", err)
		}
	}
	trackId := findTrackByAudioId(db, e.AudioId)
	updateTrack(db, e, trackId)
	insertPlay(db, ec, e, trackId)
}

func updateTrack(db *sql.DB, e Entry, id int) {
	res, err := db.Exec(updateTrackStatment(), e.Genre(), e.Bpm(), e.Key(), e.Length(), id)
	if err != nil {
		fmt.Println("Error:\n", err)
		matched, _ := regexp.MatchString("UNIQUE constraint", err.Error())
		if !matched {
			fmt.Println("Error:\n", err)
		}
	}

	affected, _ := res.RowsAffected()
	fmt.Println("Hi:", affected)
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
