package ui

import (
	"fmt"
	"strings"
	"tli/backend"
	"tli/backend/entities"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

type model struct {
	logEntries     []entities.LogEntry
	currentContent string
	ready          bool
	viewport       viewport.Model
	textInput      textinput.Model
}

func (m model) headerView() string {
	return mainTitleStyle.Width(m.viewport.Width).Render("TLI")
}

func (m model) footerView() string {
	if m.textInput.Focused() {
		return inputStyle.Width(m.viewport.Width).Render(m.textInput.View())
	}
	return footerStyle.Width(m.viewport.Width).Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
}

func (m model) handleTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" {
		if m.textInput.Value() == "" {
			m.textInput.Blur()
			return m, nil
		}
		return m, tea.Cmd(m.updateViewportContent)
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func initModel() model {
	return model{logEntries: backend.LoadFile("logs/test2.log"), textInput: textinput.New()}
}

func (m model) getViewportContent() string {
	var sb strings.Builder
	for _, logEntry := range m.logEntries {
		sb.WriteString(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				timestampStyle.Render(logEntry.Timestamp),
				logEntryBaseStyle.Copy().Foreground(GetColorForEntry(logEntry.Level)).Render(logEntry.Message),
			) + "\n",
		)
	}
	return sb.String()
}

func (m model) updateViewportContent() tea.Msg {
	var sb strings.Builder
	for _, entry := range strings.Split(m.currentContent, "\n") {
		if strings.Contains(entry, m.textInput.Value()) {
			sb.WriteString(entry + "\n")
		}
	}
	return updatedContents(sb.String())
}
