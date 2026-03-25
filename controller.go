package main

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		if key.Matches(msg, m.keys.Quit) {
			m.quitting = true
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.Backspace) {
			if wr := []rune(m.typoSequence); len(wr) > 0 {
				m.typoSequence = string(wr[:len(wr)-1])
				if m.typoSequence == "" {
					delete(m.mistypes, len([]rune(m.input)))
				}
			} else if runes := []rune(m.input); len(runes) > 0 {
				delete(m.mistypes, len(runes)-1)
				m.input = string(runes[:len(runes)-1])
			}
			break
		}
		prompt := []rune(m.currentPrompt)
		pos := len([]rune(m.input))
		if pos < len(prompt) {
			if text := msg.Key().Text; text != "" {
				if text == string(prompt[pos]) && m.typoSequence == "" {
					m.input += text
				} else {
					currentWord := m.input[strings.LastIndex(m.input, " ")+1:]
					if len([]rune(m.typoSequence)) < m.width-2-len([]rune(currentWord)) {
						m.typoSequence += text
						m.mistypes[pos] = true
					}
				}
			}
		}
	}
	return m, nil
}
