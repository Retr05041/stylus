package api

import (
	"context"
	"github.com/machinebox/graphql"
)

func Login(email string, password string) (*Session, error) {
	client := graphql.NewClient("https://api.codesociety.xyz/api")

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

func (s *Session) GetNotebooks() error {
    notebooksReq.Header.Add("Authorization", "Bearer " + s.Login.Token)    

    if err := s.Client.Run(context.Background(), notebooksReq, &s); err != nil {
        return err 
    }
    return nil
}
