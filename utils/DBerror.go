package utils

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// Postgres error codes
const (
	ErrCodeUniqueViolation     = "23505"
	ErrCodeForeignKeyViolation = "23503"
	ErrCodeUndefinedTable      = "42P01"
	ErrCodeDataNotFound        = "P0002"
	// Add more error codes as needed.
)

// DBError is a custom error type that provides additional context beyond the standard error message.
// It includes an error code, a human-readable message, and the original error that occurred.
type DBError struct {
	Code    string // Code is a string that represents the error code associated with the error.
	Message string // Message is a string that contains a human-readable message describing the error.
	Err     error  // Err is the original error that is being wrapped by DBError.
}

// Error implements the error interface for DBError. It formats the error message to include
// the human-readable message, the error code, and the original error. This method allows DBError
// to be used wherever the error interface is expected.
func (e DBError) Error() string {
	return fmt.Sprintf("%s (code: %s): %v", e.Message, e.Code, e.Err)
}

// Unwrap returns the original error wrapped by DBError. This method allows users of DBError
// to retrieve the underlying error for further inspection or handling. It is particularly useful
// when using Go's error wrapping and unwrapping features to check for specific error types.
func (e DBError) Unwrap() error {
	return e.Err
}

// Is reports whether the provided error is equivalent to the current DBError. It is used
// in error comparisons, typically with errors.Is. This method allows for comparing an error
// against a known error value or type, based on the error code.
func (e DBError) Is(target error) bool {
	t, ok := target.(DBError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// NewDBError is a constructor function for DBError. It creates and returns a new DBError
// with the provided error code, message, and the original error. This function is useful
// for creating a DBError to return from functions that encounter errors, providing a consistent
// error handling mechanism across the application.
func NewDBError(code, message string, err error) DBError {
	return DBError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// HandlePostgresError interprets a given error as a PostgreSQL error and returns a new error with a detailed message.
// It unwraps the original error and checks for specific PostgreSQL error codes to provide context-specific responses.
func HandlePostgresError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ErrCodeUniqueViolation:
			// Wrap the original error with a new message indicating a unique constraint violation.
			return NewDBError(ErrCodeUniqueViolation, "user with a given email or username already exist", err)
		case ErrCodeForeignKeyViolation:
			// Wrap the original error with a new message indicating a foreign key violation.
			return NewDBError(ErrCodeForeignKeyViolation, "invalid refrerrence code", err)
		case ErrCodeUndefinedTable:
			// Wrap the original error with a new message indicating an undefined table error.
			return NewDBError(ErrCodeUndefinedTable, "database error: invalid table", err)
		case ErrCodeDataNotFound:
			// Wrap the original error with a new message indicating no data was found.
			return NewDBError(ErrCodeDataNotFound, "data not found", err)
		// ... add more cases as needed for other PostgreSQL error codes.

		default:
			// Wrap the original error with a generic database error message.
			return NewDBError(pgErr.Code, "database error", err)
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("data not found: %w", err)
	}

	// If the error is not a PgError, return the original error.
	return err
}
