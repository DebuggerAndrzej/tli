package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
}

func (m model) headerView() string {
	return mainTitleStyle.Width(m.viewport.Width).Render("TLI")
}

func (m model) footerView() string {
	return footerStyle.Width(m.viewport.Width).Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
}
