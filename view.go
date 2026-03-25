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

	untypedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	typedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#D4A017"))
	currentWordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#89CFF0"))
	cursorStyle      = currentWordStyle.Underline(true)
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC5555"))
	errorCursorStyle = errorStyle.Underline(true)
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
	currentWord := ""
	if len(m.input) < len(m.currentPrompt) {
		currentWord = m.input[strings.LastIndex(m.input, " ")+1:] + m.typoSequence
	}
	padding := max(m.width-2-len(currentWord), 0)
	screen := inputBorderStyle.Render(untypedStyle.Render(currentWord)+strings.Repeat(" ", padding)) + "\n"
	return addBorderTitle(screen, "Input", inputStyle, inputStyle)
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

	typoRunes := []rune(m.typoSequence)
	overflowStart := min(len(typoRunes), wordEnd-cursorPos)
	overflow := typoRunes[overflowStart:]

	activeWordStyle := currentWordStyle
	if m.typoSequence != "" {
		activeWordStyle = errorStyle
	}

	content := ""
	for i, ch := range prompt {
		if i == wordEnd {
			for _, och := range overflow {
				content += errorStyle.Render(string(och))
			}
		}
		switch {
		case i < cursorPos:
			if m.mistypes[i] {
				content += errorStyle.Render(string(ch))
			} else {
				content += typedStyle.Render(string(ch))
			}
		case i == cursorPos:
			if ch == ' ' {
				content += string(ch)
			} else if m.typoSequence != "" {
				content += errorCursorStyle.Render(string(ch))
			} else {
				content += cursorStyle.Render(string(ch))
			}
		case i < wordEnd:
			content += activeWordStyle.Render(string(ch))
		default:
			content += untypedStyle.Render(string(ch))
		}
	}
	if wordEnd == len(prompt) {
		for _, och := range overflow {
			content += errorStyle.Render(string(och))
		}
	}

	return addBorderTitle(promptBorderStyle.Render(content), "Prompt", promptStyle, promptStyle)
}

// Utility function to add title text to rendered style
func addBorderTitle(renderedText string, title string, borderStyle lipgloss.Style, titleStyle lipgloss.Style) string {
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

	// Re-apply the borderStyle to the non-title parts, titleStyle to the title
	newTop := borderStyle.Render(prefix) + titleStyle.Render(title) + borderStyle.Render(suffix)
	lines[0] = newTop

	return strings.Join(lines, "\n")
}
