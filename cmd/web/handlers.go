package main

import (
	"html/template"
	"net/http"
	"os"
)

func (app *application) listPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	site := app.loveConfig.CurrentSite

	app.logger.Debug("Got site", "site", site)
	// iterate through the content folder and get markdown files.
	entries, err := os.ReadDir(site)
	if err != nil {
		app.logger.Error("error reading dir", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	parsed, err := template.ParseFiles("./ui/html/home.tmpl")
	if err != nil {
		app.logger.Error("template parsing error", "error", err)
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	err = parsed.Execute(w, entries)
	if err != nil {
		app.logger.Error("template executing error", "error", err)
		http.Error(w, "Template executing error", http.StatusInternalServerError)
	}
}
