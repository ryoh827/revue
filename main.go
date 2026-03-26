package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	gh "github.com/ryoh827/revue/github"
	"github.com/ryoh827/revue/tui"
)

func main() {
	client, err := gh.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	model := tui.New(client)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
