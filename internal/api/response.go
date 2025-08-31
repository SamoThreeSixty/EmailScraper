package api

import (
	"bytes"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/samothreesixty/EmailScraper/internal/db"
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

func replaceCidImages(htmlBody string, attachments []db.Attachment) string {
	// Compile a regex to match any "cid:..." references in the HTML
	re := regexp.MustCompile(`cid:[^'">]+`)

	// Define a function to handle each match
	var replaceFunc = func(match string) string {
		// Remove the "cid:" prefix from the matched string
		cidValue := strings.TrimPrefix(match, "cid:")

		// Loop through all attachments
		for i := 0; i < len(attachments); i++ {
			att := attachments[i]

			if att.Cid.Valid {
				cleanCid := strings.Trim(att.Cid.String, "<>")

				if cleanCid == cidValue {
					path := att.Path

					// Make sure the path starts with "/"
					if !strings.HasPrefix(path, "/") {
						path = "/" + path
					}

					// Return the path to replace the "cid:..." in HTML
					return path
				}
			}
		}

		// If no matching attachment is found, leave the HTML unchanged
		return match
	}

	// Use the regex to find all "cid:..." references and replace them
	result := re.ReplaceAllStringFunc(htmlBody, replaceFunc)

	// Return the modified HTML
	return result
}
