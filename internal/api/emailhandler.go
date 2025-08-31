package api

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/samothreesixty/EmailScraper/internal/db"
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
	email, err := dbConn.GetEmail(context.Background(), int32(emailID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch attachments
	attachments, err := dbConn.GetEmailAttachments(context.Background(), email.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Replace backslashes with forward slashes
	for i := range attachments {
		attachments[i].Path = strings.ReplaceAll(attachments[i].Path, "\\", "/")
	}

	bodyWithImages := replaceCidImages(email.HtmlBody, attachments)

	// Create a new slice of attachments that only contains the ones that are downloadable
	downloadableAttachments := make([]db.Attachment, 0, len(attachments))
	for _, att := range attachments {
		if !att.Cid.Valid {
			downloadableAttachments = append(downloadableAttachments, att)
		}
	}

	emailWithAttachments := EmailWithAttachments{
		Email:       models.ReturnEmailToEmail(email),
		BodyHtml:    template.HTML(bodyWithImages),
		Attachments: models.ReturnAttachmentsFromAttachments(downloadableAttachments),
	}

	ReturnView(w, "templates/email.html", emailWithAttachments)
}

func GetEmails(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, r, "ok")
}
