package posts

import (
	"bytes"
	"errors"
	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	"os"
	"path/filepath"
	"sort"
)

func GetPost(postPath string) (Post, error) {
	// read file contents
	content, err := os.Open(postPath)

	if err != nil {
		return Post{}, err
	}

	// parse front matter
	frontMatter := FrontMatter{}
	rest, err := frontmatter.Parse(content, &frontMatter)
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
		Title:    frontMatter.Title,
		FilePath: postPath,
		Date:     frontMatter.Date,
		Tags:     frontMatter.Tags,
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

	// sort posts by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})
	return posts, nil
}
