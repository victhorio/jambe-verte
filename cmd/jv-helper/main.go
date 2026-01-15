package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: jv-helper <post|page> <slug>")
		os.Exit(1)
	}

	contentType := os.Args[1]
	slug := os.Args[2]

	if contentType != "post" && contentType != "page" {
		fmt.Println("Error: content type must be 'post' or 'page'")
		os.Exit(1)
	}

	if err := createContent(contentType, slug); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s: %s\n", contentType, slug)
}

func createContent(contentType, slug string) error {
	var dir, filename, content string
	today := time.Now().Format("2006-01-02")

	switch contentType {
	case "post":
		dir = "content/posts"
		filename = fmt.Sprintf("%s-%s.md", today, slug)
		content = createPostContent(slug, today)
	case "page":
		dir = "content/pages"
		filename = fmt.Sprintf("%s.md", slug)
		content = createPageContent(slug, today)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	outputPath := filepath.Join(dir, filename)
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("file already exists: %s", outputPath)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	return nil
}

func createPostContent(slug, date string) string {
	return fmt.Sprintf(`---
title: "%s"
date: "%s"
tags: []
description: "Lorem ipsum dolor sit amet."
draft: true
---

# %s

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
`, slug, date, slug)
}

func createPageContent(slug, date string) string {
	return fmt.Sprintf(`---
title: "%s"
date: "%s"
---

# %s

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
`, slug, date, slug)
}
