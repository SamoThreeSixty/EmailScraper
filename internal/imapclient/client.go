package imapclient

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap/client"
	"github.com/samothreesixty/EmailScraper/internal/config"
)

type MailClient interface {
	Connect()
	Login(token string)
	GetClient() *client.Client
}

type GmailClient struct {
	Config *config.Config
	client *client.Client
}

func (ic *GmailClient) Connect() {
	address := fmt.Sprintf("%s:%s", ic.Config.Host, ic.Config.Port)

	// Connect to server
	c, err := client.DialTLS(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	ic.client = c
}

func (ic *GmailClient) Login(token string) {
	if token == "" {
		// fallback to username/password
		if err := ic.client.Login(ic.Config.Username, ic.Config.Password); err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Fatal("Not implemented")
}

func (ic *GmailClient) GetClient() *client.Client {
	return ic.client
}
