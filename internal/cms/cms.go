package cms

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

type Cms struct {
	config    *viper.Viper
	logger    *slog.Logger
	templates embed.FS
	tags      map[string][]Post
}

type HomeData struct {
	Posts []Post
	Tags  []string
}

type EditorData struct {
	Post   Post
	Tags   []string
	IsEdit bool
}

func New(v *viper.Viper, l *slog.Logger, t embed.FS) *Cms {
	return &Cms{
		config:    v,
		logger:    l,
		templates: t,
		tags:      make(map[string][]Post),
	}
}

func (c *Cms) getTags() []string {
	var tags = make([]string, 0)
	for k := range c.tags {
		tags = append(tags, k)
	}
	return tags

}

func (c *Cms) updateTags(p Post) {
	for _, t := range p.Metadata.Tags {
		val, ok := c.tags[t]
		if ok {
			c.tags[t] = append(val, p)
		} else {
			c.tags[t] = make([]Post, 0)
			c.tags[t] = append(c.tags[t], p)
		}
	}
}

// GET /
func (c *Cms) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// get sitepath
	var sitePath = c.config.GetString("SitePath")

	// if we don't have a site defined, redirect to setup
	if sitePath == "" {
		c.logger.Info("No current site defined. Redirecting to setup")
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}

	contentPath := filepath.Join(sitePath, "content")
	// Load our application
	posts, err := getPosts(contentPath)
	if err != nil {
		c.logger.Error("Error loading application: ", "err", err)
		http.Error(w, "Error parsing web", http.StatusInternalServerError)
		return
	}

	// Build our tags repo
	for _, p := range posts {
		c.updateTags(p)
	}

	data := HomeData{
		Posts: posts,
		Tags:  c.getTags(),
	}

	ts, err := template.ParseFS(c.templates, "templates/base.go.html", "templates/home.go.html")

	if err != nil {
		c.logger.Error("Error parsing templates", "err", err)
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		c.logger.Error("Error executing template", "err", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// GET /editor/?path=Foo
func (c *Cms) EditorHandler(w http.ResponseWriter, r *http.Request) {
	// get sitepath
	var sitePath = c.config.GetString("SitePath")

	// if we don't have a site defined, redirect to setup
	if sitePath == "" {
		c.logger.Info("No current site defined. Redirecting to setup")
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
		post, err = getPost(postPath)
		if err != nil {
			c.logger.Error("Error loading post", "err", err)
			http.Error(w, "Error loading post", http.StatusInternalServerError)
			return
		}
		isEdit = true
	}

	data := EditorData{
		Post:   post,
		Tags:   c.getTags(),
		IsEdit: isEdit,
	}

	ts, err := template.New("base").Funcs(template.FuncMap{
		"join": join,
	}).ParseFS(c.templates, "templates/base.go.html", "templates/editor.go.html")

	if err != nil {
		c.logger.Error("Error parsing templates", "err", err)
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		c.logger.Error("Error executing template", "err", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// POST /editor/
func (c *Cms) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// get sitepath
	var sitePath = c.config.GetString("SitePath")

	// if we don't have a site defined, redirect to setup
	if sitePath == "" {
		c.logger.Info("No current site defined. Redirecting to setup")
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		c.logger.Error("Error parsing form", "err", err)
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
	err = createPost(contentPath, post)
	if err != nil {
		c.logger.Error("Error creating post", "err", err)
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func join(sep string, s []string) string {
	return strings.Join(s, sep)
}
