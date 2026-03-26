package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/NimbleMarkets/ntcharts/v2/canvas"
	"github.com/NimbleMarkets/ntcharts/v2/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/v2/linechart/wavelinechart"

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
	if m.completed {
		return tea.NewView(m.printReport())
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
	padding := max(m.width-2-len([]rune(currentWord)), 0)
	screen := inputBorderStyle.Render(untypedStyle.Render(currentWord)+strings.Repeat(" ", padding)) + "\n"
	return addBorderTitle(screen, "Input", inputStyle, inputStyle)
}

func (m model) printPrompt(inputHeight int) string {
	promptHeight := max(m.height-inputHeight-1, 0)
	promptBorderStyle = promptBorderStyle.Width(m.width).Height(promptHeight)

	cursorPos := len([]rune(m.input))
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

	content += "\n\n" + untypedStyle.Italic(true).Render(m.currentSource)
	return addBorderTitle(promptBorderStyle.Render(content), "Prompt", promptStyle, promptStyle)
}

func (m model) printReport() string {
	elapsed := m.endTime.Sub(m.startTime)
	elapsedMin := elapsed.Minutes()
	if elapsedMin < 0.0001 {
		elapsedMin = 0.0001
	}
	correctChars := float64(len([]rune(m.currentPrompt)))
	totalKeys := float64(m.totalKeypresses)
	if totalKeys < 1 {
		totalKeys = 1
	}
	adjWPM := (correctChars / 5.0) / elapsedMin
	rawWPM := (totalKeys / 5.0) / elapsedMin
	accuracy := correctChars / totalKeys * 100.0
	if accuracy > 100 {
		accuracy = 100
	}

	label := func(s string) string { return untypedStyle.Render(fmt.Sprintf("%-16s", s)) }
	value := func(s string) string { return currentWordStyle.Render(s) }

	// Overview panel (left half)
	halfW := m.width / 2
	overviewContent := "\n" +
		label("Adjusted WPM:") + value(fmt.Sprintf("%.1f", adjWPM)) + "\n" +
		label("Accuracy:") + value(fmt.Sprintf("%.1f%%", accuracy)) + "\n" +
		label("Raw WPM:") + value(fmt.Sprintf("%.1f", rawWPM)) + "\n" +
		label("Correct Keys:") + value(fmt.Sprintf("%d/%d", int(correctChars), m.totalKeypresses)) + "\n"
	overviewStyle := promptBorderStyle.Width(halfW - 2).UnsetHeight()
	overviewPanel := addBorderTitle(overviewStyle.Render(overviewContent), "Overview", promptStyle, promptStyle)
	overviewH := lipgloss.Height(overviewPanel)

	// Worst Keys panel (right half)
	type runeStats struct {
		ch          rune
		appearances int
		errors      int
	}
	statsMap := map[rune]*runeStats{}
	for _, ch := range m.currentPrompt {
		if _, ok := statsMap[ch]; !ok {
			statsMap[ch] = &runeStats{ch: ch}
		}
		statsMap[ch].appearances++
	}
	for ch, errs := range m.keyErrors {
		if s, ok := statsMap[ch]; ok {
			s.errors += errs
		}
	}
	var runeList []*runeStats
	for _, s := range statsMap {
		if s.errors > 0 {
			runeList = append(runeList, s)
		}
	}
	sort.Slice(runeList, func(i, j int) bool {
		ai := max(float64(runeList[i].appearances-runeList[i].errors)/float64(runeList[i].appearances), 0.0)
		aj := max(float64(runeList[j].appearances-runeList[j].errors)/float64(runeList[j].appearances), 0.0)
		return ai < aj
	})
	worstContent := "\n"
	shown := 0
	for _, s := range runeList {
		if shown >= 5 {
			break
		}
		acc := max(float64(s.appearances-s.errors)/float64(s.appearances)*100.0, 0.0)
		ch := string(s.ch)
		if s.ch == ' ' {
			ch = "space"
		}
		worstContent += untypedStyle.Render("- ") +
			currentWordStyle.Render(ch) +
			untypedStyle.Render(fmt.Sprintf(" at %.1f%% accuracy", acc)) + "\n"
		shown++
	}
	if shown == 0 {
		worstContent += typedStyle.Render("  no errors!") + "\n"
	}
	rightW := m.width - halfW
	worstStyle := promptBorderStyle.Width(rightW - 2).UnsetHeight()
	// Pad worstContent with newlines to match the overview panel height
	worstDraft := addBorderTitle(worstStyle.Render(worstContent), "Worst Keys", promptStyle, promptStyle)
	heightDiff := overviewH - lipgloss.Height(worstDraft)
	if heightDiff > 0 {
		worstContent += strings.Repeat("\n", heightDiff)
	}
	worstPanel := addBorderTitle(worstStyle.Render(worstContent), "Worst Keys", promptStyle, promptStyle)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, overviewPanel, worstPanel)

	// Chart — skip the first (window-1) points; they use incomplete rolling windows
	// and produce unreliable WPM spikes that blow out the Y scale.
	const rollingWindow = 10
	stableHistory := m.wpmHistory
	if len(stableHistory) > rollingWindow-1 {
		stableHistory = stableHistory[rollingWindow-1:]
	}

	chartSection := ""
	chartH := m.height - lipgloss.Height(topRow) - 5 // -1 "Chart", -1 yLabel, -1 xLabel, -2 hint+newline
	if chartH > 3 && len(stableHistory) > 1 {
		chartW := m.width - 2

		minY, maxY := stableHistory[0], stableHistory[0]
		for _, v := range stableHistory[1:] {
			if v < minY {
				minY = v
			}
			if v > maxY {
				maxY = v
			}
		}
		padding := (maxY - minY) * 0.15
		if padding < 5 {
			padding = 5
		}

		wlc := wavelinechart.New(chartW, chartH)
		wlc.SetStyles(runes.LineUpHeavyDown, lipgloss.NewStyle().Foreground(inputBorderColor))
		wlc.AxisStyle = untypedStyle
		wlc.LabelStyle = untypedStyle
		wlc.SetViewYRange(minY-padding, maxY+padding)
		for i, v := range stableHistory {
			wlc.Plot(canvas.Float64Point{X: float64(i), Y: v})
		}
		wlc.Draw()
		yAxisOffset := len(fmt.Sprintf("%.0f", maxY+padding))
		yLabel := strings.Repeat(" ", yAxisOffset) + untypedStyle.Render("│ WPM (10-keypress rolling average)")
		xLabel := lipgloss.NewStyle().Width(chartW).Align(lipgloss.Right).
			Render(untypedStyle.Render("Keypresses"))
		chartSection = "\n" + promptStyle.Render("Chart") + "\n" +
			yLabel + "\n" + wlc.View() + "\n" + xLabel
	}

	hint := "\n" + typedStyle.Italic(true).Render("Press 'esc' to quit")
	return topRow + chartSection + hint
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
