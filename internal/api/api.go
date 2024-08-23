package api

import (
	"context"
	"github.com/machinebox/graphql"
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
}

type Notebook struct {

}

func Login(email string, password string) (*Session, error) {
	client := graphql.NewClient("https://api.codesociety.xyz/api")

	loginReq := graphql.NewRequest(`
		mutation Login($email: String!, $password: String!) {
			login(email: $email, password: $password) {
				token
				user {
					id
					username
				}
			}
		}`)

	loginReq.Var("email", email)
	loginReq.Var("password", password)

	var loginResp Session
	ctx := context.Background()

	if err := client.Run(ctx, loginReq, &loginResp); err != nil {
		return nil,err
	} 
    loginResp.Client = client
	return &loginResp, nil
}

func (s *Session) GetNotebooks() []Notebook {
    return []Notebook{}
}
