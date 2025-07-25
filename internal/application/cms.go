package application

import (
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/jaredhaight/lovecms/internal/templates"
	"github.com/jaredhaight/lovecms/internal/types"
	"github.com/spf13/viper"
)

type CmsHandler struct {
	config *viper.Viper
	logger *slog.Logger
}

func NewCmsHandler(v *viper.Viper, l *slog.Logger) *CmsHandler {
	return &CmsHandler{
		config: v,
		logger: l,
	}
}

// GET /
func (h *CmsHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	// get sitepath
	var sitePath = h.config.GetString("SitePath")

	// if we don't have a site defined, redirect to setup
	if sitePath == "" {
		h.logger.Info("No current site defined. Redirecting to setup")
		http.Redirect(w, r, "/setup", http.StatusFound)
	}

	contentPath := filepath.Join(sitePath, "content")
	// Load our application
	p, err := GetPosts(contentPath)
	if err != nil {
		h.logger.Error("Error loading application: ", "err", err)
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

// GET /editor/?path=Foo
func (h *CmsHandler) GetEditor(w http.ResponseWriter, r *http.Request) {
	// get sitepath
	var sitePath = h.config.GetString("SitePath")

	// if we don't have a site defined, redirect to setup
	if sitePath == "" {
		h.logger.Info("No current site defined. Redirecting to setup")
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}

	// if we're at /post/new, we need to just render the CMS
	if r.URL.Path == "/post/new" {
		// Create a new post form
		c := templates.Editor(types.Post{}, false)
		err := templates.Layout(c).Render(r.Context(), w)

		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	// if we're at /post/edit, we need to load the existing post
	postPath := r.URL.Query().Get("path")

	// if no path is provided, return an error
	if postPath == "" {
		http.Error(w, "Post path required", http.StatusBadRequest)
		return
	}

	// Load the existing post
	post, err := GetPost(postPath)

	if err != nil {
		h.logger.Error("Error loading post", "err", err)
		http.Error(w, "Error loading post", http.StatusInternalServerError)
		return
	}

	c := templates.Editor(post, true)
	err = templates.Layout(c).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// POST /editor/
func (h *CmsHandler) PostEditor(w http.ResponseWriter, r *http.Request) {
	// get sitepath
	var sitePath = h.config.GetString("SitePath")

	// if we don't have a site defined, redirect to setup
	if sitePath == "" {
		h.logger.Info("No current site defined. Redirecting to setup")
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		h.logger.Error("Error parsing form", "err", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	slug := r.FormValue("slug")
	tags := r.FormValue("tags")

	// Parse tags (comma-separated)
	var tagList []string
	if tags != "" {
		for _, tag := range strings.Split(tags, ",") {
			tagList = append(tagList, strings.TrimSpace(tag))
		}
	}

	// Create post
	post := types.Post{
		Metadata: types.FrontMatter{
			Title:       title,
			Date:        time.Now().Format("2006-01-02T15:04:05Z07:00"),
			Draft:       r.FormValue("draft") == "on",
			PublishDate: time.Now().Format("2006-01-02T15:04:05Z07:00"),
			Slug:        slug,
			Tags:        tagList,
		},
		Content: content,
	}

	// Create the post file
	contentPath := filepath.Join(sitePath, "content")
	err = CreatePost(contentPath, post)
	if err != nil {
		h.logger.Error("Error creating post", "err", err)
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
