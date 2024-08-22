package tui

import (
	"fmt"
	"log"
	"stylus/internal/api"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://codesociety.xyz/"
// Api url: https://api.codesociety.xyz/api
// Api Playground: https://api.codesociety.xyz/graphiql
// Api Docs: https://graphdoc.io/preview/?endpoint=https://api.codesociety.xyz/api
// GraphQL Docs: https://graphql.org/learn/

type model struct {
	loginResp string
	err    error
}

type loginMsg string

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func Start() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return login
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}

	case loginMsg:
		m.loginResp = string(msg)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m model) View() string {
	s := fmt.Sprintf("Logging in...")
	if m.err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.err)
	} else {
		s += fmt.Sprintf("%s", m.loginResp)
	}
	return s + "\n"
}

func login() tea.Msg {
	loginResp, err := api.Init("", "")
	if  err != nil {
		return errMsg{err}
	}	
	return loginMsg(loginResp)
}
