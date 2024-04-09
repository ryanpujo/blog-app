package models

import "time"

// User represents a registered user within the system.
type User struct {
	ID        uint      `json:"id"`         // Unique identifier for the user.
	FirstName string    `json:"first_name"` // User's first name.
	LastName  string    `json:"last_name"`  // User's last name.
	Username  string    `json:"username"`   // Unique username for login or display.
	Password  string    `json:"-"`          // User's password (not exposed in JSON responses).
	Email     string    `json:"email"`      // User's email address.
	CreatedAt time.Time `json:"created_at"` // Timestamp when the user record was created.
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the user record was last updated.
}

// UserPayload represents the data expected for creating or updating a user.
type UserPayload struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name" binding:"required,min=3"` // User's first name.
	LastName  string `json:"last_name" binding:"required,min=3"`  // User's last name.
	Username  string `json:"username" binding:"required,min=6"`   // Unique username for login.
	Password  string `json:"password" binding:"required,min=7"`   // User's password.
	Email     string `json:"email" binding:"required,email"`      // User's email address.
}
