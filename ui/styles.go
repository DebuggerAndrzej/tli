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
)
