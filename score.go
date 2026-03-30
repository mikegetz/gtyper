package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "charm.land/bubbletea/v2"
)

const clientVersion = "1.0.0"

// Message types

type challengeReceivedMsg struct {
	sessionID string
	token     string
	expiresAt time.Time
	err       error
}

type scoreSubmittedMsg struct {
	scoreID  string
	adjWPM   float64
	rank     int
	eligible bool
	err      error
}

type leaderboardEntry struct {
	Rank         int     `json:"rank"`
	Username     string  `json:"username"`
	AdjWPM       float64 `json:"adj_wpm"`
	Accuracy     float64 `json:"accuracy"`
	PromptSource string  `json:"prompt_source"`
}

type leaderboardFetchedMsg struct {
	entries []leaderboardEntry
	err     error
}

// sha256Hex returns the hex-encoded SHA-256 digest of text.
func sha256Hex(text string) string {
	sum := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", sum)
}

// challengeCmd posts to /v1/challenge and returns a challengeReceivedMsg.
func challengeCmd(serverURL, username, promptHash string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{
			"prompt_hash": promptHash,
			"username":    username,
		})
		resp, err := http.Post(serverURL+"/v1/challenge", "application/json", bytes.NewReader(body))
		if err != nil {
			return challengeReceivedMsg{err: err}
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return challengeReceivedMsg{err: fmt.Errorf("challenge: HTTP %d", resp.StatusCode)}
		}
		var result struct {
			SessionID string `json:"session_id"`
			Token     string `json:"token"`
			ExpiresAt string `json:"expires_at"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return challengeReceivedMsg{err: err}
		}
		exp, err := time.Parse(time.RFC3339, result.ExpiresAt)
		if err != nil {
			return challengeReceivedMsg{err: err}
		}
		return challengeReceivedMsg{
			sessionID: result.SessionID,
			token:     result.Token,
			expiresAt: exp,
		}
	}
}

// submitScoreCmd posts to /v1/scores and returns a scoreSubmittedMsg.
func submitScoreCmd(serverURL string, m model) tea.Cmd {
	return func() tea.Msg {
		times := m.keypressTimes
		if len(times) < 2 {
			return scoreSubmittedMsg{err: fmt.Errorf("insufficient keypress data")}
		}

		ms := make([]int64, len(times))
		for i, t := range times {
			ms[i] = t.UnixMilli()
		}

		elapsedMs := float64(ms[len(ms)-1] - ms[0])
		if elapsedMs < 0.1 {
			elapsedMs = 0.1
		}
		elapsedMin := elapsedMs / 60000.0
		clientWPM := (float64(len([]rune(m.currentPrompt))) / 5.0) / elapsedMin

		body, _ := json.Marshal(map[string]any{
			"session_id":        m.sessionID,
			"token":             m.scoreToken,
			"username":          m.username,
			"prompt":            m.currentPrompt,
			"prompt_source":     m.currentSource,
			"keypress_times_ms": ms,
			"total_keypresses":  m.totalKeypresses,
			"client_wpm":        clientWPM,
			"client_version":    clientVersion,
		})

		resp, err := http.Post(serverURL+"/v1/scores", "application/json", bytes.NewReader(body))
		if err != nil {
			return scoreSubmittedMsg{err: err}
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return scoreSubmittedMsg{err: fmt.Errorf("submit: HTTP %d", resp.StatusCode)}
		}
		var result struct {
			ScoreID             string  `json:"score_id"`
			AdjWPM              float64 `json:"adj_wpm"`
			Rank                int     `json:"rank"`
			LeaderboardEligible bool    `json:"leaderboard_eligible"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return scoreSubmittedMsg{err: err}
		}
		return scoreSubmittedMsg{
			scoreID:  result.ScoreID,
			adjWPM:   result.AdjWPM,
			rank:     result.Rank,
			eligible: result.LeaderboardEligible,
		}
	}
}

// leaderboardCmd fetches the top 25 scores for the given prompt from /v1/leaderboard.
func leaderboardCmd(serverURL, promptHash string) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(serverURL + "/v1/leaderboard?limit=25&prompt_hash=" + promptHash)
		if err != nil {
			return leaderboardFetchedMsg{err: err}
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return leaderboardFetchedMsg{err: fmt.Errorf("leaderboard: HTTP %d", resp.StatusCode)}
		}
		var result struct {
			Entries []leaderboardEntry `json:"entries"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return leaderboardFetchedMsg{err: err}
		}
		return leaderboardFetchedMsg{entries: result.Entries}
	}
}
