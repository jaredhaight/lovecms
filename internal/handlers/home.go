package handlers

import (
	"github.com/jaredhaight/lovecms/internal/posts"
	"github.com/jaredhaight/lovecms/internal/templates"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"path/filepath"
)

type HomeHandler struct {
	config *viper.Viper
	logger *slog.Logger
}

func NewHomeHandler(v *viper.Viper, l *slog.Logger) *HomeHandler {
	return &HomeHandler{
		config: v,
		logger: l,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get sitepath
	var sitePath = h.config.GetString("sitePath")

	// if we don't have a site defined, redirect to setup
	if sitePath == "" {
		h.logger.Info("No current site defined. Redirecting to setup")
		http.Redirect(w, r, "/setup", http.StatusFound)
	}

	contentPath := filepath.Join(sitePath, "content")
	// Load our posts
	p, err := posts.GetPosts(contentPath)
	if err != nil {
		h.logger.Error("Error loading posts: ", "err", err)
		http.Error(w, "Error parsing web", http.StatusInternalServerError)
		return
	}

	c := templates.Home(p)
	err = templates.Layout(c).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
