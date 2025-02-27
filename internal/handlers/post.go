package handlers

import (
	"github.com/jaredhaight/lovecms/internal/application"
	"github.com/jaredhaight/lovecms/internal/posts"
	"github.com/jaredhaight/lovecms/internal/templates"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type PostHandler struct {
	config application.Config
	logger slog.Logger
}

func NewPostHandler(config application.Config, logger slog.Logger) *PostHandler {
	return &PostHandler{
		config: config,
		logger: logger,
	}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get our post path
	postPath, ok := r.URL.Query()["path"]

	if !ok || postPath[0] == "" {
		h.logger.Error("Path parameter not found")
		http.Error(w, "Path parameter not found", http.StatusBadRequest)
		return
	}

	// if we don't have a site defined, redirect to setup
	if h.config.SitePath == "" {
		h.logger.Info("No current site defined. Redirecting to setup")
		http.Redirect(w, r, "/setup", http.StatusFound)
	}

	postFullPath := filepath.Join(h.config.SitePath, "content", postPath[0])

	// check if file exists
	_, err := os.Stat(postFullPath)

	if err != nil {
		h.logger.Error("File does not exist", "err", err)
		http.Error(w, "File does not exist", http.StatusInternalServerError)
		return
	}

	// Load our posts
	p, err := posts.GetPost(postFullPath)
	if err != nil {
		h.logger.Error("Error loading post", "err", err)
		http.Error(w, "Error loading post", http.StatusInternalServerError)
		return
	}

	c := templates.Post(p)
	err = templates.Layout(c).Render(r.Context(), w)

	if err != nil {
		h.logger.Error("Error rendering template", "err", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
