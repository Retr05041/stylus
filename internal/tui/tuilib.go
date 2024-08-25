package tui

import (
	"log"
	"stylus/internal/api"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Custom types
type programState uint
type loginState uint
type loginSuccessMsg struct{
    successfulSession *api.Session 
}

// For handling errors in our model
type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// Global consts for states
const (
	stateLogin programState = iota
	stateNotebooks

	stateEmail loginState = iota
	statePassword
)

var (
	programWidth  int // 100
	programHeight int // 40
	signInWidth   = 30
	signInHeight  = 1

	// Program
	programStyle lipgloss.Style
	bannerStyle  lipgloss.Style
	centerStyle  lipgloss.Style

	// Login
	focusedSignInStyle   lipgloss.Style
	unfocusedSignInStyle lipgloss.Style

    // Notebooks
    notebookListStyle lipgloss.Style

    // Utils
    errorStyle lipgloss.Style
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
	Session         api.Session
	ProgramState    programState
	ProgramViewport viewport.Model

	// Login
	LoginState       loginState
	EmailTextArea    textarea.Model
	PasswordTextArea textarea.Model

	// Notebooks
	CachedNotebooks list.Model

	// Utils
	err error
}

func (m *model) SetNotebooks() {
    cachedNotebooks := []list.Item{}  

    for i := range m.Session.Notebooks {
        cachedNotebooks = append(cachedNotebooks, cachedNotebook{title: m.Session.Notebooks[i].Title, desc: m.Session.Notebooks[i].Description, id: m.Session.Notebooks[i].ID})
    }

    m.CachedNotebooks = list.New(cachedNotebooks, list.NewDefaultDelegate(), m.ProgramViewport.Width, m.ProgramViewport.Height/2)
    m.CachedNotebooks.Title = m.Session.Login.User.Username + "'s Notebooks."
}

type cachedNotebook struct {
	title, desc, id string
}

func (n cachedNotebook) Title() string       { return n.title }
func (n cachedNotebook) Description() string { return n.desc }
func (n cachedNotebook) FilterValue() string { return n.title }

// Initialize all global variables then return the model
func InitModel() model {
	termWidth, termHeight, err := term.GetSize(0)
	if err != nil {
		log.Fatal(err)
	}
	programWidth = termWidth - 2
	programHeight = termHeight - 2

	// Program
	programStyle = lipgloss.NewStyle().
		Width(programWidth).
		Height(programHeight).
		Align(lipgloss.Left, lipgloss.Top).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))
	bannerStyle = lipgloss.NewStyle().
		Width(programWidth).
		Height(programHeight/4).
		Align(lipgloss.Center, lipgloss.Center)
	centerStyle = lipgloss.NewStyle().
		Width(programWidth).
		Height(programHeight/2).
		Align(lipgloss.Center, lipgloss.Center)

	// Login
	focusedSignInStyle = lipgloss.NewStyle().
		Width(signInWidth).
		Height(signInHeight).
		Align(lipgloss.Left, lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#69"))
	unfocusedSignInStyle = lipgloss.NewStyle().
		Width(signInWidth).
		Height(signInHeight).
		Align(lipgloss.Left, lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))

    // Notebooks
    notebookListStyle = lipgloss.NewStyle().Margin(1,2)

    // Utils
    errorStyle = lipgloss.NewStyle().
        Width(programWidth/12).
        Height(programHeight/16).
        Align(lipgloss.Right, lipgloss.Bottom).
        Foreground(lipgloss.Color("#c90025"))

	return newModel()
}

// Creates a model using the global variables provided by InitModel()
func newModel() model {
	// Program
	ProgramVp := viewport.New(programWidth, programHeight)

	// Login
	emailTa := textarea.New()
	emailTa.Placeholder = "Email"
	emailTa.Prompt = ""
	emailTa.Focus()
	emailTa.SetWidth(signInWidth)
	emailTa.SetHeight(signInHeight)
	emailTa.ShowLineNumbers = false
	emailTa.KeyMap.InsertNewline.SetEnabled(false)

	passwordTa := textarea.New()
	passwordTa.Placeholder = "Password"
	passwordTa.Prompt = ""
	passwordTa.SetWidth(signInWidth)
	passwordTa.SetHeight(signInHeight)
	passwordTa.ShowLineNumbers = false
	passwordTa.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		ProgramViewport:  ProgramVp,
		ProgramState:     stateLogin,
		EmailTextArea:    emailTa,
		PasswordTextArea: passwordTa,
		LoginState:       stateEmail,
		err:              nil,
	}
}
