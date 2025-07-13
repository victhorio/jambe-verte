package content

import (
	"html/template"
	"time"
)

type Post struct {
	Slug        string
	Title       string
	Date        time.Time
	Tags        []string
	Description string
	HTML        template.HTML
}

type PostFrontmatter struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
	Draft       bool     `yaml:"draft"`
}
