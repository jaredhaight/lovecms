package posts

type Post struct {
	Title    string
	FilePath string
	Date     string
	Tags     []string
	Content  string
}

type FrontMatter struct {
	Title        string   `yaml:"title"`
	Date         string   `yaml:"date"`
	Draft        bool     `yaml:"draft"`
	LastModified string   `yaml:"lastmod"`
	PublishDate  string   `yaml:"publishDate"`
	Slug         string   `yaml:"slug"`
	Tags         []string `yaml:"tags"`
}
