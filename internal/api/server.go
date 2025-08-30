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

	r.Get("/api/v1/email/{id}", GetEmail)
	r.Get("/api/v1/emails", GetEmails)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func methodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}
