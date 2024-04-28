package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/services"
)

type StoryController interface {
	Create(c *gin.Context)
}

type storyController struct {
	s services.StoryService
}

func NewStoryController(s services.StoryService) *storyController {
	return &storyController{
		s: s,
	}
}

func (b *storyController) Create(c *gin.Context) {

}
