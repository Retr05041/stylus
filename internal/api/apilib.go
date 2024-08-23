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
	Description string `json:"description"`
	ID          string `json:"id"`
	Pages       []Page `json:"pages"`
	Title       string `json:"title"`
	UpdatedAt   string `json:"updatedAt"`
}

type Page struct {
	ID        string `json:"id"`
	ParentId  string `json:"parentId"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updatedAt"`
}
