package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func InitTui(filePath string) {
	model := initModel(filePath)
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
