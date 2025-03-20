package handlers

import (
	"github.com/jaredhaight/lovecms/internal/posts"
	"github.com/jaredhaight/lovecms/internal/templates"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"
)

type PostHandler struct {
	config *viper.Viper
	logger *slog.Logger
}

func NewPostHandler(viper *viper.Viper, logger *slog.Logger) *PostHandler {
	return &PostHandler{
		config: viper,
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
	// check if file exists
	_, err := os.Stat(postPath[0])

	if err != nil {
		h.logger.Error("File does not exist", "err", err)
		http.Error(w, "File does not exist", http.StatusInternalServerError)
		return
	}

	// Load our posts
	p, err := posts.GetPost(postPath[0])
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
