package posts

import (
	"github.com/adrg/frontmatter"
	"os"
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

	post := Post{
		Title:    frontMatter.Title,
		FilePath: postPath,
		Date:     frontMatter.Date,
		Tags:     frontMatter.Tags,
		Content:  string(rest),
	}

	return post, nil
}
