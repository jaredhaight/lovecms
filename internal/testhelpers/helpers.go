// Package testhelpers provides common utilities for testing LoveCMS
package testhelpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaredhaight/lovecms/internal/application"
)

// CreateTempSite creates a temporary directory structure that mimics a LoveCMS site
func CreateTempSite(t *testing.T) (siteDir string, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "lovecms_test_site")
	if err != nil {
		t.Fatalf("Failed to create temp site dir: %v", err)
	}

	// Create content directory
	contentDir := filepath.Join(tempDir, "content")
	err = os.MkdirAll(contentDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create content dir: %v", err)
		os.RemoveAll(tempDir)
	}

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// CreateTestPost creates a test post file in the specified directory
func CreateTestPost(t *testing.T, contentDir, filename string, post application.Post) string {
	t.Helper()

	// Set the file path
	post.FilePath = filepath.Join(contentDir, filename)

	// Create the post
	err := application.UpdatePost(post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	return post.FilePath
}

// CreateSamplePosts creates a set of sample posts for testing
func CreateSamplePosts(t *testing.T, contentDir string) []string {
	t.Helper()

	posts := []application.Post{
		{
			Metadata: application.FrontMatter{
				Title:       "First Test Post",
				Date:        "2023-01-03T00:00:00Z",
				Draft:       false,
				PublishDate: "2023-01-03T00:00:00Z",
				Slug:        "first-test-post",
				Tags:        []string{"test", "first"},
			},
			Content: "# First Test Post\n\nThis is the content of the first test post.",
		},
		{
			Metadata: application.FrontMatter{
				Title:       "Second Test Post",
				Date:        "2023-01-02T00:00:00Z",
				Draft:       false,
				PublishDate: "2023-01-02T00:00:00Z",
				Slug:        "second-test-post",
				Tags:        []string{"test", "second"},
			},
			Content: "# Second Test Post\n\nThis is the content of the second test post.",
		},
		{
			Metadata: application.FrontMatter{
				Title:       "Draft Post",
				Date:        "2023-01-01T00:00:00Z",
				Draft:       true,
				PublishDate: "2023-01-01T00:00:00Z",
				Slug:        "draft-post",
				Tags:        []string{"test", "draft"},
			},
			Content: "# Draft Post\n\nThis is a draft post.",
		},
	}

	var filePaths []string
	for i, post := range posts {
		filename := post.Metadata.Slug + ".md"
		filePath := CreateTestPost(t, contentDir, filename, post)
		filePaths = append(filePaths, filePath)
		_ = i // Avoid unused variable
	}

	return filePaths
}

// AssertPostEqual compares two posts for equality
func AssertPostEqual(t *testing.T, got, want application.Post) {
	t.Helper()

	if got.Metadata.Title != want.Metadata.Title {
		t.Errorf("Post title = %v, want %v", got.Metadata.Title, want.Metadata.Title)
	}

	if got.Metadata.Date != want.Metadata.Date {
		t.Errorf("Post date = %v, want %v", got.Metadata.Date, want.Metadata.Date)
	}

	if got.Metadata.Draft != want.Metadata.Draft {
		t.Errorf("Post draft = %v, want %v", got.Metadata.Draft, want.Metadata.Draft)
	}

	if got.Metadata.Slug != want.Metadata.Slug {
		t.Errorf("Post slug = %v, want %v", got.Metadata.Slug, want.Metadata.Slug)
	}

	if len(got.Metadata.Tags) != len(want.Metadata.Tags) {
		t.Errorf("Post tags length = %v, want %v", len(got.Metadata.Tags), len(want.Metadata.Tags))
	} else {
		for i, tag := range got.Metadata.Tags {
			if tag != want.Metadata.Tags[i] {
				t.Errorf("Post tag[%d] = %v, want %v", i, tag, want.Metadata.Tags[i])
			}
		}
	}

	if got.Content != want.Content {
		t.Errorf("Post content = %v, want %v", got.Content, want.Content)
	}
}

// AssertFileExists checks if a file exists at the given path
func AssertFileExists(t *testing.T, filePath string) {
	t.Helper()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s does not exist", filePath)
	}
}

// AssertFileNotExists checks if a file does not exist at the given path
func AssertFileNotExists(t *testing.T, filePath string) {
	t.Helper()

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("Expected file %s should not exist", filePath)
	}
}

// AssertFileContains checks if a file contains the specified content
func AssertFileContains(t *testing.T, filePath, content string) {
	t.Helper()

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read file %s: %v", filePath, err)
		return
	}

	fileStr := string(fileContent)
	if !containsString(fileStr, content) {
		t.Errorf("File %s should contain %q", filePath, content)
	}
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// indexOf returns the index of the first occurrence of substr in s, or -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
