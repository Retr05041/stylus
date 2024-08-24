package main

import (
	//"stylus/internal/tui"
	"stylus/internal/api"
    "fmt"
    "log"
)

func main() {
	session, err := api.Login("", "")
	if err != nil {
		log.Fatal(err.Error())
        return
	}
	fmt.Println("Logged in as: " + session.Login.User.Username)
	fmt.Println("ID: " + session.Login.User.ID)
	fmt.Println("Token: " + session.Login.Token)

	fmt.Println("Querying notebooks...")
    if err := session.GetNotebooks(); err != nil {
        log.Fatal(err.Error())
        return
    }

    for i := range session.Notebooks {
        fmt.Println(session.Notebooks[i].Title)
        fmt.Println(session.Notebooks[i].Description)
    }

    
	//tui.Start()
}
