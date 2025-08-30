package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Println(err)
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, nil, map[string]string{"error": message})
}
