package main

import (
	//"stylus/internal/tui"
	"fmt"
	"log"
	"stylus/internal/api"
)

func main() {
	apiTesting()
	//tui.Start()
}

func apiTesting() {
	session, err := api.Login("", "")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Token: " + session.Login.Token)

	session.GetNotebooks()
	for notebookIndex := range session.Notebooks {
		fmt.Println("Notebook: " + session.Notebooks[notebookIndex].Title)
		for _, page := range session.Notebooks[notebookIndex].Pages {
			fmt.Println("Page name: " + page.Title) 
			fmt.Println("Content: " + page.Content)
		}
	}

}
