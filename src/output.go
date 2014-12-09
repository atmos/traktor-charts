package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"
)

func displayOutput(db *sql.DB) {
	totalPlays := countForTable(db, "plays")
	totalTracks := countForTable(db, "tracks")

	fmt.Println("# My Traktor DJ Charts\n")
	fmt.Println(totalPlays, " songs played,", totalTracks, "were unique.\n")

	for year := 2016; year > 2012; year-- {
		for month := 12; month > 0; month-- {
			monthlyEntries := findChartEntriesByMonthAndYear(db, month, year)

			if len(monthlyEntries) > 0 {
				format := delimiterFormatString()

				fmt.Println("##", time.Month(month), year, "Charts", "\n")
				fmt.Println(outputTableHeader())
				for i, chartEntry := range monthlyEntries {
					title := chartEntry.Title
					artist := chartEntry.Artist

					output := fmt.Sprintf(format, i+1, artist, title, beatPortLink(chartEntry))
					fmt.Println(output)
				}
				fmt.Println("\n")
			}
		}
	}
}

func beatPortLink(ce ChartEntry) string {
	link := "http://www.beatport.com/search?query="
	return link + url.QueryEscape(ce.Artist+" "+ce.Title)
}

func delimiterFormatString() string {
	return "| %d | %s | %s | [Beatport](%s) |"
}

func outputTableHeader() string {
	return "| Number | Artist | Title | Cop It |\n|---|---|---|---|"
}
