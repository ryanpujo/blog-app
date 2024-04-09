package response

// Response is a generic response structure for API responses.
type Response struct {
	Success bool        `json:"success"` // Indicates if the request was successful.
	Message string      `json:"message"` // Contains a message for the user.
	Data    interface{} `json:"data"`    // Holds the data to be sent to the user.
}

// NewSuccessResponse creates a new success response with the provided data.
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Message: "Operation successful.",
		Data:    data,
	}
}

// NewErrorResponse creates a new error response with the provided message.
func NewErrorResponse(message string) *Response {
	return &Response{
		Success: false,
		Message: message,
		Data:    nil,
	}
}
