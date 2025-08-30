package service

import (
	"context"
	"fmt"
	"log"
	"os"
	filepath2 "path/filepath"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/jhillyerd/enmime"
	"github.com/samothreesixty/EmailScraper/internal/db"
	"github.com/samothreesixty/EmailScraper/internal/utils/format"
)

var attachmentDir = "attachments"

func StartEmailScraper(secondInterval int, c *client.Client, query *db.Queries) {
	cx := context.Context(context.Background())
	ticker := time.NewTicker(time.Second * time.Duration(secondInterval))
	defer ticker.Stop()

	for range ticker.C {
		// Search for new/unseen emails since the last interval
		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{"\\Seen"}
		criteria.Since = time.Now().Add(-24 * time.Second)

		// Make sure the inbox is selected
		_, err := c.Select("INBOX", false)
		if err != nil {
			log.Fatal(err)
		}

		ids, err := c.Search(criteria)
		if err != nil {
			fmt.Println("Search error:", err)
			continue
		}

		if len(ids) == 0 {
			continue // skip fetch entirely
		}

		fmt.Println("Found", len(ids), "new emails")

		// Add the IDs of the found emails to the sequence set for fetching
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		// Fetch the messages concurrently; send errors to the 'done' channel
		messages := make(chan *imap.Message, len(ids))
		done := make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchRFC822}, messages)
		}()

		for msg := range messages {
			section := imap.BodySectionName{}
			var body string
			r := msg.GetBody(&section)
			if r != nil {
				env, err := enmime.ReadEnvelope(r)
				if err != nil {
					log.Println("Failed to parse MIME:", err)
				} else {
					body = env.HTML
				}

				// Save the email into the database
				email, err := query.InsertEmail(cx, db.InsertEmailParams{
					Subject:   msg.Envelope.Subject,
					FromEmail: format.EmailAddressList(msg.Envelope.From),
					ToEmail:   format.EmailAddressList(msg.Envelope.To),
					DateSent:  msg.Envelope.Date,
					Body:      body,
					CreatedAt: time.Now(),
				})
				if err != nil {
					log.Println("Failed to insert email:", err)
				} else {
					fmt.Println("Email saved:", msg.Envelope.Subject)
				}

				// Make sure there is an adequate folder for the attachments
				err = os.MkdirAll(attachmentDir, 0755)
				if err != nil {
					log.Fatal("Failed to create attachment directory:", err)
				}

				// Save the attachments
				for _, att := range env.Attachments {
					_, _ = saveAttachments(query, email, att, cx)
				}

				// Save the inline attachments
				for _, att := range env.Inlines {
					contentId, filename := saveAttachments(query, email, att, cx)

					body = strings.ReplaceAll(body, "cid:"+contentId, "attachments/"+filename)
				}

			}

		}

		// Wait for the fetch operation to complete and log any errors
		if err := <-done; err != nil {
			log.Println("Fetch error:", err)
		}
	}
}

func saveAttachments(query *db.Queries, email db.Email, att *enmime.Part, cx context.Context) (contentId, fileName string) {
	filePath := filepath2.Join(attachmentDir, att.FileName)
	err := os.WriteFile(filePath, att.Content, 0644)
	if err != nil {
		log.Println("Failed to save attachment:", err)
	}
	absPath, _ := filepath2.Abs(filePath)
	_, err = query.SaveAttachment(cx, db.SaveAttachmentParams{
		EmailID:  email.ID,
		Type:     att.ContentType,
		Filename: att.FileName,
		Path:     absPath,
	})
	if err != nil {
		log.Println("Failed to save attachment:", err)
	} else {
		fmt.Println("Attachment saved at:", absPath)
	}

	return att.ContentID, att.FileName
}
