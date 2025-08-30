package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/jhillyerd/enmime"
	"github.com/samothreesixty/EmailScraper/internal/config"
	"github.com/samothreesixty/EmailScraper/internal/db"
	"github.com/samothreesixty/EmailScraper/internal/imapclient"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil || conf == nil {
		log.Fatal(err)
	}

	dbConn, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to server...")
	c, err := imapclient.Connect(conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()
	fmt.Println("Connected and logged in!")

	messages, err := imapclient.FetchLastUnseen(c, conf.Inbox, 4)
	if err != nil {
		log.Fatal(err)
	}

	if len(messages) == 0 {
		fmt.Println("No unseen messages")
		return
	}

	section := &imap.BodySectionName{}

	fmt.Println("Last unseen messages:")
	for _, msg := range messages {
		var body string
		r := msg.GetBody(section)
		if r != nil {
			env, err := enmime.ReadEnvelope(r)
			if err != nil {
				log.Println("Failed to parse MIME:", err)
			} else {
				body = env.HTML
			}
		}

		err := dbConn.InsertEmail(context.Background(), db.InsertEmailParams{
			Subject:   msg.Envelope.Subject,
			FromEmail: formatAddressList(msg.Envelope.From),
			ToEmail:   formatAddressList(msg.Envelope.To),
			DateSent:  msg.Envelope.Date,
			Body:      body,
			CreatedAt: time.Now(),
		})
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func formatAddressList(list []*imap.Address) string {
	var result string
	for i, addr := range list {
		if i > 0 {
			result += ", "
		}
		result += addr.MailboxName + "@" + addr.HostName
	}
	return result
}
