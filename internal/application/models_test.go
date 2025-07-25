package application

import (
	"reflect"
	"testing"
)

func TestPost_NewPost(t *testing.T) {
	tests := []struct {
		name     string
		metadata FrontMatter
		content  string
		want     Post
	}{
		{
			name: "valid post with all fields",
			metadata: FrontMatter{
				Title:        "Test Post",
				Date:         "2023-01-01T00:00:00Z",
				Draft:        false,
				LastModified: "2023-01-01T00:00:00Z",
				PublishDate:  "2023-01-01T00:00:00Z",
				Slug:         "test-post",
				Tags:         []string{"tech", "golang"},
			},
			content: "This is test content",
			want: Post{
				FileName: "",
				FilePath: "",
				Metadata: FrontMatter{
					Title:        "Test Post",
					Date:         "2023-01-01T00:00:00Z",
					Draft:        false,
					LastModified: "2023-01-01T00:00:00Z",
					PublishDate:  "2023-01-01T00:00:00Z",
					Slug:         "test-post",
					Tags:         []string{"tech", "golang"},
				},
				Content: "This is test content",
			},
		},
		{
			name: "post with empty tags",
			metadata: FrontMatter{
				Title: "Another Test",
				Tags:  []string{},
			},
			content: "Content without tags",
			want: Post{
				Metadata: FrontMatter{
					Title: "Another Test",
					Tags:  []string{},
				},
				Content: "Content without tags",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Post{
				Metadata: tt.metadata,
				Content:  tt.content,
			}
			if !reflect.DeepEqual(got.Metadata, tt.want.Metadata) {
				t.Errorf("Post metadata = %v, want %v", got.Metadata, tt.want.Metadata)
			}
			if got.Content != tt.want.Content {
				t.Errorf("Post content = %v, want %v", got.Content, tt.want.Content)
			}
		})
	}
}

func TestFrontMatter_ValidateFields(t *testing.T) {
	tests := []struct {
		name     string
		fm       FrontMatter
		hasTitle bool
		hasTags  bool
	}{
		{
			name: "complete frontmatter",
			fm: FrontMatter{
				Title:        "Complete Post",
				Date:         "2023-01-01",
				Draft:        false,
				LastModified: "2023-01-01",
				PublishDate:  "2023-01-01",
				Slug:         "complete-post",
				Tags:         []string{"complete", "test"},
			},
			hasTitle: true,
			hasTags:  true,
		},
		{
			name: "minimal frontmatter",
			fm: FrontMatter{
				Title: "Minimal Post",
			},
			hasTitle: true,
			hasTags:  false,
		},
		{
			name:     "empty frontmatter",
			fm:       FrontMatter{},
			hasTitle: false,
			hasTags:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.fm.Title != "") != tt.hasTitle {
				t.Errorf("FrontMatter title presence = %v, want %v", tt.fm.Title != "", tt.hasTitle)
			}
			if (len(tt.fm.Tags) > 0) != tt.hasTags {
				t.Errorf("FrontMatter tags presence = %v, want %v", len(tt.fm.Tags) > 0, tt.hasTags)
			}
		})
	}
}
