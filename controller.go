package main

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case challengeReceivedMsg:
		if msg.err == nil {
			m.sessionID = msg.sessionID
			m.scoreToken = msg.token
			m.tokenExpires = msg.expiresAt
		}

	case scoreSubmittedMsg:
		m.scorePending = false
		if msg.err != nil {
			m.scoreErr = msg.err
		} else {
			m.scoreResult = &msg
		}

	case leaderboardFetchedMsg:
		m.leaderboardLoading = false
		if msg.err != nil {
			m.leaderboardErr = msg.err
		} else {
			m.leaderboardEntries = msg.entries
			m.leaderboardTable = buildLeaderboardTable(msg.entries, m.width, m.height)
		}

	case promptFetchedMsg:
		m.loading = false
		if msg.err == nil {
			m.currentPrompt = msg.p.text
			m.currentSource = msg.p.source
			if m.scoreServer != "" && m.username != "" {
				return m, challengeCmd(m.scoreServer, m.username, sha256Hex(m.currentPrompt))
			}
		} else {
			m.gutenbergFailed = true
		}

	case tea.KeyPressMsg:
		if m.loading {
			if key.Matches(msg, m.keys.Quit) {
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}
		if m.completed {
			var tableCmd tea.Cmd
			if m.reportView == 2 {
				m.leaderboardTable, tableCmd = m.leaderboardTable.Update(msg)
			}
			switch {
			case key.Matches(msg, m.keys.Left):
				if m.reportView > 0 {
					m.reportView--
				}
			case key.Matches(msg, m.keys.Right):
				if m.reportView < 2 {
					m.reportView++
				}
			case key.Matches(msg, m.keys.Restart):
				fresh := initialModel(!m.gutenbergMode, m.scoreServer, m.username)
				fresh.width = m.width
				fresh.height = m.height
				return fresh, fresh.Init()
			case key.Matches(msg, m.keys.Quit):
				m.quitting = true
				return m, tea.Quit
			}
			return m, tableCmd
		}
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
				if runes[len(runes)-1] != ' ' {
					delete(m.mistypes, len(runes)-1)
					m.input = string(runes[:len(runes)-1])
				}
			}
			break
		}
		prompt := []rune(m.currentPrompt)
		pos := len([]rune(m.input))
		if pos < len(prompt) {
			if text := msg.Key().Text; text != "" {
				if text == string(prompt[pos]) && m.typoSequence == "" {
					m.input += text
					m.keypressTimes = append(m.keypressTimes, time.Now())
					m.totalKeypresses++
					if len([]rune(m.input)) == 1 && m.startTime.IsZero() {
						m.startTime = m.keypressTimes[0]
					}
					if len([]rune(m.input)) == len(prompt) {
						m.endTime = time.Now()
						m.completed = true
						m.wpmHistory = computeWPMHistory(m.keypressTimes)
						var cmds []tea.Cmd
						if m.scoreServer != "" && m.sessionID != "" {
							m.scorePending = true
							cmds = append(cmds, submitScoreCmd(m.scoreServer, m))
						}
						if m.scoreServer != "" {
							m.leaderboardLoading = true
							cmds = append(cmds, leaderboardCmd(m.scoreServer, sha256Hex(m.currentPrompt)))
						}
						if len(cmds) > 0 {
							return m, tea.Batch(cmds...)
						}
					}
				} else {
					currentWord := m.input[strings.LastIndex(m.input, " ")+1:]
					if len([]rune(m.typoSequence)) < m.width-2-len([]rune(currentWord)) {
						m.typoSequence += text
						m.totalKeypresses++
						if !m.mistypes[pos] {
							m.totalMistypes++
						}
						m.keyErrors[prompt[pos]]++
						m.mistypes[pos] = true
					}
				}
			}
		}
	}
	return m, nil
}

func computeWPMHistory(times []time.Time) []float64 {
	const window = 10
	hist := make([]float64, len(times))
	for i := range times {
		start := max(i-window+1, 0)
		var elapsed float64
		if i == start {
			elapsed = 0.0001
		} else {
			elapsed = times[i].Sub(times[start]).Minutes()
			if elapsed < 0.0001 {
				elapsed = 0.0001
			}
		}
		chars := float64(i - start + 1)
		hist[i] = (chars / 5.0) / elapsed
	}
	return hist
}
