package handlers

import (
	"github.com/jaredhaight/lovecms/internal/models"
	"html/template"
	"log/slog"
	"net/http"
)

type HomeHandler struct {
	logger *slog.Logger
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	posts := make([]models.Post, 0)

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
