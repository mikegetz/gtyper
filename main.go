package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

var Version = "dev"

func main() {
	offlineMode := flag.Bool("o", false, "use offline prompts instead of Project Gutenberg")
	usernameFlag := flag.String("u", "", "username for leaderboard submission")
	flag.Parse()

	cfg := loadUserConfig()
	scoreServer := ""
	username := cfg.Username
	if cfg != nil {
		scoreServer = cfg.ScoreServer
	}
	if *usernameFlag != "" {
		username = *usernameFlag
	}

	p := tea.NewProgram(initialModel(*offlineMode, scoreServer, username))
	if _, err := p.Run(); err != nil {
		fmt.Printf("error starting bubbletea: %v\n", err)
		os.Exit(1)
	}
}
