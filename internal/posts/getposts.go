package posts

import (
	"errors"
	"github.com/adrg/frontmatter"
	"os"
	"path/filepath"
)

type FrontMatter struct {
	Title        string   `yaml:"title"`
	Date         string   `yaml:"date"`
	Draft        bool     `yaml:"draft"`
	LastModified string   `yaml:"lastmod"`
	PublishDate  string   `yaml:"publishDate"`
	Slug         string   `yaml:"slug"`
	Tags         []string `yaml:"tags"`
}

func GetPosts(directoryPath string) (map[string]Post, error) {
	if directoryPath == "" {
		return nil, errors.New("directoryPath is empty")
	}

	// get files from the content folder
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// create a post object for each item
	posts := make(map[string]Post)
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())

		if ext == ".md" || ext == ".markdown" {
			postPath := filepath.Join(directoryPath, entry.Name())
			// read file contents
			content, err := os.Open(postPath)

			if err != nil {
				return nil, err
			}

			// parse front matter
			frontMatter := FrontMatter{}
			rest, err := frontmatter.Parse(content, &frontMatter)
			if err != nil {
				return nil, err
			}

			// NOTE: If a front matter must be present in the input data, use
			//       frontmatter.MustParse instead.

			post := Post{
				Title:     frontMatter.Title,
				FilePath:  postPath,
				Published: frontMatter.PublishDate,
				Tags:      frontMatter.Tags,
				Content:   string(rest),
			}

			posts[entry.Name()] = post
		}
	}

	return posts, nil
}
