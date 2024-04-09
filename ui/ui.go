package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func InitTui(filePath, logFormat, pipedInput, warningIndicator, errorIndicator string) {
	model := initModel(filePath, logFormat, pipedInput, warningIndicator, errorIndicator)
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
