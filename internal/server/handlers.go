package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TheEditor/keyp/internal/store"
	"github.com/TheEditor/keyp/internal/vault"
)

// Public endpoints

// handleHealth returns health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, SuccessResponse(HealthResponse{Status: "ok"}))
}

// handleVersion returns version information
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, SuccessResponse(VersionResponse{
		Version: "0.1.0",
	}))
}

// Auth endpoints

// handleUnlock unlocks the vault and creates a session
func (s *Server) handleUnlock(w http.ResponseWriter, r *http.Request) {
	var req UnlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse(ErrCodeBadRequest, "Invalid JSON"))
		return
	}

	// Create vault handle and unlock
	handle := vault.NewHandle(s.vaultPath)
	if err := handle.Unlock(req.Password, 0); err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Invalid password"))
		return
	}

	// Create session
	session, err := s.sessions.Create(handle, s.sessionTimeout)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to create session"))
		return
	}

	writeJSON(w, http.StatusOK, SuccessResponse(UnlockResponse{
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
	}))
}

// handleLock locks the vault and invalidates the session
func (s *Server) handleLock(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	// Lock the vault handle
	handle := session.Handle.(*vault.VaultHandle)
	handle.Lock()

	// Delete session
	s.sessions.Delete(session.Token)

	writeJSON(w, http.StatusOK, SuccessResponse(nil))
}

// handleRefresh extends session expiry
func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	// Extend session expiry
	if err := s.sessions.Refresh(session.Token, s.sessionTimeout); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to refresh session"))
		return
	}

	// Get updated session
	session, _ = s.sessions.Get(session.Token)

	writeJSON(w, http.StatusOK, SuccessResponse(RefreshResponse{
		ExpiresAt: session.ExpiresAt,
	}))
}

// Protected endpoints

// handleListSecrets lists all secrets
func (s *Server) handleListSecrets(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	handle := session.Handle.(*vault.VaultHandle)
	st := handle.Store()
	if st == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Vault locked"))
		return
	}

	// List secrets
	secrets, err := st.List(r.Context(), nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to list secrets"))
		return
	}

	// Convert to API types
	items := make([]SecretListItem, len(secrets))
	for i, sec := range secrets {
		items[i] = ToSecretListItem(sec)
	}

	writeJSON(w, http.StatusOK, SuccessResponse(items))
}

// handleCreateSecret creates a new secret
func (s *Server) handleCreateSecret(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	handle := session.Handle.(*vault.VaultHandle)
	st := handle.Store()
	if st == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Vault locked"))
		return
	}

	var req CreateSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse(ErrCodeBadRequest, "Invalid JSON"))
		return
	}

	// Validate
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse(ErrCodeBadRequest, "Name is required"))
		return
	}
	if len(req.Fields) == 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse(ErrCodeBadRequest, "At least one field is required"))
		return
	}

	// Convert and create
	secret := req.ToSecretObject()
	if err := st.Create(r.Context(), secret); err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeJSON(w, http.StatusConflict, ErrorResponse(ErrCodeConflict, "Secret already exists"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to create secret"))
		return
	}

	detail := ToSecretDetail(secret, true)
	writeJSON(w, http.StatusCreated, SuccessResponse(detail))
}

// handleGetSecret retrieves a secret by name
func (s *Server) handleGetSecret(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	handle := session.Handle.(*vault.VaultHandle)
	st := handle.Store()
	if st == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Vault locked"))
		return
	}

	name := r.PathValue("name")

	secret, err := st.GetByName(r.Context(), name)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse(ErrCodeNotFound, "Secret not found"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to get secret"))
		return
	}

	// Return redacted by default
	detail := ToSecretDetail(secret, true)
	writeJSON(w, http.StatusOK, SuccessResponse(detail))
}

// handleUpdateSecret updates an existing secret
func (s *Server) handleUpdateSecret(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	handle := session.Handle.(*vault.VaultHandle)
	st := handle.Store()
	if st == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Vault locked"))
		return
	}

	name := r.PathValue("name")

	// Get existing secret
	secret, err := st.GetByName(r.Context(), name)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse(ErrCodeNotFound, "Secret not found"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to get secret"))
		return
	}

	// Parse update request
	var req UpdateSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse(ErrCodeBadRequest, "Invalid JSON"))
		return
	}

	// Apply updates
	if req.Tags != nil {
		secret.Tags = *req.Tags
	}
	if req.Notes != nil {
		secret.Notes = *req.Notes
	}

	// Update
	if err := st.Update(r.Context(), secret); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to update secret"))
		return
	}

	detail := ToSecretDetail(secret, true)
	writeJSON(w, http.StatusOK, SuccessResponse(detail))
}

// handleDeleteSecret deletes a secret
func (s *Server) handleDeleteSecret(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	handle := session.Handle.(*vault.VaultHandle)
	st := handle.Store()
	if st == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Vault locked"))
		return
	}

	name := r.PathValue("name")

	if err := st.Delete(r.Context(), name); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse(ErrCodeNotFound, "Secret not found"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to delete secret"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleSearch searches for secrets
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	handle := session.Handle.(*vault.VaultHandle)
	st := handle.Store()
	if st == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Vault locked"))
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse(ErrCodeBadRequest, "Query parameter 'q' is required"))
		return
	}

	// Search
	results, err := st.Search(r.Context(), query, nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Search failed"))
		return
	}

	// Convert to list items
	items := make([]SecretListItem, len(results))
	for i, sec := range results {
		items[i] = ToSecretListItem(sec)
	}

	writeJSON(w, http.StatusOK, SuccessResponse(items))
}

// handleClipboard copies a secret field to clipboard
func (s *Server) handleClipboard(w http.ResponseWriter, r *http.Request) {
	session := getSessionFromContext(r.Context())
	if session == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Not authenticated"))
		return
	}

	handle := session.Handle.(*vault.VaultHandle)
	st := handle.Store()
	if st == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse(ErrCodeUnauthorized, "Vault locked"))
		return
	}

	name := r.PathValue("name")

	// Get secret
	secret, err := st.GetByName(r.Context(), name)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse(ErrCodeNotFound, "Secret not found"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse(ErrCodeInternalError, "Failed to get secret"))
		return
	}

	// Parse request
	var req ClipboardRequest
	json.NewDecoder(r.Body).Decode(&req) // OK if empty

	// Find field value
	var value string
	if req.Field != "" {
		for _, f := range secret.Fields {
			if f.Label == req.Field {
				value = f.Value
				break
			}
		}
		if value == "" {
			writeJSON(w, http.StatusNotFound, ErrorResponse(ErrCodeNotFound, "Field not found"))
			return
		}
	} else if len(secret.Fields) > 0 {
		value = secret.Fields[0].Value
	} else {
		writeJSON(w, http.StatusBadRequest, ErrorResponse(ErrCodeBadRequest, "Secret has no fields"))
		return
	}

	// Copy to clipboard (server-side - just acknowledge for now)
	// Note: Server-side clipboard is typically done via OS clipboard integration
	// For now, we acknowledge the request succeeded
	writeJSON(w, http.StatusOK, SuccessResponse(nil))
}
