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
	username := generateUsername()
	if cfg != nil {
		scoreServer = cfg.ScoreServer
		if cfg.Username != "" {
			username = cfg.Username
		}
	}
	if *usernameFlag != "" {
		username = *usernameFlag
	}
	if *offlineMode {
		scoreServer = ""
	}

	p := tea.NewProgram(initialModel(*offlineMode, scoreServer, username))
	if _, err := p.Run(); err != nil {
		fmt.Printf("error starting bubbletea: %v\n", err)
		os.Exit(1)
	}
}
