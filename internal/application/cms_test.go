package application

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

// Helper function to create a mock viper config
func createMockConfig(sitePath string) *viper.Viper {
	v := viper.New()
	v.Set("SitePath", sitePath)
	return v
}

// Helper function to create a mock logger that discards output
func createMockLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Only log errors to reduce test noise
	}))
}

// Helper function to create test content directory with posts only
func createTestContentDir(t testing.TB) (string, func()) {
	tempDir, err := os.MkdirTemp("", "lovecms_cms_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	contentDir := filepath.Join(tempDir, "content")
	err = os.MkdirAll(contentDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create content dir: %v", err)
	}

	// Create a test post
	testPost := `---
title: "Test Post"
date: "2024-01-01T00:00:00Z"
tags: ["test", "golang"]
slug: "test-post"
draft: false
---
# Test Post Content

This is a test post for unit testing.`

	err = os.WriteFile(filepath.Join(contentDir, "test-post.md"), []byte(testPost), 0644)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestNewCmsHandler(t *testing.T) {
	config := createMockConfig("/test/site")
	logger := createMockLogger()

	handler := NewCmsHandler(config, logger)

	if handler == nil {
		t.Fatal("NewCmsHandler returned nil")
	}

	if handler.config != config {
		t.Error("Config not set correctly")
	}

	if handler.logger != logger {
		t.Error("Logger not set correctly")
	}

	// Templates might be nil if template files don't exist, which is okay
}

func TestCmsHandler_GetHome(t *testing.T) {
	tests := []struct {
		name             string
		sitePath         string
		setupContent     bool
		expectedStatus   int
		expectedRedirect string
		expectContent    string
	}{
		{
			name:             "no site path configured",
			sitePath:         "",
			setupContent:     false,
			expectedStatus:   http.StatusFound,
			expectedRedirect: "/setup",
		},
		{
			name:         "valid site with content - no templates",
			sitePath:     "", // Will be set dynamically in test
			setupContent: true,
			// Expect 500 if templates don't exist
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config *viper.Viper
			var cleanup func()

			if tt.setupContent {
				sitePath, cleanupFunc := createTestContentDir(t)
				cleanup = cleanupFunc
				config = createMockConfig(sitePath)
			} else {
				config = createMockConfig(tt.sitePath)
			}

			if cleanup != nil {
				defer cleanup()
			}

			logger := createMockLogger()
			handler := NewCmsHandler(config, logger)

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handler.GetHome(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedRedirect != "" {
				location := w.Header().Get("Location")
				if location != tt.expectedRedirect {
					t.Errorf("Expected redirect to %s, got %s", tt.expectedRedirect, location)
				}
			}

			if tt.expectContent != "" && w.Code == http.StatusOK {
				body := w.Body.String()
				if !strings.Contains(body, tt.expectContent) {
					t.Errorf("Expected body to contain %s, got: %s", tt.expectContent, body)
				}
			}
		})
	}
}

func TestCmsHandler_GetEditor(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		queryParams      string
		sitePath         string
		setupContent     bool
		expectedStatus   int
		expectedRedirect string
		expectContent    string
	}{
		{
			name:             "no site path configured",
			path:             "/editor/",
			sitePath:         "",
			setupContent:     false,
			expectedStatus:   http.StatusFound,
			expectedRedirect: "/setup",
		},
		{
			name:         "new post form - no templates",
			path:         "/post/new",
			sitePath:     "/test/site",
			setupContent: false,
			// Expect 500 if templates don't exist
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "edit post without path parameter",
			path:           "/post/edit",
			sitePath:       "/test/site",
			setupContent:   false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config *viper.Viper
			var cleanup func()

			if tt.setupContent {
				sitePath, cleanupFunc := createTestContentDir(t)
				cleanup = cleanupFunc
				config = createMockConfig(sitePath)
			} else {
				config = createMockConfig(tt.sitePath)
			}

			if cleanup != nil {
				defer cleanup()
			}

			logger := createMockLogger()
			handler := NewCmsHandler(config, logger)

			url := tt.path
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.GetEditor(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedRedirect != "" {
				location := w.Header().Get("Location")
				if location != tt.expectedRedirect {
					t.Errorf("Expected redirect to %s, got %s", tt.expectedRedirect, location)
				}
			}

			if tt.expectContent != "" && w.Code == http.StatusOK {
				body := w.Body.String()
				if !strings.Contains(body, tt.expectContent) {
					t.Errorf("Expected body to contain %s, got: %s", tt.expectContent, body)
				}
			}
		})
	}

	// Test edit existing post separately since it needs dynamic path
	t.Run("edit existing post", func(t *testing.T) {
		sitePath, cleanup := createTestContentDir(t)
		defer cleanup()

		config := createMockConfig(sitePath)
		logger := createMockLogger()
		handler := NewCmsHandler(config, logger)

		// Use the full path to the test file
		testFilePath := filepath.Join(sitePath, "content", "test-post.md")
		url := "/post/edit?path=" + testFilePath

		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		handler.GetEditor(w, req)

		// Expect 500 if templates don't exist
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestCmsHandler_PostEditor(t *testing.T) {
	tests := []struct {
		name             string
		sitePath         string
		setupContent     bool
		formData         url.Values
		expectedStatus   int
		expectedRedirect string
	}{
		{
			name:             "no site path configured",
			sitePath:         "",
			setupContent:     false,
			formData:         url.Values{},
			expectedStatus:   http.StatusFound,
			expectedRedirect: "/setup",
		},
		{
			name:         "valid post creation",
			sitePath:     "", // Will be set dynamically
			setupContent: true,
			formData: url.Values{
				"title":   []string{"New Test Post"},
				"content": []string{"This is test content"},
				"slug":    []string{"new-test-post"},
				"tags":    []string{"test, new"},
				"draft":   []string{"off"},
			},
			expectedStatus:   http.StatusSeeOther,
			expectedRedirect: "/",
		},
		{
			name:         "post creation with draft",
			sitePath:     "", // Will be set dynamically
			setupContent: true,
			formData: url.Values{
				"title":   []string{"Draft Post"},
				"content": []string{"Draft content"},
				"slug":    []string{"draft-post"},
				"tags":    []string{"draft"},
				"draft":   []string{"on"},
			},
			expectedStatus:   http.StatusSeeOther,
			expectedRedirect: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config *viper.Viper
			var cleanup func()

			if tt.setupContent {
				sitePath, cleanupFunc := createTestContentDir(t)
				cleanup = cleanupFunc
				config = createMockConfig(sitePath)
			} else {
				config = createMockConfig(tt.sitePath)
			}

			if cleanup != nil {
				defer cleanup()
			}

			logger := createMockLogger()
			handler := NewCmsHandler(config, logger)

			// Create form data
			formData := tt.formData.Encode()
			req := httptest.NewRequest("POST", "/editor/", strings.NewReader(formData))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.PostEditor(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedRedirect != "" {
				location := w.Header().Get("Location")
				if location != tt.expectedRedirect {
					t.Errorf("Expected redirect to %s, got %s", tt.expectedRedirect, location)
				}
			}

			// If we successfully created a post, verify it exists
			if tt.expectedStatus == http.StatusSeeOther && tt.setupContent {
				sitePath := config.GetString("SitePath")
				contentPath := filepath.Join(sitePath, "content")
				slug := tt.formData.Get("slug")
				if slug != "" {
					expectedFile := filepath.Join(contentPath, slug+".md")
					if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
						t.Errorf("Expected post file %s to be created", expectedFile)
					}
				}
			}
		})
	}
}

func TestCmsHandler_PostEditor_InvalidForm(t *testing.T) {
	sitePath, cleanup := createTestContentDir(t)
	defer cleanup()

	config := createMockConfig(sitePath)
	logger := createMockLogger()
	handler := NewCmsHandler(config, logger)

	// Create a request with malformed form data
	req := httptest.NewRequest("POST", "/editor/", bytes.NewReader([]byte("invalid%form%data")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.PostEditor(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid form data, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCmsHandler_GetEditor_NonexistentPost(t *testing.T) {
	sitePath, cleanup := createTestContentDir(t)
	defer cleanup()

	config := createMockConfig(sitePath)
	logger := createMockLogger()
	handler := NewCmsHandler(config, logger)

	req := httptest.NewRequest("GET", "/post/edit?path=nonexistent.md", nil)
	w := httptest.NewRecorder()

	handler.GetEditor(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d for nonexistent post, got %d", http.StatusInternalServerError, w.Code)
	}
}

// Benchmark tests
func BenchmarkCmsHandler_GetHome(b *testing.B) {
	sitePath, cleanup := createTestContentDir(b)
	defer cleanup()

	config := createMockConfig(sitePath)
	logger := createMockLogger()
	handler := NewCmsHandler(config, logger)

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.GetHome(w, req)
	}
}

func BenchmarkCmsHandler_GetEditor_NewPost(b *testing.B) {
	config := createMockConfig("/test/site")
	logger := createMockLogger()
	handler := NewCmsHandler(config, logger)

	req := httptest.NewRequest("GET", "/post/new", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.GetEditor(w, req)
	}
}
