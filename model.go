package main

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	input string
	// terminal dimensions
	width  int
	height int

	// key map
	keys keyMap
}

type keyMap struct {
	Quit    key.Binding
	Letters [26]key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{}}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "quit"),
	),
	Letters: func() [26]key.Binding {
		var bindings [26]key.Binding
		for i := range 26 {
			letter := string(rune('a' + i))
			bindings[i] = key.NewBinding(key.WithKeys(letter))
		}
		return bindings
	}(),
}

func initialModel() model {
	m := model{
		keys: keys,
	}
	return m
}

func (m model) Init() tea.Cmd {
	// nothing to init
	return nil
}
