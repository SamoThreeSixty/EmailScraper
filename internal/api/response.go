package api

import (
	"bytes"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/samothreesixty/EmailScraper/internal/models"
)

func ReturnView(w http.ResponseWriter, tmplPath string, data interface{}) {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return
	}
}

func replaceCidImages(htmlBody string, attachments []models.Attachment) string {
	re := regexp.MustCompile(`cid:[^'">]+`)

	i := 0
	return re.ReplaceAllStringFunc(htmlBody, func(match string) string {
		if i < len(attachments) {
			path := attachments[i].Path

			// Ensure it starts with "/"
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}

			i++
			return path
		}
		return match
	})
}
