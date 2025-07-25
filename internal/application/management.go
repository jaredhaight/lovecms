package application

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/jaredhaight/lovecms/internal/types"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

func GetPost(postPath string) (types.Post, error) {
	// read file contents
	content, err := os.Open(postPath)

	if err != nil {
		return types.Post{}, err
	}

	// parse front matter
	fm := types.FrontMatter{}
	rest, err := frontmatter.Parse(content, &fm)
	if err != nil {

		return types.Post{}, err
	}

	// Markdown -> HTML
	var buf bytes.Buffer
	md := goldmark.New()

	err = md.Convert(rest, &buf)
	if err != nil {
		return types.Post{}, err
	}

	post := types.Post{
		Metadata: fm,
		FileName: filepath.Base(postPath),
		FilePath: postPath,
		Content:  buf.String(),
	}

	return post, nil
}

func GetPosts(directoryPath string) ([]types.Post, error) {
	if directoryPath == "" {
		return nil, errors.New("directoryPath is empty")
	}

	// get files from the content folder
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// create a post object for each item
	var post types.Post
	var posts []types.Post
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())

		if ext == ".md" || ext == ".markdown" {
			postPath := filepath.Join(directoryPath, entry.Name())
			post, err = GetPost(postPath)

			if err != nil {
				return nil, err
			}

			posts = append(posts, post)
		}
	}

	// sort application by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Metadata.Date > posts[j].Metadata.Date
	})
	return posts, nil
}

func UpdatePost(post types.Post) error {
	// get frontmatter
	meta, err := yaml.Marshal(post.Metadata)
	if err != nil {
		return err
	}

	// create file content
	sb := strings.Builder{}
	sb.WriteString("---\n")
	sb.WriteString(string(meta))
	sb.WriteString("---\n")
	sb.WriteString(post.Content)

	// write content to disk
	err = os.WriteFile(post.FilePath, []byte(sb.String()), 0644)
	if err != nil {
		return err
	}

	return nil
}

func CreatePost(contentPath string, post types.Post) error {
	// Generate filename from title if slug is empty
	var filename string
	if post.Metadata.Slug != "" {
		filename = post.Metadata.Slug + ".md"
	} else {
		// Simple slug generation from title
		slug := strings.ToLower(post.Metadata.Title)
		slug = strings.ReplaceAll(slug, " ", "-")
		// Remove special characters (basic implementation)
		slug = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
				return r
			}
			return -1
		}, slug)
		filename = slug + ".md"
		post.Metadata.Slug = slug
	}

	// Set the file path
	post.FilePath = filepath.Join(contentPath, filename)

	// Use UpdatePost to write the file
	return UpdatePost(post)
}
