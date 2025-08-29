package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
	// "github.com/emersion/go-imap"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hostString := os.Getenv("IMAP_HOST_NAME")
	if hostString == "" {
		log.Fatal("$IMAP_HOST_NAME must be set")
	}

	hostPort := os.Getenv("IMAP_PORT")
	if hostPort == "" {
		log.Fatal("$IMAP_PORT must be set")
	}

	inbox := os.Getenv("EMAIL_INBOX")
	if inbox == "" {
		log.Fatal("$EMAIL_INBOX must be set")
	}

	username := os.Getenv("IMAP_USERNAME")
	if username == "" {
		log.Fatal("$IMAP_USERNAME must be set")
	}

	password := os.Getenv("IMAP_PASSWORD")
	if password == "" {
		log.Fatal("$IMAP_PASSWORD must be set")
	}

	fmt.Println("Connecting to server...")
	address := fmt.Sprintf("%s:%s", hostString, hostPort)

	// Connect to server
	c, err := client.DialTLS(address, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in")

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// Search for unseen messages
	criteria = imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}

	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	if len(ids) == 0 {
		log.Println("No unseen messages")
		return
	}

	// Keep only the last 4 unseen messages
	if len(ids) > 4 {
		ids = ids[len(ids)-4:]
	}

	// Add the IDs to the sequence set
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...) // <-- use the IDs directly, not AddRange

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	fmt.Println("Last 4 messages:")
	for msg := range messages {
		fmt.Println("* " + msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")

}
