package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/ryanpujo/blog-app/internal/services"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
)

// StoryController defines the interface for story related operations
type StoryController interface {
	Create(c *gin.Context)
	FindById(c *gin.Context)
	FindStories(c *gin.Context)
	Update(c *gin.Context)
	DeleteById(c *gin.Context)
}

// storyController implements the StoryController interface
type storyController struct {
	service services.StoryService
}

// NewStoryController creates a new instance of storyController
func NewStoryController(s services.StoryService) *storyController {
	return &storyController{
		service: s,
	}
}

// Create implements the Create method of the StoryController interface
func (s *storyController) Create(c *gin.Context) {
	// Define a variable to hold the story payload
	var payload models.StoryPayload
	var uri models.Uri

	// Attempt to bind the request body to the payload struct
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		utils.HandleRequestError(c, err)
		return
	}
	payload.AuthorID = uri.ID
	// Call the service layer to create the story
	id, err := s.service.Create(payload)
	if err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// Build a success response with the created story ID
	response := response.NewSuccessResponse(gin.H{"id": id})
	c.JSON(http.StatusCreated, response)
}

func (s *storyController) FindById(c *gin.Context) {
	var uri models.StoryUri

	if err := c.ShouldBindUri(&uri); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	story, err := s.service.FindById(uri.StoryID)
	if err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"story": story}))
}

func (s *storyController) FindStories(c *gin.Context) {
	stories, err := s.service.FindStories()
	if err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"stories": stories}))
}

func (s *storyController) Update(c *gin.Context) {
	var uri models.Uri
	var storyUri models.StoryUri
	var payload models.StoryPayload

	if err := c.ShouldBindUri(&uri); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	if err := c.ShouldBindUri(&storyUri); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.HandleRequestError(c, err)
		return
	}
	payload.AuthorID = uri.ID
	err := s.service.Update(storyUri.StoryID, payload)
	if err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *storyController) DeleteById(c *gin.Context) {
	var uri models.StoryUri
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	if err := s.service.DeleteById(uri.StoryID); err != nil {
		utils.HandleRequestError(c, err)
	}

	c.Status(http.StatusOK)
}
