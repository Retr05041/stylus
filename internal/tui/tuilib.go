package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// Custom types
type programState uint
type loginState uint

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

const (
	stateLogin programState = iota
)

var (
	programWidth  = 100
	programHeight = 40
	signInWidth   = 30
	signInHeight  = 1

	programStyle = lipgloss.NewStyle().
			Width(programWidth).
			Height(programHeight).
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ffffff"))

	bannerStyle = lipgloss.NewStyle().
			Width(programWidth).
			Height(10).
			Align(lipgloss.Center, lipgloss.Center)

	signInStyle = lipgloss.NewStyle().
			Width(signInWidth).
			Height(signInHeight).
			Align(lipgloss.Left, lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ffffff"))

	banner = `
   _____ _         _           
  / ____| |       | |          
 | (___ | |_ _   _| |_   _ ___ 
  \___ \| __| | | | | | | / __|
  ____) | |_| |_| | | |_| \__ \
 |_____/ \__|\__, |_|\__,_|___/
              __/ |            
             |___/             
`
)

type model struct {
	// Program
	ProgramState    programState
	ProgramViewport viewport.Model

	// Login / Registration
	EmailTextArea    textarea.Model

	// Utils
	err error
}

func newModel() model {
	// Program
	ProgramVp := viewport.New(programWidth, programHeight)

	// Login
	emailTa := textarea.New()
	emailTa.Placeholder = "Email"
	emailTa.Focus()
	emailTa.Prompt = ""
	emailTa.SetWidth(signInWidth)
	emailTa.SetHeight(signInHeight)
	emailTa.ShowLineNumbers = false
	emailTa.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		ProgramViewport:  ProgramVp,
		ProgramState:     stateLogin,
		EmailTextArea:    emailTa,
		err:              nil,
	}
}
