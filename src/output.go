package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

func writeChartByYear(fp *os.File, db *sql.DB) {
	fp.WriteString("## Top tracks by year\n\n")
	for year := 2016; year > 2012; year-- {
		yearlyEntries := findChartEntriesByYear(db, year)

		if len(yearlyEntries) > 0 {
			format := delimiterFormatString()

			fp.WriteString(fmt.Sprintf("### Best of %d\n\n", year))
			fp.WriteString(outputTableHeader())
			for i, chartEntry := range yearlyEntries {
				title := chartEntry.Title
				artist := chartEntry.Artist

				output := fmt.Sprintf(format, i+1, artist, title, listenLink(chartEntry))
				fp.WriteString(output)
			}
			fp.WriteString("\n")
		}
	}
}

func writeChartByMonthAndYear(fp *os.File, db *sql.DB) {
	fp.WriteString("## Years by month\n\n")
	for year := 2016; year > 2012; year-- {
		for month := 12; month > 0; month-- {
			monthlyEntries := findChartEntriesByMonthAndYear(db, month, year)

			if len(monthlyEntries) > 0 {
				format := delimiterFormatString()

				fp.WriteString(fmt.Sprintf("### %s ", time.Month(month)))
				fp.WriteString(strconv.Itoa(year) + " Charts" + "\n\n")
				fp.WriteString(outputTableHeader())
				for i, chartEntry := range monthlyEntries {
					title := chartEntry.Title
					artist := chartEntry.Artist

					output := fmt.Sprintf(format, i+1, artist, title, listenLink(chartEntry))
					fp.WriteString(output)
				}
				fp.WriteString("\n")
			}
		}
	}
}

func writeOutputFrom(db *sql.DB) {
	totalPlays := strconv.Itoa(countForTable(db, "plays"))
	totalTracks := strconv.Itoa(countForTable(db, "tracks"))

	markdownFile := os.ExpandEnv("${HOME}/.traktor-charts.md")
	fp, err := os.Create(markdownFile)
	if err != nil {
		fmt.Println("Unable to create", markdownFile)
	}
	defer fp.Close()

	fp.WriteString("# My Traktor DJ Charts\n\n")
	fp.WriteString(totalPlays + " songs played, " + totalTracks + " were unique.\n\n")

	fmt.Println(totalPlays, "songs played,", totalTracks, "were unique.")

	writeChartByYear(fp, db)
	writeChartByMonthAndYear(fp, db)

	fp.Sync()
}

func listenLink(ce ChartEntry) string {
	link := "https://www.youtube.com/results?search_query="
	return link + url.QueryEscape(ce.Artist+" "+ce.Title)
}

func delimiterFormatString() string {
	return "| %d | %s | %s | [YouTube](%s) |\n"
}

func outputTableHeader() string {
	return "| Number | Artist | Title | Cop It |\n|---|---|---|---|\n"
}
