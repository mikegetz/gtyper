package main

import (
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
		for i, binding := range m.keys.Letters {
			if !key.Matches(msg, binding) {
				continue
			}
			letter := rune('a' + i)
			// TODO
			m.input += string(letter)
			break
		}
	}
	return m, nil
}
