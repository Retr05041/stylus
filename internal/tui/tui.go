package tui

import (
	"fmt"
	"log"
	"stylus/internal/api"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Api url: https://api.codesociety.xyz/api
// Api Playground: https://api.codesociety.xyz/graphiql
// Api Docs: https://graphdoc.io/preview/?endpoint=https://api.codesociety.xyz/api
// GraphQL Docs: https://graphql.org/learn/

func Start() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}

	m.EmailTextArea, cmd = m.EmailTextArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, nil
}

func (m model) View() string {
	var s string
	if m.ProgramState == stateLogin {
		s += programStyle.Render(lipgloss.JoinVertical(lipgloss.Center, bannerStyle.Render(fmt.Sprintf("\n%s\n", banner)), signInStyle.Render(m.EmailTextArea.View())))
	}
	return s
}

func login() tea.Msg {
	loginResp, err := api.Init("", "")
	if  err != nil {
		return errMsg{err}
	}	
	return loginResp
}
