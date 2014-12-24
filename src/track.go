package main

import (
	"database/sql"
	"fmt"
)

type Track struct {
	Id      int
	Bpm     int
	Key     string
	Title   string
	Genre   string
	Length  int
	Artist  string
	AudioId string
	Plays   []Play
}

func trackForRows(rows *sql.Rows) Track {
	var id int
	var bpm int
	var key string
	var title string
	var artist string
	var length int
	var genre string
	var audio_id string

	if err := rows.Scan(&id, &bpm, &key, &title, &genre, &artist, &length, &audio_id); err != nil {
		fmt.Println("Unable to find this track", err)
	}
	return Track{Id: id, Artist: artist, Title: title, Bpm: bpm, Key: key, Genre: genre, Length: length, AudioId: audio_id}
}

func tracksBySql(db *sql.DB, query string) []Track {
	var tracks []Track
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Unable to execute query", err)
		return tracks
	}
	defer rows.Close()

	var track Track
	for rows.Next() {
		track = trackForRows(rows)
		track.Plays = playsBySqlAndTrack(db, &track)
		fmt.Println("Found track :", track.Id)
		fmt.Println("Found plays:", len(track.Plays))
		tracks = append(tracks, track)
	}
	return tracks
}

func getAllTracks(db *sql.DB) []Track {
	var tracks []Track

	tracks = tracksBySql(db, `SELECT tracks.* from tracks`)
	fmt.Println("Found entries:", len(tracks))
	return tracks
}
