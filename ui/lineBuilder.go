package ui

import (
	"strings"

	"github.com/DebuggerAndrzej/tli/backend/entities"
	"github.com/charmbracelet/lipgloss"
)

var highlightColors = []string{"#66961f", "#5934f0", "#9228da", "#c434c4", "#174348"}

func addLineToStringBuilder(
	builder *strings.Builder,
	logEntry entities.LogEntry,
	logBaseStyle lipgloss.Style,
	highlights []string,
	searched string,
	searchedOccurances *[]int,
	lineIndex int,
) {
	messageColor := lipgloss.NewStyle().Foreground(GetColorForEntry(logEntry.Level))
	styledMessage := removeStyleEnd(messageColor.Render(logEntry.Message))
	if len(highlights) > 0 && lineContainsHighlight(logEntry.Message, highlights) {
		addHighlightToMessage(&styledMessage, messageColor, highlights)
	}
	if searched != "" && strings.Contains(logEntry.Message, searched) {
		*searchedOccurances = append(*searchedOccurances, lineIndex)
		addSearchedTextHighlight(&styledMessage, messageColor, searched)
	}
	builder.WriteString(prepareLine(styledMessage, logEntry.Timestamp, logBaseStyle, timestampStyle))
}

func prepareLine(message, timestamp string, messageStyle, timestampStyle lipgloss.Style) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		timestampStyle.Render(timestamp),
		messageStyle.Render(message),
	) + "\n"
}

func addHighlightToMessage(message *string, messageColor lipgloss.Style, highlights []string) {
	for index, highlight := range highlights {
		if index >= len(highlightColors) {
			index = index % len(highlightColors)
		}
		if strings.Contains(*message, highlight) {
			*message = strings.Replace(
				*message,
				highlight,
				lipgloss.NewStyle().
					Background(lipgloss.Color(highlightColors[index])).
					Render(highlight)+
					removeStyleEnd(
						messageColor.Render(""),
					),
				-1,
			)
		}
	}
}

func addSearchedTextHighlight(message *string, messageColor lipgloss.Style, searched string) {
	*message = strings.Replace(
		*message,
		searched,
		lipgloss.NewStyle().
			Background(lipgloss.Color("#7FBBB3")).
			Render(searched)+
			removeStyleEnd(
				messageColor.Render(""),
			),
		-1,
	)
}
func lineContainsHighlight(line string, highlights []string) bool {
	for _, textToHighligh := range highlights {
		if strings.Contains(line, textToHighligh) {
			return true
		}
	}
	return false
}

func removeStyleEnd(styledText string) string {
	return strings.Replace(styledText, "\x1b[0m", "", 1)
}
