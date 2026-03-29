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

	case promptFetchedMsg:
		m.loading = false
		if msg.err == nil {
			m.currentPrompt = msg.p.text
			m.currentSource = msg.p.source
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
			switch {
			case key.Matches(msg, m.keys.Left):
				if m.reportView > 0 {
					m.reportView--
				}
			case key.Matches(msg, m.keys.Right):
				if m.reportView < 1 {
					m.reportView++
				}
			case key.Matches(msg, m.keys.Restart):
				fresh := initialModel(!m.gutenbergMode)
				fresh.width = m.width
				fresh.height = m.height
				return fresh, fresh.Init()
			case key.Matches(msg, m.keys.Quit):
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
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
					m.keypressTimes = append(m.keypressTimes, time.Now())
					m.totalKeypresses++
					if len([]rune(m.input)) == 1 && m.startTime.IsZero() {
						m.startTime = m.keypressTimes[0]
					}
					if len([]rune(m.input)) == len(prompt) {
						m.endTime = time.Now()
						m.completed = true
						m.wpmHistory = computeWPMHistory(m.keypressTimes)
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
