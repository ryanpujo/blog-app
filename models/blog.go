package models

import (
	"time"
)

// BlogStatus is a custom type that defines the possible statuses of a blog post.
type BlogStatus string

// Define constants for BlogStatus
const (
	Draft     BlogStatus = "draft"
	Published BlogStatus = "published"
	Archived  BlogStatus = "archived"
)

// Blog represents the structure of our resource and includes validation tags for Gin binding.
type BlogPayload struct {
	ID          uint       `json:"id" binding:"required"`
	Title       string     `json:"title" binding:"required,max=255"`
	Content     string     `json:"content" binding:"required"`
	AuthorID    uint       `json:"author_id" binding:"required"`
	Slug        string     `json:"slug" binding:"required,max=255"`
	Excerpt     *string    `json:"excerpt,omitempty"`
	Status      BlogStatus `json:"status" binding:"required,oneof=draft published archived"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Type        StoryType  `json:"type" binding:"required"`
	WordCount   uint       `json:"word_count"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type StoryType int

const (
	FlashFiction StoryType = iota
	ShortStory
	Novelette
	Novella
)

func (st StoryType) String() string {
	return [4]string{"flash_fiction", "short_story", "novelette", "novella"}[st]
}

type Blog struct {
	ID          uint       `json:"id" binding:"required"`
	Title       string     `json:"title" binding:"required,max=255"`
	Content     string     `json:"content" binding:"required"`
	Author      User       `json:"author" binding:"required"`
	Slug        string     `json:"slug" binding:"required,max=255"`
	Excerpt     *string    `json:"excerpt,omitempty"`
	Status      BlogStatus `json:"status" binding:"required,oneof=draft published archived"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Type        string     `json:"type" binding:"required"`
	WordCount   uint       `json:"word_count" binding:"required"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
