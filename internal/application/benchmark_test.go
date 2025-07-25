package application

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Benchmark tests for performance-critical functions

func BenchmarkGetPosts(b *testing.B) {
	// Create a temporary directory with test posts
	tempDir, err := os.MkdirTemp("", "lovecms_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple test posts
	testPost := `---
title: "Benchmark Post %d"
date: "2023-01-01T00:00:00Z"
tags: ["benchmark", "test"]
---
# Benchmark Content %d

This is some test content for benchmarking the GetPosts function.
It includes multiple paragraphs to simulate real blog posts.

## Section 1
Lorem ipsum dolor sit amet, consectetur adipiscing elit.

## Section 2
Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`

	for i := 0; i < 10; i++ {
		content := fmt.Sprintf(testPost, i, i)
		filename := fmt.Sprintf("post%d.md", i)
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			b.Fatalf("Failed to create test post %d: %v", i, err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GetPosts(tempDir)
		if err != nil {
			b.Fatalf("GetPosts() error: %v", err)
		}
	}
}

func BenchmarkUpdatePost(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "lovecms_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	post := Post{
		FilePath: filepath.Join(tempDir, "benchmark.md"),
		Metadata: FrontMatter{
			Title: "Benchmark Post",
			Date:  "2023-01-01T00:00:00Z",
			Draft: false,
			Tags:  []string{"benchmark", "performance", "test"},
		},
		Content: "This is benchmark content for testing UpdatePost performance.",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := UpdatePost(post)
		if err != nil {
			b.Fatalf("UpdatePost() error: %v", err)
		}
	}
}

func BenchmarkCreatePost(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "lovecms_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		post := Post{
			Metadata: FrontMatter{
				Title: fmt.Sprintf("Benchmark Post %d", i),
				Date:  "2023-01-01T00:00:00Z",
				Draft: false,
				Tags:  []string{"benchmark", "test"},
			},
			Content: fmt.Sprintf("This is benchmark content %d", i),
		}

		err := CreatePost(tempDir, post)
		if err != nil {
			b.Fatalf("CreatePost() error: %v", err)
		}

		// Clean up to avoid too many files
		if post.FilePath != "" {
			os.Remove(post.FilePath)
		}
	}
}
