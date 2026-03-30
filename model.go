package main

import (
	_ "embed"
	"math/rand"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

//go:embed prompts/prompts.txt
var promptsFile string

type prompt struct {
	text   string
	source string
}

var promptList = func() []prompt {
	var list []prompt
	for _, block := range strings.Split(promptsFile, "\n\n") {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		idx := strings.LastIndex(block, "\n")
		if idx == -1 {
			continue
		}
		list = append(list, prompt{
			text:   strings.TrimSpace(block[:idx]),
			source: strings.TrimSpace(block[idx+1:]),
		})
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
	currentPrompt   string
	currentSource   string
	usingUserConfig bool
	gutenbergMode    bool
	loading          bool
	gutenbergFailed  bool
	quitting        bool
	completed     bool
	startTime     time.Time
	endTime       time.Time
	reportView    int
	// terminal dimensions
	width  int
	height int

	// key map
	keys keyMap

	// score submission
	scoreServer  string
	username     string
	sessionID    string
	scoreToken   string
	tokenExpires time.Time
	scorePending bool
	scoreResult  *scoreSubmittedMsg
	scoreErr     error

	// leaderboard view
	leaderboardLoading bool
	leaderboardEntries []leaderboardEntry
	leaderboardErr     error
	leaderboardTable   table.Model
}

type keyMap struct {
	Quit      key.Binding
	Backspace key.Binding
	Left      key.Binding
	Right     key.Binding
	Restart   key.Binding
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
	Backspace: key.NewBinding(key.WithKeys("backspace", "shift+backspace")),
	Left:      key.NewBinding(key.WithKeys("left")),
	Right:     key.NewBinding(key.WithKeys("right")),
	Restart:   key.NewBinding(key.WithKeys("r")),
}

func initialModel(offlineMode bool, scoreServer, username string) model {
	cfg := loadUserConfig()
	var userPrompts []prompt
	if cfg != nil {
		for _, p := range cfg.Prompts {
			if strings.TrimSpace(p.Content) != "" {
				userPrompts = append(userPrompts, prompt{text: p.Content, source: p.Citation})
			}
		}
	}
	usingUserConfig := len(userPrompts) > 0
	prompts := userPrompts
	if len(prompts) == 0 {
		prompts = promptList
	}
	p := prompts[rand.Intn(len(prompts))]
	gutenbergMode := !offlineMode
	m := model{
		keys:            keys,
		currentPrompt:   p.text,
		currentSource:   p.source,
		usingUserConfig: usingUserConfig,
		gutenbergMode:   gutenbergMode,
		loading:         gutenbergMode,
		mistypes:        make(map[int]bool),
		keyErrors:       make(map[rune]int),
		scoreServer:     scoreServer,
		username:        username,
	}
	return m
}

func (m model) Init() tea.Cmd {
	if m.gutenbergMode {
		return fetchGutenbergPromptCmd
	}
	if m.scoreServer != "" && m.username != "" {
		return challengeCmd(m.scoreServer, m.username, sha256Hex(m.currentPrompt))
	}
	return nil
}
