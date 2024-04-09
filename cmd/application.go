package main

import (
	"fmt"
	"net/http"
	"time"
)

// application defines the configuration for an application server.
type application struct {
	Port int // Port defines the port on which the server listens.
}

// ApplicationOption represents a function that applies a configuration option to an application.
type ApplicationOption func(*application)

// WithPort creates an ApplicationOption that sets the port for the application.
func WithPort(port int) ApplicationOption {
	return func(a *application) {
		a.Port = port // Set the application's port to the specified value.
	}
}

// Application initializes a new application with the provided options.
// The default port is set to 2000 unless overridden by an option.
func Application(opts ...ApplicationOption) *application {
	app := &application{
		Port: 2000, // Default port.
	}

	// Apply each option to the application.
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// ServerOption represents a function that applies a configuration option to an http.Server.
type ServerOption func(*http.Server)

// WithReadTimeout creates a ServerOption that sets the read timeout for the server.
func WithReadTimeout(val time.Duration) ServerOption {
	return func(s *http.Server) {
		s.ReadTimeout = val // Set the server's read timeout to the specified duration.
	}
}

// WithWriteTimeout creates a ServerOption that sets the write timeout for the server.
func WithWriteTimeout(val time.Duration) ServerOption {
	return func(s *http.Server) {
		s.WriteTimeout = val // Set the server's write timeout to the specified duration.
	}
}

// WithReadHeaderTimeout creates a ServerOption that sets the read header timeout for the server.
func WithReadHeaderTimeout(val time.Duration) ServerOption {
	return func(s *http.Server) {
		s.ReadHeaderTimeout = val // Set the server's read header timeout to the specified duration.
	}
}

// WithIdleTimeout creates a ServerOption that sets the idle timeout for the server.
func WithIdleTimeout(val time.Duration) ServerOption {
	return func(s *http.Server) {
		s.IdleTimeout = val // Set the server's idle timeout to the specified duration.
	}
}

// Serve starts the application server with the given HTTP handler and server options.
// It initializes a new http.Server with default timeouts and applies any provided options.
func (app application) Serve(mux http.Handler, options ...ServerOption) error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.Port), // Address to bind the server to.
		Handler:           mux,                          // HTTP handler to invoke.
		ReadTimeout:       5 * time.Second,              // Default read timeout.
		WriteTimeout:      10 * time.Second,             // Default write timeout.
		IdleTimeout:       20 * time.Second,             // Default idle timeout.
		ReadHeaderTimeout: 3 * time.Second,              // Default read header timeout.
	}

	// Apply each server option to the http.Server.
	for _, opt := range options {
		opt(srv)
	}

	// Start the server and listen for incoming requests.
	return srv.ListenAndServe()
}
