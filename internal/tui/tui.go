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
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
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
		case "tab":
			if m.ProgramState == stateLogin {
				if m.LoginState == stateEmail {
					m.LoginState = statePassword	
				} else {
					m.LoginState = stateEmail
				}
			}
		}
		switch m.ProgramState {
		case stateLogin:
			switch m.LoginState {
			case stateEmail:
				m.EmailTextArea, cmd = m.EmailTextArea.Update(msg)
				m.EmailTextArea.Focus()
				m.PasswordTextArea.Blur()
				cmds = append(cmds, cmd)
			case statePassword:
				m.PasswordTextArea, cmd = m.PasswordTextArea.Update(msg)
				m.PasswordTextArea.Focus()
				m.EmailTextArea.Blur()
				cmds = append(cmds, cmd)
			}
		}
	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	if m.ProgramState == stateLogin {
		if m.LoginState == stateEmail {
			s += programStyle.Render(centerStyle.Render(lipgloss.JoinVertical(lipgloss.Center, bannerStyle.Render(fmt.Sprintf("%s\n\nLogin", banner)), lipgloss.JoinVertical(lipgloss.Center, focusedSignInStyle.Render(m.EmailTextArea.View()), unfocusedSignInStyle.Render(m.PasswordTextArea.View())))))
		} else if m.LoginState == statePassword {
			s += programStyle.Render(centerStyle.Render(lipgloss.JoinVertical(lipgloss.Center, bannerStyle.Render(fmt.Sprintf("%s\n\nLogin", banner)), lipgloss.JoinVertical(lipgloss.Center, unfocusedSignInStyle.Render(m.EmailTextArea.View()), focusedSignInStyle.Render(m.PasswordTextArea.View())))))
		}
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
