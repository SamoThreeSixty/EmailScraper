package repository

import (
	"context"

	"github.com/samothreesixty/EmailScraper/internal/db"
	"github.com/samothreesixty/EmailScraper/internal/models"
)

func GetEmail(queries db.Queries, ctx context.Context) (models.Email, error) {
	res, err := queries.GetEmail(ctx, 1)
	if err != nil {
		return models.Email{}, err
	}
	return models.ReturnEmailToEmail(res), nil
}
