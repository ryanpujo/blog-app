package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanpujo/blog-app/config"
)

// db is a package-level variable that holds the database connection.
var db *sql.DB

// initDB initializes a new database connection using the configuration settings.
// It returns a pointer to the sql.DB object and an error if any occurs during the connection process.
func initDB() (*sql.DB, error) {
	// Retrieve the configuration settings for the database connection.
	config := config.Config()

	// Attempt to open a new database connection using the DSN from the configuration.
	db, err := sql.Open("pgx", config.DSN)
	if err != nil {
		// If opening the database connection fails, return the error.
		return nil, err
	}

	// Ping the database to verify the connection is established.
	if err := db.Ping(); err != nil {
		// If pinging the database fails, return the error.
		return nil, err
	}

	// Return the database connection.
	return db, nil
}

// EstablishDBConnectionWithRetry attempts to establish a database connection with retries.
// It uses a ticker to retry the connection every 2 seconds and will log a fatal error
// if the connection cannot be established after 5 attempts.
func EstablishDBConnectionWithRetry() *sql.DB {
	// Create a new ticker that triggers every 2 seconds.
	ticker := time.NewTicker(2 * time.Second)
	// Ensure the ticker is stopped to free resources.
	defer ticker.Stop()

	// Initialize a counter to track the number of connection attempts.
	count := 0
	// Declare an error variable to hold any errors that occur during connection attempts.
	var err error

	// Continue attempting to establish a database connection until successful.
	for db == nil {
		// Call initDB to attempt to initialize the database connection.
		db, err = initDB()
		if err != nil {
			// If an error occurs, log the error message.
			log.Println("postgres is not ready yet: ", err)
		}

		// Increment the connection attempt counter.
		count++
		if count > 5 {
			// If more than 5 attempts have been made, log a fatal error and exit.
			log.Fatal("something went wrong: ", err)
		}

		// Wait for the next ticker signal before retrying.
		<-ticker.C
	}

	// Return the established database connection.
	return db
}
