package models

// Uri represents the structure for URI parameters.
type Uri struct {
	ID uint `uri:"id" binding:"gt=0"` // ID must be greater than 0.
}

type StoryUri struct {
	StoryID uint `uri:"storyID" binding:"gt=0"`
}
