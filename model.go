package main

import (
	tea "charm.land/bubbletea/v2"
)

type model struct {

	//terminal dimensions
	width  int
	height int
}

func initialModel() model {
	m := model{}
	return m
}

func (m model) Init() tea.Cmd {
	//nothing to init
	return nil
}
