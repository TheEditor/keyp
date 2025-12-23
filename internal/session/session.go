package session

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// DefaultTimeout is the default session timeout (15 minutes)
	DefaultTimeout = 15 * time.Minute
	// SessionFileName is the name of the session file
	SessionFileName = "session"
)

// Manager handles session persistence
type Manager struct {
	sessionDir string
	timeout    time.Duration
}

// New creates a new session manager
func New(timeout time.Duration) *Manager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	sessionDir := filepath.Join(homeDir, ".keyp")
	return &Manager{
		sessionDir: sessionDir,
		timeout:    timeout,
	}
}

// Save writes the derived key and expiry to the session file
func (m *Manager) Save(derivedKey []byte) error {
	// Ensure session directory exists
	if err := os.MkdirAll(m.sessionDir, 0700); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	sessionPath := filepath.Join(m.sessionDir, SessionFileName)

	// Create the session file with the derived key in hex and expiry timestamp
	keyHex := hex.EncodeToString(derivedKey)
	expiry := time.Now().Add(m.timeout)
	content := fmt.Sprintf("%s\n%d", keyHex, expiry.Unix())

	// Write with restricted permissions (0600)
	if err := os.WriteFile(sessionPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// Load reads the session file and returns the derived key if valid and not expired
func (m *Manager) Load() ([]byte, error) {
	sessionPath := filepath.Join(m.sessionDir, SessionFileName)

	// Check if session file exists
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no session found")
		}
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	// Parse the session file
	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid session file format")
	}

	keyHex := strings.TrimSpace(lines[0])
	expiryStr := strings.TrimSpace(lines[1])

	// Decode the hex key
	derivedKey, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode session key: %w", err)
	}

	// Parse the expiry timestamp
	var expiry int64
	_, err = fmt.Sscanf(expiryStr, "%d", &expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session expiry: %w", err)
	}

	// Check if session has expired
	if time.Now().Unix() > expiry {
		return nil, fmt.Errorf("session expired")
	}

	return derivedKey, nil
}

// Clear deletes the session file
func (m *Manager) Clear() error {
	sessionPath := filepath.Join(m.sessionDir, SessionFileName)
	err := os.Remove(sessionPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear session: %w", err)
	}
	return nil
}
