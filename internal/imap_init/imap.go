package imap_init

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap/client"
)

func Connect(host, port, username, password string) (*client.Client, error) {
	address := fmt.Sprintf("%s:%s", host, port)

	// Connect to server
	c, err := client.DialTLS(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Login
	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}

	return c, nil
}
