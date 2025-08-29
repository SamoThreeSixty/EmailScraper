package main

import (
	"fmt"
	"log"

	"github.com/samothreesixty/EmailScraper/internal/config"
	"github.com/samothreesixty/EmailScraper/internal/imap"
	// "github.com/emersion/go-imap"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil || conf == nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to server...")
	c, err := imap.Connect(conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()
	fmt.Println("Connected and logged in!")

	messages, err := imap.FetchLastUnseen(c, conf.Inbox, 4)
	if err != nil {
		log.Fatal(err)
	}

	if len(messages) == 0 {
		fmt.Println("No unseen messages")
		return
	}

	fmt.Println("Last unseen messages:")
	for _, msg := range messages {
		fmt.Println("* " + msg.Envelope.Subject)
	}
}
