package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DebuggerAndrzej/tli/backend"
	"github.com/DebuggerAndrzej/tli/backend/entities"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	logEntries    []entities.LogEntry
	ready         bool
	viewport      viewport.Model
	textInput     textinput.Model
	weakFilters   []string
	strongFilters []string
	inputType     string
	emptyViewport bool
}

func initModel(filePath, logFormat, pipedInput string) model {
	return model{logEntries: backend.LoadData(filePath, logFormat, pipedInput), textInput: textinput.New()}
}

func (m model) headerView() string {
	return mainTitleStyle.Width(m.viewport.Width).Render("TLI")
}

func (m model) footerView() string {
	if m.textInput.Focused() {
		return inputStyle.Width(m.viewport.Width).Render(m.textInput.View())
	}
	if m.emptyViewport {
		return footerStyle.Width(m.viewport.Width).Render("Expected your logs here? Good luck next time!")
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
			if m.inputType == "Weak filter" {
				m.weakFilters = append(m.weakFilters, inputValue)
			}
			if m.inputType == "Strong filter" {
				m.strongFilters = append(m.strongFilters, inputValue)
			}
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
			m.inputType = "Weak filter"
			return m, nil
		case "F":
			m.textInput.Focus()
			m.inputType = "Strong filter"
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
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight + 1
	if !m.ready {
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.HighPerformanceRendering = true 
		m.ready = true
		m.viewport.YPosition = headerHeight + 1  
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
	}
	return m, m.updateViewportContent
}

func (m model) clearFilters() (tea.Model, tea.Cmd) {
	if !m.hasAnyFilters() {
		return m, nil
	}
	m.strongFilters = nil
	m.weakFilters = nil
	return m, tea.Cmd(m.updateViewportContent)
}

func (m model) updateViewportContent() tea.Msg {
	var sb strings.Builder
	if !m.hasAnyFilters() {
		for _, logEntry := range m.logEntries {
			sb.WriteString(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					timestampStyle.Render(logEntry.Timestamp),
					"  ",
					lipgloss.NewStyle().Foreground(GetColorForEntry(logEntry.Level)).Render(logEntry.Message),
				) + "\n",
			)
		}
	} else {
		for _, logEntry := range m.logEntries {
			if m.doesMatchFilters(logEntry.Message) {
				sb.WriteString(
					lipgloss.JoinHorizontal(
						lipgloss.Left,
						timestampStyle.Render(logEntry.Timestamp),
						"  ",
						lipgloss.NewStyle().Foreground(GetColorForEntry(logEntry.Level)).Render(logEntry.Message),
					) + "\n",
				)
			}
		}
	}

	return updatedContents(sb.String())
}

func (m model) hasAnyFilters() bool {
	return len(m.weakFilters)+len(m.strongFilters) > 0
}

func (m model) doesMatchFilters(message string) bool {
	weakFilterRegex := regexp.MustCompile(strings.Join(m.weakFilters, "|"))
	if !weakFilterRegex.MatchString(message) {
		return false
	}

	if len(m.strongFilters) > 0 && !m.doesMatchStrongFilters(message) {
		return false
	}

	return true
}

func (m model) doesMatchStrongFilters(message string) bool {
	var strongFilterRegexes []*regexp.Regexp
	for _, strongFilter := range m.strongFilters {
		strongFilterRegexes = append(strongFilterRegexes, regexp.MustCompile(strongFilter))
		if !regexp.MustCompile(strongFilter).MatchString(message) {
			return false
		}
	}
	return true
}

