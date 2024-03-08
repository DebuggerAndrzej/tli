package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func InitTui() {
	content, err := os.ReadFile("test2.log")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		model{content: string(content)},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
