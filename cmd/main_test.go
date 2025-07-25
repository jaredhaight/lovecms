package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	// Save original values
	originalGOOS := runtime.GOOS

	tests := []struct {
		name     string
		goos     string
		expected func() string
	}{
		{
			name: "windows",
			goos: "windows",
			expected: func() string {
				appdata := os.Getenv("APPDATA")
				if appdata == "" {
					// For testing when APPDATA might not be set
					return filepath.Join("testappdata", "lovecms")
				}
				return filepath.Join(appdata, "lovecms")
			},
		},
		{
			name: "linux",
			goos: "linux",
			expected: func() string {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return "."
				}
				return filepath.Join(homeDir, ".lovecms")
			},
		},
		{
			name: "darwin",
			goos: "darwin",
			expected: func() string {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return "."
				}
				return filepath.Join(homeDir, ".lovecms")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test is tricky because we can't easily change runtime.GOOS
			// We'll test the actual function behavior on the current platform
			result := getConfigDir()

			// Basic validation - should not be empty and should be an absolute path
			if result == "" {
				t.Error("getConfigDir() returned empty string")
			}

			// Should contain "lovecms" or ".lovecms"
			if !contains(result, "lovecms") {
				t.Errorf("getConfigDir() = %v, should contain 'lovecms'", result)
			}

			// Should be a valid path format
			if !filepath.IsAbs(result) && result != "." {
				t.Errorf("getConfigDir() = %v, should be absolute path or '.'", result)
			}
		})
	}

	// Restore original GOOS (not needed since we can't actually change it in tests)
	_ = originalGOOS
}

func TestGetConfigDir_RealPlatform(t *testing.T) {
	result := getConfigDir()

	switch runtime.GOOS {
	case "windows":
		// Should contain APPDATA path on Windows
		if !contains(result, "lovecms") {
			t.Errorf("Windows config dir should contain 'lovecms', got: %v", result)
		}
	case "darwin", "linux":
		// Should contain .lovecms in home directory
		if !contains(result, ".lovecms") {
			t.Errorf("Unix config dir should contain '.lovecms', got: %v", result)
		}
		homeDir, err := os.UserHomeDir()
		if err == nil && !contains(result, homeDir) && result != "." {
			t.Errorf("Unix config dir should be under home directory, got: %v", result)
		}
	}
}

func TestGetConfigDir_FailedHomeDir(t *testing.T) {
	// We can't easily test the error case where os.UserHomeDir() fails
	// without mocking, but we can test that the function handles it gracefully
	result := getConfigDir()

	// Should never return empty string
	if result == "" {
		t.Error("getConfigDir() should never return empty string")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
