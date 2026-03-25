package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

var Version = "dev"

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error starting bubbletea: %v\n", err)
		os.Exit(1)
	}
}
