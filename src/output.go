package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func displayOutput(db *sql.DB) {
	totalPlays := countForTable(db, "plays")
	totalTracks := countForTable(db, "tracks")
	fmt.Println("Found", totalTracks, "unique tracks.")
	fmt.Println("For a total of", totalPlays, "plays.\n")

	for year := 2012; year < 2016; year++ {
		for month := 1; month < 13; month++ {
			monthlyEntries := findChartEntriesByMonthAndYear(db, month, year)

			if len(monthlyEntries) > 0 {
				artistLength := longestArtistName(monthlyEntries)
				titleLength := longestTitleName(monthlyEntries)

				format := delimiterFormatString(artistLength, titleLength)

				fmt.Println(time.Month(month), year, "Charts")
				fmt.Println(lineDelimiter(artistLength, titleLength))
				for i, chartEntry := range monthlyEntries {
					output := fmt.Sprintf(format, i+1, chartEntry.Artist, chartEntry.Title, chartEntry.Count)
					fmt.Println(output)
				}
				fmt.Println(lineDelimiter(artistLength, titleLength) + "\n")
			}
		}
	}
}

func delimiterFormatString(artistLength int, titleLength int) string {
	return "| %2d | %-" + strconv.Itoa(artistLength) + "q | %-" + strconv.Itoa(titleLength) + "q | %02d |"
}

func longestArtistName(monthlyEntries []ChartEntry) int {
	max := 0
	for i := 0; i < len(monthlyEntries); i++ {
		length := len(monthlyEntries[i].Artist)
		if length > max {
			max = length
		}
	}
	return max + 2
}

func longestTitleName(monthlyEntries []ChartEntry) int {
	max := 0
	for i := 0; i < len(monthlyEntries); i++ {
		length := len(monthlyEntries[i].Title)
		if length > max {
			max = length
		}
	}
	return max + 6
}

func lineDelimiter(artistLength int, titleLength int) string {
	delimiter := ""
	totalDashes := artistLength + titleLength + 15
	for i := 0; i < totalDashes; i++ {
		delimiter += "-"
	}
	return "+" + delimiter + "+"
}
