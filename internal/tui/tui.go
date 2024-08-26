package tui

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
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
	return tea.Batch(textinput.Blink)
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
		case "enter":
			if m.ProgramState == stateLogin { // This might need to be changed to be done in a cmd so we can handle errors...
				cmds = append(cmds, LoginToApi(m.EmailTextInput.Value(), m.PasswordTextInput.Value()))
			}
		}
		switch m.ProgramState {
		case stateLogin:
			switch m.LoginState {
			case stateEmail:
				m.EmailTextInput, cmd = m.EmailTextInput.Update(msg)
				m.EmailTextInput.Focus()
				m.PasswordTextInput.Blur()
				cmds = append(cmds, cmd)
			case statePassword:
				m.PasswordTextInput, cmd = m.PasswordTextInput.Update(msg)
				m.PasswordTextInput.Focus()
				m.EmailTextInput.Blur()
				cmds = append(cmds, cmd)
			}
		case stateNotebooks:
			m.CachedNotebooks, cmd = m.CachedNotebooks.Update(msg)
			cmds = append(cmds, cmd)
		}

	case errMsg:
		m.err = msg
		m.errNotificationTime = time.Now()
		return m, nil

	case loginSuccessMsg:
		m.Session = *msg.successfulSession
		m.Session.GetNotebooks()
		m.SetNotebooks()
		m.ProgramState = stateNotebooks
		return m, nil

	default:
		return m, nil
	}

	if m.err != nil && time.Since(m.errNotificationTime) > 3*time.Second {
		m.err = nil
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var programContent string
	switch m.ProgramState {
	case stateLogin:
		switch m.LoginState {
		case stateEmail:
			programContent += lipgloss.JoinVertical(
					lipgloss.Center,
					headerBannerStyle.Render(fmt.Sprintf("%s\n", banner)),
					centerSignInStyle.Render(
						lipgloss.JoinVertical(
							lipgloss.Center,
							"Login\n"+focusedSignInStyle.Render(m.EmailTextInput.View()),
							unfocusedSignInStyle.Render(m.PasswordTextInput.View()))))
		case statePassword:
			programContent += lipgloss.JoinVertical(
					lipgloss.Center,
					headerBannerStyle.Render(fmt.Sprintf("%s\n", banner)),
					centerSignInStyle.Render(
						lipgloss.JoinVertical(
							lipgloss.Center,
							"Login\n"+unfocusedSignInStyle.Render(m.EmailTextInput.View()),
							focusedSignInStyle.Render(m.PasswordTextInput.View()))))
		}
		if m.err != nil {
			duration := time.Since(m.errNotificationTime) 
			if duration < 3*time.Second {
				programContent += "\n" + signInErrorStyle.Render("Error: " + m.err.Error())
			}
		}
	case stateNotebooks:
		programContent += centerSignInStyle.Render(notebookListStyle.Render(m.CachedNotebooks.View()))
	}


	return programStyle.Render(programContent)
}
