package application

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type CmsHandler struct {
	config    *viper.Viper
	logger    *slog.Logger
	templates embed.FS
}

type HomeData struct {
	Posts []Post
}

type EditorData struct {
	Post   Post
	IsEdit bool
}

func NewCmsHandler(v *viper.Viper, l *slog.Logger, t embed.FS) *CmsHandler {
	return &CmsHandler{
		config:    v,
		logger:    l,
		templates: t,
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
		return
	}

	contentPath := filepath.Join(sitePath, "content")
	// Load our application
	p, err := GetPosts(contentPath)
	if err != nil {
		h.logger.Error("Error loading application: ", "err", err)
		http.Error(w, "Error parsing web", http.StatusInternalServerError)
		return
	}

	data := HomeData{
		Posts: p,
	}

	ts, err := template.ParseFS(h.templates, "templates/base.go.html", "templates/home.go.html")

	if err != nil {
		h.logger.Error("Error parsing templates", "err", err)
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.logger.Error("Error executing template", "err", err)
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

	var post = Post{}
	var isEdit = false
	var err error

	// get post path
	postPath := r.URL.Query().Get("path")

	// Load the existing post is we have postpath
	if postPath != "" {
		post, err = GetPost(postPath)
		if err != nil {
			h.logger.Error("Error loading post", "err", err)
			http.Error(w, "Error loading post", http.StatusInternalServerError)
			return
		}
		isEdit = true
	}

	data := EditorData{
		Post:   post,
		IsEdit: isEdit,
	}

	ts, err := template.New("base").Funcs(template.FuncMap{
		"join": join,
	}).ParseFS(h.templates, "templates/base.go.html", "templates/editor.go.html")

	if err != nil {
		h.logger.Error("Error parsing templates", "err", err)
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.logger.Error("Error executing template", "err", err)
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
	post := Post{
		Metadata: FrontMatter{
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

func join(sep string, s []string) string {
	return strings.Join(s, sep)
}
