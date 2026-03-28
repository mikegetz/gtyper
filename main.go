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
	flag.Parse()

	p := tea.NewProgram(initialModel(*offlineMode))
	if _, err := p.Run(); err != nil {
		fmt.Printf("error starting bubbletea: %v\n", err)
		os.Exit(1)
	}
}
