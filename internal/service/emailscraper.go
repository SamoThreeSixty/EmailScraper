package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	filepath2 "path/filepath"
	"strconv"
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
		criteria.Since = time.Now().Add(-24 * time.Hour)

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
			r := msg.GetBody(&section)
			if r != nil {
				env, err := enmime.ReadEnvelope(r)
				if err != nil {
					log.Println("Failed to parse MIME:", err)
				}

				var htmlBody, txtBody string
				if env != nil {
					htmlBody = env.HTML
					txtBody = env.Text
				}

				// Save the email into the database
				email, err := query.InsertEmail(cx, db.InsertEmailParams{
					Subject:   msg.Envelope.Subject,
					FromEmail: format.EmailAddressList(msg.Envelope.From),
					ToEmail:   format.EmailAddressList(msg.Envelope.To),
					DateSent:  msg.Envelope.Date,
					HtmlBody:  htmlBody,
					TextBody:  txtBody,
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

				allAttachments := collectAttachments(env.Root)
				fmt.Println("Total attachments:", len(allAttachments))

				for _, att := range allAttachments {
					_, _ = saveAttachments(query, email, att, cx)
				}

			}

		}

		// Wait for the fetch operation to complete and log any errors
		if err := <-done; err != nil {
			log.Println("Fetch error:", err)
		}
	}
}

func collectAttachments(part *enmime.Part) []*enmime.Part {
	var atts []*enmime.Part

	// If this part has a filename or Content-ID, treat as attachment
	if part.FileName != "" || part.ContentID != "" {
		atts = append(atts, part)
	}

	// Recursively walk children via FirstChild / NextSibling
	for child := part.FirstChild; child != nil; child = child.NextSibling {
		atts = append(atts, collectAttachments(child)...)
	}

	return atts
}

func saveAttachments(query *db.Queries, email db.Email, att *enmime.Part, cx context.Context) (contentId, fileName string) {
	// Make a safe filename for storage
	timestamp := time.Now().Format("20060102150405")
	fileName = fmt.Sprintf("%s_%s", timestamp, att.FileName)

	// Save the file into the database
	attDb, err := query.SaveAttachment(cx, db.SaveAttachmentParams{
		EmailID:          email.ID,
		Type:             att.ContentType,
		OriginalFilename: att.FileName,
		SavedFilename:    "", // placeholder
		Path:             "", // placeholder
		Cid:              sql.NullString{String: att.ContentID, Valid: att.ContentID != ""},
	})
	if err != nil {
		log.Println("Failed to save attachment into database:", err)
	}

	// Make sure that there is an adequate folder for the attachments
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	day := time.Now().Format("02")
	emailIDStr := strconv.FormatInt(int64(email.ID), 10)

	err = os.MkdirAll(filepath2.Join(attachmentDir, year, month, day, emailIDStr), 0755)
	if err != nil {
		log.Fatal("Failed to create attachment directory:", err)
	}

	filePath := filepath2.Join(attachmentDir, year, month, day, emailIDStr, fileName)
	err = os.WriteFile(filePath, att.Content, 0644)
	if err != nil {
		log.Println("Failed to save attachment:", err)
	}

	// Save the file path into the database
	err = query.UpdateAttachmentPathFilenames(cx, db.UpdateAttachmentPathFilenamesParams{
		ID:               attDb.ID,
		Path:             filePath,
		SavedFilename:    fileName,
		OriginalFilename: att.FileName,
	})
	if err != nil {
		log.Println("Failed to save attachment path and filename into database:", err)
	}

	return att.ContentID, fileName
}
