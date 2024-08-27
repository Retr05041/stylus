package tui

import (
	"fmt"
	"log"
	"stylus/internal/api"
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
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
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
		case "esc": // Esc will be out "backspace" button - it will go backwords through the program
			switch m.ProgramState {
			case programStateLogin:
				return m, tea.Quit
			case programStateNotebooks:
				m.ProgramState = programStateLogin
			case programStatePages:
				switch m.PageState {
				case pageStateList:
					m.ProgramState = programStateNotebooks
				case pageStatePage:
					m.PageState = pageStateList
					m.EditablePage.Blur()
				case pageStateRender:
					m.PageState = pageStatePage
					m.EditablePage.Focus()
				}
			}
		case "ctrl+c":
			if m.ProgramState == programStatePages {
				if m.PageState == pageStatePage {
					m.SaveCachedPageContent()
					m.RenderPage()
					m.PageState = pageStateRender
				} 
			}
		case "tab":
			switch m.ProgramState {
			case programStateLogin:
				switch m.LoginState {
				case loginStateEmail:
					m.LoginState = loginStatePassword
				case loginStatePassword:
					m.LoginState = loginStateEmail
				}
			}
		case "enter":
			switch m.ProgramState {
			case programStateLogin:
				cmds = append(cmds, LoginToApi(m.EmailTextInput.Value(), m.PasswordTextInput.Value()))
			case programStateNotebooks:
				selectedNotebook, ok := m.CachedNotebooks.SelectedItem().(api.Notebook)
				if ok {
					m.SelectedNotebook = selectedNotebook
					m.ListPages()
					m.ProgramState = programStatePages
					m.PageState = pageStateList
				}
			case programStatePages:
				switch m.PageState {
				case pageStateList:
					selectedPage, ok := m.CachedPages.SelectedItem().(api.Page)
					if ok {
						m.SelectedPage = selectedPage
						m.EditablePage.Reset()
						m.DisplayEditablePage()
						m.PageState = pageStatePage
						m.EditablePage.Focus()
					}
				}
			}
		}
		switch m.ProgramState {
		case programStateLogin:
			switch m.LoginState {
			case loginStateEmail:
				m.EmailTextInput, cmd = m.EmailTextInput.Update(msg)
				m.EmailTextInput.Focus()
				m.PasswordTextInput.Blur()
				cmds = append(cmds, cmd)
			case loginStatePassword:
				m.PasswordTextInput, cmd = m.PasswordTextInput.Update(msg)
				m.PasswordTextInput.Focus()
				m.EmailTextInput.Blur()
				cmds = append(cmds, cmd)
			}
		case programStateNotebooks:
			m.CachedNotebooks, cmd = m.CachedNotebooks.Update(msg)
			cmds = append(cmds, cmd)
		case programStatePages:
			m.CachedPages, cmd = m.CachedPages.Update(msg)
			cmds = append(cmds, cmd)
			m.EditablePage, cmd = m.EditablePage.Update(msg)
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
		m.EmailTextInput.Reset()
		m.PasswordTextInput.Reset()
		m.ProgramState = programStateNotebooks
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
	case programStateLogin:
		switch m.LoginState {
		case loginStateEmail:
			programContent += lipgloss.JoinVertical(
				lipgloss.Center,
				headerBannerStyle.Render(fmt.Sprintf("%s\n", banner)),
				centerSignInStyle.Render(
					lipgloss.JoinVertical(
						lipgloss.Center,
						"Login\n"+focusedSignInStyle.Render(m.EmailTextInput.View()),
						unfocusedSignInStyle.Render(m.PasswordTextInput.View()))))
		case loginStatePassword:
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
				programContent += "\n" + signInErrorStyle.Render("Error: "+m.err.Error())
			}
		}
	case programStateNotebooks:
		programContent += notebookListStyle.Render(m.CachedNotebooks.View())
	case programStatePages:
		switch m.PageState {
		case pageStatePage, pageStateList:
			programContent += lipgloss.JoinHorizontal(lipgloss.Center, pageListStyle.Render(m.CachedPages.View()), pageStyle.Render(m.EditablePage.View()))
		case pageStateRender:
			programContent += lipgloss.JoinHorizontal(lipgloss.Center, pageListStyle.Render(m.CachedPages.View()), pageStyle.Render(m.RenderedPage.View()))
		}
	}

	return programStyle.Render(programContent)
}
