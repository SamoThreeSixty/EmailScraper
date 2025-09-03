package imapclient

import (
	"github.com/emersion/go-imap/client"
)

type MailClient interface {
	init()
	connect()
	login(token string)
	getClient() *client.Client
}
