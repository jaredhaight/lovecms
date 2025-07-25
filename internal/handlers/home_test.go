package handlers

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestHomeHandler_Get(t *testing.T) {
	// Create a temporary directory for testing
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

	// Create a test post
	testPost := `---
title: "Test Post"
date: "2023-01-01T00:00:00Z"
tags: ["test"]
---
# Test Content`

	postPath := filepath.Join(contentDir, "test.md")
	err = os.WriteFile(postPath, []byte(testPost), 0644)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

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
			// Setup config
			v := viper.New()
			v.Set("SitePath", tt.sitePath)

			// Setup logger
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			// Create handler
			handler := NewHomeHandler(v, logger)

			// Create request
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.Get(w, req)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("HomeHandler.Get() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// Check redirect if expected
			if tt.wantRedirect {
				location := w.Header().Get("Location")
				if location != tt.redirectTarget {
					t.Errorf("HomeHandler.Get() redirect location = %v, want %v", location, tt.redirectTarget)
				}
			}
		})
	}
}

func TestHomeHandler_Get_InvalidContentPath(t *testing.T) {
	// Setup config with invalid path
	v := viper.New()
	v.Set("SitePath", "/nonexistent/path")

	// Setup logger
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	// Create handler
	handler := NewHomeHandler(v, logger)

	// Create request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Execute request
	handler.Get(w, req)

	// Should return internal server error
	if w.Code != http.StatusInternalServerError {
		t.Errorf("HomeHandler.Get() with invalid path status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestNewHomeHandler(t *testing.T) {
	v := viper.New()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	handler := NewHomeHandler(v, logger)

	if handler == nil {
		t.Error("NewHomeHandler() returned nil")
	}

	if handler.config != v {
		t.Error("NewHomeHandler() config not set correctly")
	}

	if handler.logger != logger {
		t.Error("NewHomeHandler() logger not set correctly")
	}
}
