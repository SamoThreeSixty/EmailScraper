package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

func AttachmentHandler(w http.ResponseWriter, r *http.Request) {
	baseFolderPath := chi.URLParam(r, "baseFolderPath")
	year := chi.URLParam(r, "year")
	month := chi.URLParam(r, "month")
	day := chi.URLParam(r, "day")
	emailId := chi.URLParam(r, "emailId")
	savedFilename := chi.URLParam(r, "savedFilename")

	folderPath := filepath.Join(baseFolderPath, year, month, day, emailId)

	filePath := filepath.Join(folderPath, savedFilename)
	_, err := os.Open(filePath)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "File not found")
		return
	}

	http.ServeFile(w, r, filePath)
}
