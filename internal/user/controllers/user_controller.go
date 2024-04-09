package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/ryanpujo/blog-app/internal/user/services"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
)

type UserController interface {
	Create(x *gin.Context)
}

type userController struct {
	s services.UserService
}

func NewUserController(s services.UserService) *userController {
	return &userController{
		s: s,
	}
}

// Create handles the creation of a new user. It binds the incoming JSON to a UserPayload struct,
// calls the service to create a user, and returns the result in a JSON response.
func (uc *userController) Create(c *gin.Context) {
	// Initialize a new instance of UserPayload.
	var payload models.UserPayload

	// Bind the incoming JSON to the payload. If there's an error, handle it and return.
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// Call the service to create a new user with the provided payload.
	id, err := uc.s.Create(payload)
	if err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// If the user is created successfully, return the user ID in a success response.
	c.JSON(http.StatusCreated, response.NewSuccessResponse(gin.H{"id": id}))
}
