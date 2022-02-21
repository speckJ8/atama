package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speckJ8/atama/display"
)

func main() {
	p := tea.NewProgram(&display.DisplayModel{}, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		log.Fatalf("%s", err.Error())
	}
	/*
		mem := NewMemory(1 << 16)
		proc := NewProcessor(&mem)
		proc.Start()
	*/
}
