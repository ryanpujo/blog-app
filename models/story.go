// Package models contains the data structures used in our application.
package models

import (
	"fmt"
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
	ID          uint        `json:"id"`                               // Unique identifier for the story
	Title       string      `json:"title" binding:"required,max=255"` // Title of the story
	Content     string      `json:"content" binding:"required"`       // Content of the story
	AuthorID    uint        `json:"author_id"`                        // Unique identifier for the author
	Slug        string      `json:"slug" binding:"required,max=255"`  // URL-friendly version of the story title
	Excerpt     *string     `json:"excerpt,omitempty"`                // Short summary of the story
	Status      StoryStatus `json:"status" default:"1"`               // Status of the story
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
	_ StoryType = iota
	FlashFiction
	ShortStory
	Novelette
	Novella
)

// String returns the string representation of the StoryType.
func (st StoryType) String() string {
	return [...]string{"", "flash_fiction", "short_story", "novelette", "novella"}[st]
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

// IsValidWordCountForStoryType checks if the word count of a story falls within the typical range for its type.
// It takes two parameters: storyType and wordCount. storyType is the type of the story and wordCount is the number of words in the story.
// It returns nil if the word count is within the typical range for the given story type, and a StoryError otherwise.
func IsValidWordCountForStoryType(storyType StoryType, wordCount uint) error {
	switch storyType {
	case FlashFiction:
		if wordCount <= 100 || wordCount > 1000 {
			return &StoryError{storyType, wordCount, "word count for flash fiction should be between 100 and 1000"}
		}
	case ShortStory:
		if wordCount <= 1000 || wordCount > 7500 {
			return &StoryError{storyType, wordCount, "word count for short story should be between 1000 and 7500"}
		}
	case Novelette:
		if wordCount <= 7500 || wordCount > 20_000 {
			return &StoryError{storyType, wordCount, "word count for novelette should be between 7500 and 20,000"}
		}
	case Novella:
		if wordCount <= 20_000 || wordCount > 40_000 {
			return &StoryError{storyType, wordCount, "word count for novella should be between 20,000 and 40,000"}
		}
	default:
		return &StoryError{storyType, wordCount, "invalid story type"}
	}
	return nil
}

// StoryError represents an error that occurs during story operations.
type StoryError struct {
	StoryType StoryType
	WordCount uint
	Message   string
}

// Error implements the error interface.
func (e StoryError) Error() string {
	return fmt.Sprintf("story error: %s (story type: %s, word count: %d)", e.Message, e.StoryType, e.WordCount)
}

func (e StoryError) Is(target error) bool {
	err, ok := target.(StoryError)
	if !ok {
		return false
	}
	return e.StoryType == err.StoryType
}

// As implements the As method for the error interface.
func (e *StoryError) As(target interface{}) bool {
	t, ok := target.(*StoryError)
	if !ok {
		return false
	}
	*t = *e
	return true
}
