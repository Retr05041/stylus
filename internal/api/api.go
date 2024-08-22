package api

import (
	"context"
	"github.com/machinebox/graphql"
)


type Session struct {
	Login struct {
		Token string `json:"token"`
		User  struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"login"`
}

func Init(email string, password string) (string, error) {
	client := graphql.NewClient("https://api.codesociety.xyz/api")

	loginReq := graphql.NewRequest(`
		mutation ($email: String!, $password: String!) {
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
		return "",err
	} 
	return "successfully logged in as " + loginResp.Login.User.Username, nil
}
