package internal

import (
	"regexp"
	"strings"
	"time"
)

func CreateDumpFilenameFromPatternAndDate(pattern string, date time.Time) string {
	replacements := map[string]string{
		"{yyyy}": date.Format("2006"),
		"{yy}":   date.Format("06"),
		"{MM}":   date.Format("01"),
		"{M}":    date.Format("1"),
		"{dd}":   date.Format("02"),
		"{d}":    date.Format("2"),
		"{hh}":   date.Format("15"),
		"{mm}":   date.Format("04"),
		"{m}":    date.Format("4"),
		"{ss}":   date.Format("05"),
		"{s}":    date.Format("5"),
	}

	for oldString, newString := range replacements {
		pattern = strings.Replace(pattern, oldString, newString, -1)
	}

	return pattern
}

func FindAllCommandOptionsFromString(optionsString string) []string {
	return regexp.MustCompile(`-{1,2}\S+((=|\s)[^-]\S+)?`).FindAllString(optionsString, -1)
}
