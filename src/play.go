package main

import (
	"database/sql"
	"fmt"
)

type Play struct {
	Id      int
	TrackId int
	Year    int
	Month   int
	Day     int
	Hour    int
	Minute  int
}

func playForRows(rows *sql.Rows) Play {
	var id int
	var track_id int
	var year int
	var month int
	var day int
	var hour int
	var minute int

	if err := rows.Scan(&id, &track_id, &year, &month, &day, &hour, &minute); err != nil {
		fmt.Println("Unable to find this play", err)
	}
	return Play{Id: id, TrackId: track_id, Year: year, Month: month, Day: day, Hour: hour, Minute: minute}
}

func playsByTrackId(db *sql.DB, track_id int) []Play {
	var plays []Play
	statement := `SELECT plays.* from plays where track_id = ?`

	rows, err := db.Query(statement, track_id)
	if err != nil {
		fmt.Println("Unable to execute query", track_id, err)
		return plays
	}
	defer rows.Close()

	for rows.Next() {
		plays = append(plays, playForRows(rows))
	}
	return plays
}
