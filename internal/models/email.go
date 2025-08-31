package models

import (
	"time"

	"github.com/samothreesixty/EmailScraper/internal/db"
)

type Email struct {
	ID        int32     `json:"id"`
	Subject   string    `json:"subject"`
	FromEmail string    `json:"from_email"`
	ToEmail   string    `json:"to_email"`
	DateSent  time.Time `json:"date_sent"`
	HtmlBody  string    `json:"html_body"`
	TextBody  string    `json:"text_body"`
}

func ReturnEmailToEmail(email db.Email) Email {
	return Email{
		ID:        email.ID,
		Subject:   email.Subject,
		FromEmail: email.FromEmail,
		ToEmail:   email.ToEmail,
		DateSent:  email.DateSent,
		HtmlBody:  email.HtmlBody,
		TextBody:  email.TextBody,
	}
}
