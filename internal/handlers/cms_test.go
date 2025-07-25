package handlers

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestPostHandler_GetNew(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name           string
		sitePath       string
		wantStatus     int
		wantRedirect   bool
		redirectTarget string
	}{
		{
			name:         "valid site path",
			sitePath:     tempDir,
			wantStatus:   http.StatusOK,
			wantRedirect: false,
		},
		{
			name:           "empty site path - should redirect",
			sitePath:       "",
			wantStatus:     http.StatusFound,
			wantRedirect:   true,
			redirectTarget: "/setup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.Set("SitePath", tt.sitePath)
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			handler := NewPostHandler(v, logger)
			req := httptest.NewRequest("GET", "/posts/new", nil)
			w := httptest.NewRecorder()

			handler.GetNew(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("PostHandler.GetNew() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantRedirect {
				location := w.Header().Get("Location")
				if location != tt.redirectTarget {
					t.Errorf("PostHandler.GetNew() redirect location = %v, want %v", location, tt.redirectTarget)
				}
			}
		})
	}
}

func TestPostHandler_PostNew(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create content directory
	contentDir := filepath.Join(tempDir, "content")
	err = os.MkdirAll(contentDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create content dir: %v", err)
	}

	tests := []struct {
		name         string
		sitePath     string
		formData     url.Values
		wantStatus   int
		wantRedirect bool
		wantFile     string
	}{
		{
			name:     "valid post creation",
			sitePath: tempDir,
			formData: url.Values{
				"title":   {"Test Post"},
				"content": {"This is test content"},
				"slug":    {"test-post"},
				"tags":    {"tech, golang"},
				"draft":   {""},
			},
			wantStatus:   http.StatusSeeOther,
			wantRedirect: true,
			wantFile:     "test-post.md",
		},
		{
			name:     "post with auto-generated slug",
			sitePath: tempDir,
			formData: url.Values{
				"title":   {"Auto Generated Post"},
				"content": {"Content for auto-generated post"},
				"tags":    {"auto"},
			},
			wantStatus:   http.StatusSeeOther,
			wantRedirect: true,
			wantFile:     "auto-generated-post.md",
		},
		{
			name:     "post marked as draft",
			sitePath: tempDir,
			formData: url.Values{
				"title":   {"Draft Post"},
				"content": {"Draft content"},
				"draft":   {"on"},
			},
			wantStatus:   http.StatusSeeOther,
			wantRedirect: true,
			wantFile:     "draft-post.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.Set("SitePath", tt.sitePath)
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			handler := NewPostHandler(v, logger)

			// Create form request
			req := httptest.NewRequest("POST", "/posts/new", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.PostNew(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("PostHandler.PostNew() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantRedirect {
				location := w.Header().Get("Location")
				if location != "/" {
					t.Errorf("PostHandler.PostNew() redirect location = %v, want %v", location, "/")
				}
			}

			// Check if file was created
			if tt.wantFile != "" {
				filePath := filepath.Join(contentDir, tt.wantFile)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not created", filePath)
				} else {
					// Verify file content
					content, err := os.ReadFile(filePath)
					if err != nil {
						t.Errorf("Failed to read created file: %v", err)
					} else {
						contentStr := string(content)
						if !strings.Contains(contentStr, tt.formData.Get("title")) {
							t.Errorf("File should contain post title")
						}
						if !strings.Contains(contentStr, tt.formData.Get("content")) {
							t.Errorf("File should contain post content")
						}
					}
					// Clean up
					os.Remove(filePath)
				}
			}
		})
	}
}

