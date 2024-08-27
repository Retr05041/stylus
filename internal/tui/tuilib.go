package tui

import (
	"log"
	"stylus/internal/api"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Global consts for states
const (
	stateLogin programState = iota
	stateNotebooks
	statePages

	stateEmail loginState = iota
	statePassword

	statePageList pageState = iota
	statePage
	statePageRender
)

var (
	programWidth  int // 100
	programHeight int // 40
	signInWidth   = 30

	// Program
	programStyle lipgloss.Style

	// Login
	headerBannerStyle    lipgloss.Style
	centerSignInStyle    lipgloss.Style
	focusedSignInStyle   lipgloss.Style
	unfocusedSignInStyle lipgloss.Style
	signInErrorStyle     lipgloss.Style

	// Notebooks
	notebookListStyle lipgloss.Style
	pageListStyle     lipgloss.Style
	pageStyle         lipgloss.Style

	// Utils
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
	LoginState        loginState
	EmailTextInput    textinput.Model
	PasswordTextInput textinput.Model

	// Notebooks
	CachedNotebooks    list.Model
	SelectedNotebookID string
	SelectedNotebook   api.Notebook

	// Pages
	PageState      pageState
	CachedPages    list.Model
	SelectedPageID string
	EditablePage   textarea.Model
	RenderedPage   viewport.Model

	// Utils
	err                 error
	errNotificationTime time.Time
}

// Custom types
type programState uint
type loginState uint
type loginSuccessMsg struct {
	successfulSession *api.Session
}
type pageState uint

// For handling errors in our model
type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// Lists for Cached Notebooks and Pages
type cachedNotebook struct {
	title, desc, id string
}

func (n cachedNotebook) Title() string       { return n.title }
func (n cachedNotebook) Description() string { return n.desc }
func (n cachedNotebook) FilterValue() string { return n.title }

type cachedPage struct {
	title, id, updatedAt string
}

func (p cachedPage) Title() string       { return p.title }
func (p cachedPage) Description() string { return p.updatedAt }
func (p cachedPage) FilterValue() string { return p.id }

func (m *model) SetNotebooks() {
	cachedNotebooks := []list.Item{}

	for i := range m.Session.Notebooks {
		cachedNotebooks = append(cachedNotebooks, cachedNotebook{title: m.Session.Notebooks[i].Title, desc: m.Session.Notebooks[i].Description, id: m.Session.Notebooks[i].ID})
	}

	m.CachedNotebooks = list.New(cachedNotebooks, list.NewDefaultDelegate(), m.ProgramViewport.Width/2, m.ProgramViewport.Height/2)
	m.CachedNotebooks.Title = m.Session.Login.User.Username + "'s Notebooks."
	m.CachedNotebooks.SetShowHelp(false)
	m.CachedNotebooks.DisableQuitKeybindings()
}

func (m *model) ListPages() {
	cachedPages := []list.Item{}
	var chosenNotebook api.Notebook

	for notebookIndex := range m.Session.Notebooks {
		if m.Session.Notebooks[notebookIndex].ID == m.SelectedNotebookID {
			chosenNotebook = m.Session.Notebooks[notebookIndex]
			m.SelectedNotebook = chosenNotebook
			for _, page := range m.Session.Notebooks[notebookIndex].Pages {
				cachedPages = append(cachedPages, cachedPage{title: page.Title, id: page.ID, updatedAt: page.UpdatedAt})
			}
			break
		}
	}

	m.CachedPages = list.New(cachedPages, list.NewDefaultDelegate(), programWidth/4, programHeight/2)
	m.CachedPages.Title = chosenNotebook.Title
	m.CachedPages.SetShowHelp(false)
	m.CachedPages.DisableQuitKeybindings()

	// So we display nothing... once listing is done
	m.EditablePage = textarea.New()
	m.EditablePage.SetWidth(programWidth - (programWidth / 4) - 10)
	m.EditablePage.SetHeight(programHeight - 2)
}

func (m *model) DisplayEditablePage() {
	m.EditablePage.CursorStart()
	m.EditablePage.InsertString("")

	for _, page := range m.SelectedNotebook.Pages {
		if page.ID == m.SelectedPageID {
			m.EditablePage.InsertString(page.Content)
			break
		}
	}
}

func (m *model) RenderPage() {
	m.RenderedPage = viewport.New(programWidth-(programWidth/4)-10, programHeight-2)
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle())

	for _, page := range m.SelectedNotebook.Pages {
		if page.ID == m.SelectedPageID {
			Rstr, _:= r.Render(page.Content)
			m.RenderedPage.SetContent(Rstr)
			break
		}
	}

}

func (m *model) SavePageContent() {
	for pageIndex, page := range m.SelectedNotebook.Pages {
		if page.ID == m.SelectedPageID {
			m.SelectedNotebook.Pages[pageIndex].Content = m.EditablePage.Value()
			break
		}
	}
}

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
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))

	// Sign In (Login/Register)
	headerBannerStyle = lipgloss.NewStyle().
		Width(programWidth).
		Height(programHeight/4). // Banner block takes up 1/4 of the program window
		Align(lipgloss.Center, lipgloss.Center)
	centerSignInStyle = lipgloss.NewStyle().
		Width(programWidth).
		Height(programHeight/4). // Anything using this style gets a block height of half the program window (plus banner makes total block size 3/4) - the login fields use this in the same block
		Align(lipgloss.Center, lipgloss.Bottom)
	focusedSignInStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#69"))
	unfocusedSignInStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))
	signInErrorStyle = lipgloss.NewStyle().
		Width(programWidth).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#ff0000"))

	// Notebooks + Pages
	notebookListStyle = lipgloss.NewStyle().
		Width(programWidth).
		Height(programHeight-2).
		Align(lipgloss.Center, lipgloss.Center).
		Margin(1, 2)
	pageListStyle = lipgloss.NewStyle().
		Width(programWidth/4).
		Height(programHeight-2).
		Align(lipgloss.Left, lipgloss.Top).
		Margin(1, 2)
	pageStyle = lipgloss.NewStyle().
		Width(programWidth - (programWidth / 4)).
		Height(programHeight - 2)

	return newModel()
}

// Creates a model using the global variables provided by InitModel()
func newModel() model {
	// Program
	ProgramVp := viewport.New(programWidth, programHeight)

	// Login
	emailTi := textinput.New()
	emailTi.Placeholder = "Email"
	emailTi.Prompt = ""
	emailTi.Focus()
	emailTi.Width = signInWidth

	passwordTi := textinput.New()
	passwordTi.Placeholder = "Password"
	passwordTi.Prompt = ""
	passwordTi.Width = signInWidth
	passwordTi.EchoMode = textinput.EchoPassword
	passwordTi.EchoCharacter = '*'

	return model{
		ProgramViewport:   ProgramVp,
		ProgramState:      stateLogin,
		EmailTextInput:    emailTi,
		PasswordTextInput: passwordTi,
		LoginState:        stateEmail,
		err:               nil,
	}
}
