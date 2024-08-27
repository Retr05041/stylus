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
	programStateLogin programState = iota
	programStateNotebooks
	programStatePages

	loginStateEmail loginState = iota
	loginStatePassword

	pageStateList pageState = iota
	pageStatePage
	pageStateRender
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
	CachedNotebooks  list.Model
	SelectedNotebook api.Notebook

	// Pages
	PageState    pageState
	CachedPages  list.Model
	SelectedPage api.Page
	EditablePage textarea.Model
	Renderer     glamour.TermRenderer
	RenderedPage viewport.Model

	// Utils
	err                 error
	errNotificationTime time.Time
}

// Custom types
// States
type programState uint
type loginState uint
type pageState uint

// Tea.msg
type loginSuccessMsg struct {
	successfulSession *api.Session
}
type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// Given the list of notebooks from the API, set them in the model for use
func (m *model) SetNotebooks() {
	cachedNotebooks := []list.Item{}

	for _, notebook := range m.Session.Notebooks {
		cachedNotebooks = append(cachedNotebooks, notebook)
	}

	m.CachedNotebooks = list.New(cachedNotebooks, list.NewDefaultDelegate(), m.ProgramViewport.Width/2, m.ProgramViewport.Height/2)
	m.CachedNotebooks.Title = m.Session.Login.User.Username + "'s Notebooks."
	m.CachedNotebooks.SetShowHelp(false)
	m.CachedNotebooks.DisableQuitKeybindings()
}

// Given a specified chachedNotebook, cache all the pages that notebook has
func (m *model) ListPages() {
	cachedPages := []list.Item{}
	var chosenNotebook api.Notebook

	for notebookIndex := range m.Session.Notebooks {
		if m.Session.Notebooks[notebookIndex].ID == m.SelectedNotebook.ID {
			chosenNotebook = m.Session.Notebooks[notebookIndex]
			m.SelectedNotebook = chosenNotebook
			for _, page := range m.Session.Notebooks[notebookIndex].Pages {
				cachedPages = append(cachedPages, page)
			}
			break
		}
	}

	m.CachedPages = list.New(cachedPages, list.NewDefaultDelegate(), programWidth/4, programHeight/2)
	m.CachedPages.Title = chosenNotebook.NotebookTitle
	m.CachedPages.SetShowHelp(false)
	m.CachedPages.DisableQuitKeybindings()

}

// Once selected a page to edit, this sets the contents of the textarea
func (m *model) DisplayEditablePage() {
	m.EditablePage.Reset()
	m.EditablePage.CursorStart()
	m.EditablePage.InsertString("")

	for _, page := range m.SelectedNotebook.Pages {
		if page.ID == m.SelectedPage.ID {
			m.EditablePage.InsertString(page.Content)
			break
		}
	}
}

// Post page selection, called if you want to render the page in markdown
func (m *model) RenderPage() {
	Rstr, _ := m.Renderer.Render(m.SelectedPage.Content)
	m.RenderedPage.SetContent(Rstr)
}

// Sets the selected cached page's contents to that of whats in the textarea - for saving and rendering - since it's just editing the cached pages, leaving the notebook will revert the changes
func (m *model) SaveCachedPageContent() {
	m.SelectedPage.Content = m.EditablePage.Value()
}

// Creates a model
func newModel() model {
	// Global Variable Assignment
	termWidth, termHeight, err := term.GetSize(0)
	if err != nil {
		log.Fatal(err)
	}
	programWidth = termWidth - 2
	programHeight = termHeight - 2

	// Program vars
	programStyle = lipgloss.NewStyle().
		Width(programWidth).
		Height(programHeight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff"))

	// Sign In (Login/Register) vars
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

	// **************************************

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

	// Pages
	EditPageTa := textarea.New()
	EditPageTa.SetWidth(programWidth - (programWidth / 4) - 10)
	EditPageTa.SetHeight(programHeight - 2)

	RenderedPageRenderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
	if err != nil {
		log.Fatal("Error: Markdown renderer failed to init")
	}
	RenderedPageVp := viewport.New(programWidth-(programWidth/4)-10, programHeight-2)

	return model{
		ProgramViewport:   ProgramVp,
		ProgramState:      programStateLogin,
		EmailTextInput:    emailTi,
		PasswordTextInput: passwordTi,
		LoginState:        loginStateEmail,
		EditablePage:      EditPageTa,
		Renderer:          *RenderedPageRenderer,
		RenderedPage:      RenderedPageVp,
		err:               nil,
	}
}