func TestPostHandler_PostNew_EmptySitePath(t *testing.T) {
	v := viper.New()
	v.Set("SitePath", "")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	handler := NewPostHandler(v, logger)

	formData := url.Values{
		"title":   {"Test Post"},
		"content": {"Content"},
	}

	req := httptest.NewRequest("POST", "/posts/new", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.PostNew(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("PostHandler.PostNew() with empty site path status = %v, want %v", w.Code, http.StatusFound)
	}

	location := w.Header().Get("Location")
	if location != "/setup" {
		t.Errorf("PostHandler.PostNew() redirect location = %v, want %v", location, "/setup")
	}
}

func TestPostHandler_Get(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test post
	testPost := `---
title: "Test Post"
date: "2023-01-01T00:00:00Z"
tags: ["test"]
---
# Test Content`

	postPath := filepath.Join(tempDir, "test.md")
	err = os.WriteFile(postPath, []byte(testPost), 0644)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name       string
		postPath   string
		wantStatus int
	}{
		{
			name:       "valid post path",
			postPath:   postPath,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing post path",
			postPath:   "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "nonexistent post",
			postPath:   "/nonexistent/post.md",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			handler := NewPostHandler(v, logger)

			req := httptest.NewRequest("GET", "/post?path="+tt.postPath, nil)
			w := httptest.NewRecorder()

			handler.Get(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("PostHandler.Get() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestPostHandler_PostNew_InvalidForm(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	v := viper.New()
	v.Set("SitePath", tempDir)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	handler := NewPostHandler(v, logger)

	// Create request with invalid form data (malformed)
	req := httptest.NewRequest("POST", "/posts/new", bytes.NewReader([]byte("invalid%form%data")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.PostNew(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("PostHandler.PostNew() with invalid form status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestNewPostHandler(t *testing.T) {
	v := viper.New()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	handler := NewPostHandler(v, logger)

	if handler == nil {
		t.Error("NewPostHandler() returned nil")
	}

	if handler.config != v {
		t.Error("NewPostHandler() config not set correctly")
	}

	if handler.logger != logger {
		t.Error("NewPostHandler() logger not set correctly")
	}
}

func TestPostHandler_PostEdit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lovecms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test post file
	testPost := `---
title: "Original Title"
date: "2023-01-01T00:00:00Z"
tags: ["original", "test"]
draft: false
slug: "original-post"
---
# Original Content`

	postPath := filepath.Join(tempDir, "test.md")
	err = os.WriteFile(postPath, []byte(testPost), 0644)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name         string
		postPath     string
		formData     url.Values
		wantStatus   int
		wantRedirect bool
	}{
		{
			name:     "valid edit",
			postPath: postPath,
			formData: url.Values{
				"title":   {"Updated Title"},
				"content": {"Updated content"},
				"slug":    {"updated-post"},
				"tags":    {"updated, test"},
				"date":    {"2023-01-02T00:00:00Z"},
			},
			wantStatus:   http.StatusSeeOther,
			wantRedirect: true,
		},
		{
			name:     "edit with draft status",
			postPath: postPath,
			formData: url.Values{
				"title":   {"Draft Title"},
				"content": {"Draft content"},
				"draft":   {"on"},
			},
			wantStatus:   http.StatusSeeOther,
			wantRedirect: true,
		},
		{
			name:     "missing post path",
			postPath: "",
			formData: url.Values{
				"title":   {"Test"},
				"content": {"Test"},
			},
			wantStatus:   http.StatusBadRequest,
			wantRedirect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			handler := NewPostHandler(v, logger)

			// Create form request
			var url string
			if tt.postPath != "" {
				url = "/posts/edit?path=" + tt.postPath
			} else {
				url = "/posts/edit"
			}

			req := httptest.NewRequest("POST", url, strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.PostEdit(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("PostHandler.PostEdit() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantRedirect {
				location := w.Header().Get("Location")
				expectedLocation := "/post?path=" + tt.postPath
				if location != expectedLocation {
					t.Errorf("PostHandler.PostEdit() redirect location = %v, want %v", location, expectedLocation)
				}
			}

			// Verify file was updated if successful
			if tt.wantStatus == http.StatusSeeOther && tt.postPath != "" {
				content, err := os.ReadFile(tt.postPath)
				if err != nil {
					t.Errorf("Failed to read updated file: %v", err)
				} else {
					contentStr := string(content)
					if !strings.Contains(contentStr, tt.formData.Get("title")) {
						t.Errorf("File should contain updated title")
					}
					if !strings.Contains(contentStr, tt.formData.Get("content")) {
						t.Errorf("File should contain updated content")
					}
				}
			}
		})
	}
}

func TestPostHandler_PostEdit_NonexistentPost(t *testing.T) {
	v := viper.New()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	handler := NewPostHandler(v, logger)

	formData := url.Values{
		"title":   {"Test"},
		"content": {"Test"},
	}

	req := httptest.NewRequest("POST", "/posts/edit?path=/nonexistent/post.md", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.PostEdit(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("PostHandler.PostEdit() with nonexistent post status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
