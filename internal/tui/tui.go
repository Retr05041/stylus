package tui

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func Start() {
	// Log to a file. Useful in debugging since you can't really log to stdout.
	// Not required.
	logfilePath := os.Getenv("BUBBLETEA_LOG")
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "simple"); err != nil {
			log.Fatal(err)
		}
	}

	// Initialize our program
	p := tea.NewProgram(model(5))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model int

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}

	case tickMsg:
		m--
		if m <= 0 {
			return m, tea.Quit
		}
		return m, tick
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Hi. This program will exit in %d seconds.\n\nTo quit sooner press q...\n", m)
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}
