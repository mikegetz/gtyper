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
	inputBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(inputBorderColor)

	promptBorderColor = lipgloss.Color("#D4A017")
	promptStyle       = lipgloss.NewStyle().Foreground(promptBorderColor)
	promptBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(promptBorderColor)

	typedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#D4A017"))
	currentWordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#89CFF0"))
	cursorStyle      = currentWordStyle.Underline(true)
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
	currentWord := m.input[strings.LastIndex(m.input, " ")+1:]
	padding := max(m.width-2-len(currentWord), 0)
	screen := inputBorderStyle.Render(currentWord+strings.Repeat(" ", padding)) + "\n"
	return addBorderTitle(screen, "Input", inputStyle)
}

func (m model) printPrompt(inputHeight int) string {
	promptHeight := max(m.height-inputHeight, 0)
	promptBorderStyle = promptBorderStyle.Width(m.width).Height(promptHeight)

	cursorPos := len(m.input)
	prompt := []rune(m.currentPrompt)

	wordEnd := len(prompt)
	for i := cursorPos; i < len(prompt); i++ {
		if prompt[i] == ' ' {
			wordEnd = i
			break
		}
	}

	content := ""
	for i, ch := range prompt {
		switch {
		case i < cursorPos:
			content += typedStyle.Render(string(ch))
		case i == cursorPos:
			if ch == ' ' {
				content += string(ch)
			} else {
				content += cursorStyle.Render(string(ch))
			}
		case i < wordEnd:
			content += currentWordStyle.Render(string(ch))
		default:
			content += string(ch)
		}
	}

	return addBorderTitle(promptBorderStyle.Render(content), "Prompt", promptStyle)
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
