package main

import (
	_ "embed"
	"math/rand"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

//go:embed prompts/prompts.txt
var promptsFile string

var promptList = func() []string {
	var list []string
	for _, p := range strings.Split(promptsFile, "\n\n") {
		p = strings.TrimSpace(p)
		if p != "" {
			list = append(list, p)
		}
	}
	return list
}()

type model struct {
	input         string
	typoSequence  string
	mistypes       map[int]bool
	totalMistypes  int
	keypressTimes  []time.Time
	totalKeypresses int
	keyErrors      map[rune]int
	wpmHistory     []float64
	currentPrompt string
	quitting      bool
	completed     bool
	startTime     time.Time
	endTime       time.Time
	// terminal dimensions
	width  int
	height int

	// key map
	keys keyMap
}

type keyMap struct {
	Quit      key.Binding
	Backspace key.Binding
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
	Backspace: key.NewBinding(key.WithKeys("backspace")),
}

func initialModel() model {
	m := model{
		keys:          keys,
		currentPrompt: promptList[rand.Intn(len(promptList))],
		mistypes:      make(map[int]bool),
		keyErrors:     make(map[rune]int),
	}
	return m
}

func (m model) Init() tea.Cmd {
	// nothing to init
	return nil
}
