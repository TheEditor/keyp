package server

import (
	"context"
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
