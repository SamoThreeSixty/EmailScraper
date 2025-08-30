package models

import "github.com/samothreesixty/EmailScraper/internal/db"

type Attachment struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

func ReturnAttachmentFromAttachment(attachment db.Attachment) Attachment {
	return Attachment{
		Type:     attachment.Type,
		Filename: attachment.Filename,
		Path:     attachment.Path,
	}
}

func ReturnAttachmentsFromAttachments(dbAttachments []db.Attachment) []Attachment {
	result := make([]Attachment, 0, len(dbAttachments))
	for _, a := range dbAttachments {
		result = append(result, ReturnAttachmentFromAttachment(a))
	}
	return result
}
