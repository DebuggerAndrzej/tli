package backend

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/DebuggerAndrzej/tli/backend/entities"
)

func LoadData(filePath, logFormat, pipedInput string) []entities.LogEntry {
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
		logEntries = append(logEntries, getLogEntryForLine(entry, logFormat))
	}

	return logEntries
}

func getLogEntryForLine(entry, logFormat string) entities.LogEntry {
	sliced := strings.Split(entry, " ")
	format := strings.Split(logFormat, " ")
	return entities.LogEntry{
		Timestamp: getPartOfEntry(sliced, format, "T"),
		Level:     getPartOfEntry(sliced, format, "S"),
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
