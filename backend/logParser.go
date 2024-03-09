package backend

import (
	//"tli/backend/entities"
	"fmt"
	"os"
	"strings"
	"tli/backend/entities"
)

func LoadFile(filePath string) []entities.LogEntry {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}
	var logEntries []entities.LogEntry
	for _, entry := range strings.Split(string(content), "\n") {
		sliced := strings.Split(entry, " ")
		if len(sliced) >= 4 {
			logEntries = append(
				logEntries,
				entities.LogEntry{
					Timestamp: strings.Join(sliced[0:2], " "),
					Level:     sliced[2],
					Message:   strings.Join(sliced[3:], " "),
				},
			)
		}
	}

	return logEntries
}
