package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server represents the HTTP API server
type Server struct {
	httpServer     *http.Server
	mux            *http.ServeMux
	address        string
	vaultPath      string
	sessions       SessionStore
	sessionTimeout time.Duration
}

// NewServer creates a new server instance
func NewServer(address, vaultPath string) *Server {
	return &Server{
		address:        address,
		vaultPath:      vaultPath,
		mux:            http.NewServeMux(),
		sessions:       NewSessionStore(),
		sessionTimeout: 15 * time.Minute, // Default 15 min
	}
}

// SetSessionTimeout sets the session expiry duration
func (s *Server) SetSessionTimeout(timeout time.Duration) {
	s.sessionTimeout = timeout
}

// Handler returns the HTTP handler for the server
func (s *Server) Handler() http.Handler {
	return s.mux
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:    s.address,
		Handler: s.mux,
	}

	log.Printf("Starting server on %s", s.address)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	// Lock all active sessions
	s.sessions.LockAll()

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Println("Server stopped")
	return nil
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// contextKey type for context values
type contextKey string

const sessionKey contextKey = "session"

// getSessionFromContext extracts session from request context
func getSessionFromContext(ctx context.Context) *Session {
	session := ctx.Value(sessionKey)
	if session == nil {
		return nil
	}
	return session.(*Session)
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Public routes
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /version", s.handleVersion)

	// Auth routes (no auth required)
	s.mux.HandleFunc("POST /v1/unlock", s.handleUnlock)

	// Protected routes (require auth)
	s.mux.HandleFunc("POST /v1/lock", s.withAuth(s.handleLock))
	s.mux.HandleFunc("POST /v1/refresh", s.withAuth(s.handleRefresh))

	// Secret routes (protected)
	s.mux.HandleFunc("GET /v1/secrets", s.withAuth(s.handleListSecrets))
	s.mux.HandleFunc("POST /v1/secrets", s.withAuth(s.handleCreateSecret))
	s.mux.HandleFunc("GET /v1/secrets/{name}", s.withAuth(s.handleGetSecret))
	s.mux.HandleFunc("PUT /v1/secrets/{name}", s.withAuth(s.handleUpdateSecret))
	s.mux.HandleFunc("DELETE /v1/secrets/{name}", s.withAuth(s.handleDeleteSecret))

	// Search route (protected)
	s.mux.HandleFunc("GET /v1/search", s.withAuth(s.handleSearch))
}
