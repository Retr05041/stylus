package api

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"log"
	"encoding/json"
)

type Session struct {
	Data struct {
		Login struct {
			Token string `json:"token"`
			User  struct {
				ID       string `json:"id"`
				Username string `json:"username"`
			} `json:"user"`
		} `json:"login"`
	} `json:"data"`
}

func Init(email string, password string) {
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
		log.Fatal(err)
	} 

	// Print the raw response for debugging
	respBytes, _ := json.MarshalIndent(loginResp, "", "  ")
	fmt.Println("Response:", string(respBytes))

	// Access and print the username
	if loginResp.Data.Login.User.Username != "" {
		fmt.Println("Logged in as " + loginResp.Data.Login.User.Username)
	} else {
		fmt.Println("Username field not found in response.")
	}
}
