package main

import (
	//"stylus/internal/tui"
    "stylus/internal/api"
    "log"
    "fmt"
)

func main() {
    session, err := api.Login("", "")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Logged in as: " + session.Login.User.Username)
    fmt.Println("ID: " + session.Login.User.ID)
    fmt.Println("Token: " + session.Login.Token)
   //tui.Start()
}
