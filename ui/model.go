package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"strings"
	"tli/backend"
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
		sb.WriteString(fmt.Sprintf("%s --- %s\n", logEntry.Timestamp, logEntry.Message))
	}
	return model{content: sb.String()}
}
