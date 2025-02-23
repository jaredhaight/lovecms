package application

import (
	"errors"
	"github.com/jaredhaight/lovecms/internal/models"
	"os"
	"path"
)

func LoadPosts(directoryPath string) (map[string]models.Post, error) {
	if directoryPath == "" {
		return nil, errors.New("directoryPath is empty")
	}

	// get files from the content folder
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// create a post object for each item
	posts := make(map[string]models.Post)
	for _, entry := range entries {
		ext := path.Ext(entry.Name())

		if ext == ".md" || ext == ".markdown" {
			post := models.Post{
				Title:    entry.Name(),
				FilePath: path.Join(directoryPath, entry.Name()),
				Tags:     nil,
			}

			posts[entry.Name()] = post
		}
	}

	return posts, nil
}
