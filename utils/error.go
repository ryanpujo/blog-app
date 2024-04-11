package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ryanpujo/blog-app/internal/response"
)

// GetValidationErrorMessage generates a user-friendly error message based on the validation errors.
func GetValidationErrorMessage(vErr validator.ValidationErrors) string {
	// Default error message
	errMessage := "Validation failed"

	// Handle different validation errors
	switch vErr[0].ActualTag() {
	case "min":
		errMessage = fmt.Sprintf("The %s field must be at least %s characters", vErr[0].Field(), vErr[0].Param())
	case "email":
		errMessage = fmt.Sprintf("The %s field must be a valid email address", vErr[0].Field())
	case "gt":
		errMessage = fmt.Sprintf("The %s field must be grater than %s", vErr[0].Field(), vErr[0].Param())
	}

	return errMessage
}

// HandleRequestError handles errors and sends an appropriate JSON response to the client.
func HandleRequestError(c *gin.Context, err error) {

	// Determine the type of error and respond accordingly
	var validationErrs validator.ValidationErrors
	var DBerr DBError
	if errors.As(err, &validationErrs) {
		// Handle validation errors
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewErrorResponse(GetValidationErrorMessage(validationErrs)))
	} else if errors.As(err, &DBerr) {
		// Handle database errors
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewErrorResponse(DBerr.Message))
	} else if errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatusJSON(http.StatusNotFound, response.NewErrorResponse("data not found"))
	} else {
		// Handle other types of errors
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewErrorResponse("An unexpected error occurred"))
	}
}
