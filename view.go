package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var (
	inputBorderColor = lipgloss.Color("#99cc00")
	inputStyle       = lipgloss.NewStyle().Foreground(inputBorderColor)
	inputBorderStyle = inputStyle.Border(lipgloss.RoundedBorder()).BorderForeground(inputBorderColor)

	promptBorderColor = lipgloss.Color("#D4A017")
	promptStyle       = lipgloss.NewStyle().Foreground(promptBorderColor)
	promptBorderStyle = promptStyle.Border(lipgloss.RoundedBorder()).BorderForeground(promptBorderColor)
)

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	screen := ""
	input := m.printInput()
	inputHeight := lipgloss.Height(input)

	screen += input
	screen += m.printPrompt(inputHeight)

	return tea.NewView(screen)
}

func (m model) printInput() string {
	inputBorderStyle = inputBorderStyle.Width(m.width)
	screen := ""
	padding := max(m.width-2, 0) - len(m.input)
	screen += inputBorderStyle.Render(m.input + strings.Repeat(" ", padding))
	screen += "\n"
	return addBorderTitle(screen, "Input", inputStyle)
}

func (m model) printPrompt(inputHeight int) string {
	promptHeight := max(m.height-inputHeight, 0)
	promptBorderStyle = promptBorderStyle.Width(m.width).Height(promptHeight)
	screen := ""
	// TODO add prompt

	// pad empty space between prompt and term bottom
	padding := max(promptHeight-lipgloss.Height(screen)-2, 0)
	screen += promptBorderStyle.Render(strings.Repeat(strings.Repeat(" ", max(m.width-2, 0))+"\n", padding))
	return addBorderTitle(screen, "Prompt", promptStyle)
}

// Utility function to add title text to rendered style
func addBorderTitle(renderedText string, title string, renderedTextStyle lipgloss.Style) string {
	lines := strings.Split(renderedText, "\n")
	if len(lines) == 0 {
		return renderedText
	}

	topBorder := lines[0]
	// Strip ANSI codes to get the plain border characters
	plainTop := ansi.Strip(topBorder)
	topWidth := ansi.StringWidth(plainTop)
	titleWidth := ansi.StringWidth(title)

	// +2 accounts for the corner character on each end
	if titleWidth+2 > topWidth {
		return renderedText
	}

	// Slice the plain string safely
	runes := []rune(plainTop)
	prefix := string(runes[:1])
	suffix := string(runes[1+titleWidth:])

	// Re-apply the renderedTextStyle (this style should only be foreground color) to the non-title parts
	newTop := renderedTextStyle.Render(prefix) + title + renderedTextStyle.Render(suffix)
	lines[0] = newTop

	return strings.Join(lines, "\n")
}
