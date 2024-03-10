package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	mainTitleStyle = lipgloss.NewStyle().
			Height(1).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderBottom(true)
	footerStyle = lipgloss.NewStyle().
			Height(1).
			Bold(false).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderTop(true)
	logEntryBaseStyle = lipgloss.NewStyle().
				Width(180)
	// https://github.com/charmbracelet/lipgloss/issues/85 formatting lost on line to long
	timestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A9A9A9")).Width(20)
)

func GetColorForEntry(level string) lipgloss.Color {
	if level == "WARN" {
		return lipgloss.Color("#efef8d")
	}
	if level == "ERROR" {
		return lipgloss.Color("#d19191")
	}
	return lipgloss.Color("#D3D3D3")
}
