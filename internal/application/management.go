package application

import (
	"bytes"
	"errors"
	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetPost(postPath string) (Post, error) {
	// read file contents
	content, err := os.Open(postPath)

	if err != nil {
		return Post{}, err
	}

	// parse front matter
	fm := FrontMatter{}
	rest, err := frontmatter.Parse(content, &fm)
	if err != nil {

		return Post{}, err
	}

	// Markdown -> HTML
	var buf bytes.Buffer
	md := goldmark.New()

	err = md.Convert(rest, &buf)
	if err != nil {
		return Post{}, err
	}

	post := Post{
		Metadata: fm,
		FileName: filepath.Base(postPath),
		FilePath: postPath,
		Content:  buf.String(),
	}

	return post, nil
}

func GetPosts(directoryPath string) ([]Post, error) {
	if directoryPath == "" {
		return nil, errors.New("directoryPath is empty")
	}

	// get files from the content folder
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// create a post object for each item
	var post Post
	var posts []Post
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())

		if ext == ".md" || ext == ".markdown" {
			postPath := filepath.Join(directoryPath, entry.Name())
			post, err = GetPost(postPath)
			posts = append(posts, post)
		}
	}

	// sort application by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Metadata.Date > posts[j].Metadata.Date
	})
	return posts, nil
}

func UpdatePost(post Post) error {
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
