package imapclient

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap/client"
	"github.com/samothreesixty/EmailScraper/internal/config"
)

type GmailClient struct {
	config *config.Config
	client *client.Client
}

func NewGmailClient() (client *client.Client, error error) {
	c := &GmailClient{}

	// Start the IMAP client
	fmt.Println("Connecting to server...")

	err := c.init()
	if err != nil {
		return nil, err
	}

	err = c.connect()
	if err != nil {
		return nil, err
	}

	err = c.login("")
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected and logged in!")

	return c.GetClient(), nil
}

func (ic *GmailClient) init() error {
	conf, err := config.NewConfig()
	if err != nil || conf == nil {
		return err
	}

	ic.config = conf
	return nil
}

func (ic *GmailClient) connect() error {
	address := fmt.Sprintf("%s:%s", ic.config.Host, ic.config.Port)

	// connect to server
	c, err := client.DialTLS(address, nil)
	if err != nil {
		return err
	}

	ic.client = c
	return nil
}

func (ic *GmailClient) login(token string) error {
	if token == "" {
		// fallback to username/password
		if err := ic.client.Login(ic.config.Username, ic.config.Password); err != nil {
			return err
		}

		return nil
	}

	log.Fatal("Not implemented")
	return nil
}

func (ic *GmailClient) GetClient() *client.Client {
	return ic.client
}
