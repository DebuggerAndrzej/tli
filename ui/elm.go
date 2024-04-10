package ui

import (

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
type syncMsg bool
type requiresOffsetCalculation bool

type updatedContents struct {
	Content         string
	searchedIndexes []int
	maxIndex        int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyInput(msg)
	case tea.MouseMsg:
		return m.viewPortUpdate(msg)
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	case updatedContents:
		if msg.Content == "" {
			m.emptyViewport = true
		} else {
			m.emptyViewport = false
		}
		m.viewport.SetContent(msg.Content)
		m.searchedOccurances = msg.searchedIndexes
		m.visibleLogEntriesAmount = msg.maxIndex
		return m, requiresSync
	case syncMsg:
		return m, viewport.Sync(m.viewport)
	case requiresOffsetCalculation:
		return m.updateYOffset()
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	var mainContent string
	if m.emptyViewport {
		height := m.viewport.Height
		width := m.viewport.Width
		emptyViewStyle := lipgloss.NewStyle().
			Width(width).
			Height(height).
			PaddingTop(height/2 - 1).
			PaddingLeft(width/2 - 8)
		mainContent = emptyViewStyle.Render("No results left...")
	} else {
		mainContent = m.viewport.View()
		//mainContent = "TEST"
	}
	return lipgloss.JoinVertical(lipgloss.Top, m.headerView(), mainContent, m.footerView())
}

func requiresSync() tea.Msg {
	return syncMsg(true)
}
