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
	logEntries              []entities.LogEntry
	ready                   bool
	viewport                viewport.Model
	textInput               textinput.Model
	weakFilters             []string
	strongFilters           []string
	highlights              []string
	inputType               string
	emptyViewport           bool
	searched                string
	searchedOccurances      []int
	visibleLogEntriesAmount int
	currentSearchIndex      int
	minimalSeverity         string
}

func initModel(filePath, logFormat, pipedInput, warningIndicator, errorIndicator string) model {
	return model{
		logEntries: backend.LoadData(filePath, logFormat, pipedInput, warningIndicator, errorIndicator),
		textInput:  textinput.New(),
	}
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
		switch keypress := msg.String(); keypress {
		case "esc":
			m.textInput.Blur()
			m.textInput.Reset()
			return m, nil
		case "enter":
			inputValue := m.textInput.Value()
			m.textInput.Blur()
			m.textInput.Reset()
			if inputValue == "" {
				return m, nil
			}
			switch m.inputType {
			case "Weak filter":
				m.weakFilters = append(m.weakFilters, inputValue)
			case "Strong filter":
				m.strongFilters = append(m.strongFilters, inputValue)
			case "Highlight":
				m.highlights = append(m.highlights, inputValue)
			case "Search":
				m.searched = inputValue
				m.currentSearchIndex = 0
			}
			return m, tea.Sequence(m.updateViewportContent, m.requiresOffsetCalculation)
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
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
		case "w":
			m.minimalSeverity = "WARN"
			return m, m.updateViewportContent
		case "e":
			m.minimalSeverity = "ERROR"
			return m, m.updateViewportContent
		case "h":
			m.textInput.Focus()
			m.inputType = "Highlight"
			return m, nil
		case "/":
			m.textInput.Focus()
			m.inputType = "Search"
			return m, nil
		case "n":
			if m.currentSearchIndex != len(m.searchedOccurances)-1 {
				m.currentSearchIndex++
			} else {
				m.currentSearchIndex = 0
			}
			return m, m.requiresOffsetCalculation
		case "N":
			if m.currentSearchIndex != 0 {
				m.currentSearchIndex--
			} else {
				m.currentSearchIndex = len(m.searchedOccurances) - 1
			}
			return m, m.requiresOffsetCalculation
		case "r":
			return m.resetModifiers()
		default:
			return m.viewPortUpdate(msg)
		}
	}
	return m, nil
}

func (m model) requiresOffsetCalculation() tea.Msg {
	return requiresOffsetCalculation(true)
}

func (m model) updateYOffset() (model, tea.Cmd) {
	if m.searchedOccurances == nil {
		return m, nil
	}
	offset := m.caculateViewportOffsetForSearchHit()
	m.viewport.SetYOffset(offset)
	return m, viewport.Sync(m.viewport) 
}

func (m model) caculateViewportOffsetForSearchHit() int {
	return int(
		float64(
			m.searchedOccurances[m.currentSearchIndex],
		) * float64(
			m.viewport.TotalLineCount()-m.viewport.Height,
		) / float64(
			m.visibleLogEntriesAmount,
		))
}

func (m model) viewPortUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	if !m.ready {
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight +1
		m.viewport.HighPerformanceRendering= true
		m.ready = true
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
	}
	return m, m.updateViewportContent
}

func (m model) resetModifiers() (tea.Model, tea.Cmd) {
	if !m.hasAnyFilters() && len(m.highlights) == 0 && m.searched == "" {
		return m, nil
	}
	m.strongFilters = nil
	m.weakFilters = nil
	m.minimalSeverity = ""
	m.highlights = nil
	m.searched = ""
	m.searchedOccurances = nil
	return m, tea.Cmd(m.updateViewportContent)
}

func (m model) updateViewportContent() tea.Msg {
	var builder strings.Builder
	var logBaseStyle lipgloss.Style
	var searchedOccurances []int
	var maxIndex int
	if m.logEntries[0].Timestamp == "" {
		logBaseStyle = lipgloss.NewStyle().Width(m.viewport.Width)
	} else {
		logBaseStyle = lipgloss.NewStyle().
			Width(m.viewport.Width - lipgloss.Width(timestampStyle.Render(m.logEntries[0].Timestamp)) - 3).
			MarginLeft(3)
	}
	if !m.hasAnyFilters() {
		for index, logEntry := range m.logEntries {
			addLineToStringBuilder(
				&builder,
				logEntry,
				logBaseStyle,
				m.highlights,
				m.searched,
				&searchedOccurances,
				index,
			)
		}
		maxIndex = len(m.logEntries)
	} else {
		var relativeIndex int
		for _, logEntry := range m.logEntries {
			if m.doesMatchFilters(logEntry.Message, logEntry.Level) {
				addLineToStringBuilder(&builder, logEntry, logBaseStyle, m.highlights, m.searched, &searchedOccurances, relativeIndex)
				relativeIndex++
			}
		}
		maxIndex = relativeIndex
	}

	entries := builder.String()
	return updatedContents{entries, searchedOccurances, maxIndex}
}

func (m model) hasAnyFilters() bool {
	return len(m.weakFilters)+len(m.strongFilters) > 0 || m.minimalSeverity != ""
}

func (m model) doesMatchFilters(message, level string) bool {
	weakFilterRegex := regexp.MustCompile(strings.Join(m.weakFilters, "|"))
	if !weakFilterRegex.MatchString(message) {
		return false
	}

	if len(m.strongFilters) > 0 && !m.doesMatchStrongFilters(message) {
		return false
	}

	if m.minimalSeverity != "" && level != "" {
		if m.minimalSeverity == "WARN" {
			if level == "WARN" || level == "ERROR" {
				return true
			}
		} else {
			if level == "ERROR" {
				return true
			}
		}
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
