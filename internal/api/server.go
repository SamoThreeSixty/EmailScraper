package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/samothreesixty/EmailScraper/internal/db"
)

var dbConn *db.Queries

func InitAPI(db *db.Queries) {
	dbConn = db
}

func StartAPIService() {
	fmt.Println("Listening to api requests on port 8080")
	r := chi.NewRouter()

	r.Get("/api/v1/email/{id}", GetEmailView)
	r.Get("/api/v1/email/{id}/data", GetEmail)
	r.Get("/api/v1/emails", GetEmails)

	r.Handle("/attachments/*", http.StripPrefix("/attachments/", http.FileServer(http.Dir("./attachments"))))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
