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
			fp.WriteString(fmt.Sprintf("#### Best of %d\n\n", year))
			fp.WriteString(outputTableHeader())
			for i, chartEntry := range yearlyEntries {
				output := chartEntry.toMarkdown(i+1, listenLink(chartEntry))
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
				fp.WriteString(fmt.Sprintf("#### %s ", time.Month(month)))
				fp.WriteString(strconv.Itoa(year) + " Charts" + "\n\n")
				fp.WriteString(outputTableHeader())
				for i, chartEntry := range monthlyEntries {
					output := chartEntry.toMarkdown(i+1, listenLink(chartEntry))
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
	stats := totalPlays + " songs played, " + totalTracks + " were unique."

	markdownFile := os.ExpandEnv("${HOME}/.traktor-charts.md")
	fp, err := os.Create(markdownFile)
	if err != nil {
		fmt.Println("Unable to create", markdownFile)
	}
	defer fp.Close()

	fp.WriteString("# My Traktor DJ Charts\n\n")
	fp.WriteString(stats + "\n\n")

	fmt.Println(stats)

	writeChartByYear(fp, db)
	writeChartByMonthAndYear(fp, db)

	fp.Sync()
}

func listenLink(ce ChartEntry) string {
	link := "https://www.youtube.com/results?search_query="
	return link + url.QueryEscape(ce.Artist+" "+ce.Title)
}

func (ce ChartEntry) toMarkdown(num int, link string) string {
	format := "| %d | %s | %s | %s | %d | %s | [YouTube](%s) |\n"

	return fmt.Sprintf(format, num, ce.Artist, ce.Title, ce.Genre, ce.Bpm, ce.Key, link)
}

func outputTableHeader() string {
	return "| Number | Artist | Title | Genre | BPM | Key | Check It |\n|---|---|---|---|---|---|---|\n"
}
