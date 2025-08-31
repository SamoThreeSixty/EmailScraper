package api

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/samothreesixty/EmailScraper/internal/models"
	"github.com/samothreesixty/EmailScraper/internal/repository"
)

type EmailWithAttachments struct {
	Email       models.Email        `json:"email"`
	BodyHtml    template.HTML       `json:"body_html"`
	Attachments []models.Attachment `json:"attachments"`
}

func GetEmail(w http.ResponseWriter, r *http.Request) {
	emailIDStr := chi.URLParam(r, "id")
	emailID, err := strconv.Atoi(emailIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid email ID")
		return
	}

	// Fetch email
	email, err := repository.GetEmail(dbConn, context.Background(), int32(emailID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch attachments
	attachments, err := repository.GetEmailAttachments(dbConn, int(email.ID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	emailWithAttachments := EmailWithAttachments{
		Email:       email,
		Attachments: attachments,
	}

	RespondWithJSON(w, r, emailWithAttachments)
}

func GetEmailView(w http.ResponseWriter, r *http.Request) {
	emailIDStr := chi.URLParam(r, "id")
	emailID, err := strconv.Atoi(emailIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid email ID")
		return
	}

	// Fetch email
	email, err := repository.GetEmail(dbConn, context.Background(), int32(emailID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch attachments
	attachments, err := repository.GetEmailAttachments(dbConn, int(email.ID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	emailWithAttachments := EmailWithAttachments{
		Email:       email,
		BodyHtml:    template.HTML(email.HtmlBody),
		Attachments: attachments,
	}

	ReturnView(w, "templates/email.html", emailWithAttachments)
}

func GetEmails(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, r, "ok")
}
