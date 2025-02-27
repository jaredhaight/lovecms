package handlers

import (
	"github.com/jaredhaight/lovecms/internal/application"
	"github.com/jaredhaight/lovecms/internal/posts"
	"github.com/jaredhaight/lovecms/internal/templates"
	"log/slog"
	"net/http"
	"path/filepath"
)

type HomeHandler struct {
	config application.Config
	logger slog.Logger
}

func NewHomeHandler(config application.Config, logger slog.Logger) *HomeHandler {
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

	contentPath := filepath.Join(h.config.SitePath, "content")
	// Load our posts
	p, err := posts.GetPosts(contentPath)
	if err != nil {
		h.logger.Error("Error loading posts: ", "err", err)
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	c := templates.Home(p)
	err = templates.Layout(c).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
