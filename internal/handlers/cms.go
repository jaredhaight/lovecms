package handlers

import (
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/jaredhaight/lovecms/internal/application"
	"github.com/jaredhaight/lovecms/internal/templates"
	"github.com/spf13/viper"
)

type PostHandler struct {
	config *viper.Viper
	logger *slog.Logger
}

func NewPostHandler(v *viper.Viper, l *slog.Logger) *PostHandler {
	return &PostHandler{
		config: v,
		logger: l,
	}
}

// GET /post/new
// GET /post/edit?path=Foo
func (h *PostHandler) GetCMS(w http.ResponseWriter, r *http.Request) {
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
		c := templates.PostForm(application.Post{}, false)
		err := templates.Layout(c).Render(r.Context(), w)

		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	// if we're at /post/edit, we need to load the existing post
	postPath := r.URL.Query().Get("path")
	if postPath != "" {
		// Load the existing post
		post, err := application.GetPost(postPath)
		if err != nil {
			h.logger.Error("Error loading post", "err", err)
			http.Error(w, "Error loading post", http.StatusInternalServerError)
			return
		}
		// Create a post for the form
		post := application.Post{
			Metadata: application.FrontMatter{
				Title:       post.Metadata.Title,
				Date:        post.Metadata.Date,
				Draft:       post.Metadata.Draft,
				PublishDate: post.Metadata.PublishDate,
				Tags:        post.Metadata.Tags,
			},
			Content: post.Content,
	}

	c := templates.PostForm(post, false)
	err := templates.Layout(c).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// POST /post/new
func (h *PostHandler) PostNew(w http.ResponseWriter, r *http.Request) {
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
	post := application.Post{
		Metadata: application.FrontMatter{
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
	err = application.CreatePost(contentPath, post)
	if err != nil {
		h.logger.Error("Error creating post", "err", err)
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	postPath := r.URL.Query().Get("path")
	if postPath == "" {
		http.Error(w, "Post path required", http.StatusBadRequest)
		return
	}

	post, err := application.GetPost(postPath)
	if err != nil {
		h.logger.Error("Error loading post", "err", err)
		http.Error(w, "Error loading post", http.StatusInternalServerError)
		return
	}

	c := templates.Post(post)
	err = templates.Layout(c).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) PostEdit(w http.ResponseWriter, r *http.Request) {
	postPath := r.URL.Query().Get("path")
	if postPath == "" {
		http.Error(w, "Post path required", http.StatusBadRequest)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		h.logger.Error("Error parsing form", "err", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Load the existing post to get current data
	existingPost, err := application.GetPost(postPath)
	if err != nil {
		h.logger.Error("Error loading existing post", "err", err)
		http.Error(w, "Error loading post", http.StatusInternalServerError)
		return
	}

	// Extract form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	slug := r.FormValue("slug")
	tags := r.FormValue("tags")
	date := r.FormValue("date")
	lastmod := r.FormValue("lastmod")

	// Parse tags (comma-separated)
	var tagList []string
	if tags != "" {
		for _, tag := range strings.Split(tags, ",") {
			tagList = append(tagList, strings.TrimSpace(tag))
		}
	}

	// Set last modified to current time if not provided
	if lastmod == "" {
		lastmod = time.Now().Format("2006-01-02T15:04:05Z07:00")
	}

	// Use existing date if not provided
	if date == "" {
		date = existingPost.Metadata.Date
	}

	// Create updated post
	updatedPost := application.Post{
		FilePath: postPath,
		FileName: existingPost.FileName,
		Metadata: application.FrontMatter{
			Title:        title,
			Date:         date,
			Draft:        r.FormValue("draft") == "on",
			LastModified: lastmod,
			PublishDate:  existingPost.Metadata.PublishDate, // Keep original publish date
			Slug:         slug,
			Tags:         tagList,
		},
		Content: content,
	}

	// Update the post
	err = application.UpdatePost(updatedPost)
	if err != nil {
		h.logger.Error("Error updating post", "err", err)
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	// Redirect back to the post view
	http.Redirect(w, r, "/post?path="+postPath, http.StatusSeeOther)
}
