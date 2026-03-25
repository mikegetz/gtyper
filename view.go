package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var (
	primaryColor     = lipgloss.Color("#99cc00")
	inputStyle       = lipgloss.NewStyle().Foreground(primaryColor)
	inputBorderStyle = inputStyle.Border(lipgloss.RoundedBorder()).BorderForeground(primaryColor)
)

func (m model) View() tea.View {
	screen := ""

	screen += m.printInput()
	return tea.NewView(screen)
}

func (m model) printInput() string {
	inputBorderStyle.Width(m.width)
	screen := ""
	padding := max(m.width-2, 0) - len(m.input)
	screen += inputBorderStyle.Render(m.input + strings.Repeat(" ", padding))
	return addBorderTitle(screen, "Input", inputStyle)
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

	// +4 accounts for the corner character on each end (2) and a space on each side of the title (2)
	if titleWidth+4 > topWidth {
		return renderedText
	}

	// Slice the plain string safely
	runes := []rune(plainTop)
	prefix := string(runes[:1])
	suffix := string(runes[1+1+titleWidth+1:])

	// Re-apply the renderedTextStyle (this style should only be foreground color) to the non-title parts
	newTop := renderedTextStyle.Render(prefix) + " " + title + " " + renderedTextStyle.Render(suffix)
	lines[0] = newTop

	return strings.Join(lines, "\n")
}
