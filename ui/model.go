package ui

import (
	"fmt"
	"regexp"
	"strings"
	"tli/backend"
	"tli/backend/entities"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

type model struct {
	logEntries []entities.LogEntry
	ready      bool
	viewport   viewport.Model
	textInput  textinput.Model
	filters    []string
}

func initModel(filePath string) model {
	return model{logEntries: backend.LoadFile(filePath), textInput: textinput.New()}
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

func (m model) handleKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.textInput.Focused() {
	case true:
		if msg.String() == "enter" {
			inputValue := m.textInput.Value()
			m.textInput.Blur()
			m.textInput.Reset()
			if inputValue == "" {
				return m, nil
			}
			m.filters = append(m.filters, inputValue)
			return m, tea.Cmd(m.updateViewportContent)
		}
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	case false:
		switch keypress := msg.String(); keypress {
		case "q":
			return m, tea.Quit
		case "f":
			m.textInput.Focus()
			return m, nil
		case "r":
			return m.clearFilters()
		default:
			return m.viewPortUpdate(msg)
		}
	}
	return m, nil
}

func (m model) viewPortUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	if !m.ready {
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
		m.ready = true
		m.viewport.YPosition = headerHeight + 1
		cmds = append(cmds, m.updateViewportContent)
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
	}
	if useHighPerformanceRenderer {
		cmds = append(cmds, viewport.Sync(m.viewport))
	}
	return m, tea.Batch(cmds...)
}

func (m model) clearFilters() (tea.Model, tea.Cmd) {
	if len(m.filters) == 0 {
		return m, nil
	}
	m.filters = nil
	return m, tea.Cmd(m.updateViewportContent)
}

func (m model) updateViewportContent() tea.Msg {
	var sb strings.Builder
	if len(m.filters) == 0 {
		for _, logEntry := range m.logEntries {
			sb.WriteString(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					timestampStyle.Render(logEntry.Timestamp),
					logEntryBaseStyle.Copy().Foreground(GetColorForEntry(logEntry.Level)).Render(logEntry.Message),
				) + "\n",
			)
		}
	} else {
		var filterRegex = regexp.MustCompile(strings.Join(m.filters, "|"))
		for _, logEntry := range m.logEntries {
			if filterRegex.MatchString(logEntry.Message) {
				sb.WriteString(
					lipgloss.JoinHorizontal(
						lipgloss.Left,
						timestampStyle.Render(logEntry.Timestamp),
						logEntryBaseStyle.Copy().Foreground(GetColorForEntry(logEntry.Level)).Render(logEntry.Message),
					) + "\n",
				)
			}
		}
	}

	return updatedContents(sb.String())
}
