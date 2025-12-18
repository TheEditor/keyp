package vault

import (
	"sync"
	"time"

	"github.com/TheEditor/keyp/internal/store"
)

// VaultHandle represents an unlocked vault that can be reused
// without re-entering the password repeatedly
type VaultHandle struct {
	mu         sync.RWMutex
	store      *store.Store
	key        []byte        // Derived encryption key
	unlockedAt time.Time
	timeout    time.Duration
	path       string
	password   string // Keep password for re-unlocking after auto-lock
}

// NewHandle creates a new vault handle (initially locked)
func NewHandle(path string) *VaultHandle {
	return &VaultHandle{
		path:    path,
		timeout: 30 * time.Minute, // Default 30 min timeout
	}
}

// Path returns the vault file path
func (h *VaultHandle) Path() string {
	return h.path
}

// Store returns the underlying store if unlocked, nil if locked
func (h *VaultHandle) Store() *store.Store {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.store
}

// IsUnlocked returns true if vault is currently unlocked
func (h *VaultHandle) IsUnlocked() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.store != nil
}

// IsExpired returns true if unlock timeout has elapsed
func (h *VaultHandle) IsExpired() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.store == nil {
		return true // Already locked
	}
	return time.Since(h.unlockedAt) > h.timeout
}

// Unlock opens the vault and keeps it open in the handle
func (h *VaultHandle) Unlock(password string, timeout time.Duration) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Open the vault
	v, err := Open(h.path, password)
	if err != nil {
		return err
	}

	// Extract store from vault (close vault, keep store)
	h.store = v.store
	h.key = v.key
	h.password = password
	h.unlockedAt = time.Now()

	if timeout > 0 {
		h.timeout = timeout
	}

	return nil
}

// Lock explicitly locks the vault
func (h *VaultHandle) Lock() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.store != nil {
		h.store.Close()
	}

	h.store = nil
	h.key = nil
	h.password = ""
	h.unlockedAt = time.Time{}
}

// UnlockedTime returns when the vault was unlocked
func (h *VaultHandle) UnlockedTime() time.Time {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.unlockedAt
}

// TimeUntilExpire returns time remaining until auto-lock
func (h *VaultHandle) TimeUntilExpire() time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.store == nil {
		return 0
	}

	elapsed := time.Since(h.unlockedAt)
	if elapsed >= h.timeout {
		return 0
	}

	return h.timeout - elapsed
}

// Timeout returns the current timeout setting
func (h *VaultHandle) Timeout() time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.timeout
}

// SetTimeout updates the timeout duration
func (h *VaultHandle) SetTimeout(timeout time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.timeout = timeout
}
