package handlers

import (
	"github.com/jaredhaight/lovecms/internal/application"
	"html/template"
	"log/slog"
	"net/http"
	"path"
)

type HomeHandler struct {
	config *application.Config
	logger *slog.Logger
}

func NewHomeHandler(config *application.Config, logger *slog.Logger) *HomeHandler {
	return &HomeHandler{
		config: config,
		logger: logger,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if we don't have a site defined, redirect to setup
	if h.config.SitePath == "" {
		h.logger.Info("No current site defined. Redirecting to setup")
		http.Redirect(w, r, "/setup", http.StatusFound)
	}

	contentPath := path.Join(h.config.SitePath, "content")
	// Load our posts
	posts, err := application.LoadPosts(contentPath)
	if err != nil {
		h.logger.Error("Error loading posts: ", "err", err)
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	templates := []string{
		"./templates/base.gohtml",
		"./templates/home.gohtml",
	}

	ts, err := template.ParseFiles(templates...)

	if err != nil {
		h.logger.Error("Error parsing templates", "err", err)
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", posts)
	if err != nil {
		h.logger.Error("Error parsing templates", "err", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
