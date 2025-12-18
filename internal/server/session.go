package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/TheEditor/keyp/internal/vault"
)

const TokenBytes = 32 // 256 bits

// SessionStore interface for session management
type SessionStore interface {
	Create(handle *vault.VaultHandle, expiry time.Duration) (*Session, error)
	Get(token string) (*Session, error)
	Delete(token string) error
	Refresh(token string, expiry time.Duration) error
	Cleanup()
	LockAll()
}

// MemorySessionStore implements SessionStore with in-memory storage
type MemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewSessionStore creates a new memory-based session store
func NewSessionStore() *MemorySessionStore {
	store := &MemorySessionStore{
		sessions: make(map[string]*Session),
	}

	// Start cleanup goroutine
	go store.cleanupLoop()

	return store
}

// cleanupLoop periodically removes expired sessions
func (m *MemorySessionStore) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.Cleanup()
	}
}

// generateToken creates a cryptographically secure token
func generateToken() (string, error) {
	bytes := make([]byte, TokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// Create creates a new session
func (m *MemorySessionStore) Create(handle *vault.VaultHandle, expiry time.Duration) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:     token,
		Handle:    handle,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expiry),
	}

	m.mu.Lock()
	m.sessions[token] = session
	m.mu.Unlock()

	return session, nil
}

// Get retrieves a session by token
func (m *MemorySessionStore) Get(token string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[token]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	// Check expiry
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

// Delete removes a session
func (m *MemorySessionStore) Delete(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[token]
	if !ok {
		return nil // Already deleted or never existed
	}

	// Lock the vault handle
	if handle, ok := session.Handle.(*vault.VaultHandle); ok {
		handle.Lock()
	}

	delete(m.sessions, token)
	return nil
}

// Refresh extends a session's expiry
func (m *MemorySessionStore) Refresh(token string, expiry time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[token]
	if !ok {
		return fmt.Errorf("session not found")
	}

	session.ExpiresAt = time.Now().Add(expiry)
	return nil
}

// Cleanup removes expired sessions
func (m *MemorySessionStore) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for token, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			// Lock the vault handle before deleting
			if handle, ok := session.Handle.(*vault.VaultHandle); ok {
				handle.Lock()
			}
			delete(m.sessions, token)
		}
	}
}

// LockAll locks all active sessions' vault handles
func (m *MemorySessionStore) LockAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, session := range m.sessions {
		if handle, ok := session.Handle.(*vault.VaultHandle); ok {
			handle.Lock()
		}
	}
}
