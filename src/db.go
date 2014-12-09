package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"regexp"
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
  month    INTEGER,
  day      INTEGER,
  hour     INTEGER,
  minute   INTEGER
);
`
}

func insertTrackStatment() string {
	return `
INSERT INTO tracks (artist,name,audio_id) values(?,?,?)
`
}
func insertPlayStatment() string {
	return `
INSERT INTO plays (track_id, month, day, hour, minute) values(?,?,?,?,?)
`
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
	_, err := db.Exec(insertPlayStatment(), id, ec.Month, ec.Day, ec.Hour, ec.Minute)
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
	os.Remove(s)

	db, err := sql.Open("sqlite3", s)
	if err != nil {
		fmt.Println("Error opening db file:", err)
		return db, false
	}

	db.Exec(createTableStatement())

	return db, true
}
