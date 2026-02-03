package article

import (
	"errors"
	"strings"
	"time"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
)

var (
	ErrInvalidTitle   = errors.New("title is required")
	ErrInvalidContent = errors.New("markdown content is required")
	ErrInvalidHTML    = errors.New("html content is required")
)

type Article struct {
	ID              int64
	Title           string
	Slug            string
	MarkdownContent string
	HTMLContent     string
	Status          Status
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewDraft(title, markdown string, now time.Time) (*Article, error) {
	cleanTitle := strings.TrimSpace(title)
	if cleanTitle == "" {
		return nil, ErrInvalidTitle
	}
	cleanMarkdown := strings.TrimSpace(markdown)
	if cleanMarkdown == "" {
		return nil, ErrInvalidContent
	}

	slug := Slugify(cleanTitle)
	return &Article{
		Title:           cleanTitle,
		Slug:            slug,
		MarkdownContent: cleanMarkdown,
		Status:          StatusDraft,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (a *Article) Update(title, markdown string, now time.Time) error {
	cleanTitle := strings.TrimSpace(title)
	if cleanTitle == "" {
		return ErrInvalidTitle
	}
	cleanMarkdown := strings.TrimSpace(markdown)
	if cleanMarkdown == "" {
		return ErrInvalidContent
	}

	a.Title = cleanTitle
	a.Slug = Slugify(cleanTitle)
	a.MarkdownContent = cleanMarkdown
	a.UpdatedAt = now
	return nil
}

func (a *Article) Publish(html string, now time.Time) error {
	cleanHTML := strings.TrimSpace(html)
	if cleanHTML == "" {
		return ErrInvalidHTML
	}
	a.HTMLContent = cleanHTML
	a.Status = StatusPublished
	a.UpdatedAt = now
	return nil
}

func (a *Article) MarkDraft(now time.Time) {
	a.Status = StatusDraft
	a.HTMLContent = ""
	a.UpdatedAt = now
}
