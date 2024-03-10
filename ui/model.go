package ui

import (
	"fmt"
	"strings"
	"tli/backend"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
}

func (m model) headerView() string {
	return mainTitleStyle.Width(m.viewport.Width).Render("TLI")
}

func (m model) footerView() string {
	return footerStyle.Width(m.viewport.Width).Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
}

func initModel() model {
	logEntries := backend.LoadFile("logs/test2.log")

	var sb strings.Builder
	for _, logEntry := range logEntries {
		sb.WriteString(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				timestampStyle.Render(logEntry.Timestamp),
				logEntryBaseStyle.Copy().Foreground(GetColorForEntry(logEntry.Level)).Render(logEntry.Message),
			) + "\n",
		)
	}
	return model{content: sb.String()}
}
