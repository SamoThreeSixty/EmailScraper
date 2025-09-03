package imapclient

import (
	"errors"

	"github.com/emersion/go-imap/client"
	"github.com/samothreesixty/EmailScraper/internal/config"
)

type MicrosoftClient struct {
	Config *config.Config
	client *client.Client
}

func NewMicrosoftClient() (client *client.Client, error error) {
	// TODO
	return nil, errors.New("microsoft client not implemented")
}

func (ic *MicrosoftClient) connect() {

}

func (ic *MicrosoftClient) login() {

}

func (ic *MicrosoftClient) getClient() *client.Client {
	return ic.client
}
