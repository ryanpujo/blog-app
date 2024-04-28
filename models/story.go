// Package models contains the data structures used in our application.
package models

import (
	"time"
)

// StoryStatus represents the possible statuses of a story.
type StoryStatus int

// Constants for StoryStatus.
const (
	Draft StoryStatus = iota
	Published
	Archived
)

// String returns the string representation of the StoryStatus.
func (ss StoryStatus) String() string {
	return [...]string{"draft", "published", "archived"}[ss]
}

// StoryPayload represents the structure of a story resource and includes validation tags for Gin binding.
type StoryPayload struct {
	ID          uint        `json:"id" binding:"required"`            // Unique identifier for the story
	Title       string      `json:"title" binding:"required,max=255"` // Title of the story
	Content     string      `json:"content" binding:"required"`       // Content of the story
	AuthorID    uint        `json:"author_id" binding:"required"`     // Unique identifier for the author
	Slug        string      `json:"slug" binding:"required,max=255"`  // URL-friendly version of the story title
	Excerpt     *string     `json:"excerpt,omitempty"`                // Short summary of the story
	Status      StoryStatus `json:"status" binding:"required"`        // Status of the story
	PublishedAt *time.Time  `json:"published_at,omitempty"`           // Date and time when the story was published
	Type        StoryType   `json:"type" binding:"required"`          // Type of the story
	WordCount   uint        `json:"word_count"`                       // Word count of the story
	CreatedAt   time.Time   `json:"created_at,omitempty"`             // Date and time when the story was created
	UpdatedAt   *time.Time  `json:"updated_at,omitempty"`             // Date and time when the story was last updated
}

// StoryType represents the possible types of a story.
type StoryType int

// Constants for StoryType.
const (
	FlashFiction StoryType = iota
	ShortStory
	Novelette
	Novella
)

// String returns the string representation of the StoryType.
func (st StoryType) String() string {
	return [...]string{"flash_fiction", "short_story", "novelette", "novella"}[st]
}

// Story represents the structure of a story resource.
type Story struct {
	ID          uint       `json:"id" binding:"required"`                                                     // Unique identifier for the story
	Title       string     `json:"title" binding:"required,max=255"`                                          // Title of the story
	Content     string     `json:"content" binding:"required"`                                                // Content of the story
	Author      User       `json:"author" binding:"required"`                                                 // Author of the story
	Slug        string     `json:"slug" binding:"required,max=255"`                                           // URL-friendly version of the story title
	Excerpt     *string    `json:"excerpt,omitempty"`                                                         // Short summary of the story
	Status      string     `json:"status" binding:"required,oneof=draft published archived"`                  // Status of the story
	PublishedAt *time.Time `json:"published_at,omitempty"`                                                    // Date and time when the story was published
	Type        string     `json:"type" binding:"required,oneof=flash_fiction short_story novelette novella"` // Type of the story
	WordCount   uint       `json:"word_count" binding:"required"`                                             // Word count of the story
	CreatedAt   time.Time  `json:"created_at,omitempty"`                                                      // Date and time when the story was created
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`                                                      // Date and time when the story was last updated
}
