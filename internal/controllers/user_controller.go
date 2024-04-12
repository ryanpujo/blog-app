package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/ryanpujo/blog-app/internal/services"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
)

type UserController interface {
	Create(c *gin.Context)
	FindById(c *gin.Context)
	FindUsers(c *gin.Context)
	DeleteById(c *gin.Context)
	Update(c *gin.Context)
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

// FindById handles the request to find a user by their ID.
// It binds the URI parameters to a Uri struct, then retrieves the user by ID.
// If successful, it responds with the user data in JSON format.
func (uc *userController) FindById(c *gin.Context) {
	// Initialize a new Uri struct to store URI parameters.
	var uri models.Uri

	// Bind URI parameters to the struct. If there's an error, handle it and return.
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// Retrieve the user by ID using the service layer.
	user, err := uc.s.FindById(uri.ID)
	if err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// Respond with the user data in JSON format if retrieval is successful.
	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"user": user}))
}

// FindUsers handles the HTTP request for retrieving a list of users.
// It uses the userController's service to fetch the users and returns
// a JSON response with the users or an error message.
func (uc *userController) FindUsers(c *gin.Context) {
	// Attempt to find users through the service layer.
	users, err := uc.s.FindUsers()
	if err != nil {
		// If an error occurs, use the utility function to handle the error response.
		utils.HandleRequestError(c, err)
		return
	}

	// If no error occurs, respond with the list of users in a success response.
	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"users": users}))
}

// DeleteById handles the HTTP request to delete a user by ID.
// It binds the URI to a models.Uri struct, attempts to delete the user
// through the service layer, and responds with the appropriate status code.
func (uc *userController) DeleteById(c *gin.Context) {
	// Initialize a Uri struct to store the URI parameters.
	var uri models.Uri

	// Bind the URI parameters to the 'uri' struct.
	// If there's an error in binding, handle the error and return.
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// Attempt to delete the user by ID through the service layer.
	// If there's an error in deletion, handle the error and return.
	if err := uc.s.DeleteById(uri.ID); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// If the deletion is successful, respond with 'http.StatusOK'.
	c.Status(http.StatusOK)
}

// Update handles the user update request.
func (uc *userController) Update(c *gin.Context) {
	// Define the payload and URI variables to store the incoming data.
	var payload models.UserPayload
	var uri models.Uri

	// Bind the JSON body to the payload variable. If there's an error, handle it and return.
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// Bind the URI parameters to the uri variable. If there's an error, handle it and return.
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// Call the update service with the URI ID and payload. If there's an error, handle it and return.
	if err := uc.s.Update(uri.ID, &payload); err != nil {
		utils.HandleRequestError(c, err)
		return
	}

	// If the update is successful, set the status to OK.
	c.Status(http.StatusOK)
}
