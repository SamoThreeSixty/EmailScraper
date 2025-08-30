package repository

import (
	"context"

	"github.com/samothreesixty/EmailScraper/internal/db"
	"github.com/samothreesixty/EmailScraper/internal/models"
)

func GetEmailAttachments(queries *db.Queries, int int) ([]models.Attachment, error) {
	res, err := queries.GetEmailAttachments(context.Background(), int32(int))
	if err != nil {
		return []models.Attachment{}, err
	}
	return models.ReturnAttachmentsFromAttachments(res), nil
}
