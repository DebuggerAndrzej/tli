package backend

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/DebuggerAndrzej/tli/backend/entities"
)

func LoadData(filePath, logFormat, pipedInput, warningIndicator, errorIndicator string) []entities.LogEntry {
	var content string
	if filePath != "" {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Could not load file:", err)
			os.Exit(1)
		}
		content = string(fileContent)
	} else {
		content = pipedInput
	}
	var logEntries []entities.LogEntry
	for _, entry := range strings.Split(content, "\n") {
		logEntry := getLogEntryForLine(entry, logFormat, warningIndicator, errorIndicator)
		if logEntry.Message != "" {
			logEntries = append(logEntries, logEntry)
		}
	}

	return logEntries
}

func getLogEntryForLine(entry, logFormat, warningIndicator, errorIndicator string) entities.LogEntry {
	sliced := strings.Split(entry, " ")
	format := strings.Split(logFormat, " ")
	level := getLogLevel(sliced, format, warningIndicator, errorIndicator)
	return entities.LogEntry{
		Timestamp: getPartOfEntry(sliced, format, "T"),
		Level:     level,
		Message:   getPartOfEntry(sliced, format, "M"),
	}
}

func getPartOfEntry(entry, logFormat []string, entryPart string) string {
	indexes := getIndexesOfPart(logFormat, entryPart)
	isTheLastOne := slices.Contains(indexes, len(logFormat)-1)
	var part string
	for iteration, idx := range indexes {
		if isLineShortedThenExpected(idx, len(entry)) {
			continue
		}
		if iteration == 0 {
			part = part + entry[idx]
			continue
		}
		part = part + " " + entry[idx]
	}
	if isTheLastOne && len(entry) > len(logFormat) {
		part = part + strings.Join(entry[len(logFormat):], " ")
	}
	return part
}

func getIndexesOfPart(logFormat []string, entryPart string) []int {
	var indexes []int
	for index, elem := range logFormat {
		if elem == entryPart {
			indexes = append(indexes, index)
		}
	}
	return indexes
}

func isLineShortedThenExpected(expected, actual int) bool {
	return expected >= actual
}

func getLogLevel(slicedEntry, format []string, warningIndicator, errorIndicator string) string {
	level := getPartOfEntry(slicedEntry, format, "S")

	if level != "" {
		if strings.Contains(level, warningIndicator) {
			return "WARN"
		}
		if strings.Contains(level, errorIndicator) {
			return "ERROR"
		}
	}
	return level
}
