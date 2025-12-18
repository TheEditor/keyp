package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// withAuth is middleware that verifies Bearer token authentication
func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract Bearer token
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Missing or invalid Authorization header"))
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		// Lookup session
		session, err := s.sessions.Get(token)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Invalid or expired token"))
			return
		}

		// Check expiry
		if time.Now().After(session.ExpiresAt) {
			s.sessions.Delete(token)
			writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Session expired"))
			return
		}

		// Add session to context
		ctx := context.WithValue(r.Context(), sessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// withLogging logs HTTP requests with method, path, status, and duration
func (s *Server) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.status, duration)
	})
}

// withRecovery catches panics and returns 500 error response
func (s *Server) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				resp := ErrorResponse(ErrCodeInternalError, "Internal server error")
				json.NewEncoder(w).Encode(resp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
