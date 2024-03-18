package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type updatedContents string

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
		if string(msg) == "" {
			m.emptyViewport = true
		} else {
			m.emptyViewport = false
		}
		m.viewport.SetContent(string(msg))
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
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), mainContent, m.footerView())
}
