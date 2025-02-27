package posts

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
)

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
