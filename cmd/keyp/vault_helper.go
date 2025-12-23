package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/session"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var globalHandle *vault.VaultHandle
var sessionMgr *session.Manager

func init() {
	sessionMgr = session.New(session.DefaultTimeout)
}

// getOrUnlockVault returns the vault handle, unlocking if necessary
func getOrUnlockVault(cmd *cobra.Command, timeout time.Duration) (*vault.VaultHandle, error) {
	// Check process-global handle (only useful within same process)
	if globalHandle != nil && globalHandle.IsUnlocked() && !globalHandle.IsExpired() {
		return globalHandle, nil
	}

	// Try to load session from disk
	if derivedKey, err := sessionMgr.Load(); err == nil {
		// Session is valid, use it to unlock
		handle := vault.NewHandle(getVaultPath())
		if err := handle.UnlockWithKey(derivedKey, timeout); err == nil {
			globalHandle = handle
			return handle, nil
		}
		// If unlock fails, session is invalid, continue to prompt
	}

	// Need to unlock - prompt for password
	password, err := ui.PromptPassword("Vault password: ")
	if err != nil {
		return nil, err
	}

	handle := vault.NewHandle(getVaultPath())
	if err := handle.Unlock(password, timeout); err != nil {
		return nil, fmt.Errorf("failed to unlock vault: %w", err)
	}

	// Save session for future use
	if derivedKey := handle.GetDerivedKey(); derivedKey != nil {
		_ = sessionMgr.Save(derivedKey)
	}

	globalHandle = handle
	return handle, nil
}

// setVaultHandle stores a vault handle globally
func setVaultHandle(h *vault.VaultHandle) {
	globalHandle = h
}

// getVaultHandle retrieves the global vault handle
func getVaultHandle() *vault.VaultHandle {
	return globalHandle
}

// clearVaultHandle removes the global vault handle
func clearVaultHandle() {
	if globalHandle != nil {
		globalHandle.Lock()
	}
	globalHandle = nil
	// Also clear the session file
	_ = sessionMgr.Clear()
}
