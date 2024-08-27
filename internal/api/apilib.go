package api

import "github.com/machinebox/graphql"

var (
	loginReq = graphql.NewRequest(`
		mutation Login($email: String!, $password: String!) {
			login(email: $email, password: $password) {
				token
				user {
					id
					username
				}
			}
		}`)

	notebooksReq = graphql.NewRequest(`
		query GetNotebooks {
		  notebooks {
			id
			title
			description
			updatedAt
			pages {
			  id
			  title
			  parentId
			  updatedAt
			  content
			}
		  }
		}`)
)

type Session struct {
	Client *graphql.Client

	Login struct {
		Token string `json:"token"`
		User  struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"login"`

	Notebooks []Notebook `json:"notebooks"`
}

type Notebook struct {
	ID                  string `json:"id"`
	NotebookTitle       string `json:"title"`
	NotebookDescription string `json:"description"`
	UpdatedAt           string `json:"updatedAt"`
	Pages               []Page `json:"pages"`
}
type Page struct {
	ID        string `json:"id"`
	ParentId  string `json:"parentId"`
	PageTitle string `json:"title"`
	UpdatedAt string `json:"updatedAt"`
	Content   string `json:"content"`
}

func (n Notebook) Title() string       { return n.NotebookTitle }
func (n Notebook) Description() string { return n.NotebookDescription }
func (n Notebook) FilterValue() string { return n.ID }

func (p Page) Title() string       { return p.PageTitle }
func (p Page) Description() string { return p.UpdatedAt }
func (p Page) FilterValue() string { return p.ID }
