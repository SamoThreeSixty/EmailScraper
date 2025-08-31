package models

import "github.com/samothreesixty/EmailScraper/internal/db"

type Attachment struct {
	Type             string `json:"type"`
	OriginalFilename string `json:"original_filename"`
	SavedFilename    string `json:"saved_filename"`
	Path             string `json:"path"`
}

func ReturnAttachmentFromAttachment(attachment db.Attachment) Attachment {
	return Attachment{
		Type:             attachment.Type,
		OriginalFilename: attachment.OriginalFilename,
		SavedFilename:    attachment.SavedFilename,
		Path:             attachment.Path,
	}
}

func ReturnAttachmentsFromAttachments(dbAttachments []db.Attachment) []Attachment {
	result := make([]Attachment, 0, len(dbAttachments))
	for _, a := range dbAttachments {
		result = append(result, ReturnAttachmentFromAttachment(a))
	}
	return result
}
