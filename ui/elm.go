package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const useHighPerformanceRenderer = false

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
		m.viewport.SetContent(string(msg))
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
