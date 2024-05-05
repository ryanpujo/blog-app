package repositories

import (
	"context"
	"log"
	"time"

	"github.com/ryanpujo/blog-app/database"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
)

// userRepository implements the UserRepository interface for operations on the users table.
type userRepository struct {
	db database.DatabaseOperations
}

// UserRepository defines the interface for user repository operations.
type UserRepository interface {
	Create(payload models.UserPayload) (*uint, error)
	FindById(id uint) (*models.User, error)
	FindUsers() ([]*models.User, error)
	DeleteById(id uint) error
	Update(id uint, user *models.UserPayload) error
	CheckIfEmailOrUsernameExist(email, username string) bool
}

// NewUserRepository creates a new instance of a userRepository.
// It requires a database connection object (*sql.DB) to perform operations.
func NewUserRepository(db database.DatabaseOperations) *userRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database using the provided UserPayload.
// It returns the ID of the newly created user or an error if the operation fails.
func (repo *userRepository) Create(payload models.UserPayload) (*uint, error) {
	var id uint
	// Set a timeout context to avoid long-running database operations.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// SQL statement to insert a new user and return the generated ID.
	statement := `
	INSERT INTO users (first_name, last_name, username, password, email)
	VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	// Execute the SQL statement with the provided payload data.
	err := repo.db.QueryRowContext(ctx, statement,
		payload.FirstName,
		payload.LastName,
		payload.Username,
		payload.Password,
		payload.Email,
	).Scan(&id)

	// Handle any errors that occur during the insert operation.
	if err != nil {
		// Use a utility function to handle common PostgreSQL errors.
		return nil, utils.HandlePostgresError(err)
	}

	// Return the ID of the newly created user.
	return &id, nil
}

// FindById retrieves a user by their ID from the database.
// It returns a pointer to a User model and any error encountered during the operation.
func (repo *userRepository) FindById(id uint) (*models.User, error) {
	// Define the context with a timeout to avoid long-running queries.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure the context is canceled to free resources.

	// SQL statement to select a user by ID.
	stmt := `
		SELECT id, first_name, last_name, username, password, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	// Initialize an empty User model to store the query result.
	var userFound models.User

	// Execute the query with the provided context and ID.
	row := repo.db.QueryRowContext(ctx, stmt, id)

	// Scan the result into the User model.
	err := row.Scan(
		&userFound.ID,
		&userFound.FirstName,
		&userFound.LastName,
		&userFound.Username,
		&userFound.Password,
		&userFound.Email,
		&userFound.CreatedAt,
		&userFound.UpdatedAt,
	)

	// Handle any errors that occurred during the query or scanning.
	if err != nil {
		// Use a helper function to handle common database errors.
		return nil, utils.HandlePostgresError(err)
	}

	// Return the found user and nil error if the operation was successful.
	return &userFound, nil
}

// FindUsers retrieves all users from the database.
// It returns a slice of pointers to User models and any error encountered during the operation.
func (repo *userRepository) FindUsers() ([]*models.User, error) {
	// Define the context with a timeout to avoid long-running queries.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure the context is canceled to free resources.

	// SQL statement to select all users.
	stmt := `
	SELECT id, first_name, last_name, username, password, email, created_at, updated_at
	FROM users
	`

	// Execute the query with the provided context.
	rows, err := repo.db.QueryContext(ctx, stmt)
	if err != nil {
		// Handle any database-related errors.
		return nil, utils.HandlePostgresError(err)
	}
	defer rows.Close() // Ensure the rows are closed after the function returns.

	users := []*models.User{} // Initialize a slice to hold the user records.

	// Iterate over the query results.
	for rows.Next() {
		var user models.User // Initialize a User struct to hold each record.

		// Scan the result into the User struct.
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			// Handle any scanning-related errors.
			return nil, utils.HandlePostgresError(err)
		}

		users = append(users, &user) // Add the user to the slice.
	}

	// Check for any errors encountered during iteration.
	if err := rows.Err(); err != nil {
		return nil, utils.HandlePostgresError(err)
	}

	return users, nil // Return the slice of users and nil error if successful.
}

// DeleteById removes a user from the database by their ID.
// It returns an error if the delete operation fails.
func (repo *userRepository) DeleteById(id uint) error {
	// Create a context with a timeout to prevent the operation from hanging indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure that the context is canceled when the operation is complete.

	// Prepare the SQL statement for deleting a user by ID.
	stmt := "DELETE FROM users WHERE id = $1"

	// Execute the delete operation with the provided context and ID.
	result, err := repo.db.ExecContext(ctx, stmt, id)
	if err != nil {
		// Handle any errors that occur during the delete operation.
		return utils.HandlePostgresError(err)
	}

	// Check if the record was actually updated.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.HandlePostgresError(err)
	}
	if rowsAffected == 0 {
		return utils.ErrNoDataFound
	}

	// Return nil if the delete operation is successful.
	return nil
}

// UpdateUser updates an existing user's information in the database.
// It takes a user model containing the updated information and the user's ID.
func (repo *userRepository) Update(id uint, user *models.UserPayload) error {
	// Create a context with a timeout to prevent the operation from hanging indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure that the context is canceled when the operation is complete.

	// Prepare the SQL statement for updating the user.
	stmt := `
	UPDATE users
	SET first_name = $1, last_name = $2, username = $3, password = $4, email = $5
	WHERE id = $6
	`

	// Execute the update operation with the provided context and user information.
	result, err := repo.db.ExecContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Password,
		user.Email,
		id,
	)
	if err != nil {
		// Handle any errors that occur during the update operation.
		return utils.HandlePostgresError(err)
	}

	// Check if the record was actually updated.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.HandlePostgresError(err)
	}
	if rowsAffected == 0 {
		return utils.ErrNoDataFound
	}

	// Return nil if the update operation is successful.
	return nil
}

// CheckIfEmailOrUsernameExist checks if a user with the given email or username exists in the database.
// It returns true if the user exists, and false otherwise.
func (repo *userRepository) CheckIfEmailOrUsernameExist(email, username string) bool {
	// Create a context with a timeout to ensure the query does not run indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure the context is canceled to avoid resource leaks.

	// Prepare the SQL statement to check for the existence of the email or username.
	stmt := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE email = $1 OR username = $2
		)
	`

	// Variable to store the result of the query.
	var isExists bool

	// Execute the query with the provided email and username, and scan the result into the isExists variable.
	err := repo.db.QueryRowContext(ctx, stmt, email, username).Scan(&isExists)
	if err != nil {
		// Log the error and return false if there's an error executing the query or scanning the result.
		log.Printf("Error checking if email or username exists: %v", err)
		return false
	}

	// Return the result of the query.
	return isExists
}
