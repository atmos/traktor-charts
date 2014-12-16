package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

func chartByYear(ty TraktorDataByYear) string {
	result := fmt.Sprintf("#### Best of %d\n\n", ty.Year)
	result += outputTableHeader()
	for i, chartEntry := range ty.Charts {
		result += chartEntry.toMarkdown(i+1, listenLink(chartEntry))
	}
	result += "\n"
	return result
}

func chartByMonthAndYear(tbmay TraktorByYearAndMonth) string {
	result := "## Years by month\n\n"

	result += fmt.Sprintf("#### %s ", time.Month(tbmay.Month))
	result += strconv.Itoa(tbmay.Year) + " Charts" + "\n\n"
	result += outputTableHeader()

	for i, chartEntry := range tbmay.Charts {
		result += chartEntry.toMarkdown(i+1, listenLink(chartEntry))
	}
	result += "\n"
	return result
}

func writeMarkdownFile(traktorData TraktorData) {
	stats := fmt.Sprintf("%d songs playd, %d were unique.", traktorData.Plays, traktorData.Tracks)

	markdownFile := os.ExpandEnv("${HOME}/.traktor-charts.md")
	fp, err := os.Create(markdownFile)
	if err != nil {
		fmt.Println("Unable to create", markdownFile)
	}
	defer fp.Close()

	fp.WriteString("# My Traktor DJ Charts\n\n")
	fp.WriteString(stats + "\n\n")

	fmt.Println(stats)

	fp.WriteString("## Top tracks by year\n\n")
	for _, outputYear := range traktorData.ByYear {
		fp.WriteString(chartByYear(outputYear))
	}

	for _, outputYear := range traktorData.ByYear {
		for _, outputByMonth := range outputYear.ByMonth {
			fp.WriteString(chartByMonthAndYear(outputByMonth))
		}
	}
	fp.Sync()
}

func listenLink(ce ChartEntry) string {
	link := "https://www.youtube.com/results?search_query="
	return link + url.QueryEscape(ce.Artist+" "+ce.Title)
}

func (ce ChartEntry) toMarkdown(num int, link string) string {
	format := "| %d | %s | %s | %s | %d | %s | %s | [YouTube](%s) |\n"

	return fmt.Sprintf(format, num, ce.Artist, ce.Title, ce.Genre, ce.Bpm, ce.Key, ce.StringLength(), link)
}

func outputTableHeader() string {
	return "| Number | Artist | Title | Genre | BPM | Key | Length | Check It |\n|---|---|---|---|---|---|---|---|\n"
}
