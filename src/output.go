package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

func writeOutputFrom(db *sql.DB) {
	totalPlays := strconv.Itoa(countForTable(db, "plays"))
	totalTracks := strconv.Itoa(countForTable(db, "tracks"))

	markdownFile := os.ExpandEnv("${HOME}/.traktor-charts.md")
	fp, err := os.Create(markdownFile)
	if err != nil {
		fmt.Println("Unable to create", markdownFile)
	}
	defer fp.Close()

	fp.WriteString("# My Traktor DJ Charts\n")
	fp.WriteString(totalPlays + " songs played, " + totalTracks + " were unique.\n")

	fmt.Println(totalPlays, "songs played,", totalTracks, "were unique.\n")

	for year := 2016; year > 2012; year-- {
		for month := 12; month > 0; month-- {
			monthlyEntries := findChartEntriesByMonthAndYear(db, month, year)

			if len(monthlyEntries) > 0 {
				format := delimiterFormatString()

				fp.WriteString("##")
				fp.WriteString(fmt.Sprintf(" %s ", time.Month(month)))
				fp.WriteString(strconv.Itoa(year) + " Charts" + "\n\n")
				fp.WriteString(outputTableHeader())
				for i, chartEntry := range monthlyEntries {
					title := chartEntry.Title
					artist := chartEntry.Artist

					output := fmt.Sprintf(format, i+1, artist, title, beatPortLink(chartEntry))
					fp.WriteString(output)
				}
				fp.WriteString("\n")
			}
		}
	}
	fp.Sync()
}

func beatPortLink(ce ChartEntry) string {
	link := "http://www.beatport.com/search?query="
	return link + url.QueryEscape(ce.Artist+" "+ce.Title)
}

func delimiterFormatString() string {
	return "| %d | %s | %s | [Beatport](%s) |\n"
}

func outputTableHeader() string {
	return "| Number | Artist | Title | Cop It |\n|---|---|---|---|\n"
}
