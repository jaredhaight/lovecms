package handlers

import (
	"encoding/json"
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/jaredhaight/lovecms/internal/application"
	"github.com/jaredhaight/lovecms/internal/templates"
	"github.com/spf13/viper"
	"io"
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

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
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

	// Load our application
	p, err := application.GetPost(postPath[0])
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

func (h *PostHandler) Post(w http.ResponseWriter, r *http.Request) {
	// validate post path
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

	p, err := application.GetPost(postPath[0])
	if err != nil {
		h.logger.Error("Error loading post", "err", err)
		http.Error(w, "Error loading post", http.StatusInternalServerError)
		return
	}

	// get content from body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error reading request body", "err", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the JSON data
	var pu application.PostUpdate
	err = json.Unmarshal(body, &pu)
	if err != nil {
		http.Error(w, "Error parsing Post update", http.StatusBadRequest)
		return
	}

	// convert html to markdown
	md, err := htmltomarkdown.ConvertString(pu.Content)
	if err != nil {
		http.Error(w, "Error parsing Post update", http.StatusBadRequest)
		return
	}

	// update post
	p.Title = pu.Title
	p.Content = md

	// save post
	err = application.UpdatePost(p)
	if err != nil {
		h.logger.Error("Error updating post", "err", err)
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
