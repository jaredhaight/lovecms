package application

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaredhaight/lovecms/internal/types"
)

func TestGetPosts(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		files       map[string]string
		want        int
		wantErr     bool
		errContains string
	}{
		{
			name: "valid markdown files",
			files: map[string]string{
				"post1.md": `---
title: "First Post"
date: "2023-01-02T00:00:00Z"
tags: ["tech", "golang"]
---
# First Post Content`,
				"post2.md": `---
title: "Second Post"
date: "2023-01-01T00:00:00Z"
tags: ["blog"]
---
# Second Post Content`,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "mixed file types",
			files: map[string]string{
				"post1.md": `---
title: "Markdown Post"
date: "2023-01-01T00:00:00Z"
---
Content`,
				"post2.txt": "Not a markdown file",
				"post3.markdown": `---
title: "Markdown File"
date: "2023-01-01T00:00:00Z"
---
Content`,
			},
			want:    2,
			wantErr: false,
		},
		{
			name:    "empty directory",
			files:   map[string]string{},
			want:    0,
			wantErr: false,
		},
		{
			name:        "nonexistent directory",
			files:       nil,
			want:        0,
			wantErr:     true,
			errContains: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testDir string
			if tt.files != nil {
				testDir = tempDir
				// Create test files
				for filename, content := range tt.files {
					filePath := filepath.Join(testDir, filename)
					err := os.WriteFile(filePath, []byte(content), 0644)
					if err != nil {
						t.Fatalf("Failed to create test file %s: %v", filename, err)
					}
				}
			} else {
				testDir = "/nonexistent/directory"
			}

			got, err := GetPosts(testDir)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPosts() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetPosts() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GetPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.want {
				t.Errorf("GetPosts() returned %d posts, want %d", len(got), tt.want)
			}

			// Verify posts are sorted by date (newest first)
			if len(got) > 1 {
				for i := 0; i < len(got)-1; i++ {
					if got[i].Metadata.Date < got[i+1].Metadata.Date {
						t.Errorf("Posts not sorted by date correctly: %s should come before %s",
							got[i].Metadata.Date, got[i+1].Metadata.Date)
					}
				}
			}

			// Clean up test files for next iteration
			if tt.files != nil {
				for filename := range tt.files {
					os.Remove(filepath.Join(testDir, filename))
				}
			}
		})
	}
}

func TestGetPosts_EmptyDirectory(t *testing.T) {
	got, err := GetPosts("")
	if err == nil {
		t.Errorf("GetPosts(\"\") should return error for empty directory path")
	}
	if got != nil {
		t.Errorf("GetPosts(\"\") should return nil posts for empty directory path")
	}
}

func TestUpdatePost(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		post    types.Post
		wantErr bool
	}{
		{
			name: "valid post",
			post: types.Post{
				FilePath: filepath.Join(tempDir, "test.md"),
				Metadata: types.FrontMatter{
					Title: "Test Post",
					Date:  "2023-01-01T00:00:00Z",
					Draft: false,
					Tags:  []string{"test", "golang"},
				},
				Content: "This is test content",
			},
			wantErr: false,
		},
		{
			name: "post with invalid file path",
			post: types.Post{
				FilePath: "/invalid/path/that/does/not/exist/test.md",
				Metadata: types.FrontMatter{
					Title: "Invalid Post",
				},
				Content: "Content",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdatePost(tt.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created and content is correct
				content, err := os.ReadFile(tt.post.FilePath)
				if err != nil {
					t.Errorf("Failed to read created file: %v", err)
					return
				}

				contentStr := string(content)
				if !strings.Contains(contentStr, "---") {
					t.Errorf("File should contain frontmatter delimiters")
				}
				if !strings.Contains(contentStr, tt.post.Metadata.Title) {
					t.Errorf("File should contain post title")
				}
				if !strings.Contains(contentStr, tt.post.Content) {
					t.Errorf("File should contain post content")
				}
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name         string
		contentPath  string
		post         types.Post
		wantErr      bool
		wantFilename string
	}{
		{
			name:        "post with custom slug",
			contentPath: tempDir,
			post: types.Post{
				Metadata: types.FrontMatter{
					Title: "Custom Slug Post",
					Slug:  "custom-slug",
				},
				Content: "Content with custom slug",
			},
			wantErr:      false,
			wantFilename: "custom-slug.md",
		},
		{
			name:        "post without slug - auto generate",
			contentPath: tempDir,
			post: types.Post{
				Metadata: types.FrontMatter{
					Title: "Auto Generated Slug Post",
				},
				Content: "Content with auto-generated slug",
			},
			wantErr:      false,
			wantFilename: "auto-generated-slug-post.md",
		},
		{
			name:        "post with special characters in title",
			contentPath: tempDir,
			post: types.Post{
				Metadata: types.FrontMatter{
					Title: "Post With Special Characters!@#$%",
				},
				Content: "Content with special chars",
			},
			wantErr:      false,
			wantFilename: "post-with-special-characters.md",
		},
		{
			name:        "invalid content path",
			contentPath: "/invalid/path",
			post: types.Post{
				Metadata: types.FrontMatter{
					Title: "Invalid Path Post",
				},
				Content: "Content",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreatePost(tt.contentPath, tt.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created with correct name
				expectedPath := filepath.Join(tt.contentPath, tt.wantFilename)
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not created", expectedPath)
				}

				// Clean up created file
				os.Remove(expectedPath)
			}
		})
	}
}

func TestCreatePost_SlugGeneration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		title    string
		expected string
	}{
		{"Simple Title", "simple-title"},
		{"Title With Numbers 123", "title-with-numbers-123"},
		{"UPPERCASE TITLE", "uppercase-title"},
		{"Title-With-Hyphens", "title-with-hyphens"},
		{"Title With Special @#$% Characters", "title-with-special--characters"}, // Special chars removed, leaving double dash
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			post := types.Post{
				Metadata: types.FrontMatter{
					Title: tt.title,
				},
				Content: "Test content",
			}

			err := CreatePost(tempDir, post)
			if err != nil && tt.expected != "" {
				t.Errorf("CreatePost() error = %v", err)
				return
			}

			if tt.expected != "" {
				expectedPath := filepath.Join(tempDir, tt.expected+".md")
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not created", expectedPath)
				} else {
					// Clean up
					os.Remove(expectedPath)
				}
			}
		})
	}
}
