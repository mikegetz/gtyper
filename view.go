package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#99cc00"))

func (m model) View() tea.View {
	screen := ""

	screen += m.printInput()
	return tea.NewView(screen)
}

func (m model) printInput() string {
	screen := ""
	screen += inputStyle.Render(m.input)
	return screen
}
