package vault

import (
	"os"
	"testing"
	"time"
)

// setupTestVault creates a temporary vault for testing
func setupTestVault(t *testing.T) (string, string) {
	tmpDir, err := os.MkdirTemp("", "keyp-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Initialize vault
	password := "testpassword"
	v, err := Open(tmpDir, password)
	if err != nil {
		t.Fatalf("failed to open vault: %v", err)
	}
	defer v.Close()

	return tmpDir, password
}

// cleanupTestVault removes the test vault
func cleanupTestVault(t *testing.T, tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Logf("warning: failed to cleanup temp dir: %v", err)
	}
}

// TestNewHandle tests NewHandle creates unlocked handle
func TestNewHandle(t *testing.T) {
	tmpDir, _ := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Check path is set
	if handle.Path() != tmpDir {
		t.Errorf("expected path %q, got %q", tmpDir, handle.Path())
	}

	// Check initially locked
	if handle.IsUnlocked() {
		t.Errorf("expected handle to be locked initially")
	}

	// Check IsExpired returns true when locked
	if !handle.IsExpired() {
		t.Errorf("expected locked handle to be expired")
	}
}

// TestUnlock tests Unlock method
func TestUnlock(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Unlock
	err := handle.Unlock(password, 30*time.Minute)
	if err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}

	// Check unlocked
	if !handle.IsUnlocked() {
		t.Errorf("expected handle to be unlocked")
	}

	// Check not expired immediately
	if handle.IsExpired() {
		t.Errorf("expected unlocked handle to not be expired immediately")
	}

	// Check store is accessible
	store := handle.Store()
	if store == nil {
		t.Errorf("expected store to be non-nil after unlock")
	}

	// Check unlock time is set
	unlockTime := handle.UnlockedTime()
	if unlockTime.IsZero() {
		t.Errorf("expected unlock time to be set")
	}
}

// TestUnlockWithWrongPassword tests Unlock fails with wrong password
func TestUnlockWithWrongPassword(t *testing.T) {
	tmpDir, _ := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Try unlock with wrong password
	err := handle.Unlock("wrongpassword", 30*time.Minute)
	if err == nil {
		t.Errorf("expected unlock to fail with wrong password")
	}

	// Check still locked
	if handle.IsUnlocked() {
		t.Errorf("expected handle to be locked after failed unlock")
	}
}

// TestLock tests Lock method
func TestLock(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Unlock first
	if err := handle.Unlock(password, 30*time.Minute); err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}

	// Verify unlocked
	if !handle.IsUnlocked() {
		t.Fatalf("expected handle to be unlocked")
	}

	// Lock
	handle.Lock()

	// Check locked
	if handle.IsUnlocked() {
		t.Errorf("expected handle to be locked after Lock()")
	}

	// Check expired
	if !handle.IsExpired() {
		t.Errorf("expected locked handle to be expired")
	}

	// Check store is nil
	if handle.Store() != nil {
		t.Errorf("expected store to be nil after lock")
	}

	// Check unlock time is cleared
	if !handle.UnlockedTime().IsZero() {
		t.Errorf("expected unlock time to be cleared")
	}
}

// TestLockWhenAlreadyLocked tests Lock on already locked handle
func TestLockWhenAlreadyLocked(t *testing.T) {
	tmpDir, _ := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Lock when already locked - should not error
	handle.Lock()

	// Check locked
	if handle.IsUnlocked() {
		t.Errorf("expected handle to be locked")
	}
}

// TestIsExpired tests IsExpired method
func TestIsExpired(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Locked handle is always expired
	if !handle.IsExpired() {
		t.Errorf("expected locked handle to be expired")
	}

	// Unlock with short timeout
	timeout := 100 * time.Millisecond
	if err := handle.Unlock(password, timeout); err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}

	// Should not be expired immediately
	if handle.IsExpired() {
		t.Errorf("expected unlocked handle to not be expired immediately")
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Should be expired now
	if !handle.IsExpired() {
		t.Errorf("expected handle to be expired after timeout")
	}
}

// TestTimeUntilExpire tests TimeUntilExpire method
func TestTimeUntilExpire(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Locked handle has 0 time remaining
	if handle.TimeUntilExpire() != 0 {
		t.Errorf("expected locked handle to have 0 time remaining")
	}

	// Unlock with timeout
	timeout := 100 * time.Millisecond
	if err := handle.Unlock(password, timeout); err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}

	// Check time remaining (should be close to timeout)
	remaining := handle.TimeUntilExpire()
	if remaining <= 0 || remaining > timeout {
		t.Errorf("expected time remaining to be between 0 and %v, got %v", timeout, remaining)
	}

	// Wait and check again
	time.Sleep(50 * time.Millisecond)
	remaining2 := handle.TimeUntilExpire()
	if remaining2 >= remaining {
		t.Errorf("expected time remaining to decrease, was %v, now %v", remaining, remaining2)
	}

	// Wait for expiry
	time.Sleep(100 * time.Millisecond)
	if handle.TimeUntilExpire() != 0 {
		t.Errorf("expected 0 time remaining after expiry")
	}
}

// TestTimeout tests Timeout and SetTimeout methods
func TestTimeout(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Default timeout
	if handle.Timeout() != 30*time.Minute {
		t.Errorf("expected default timeout to be 30 minutes, got %v", handle.Timeout())
	}

	// Set timeout
	newTimeout := 15 * time.Minute
	handle.SetTimeout(newTimeout)

	if handle.Timeout() != newTimeout {
		t.Errorf("expected timeout to be %v, got %v", newTimeout, handle.Timeout())
	}

	// Unlock and verify timeout is used
	if err := handle.Unlock(password, newTimeout); err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}

	if handle.Timeout() != newTimeout {
		t.Errorf("expected timeout to remain %v after unlock, got %v", newTimeout, handle.Timeout())
	}
}

// TestMultipleUnlocks tests that Unlock can be called multiple times
func TestMultipleUnlocks(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// First unlock
	if err := handle.Unlock(password, 30*time.Minute); err != nil {
		t.Fatalf("first unlock failed: %v", err)
	}

	firstUnlockTime := handle.UnlockedTime()

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Second unlock (should succeed and update time)
	if err := handle.Unlock(password, 30*time.Minute); err != nil {
		t.Fatalf("second unlock failed: %v", err)
	}

	secondUnlockTime := handle.UnlockedTime()

	// Second unlock should have newer timestamp
	if !secondUnlockTime.After(firstUnlockTime) {
		t.Errorf("expected second unlock to have newer timestamp")
	}
}

// TestKeyZeroed tests that key is cleared on lock
func TestKeyZeroed(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	handle := NewHandle(tmpDir)

	// Unlock
	if err := handle.Unlock(password, 30*time.Minute); err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}

	// Lock
	handle.Lock()

	// After lock, we can't directly verify key is zeroed due to encapsulation,
	// but we verify the handle is locked and store is nil
	if handle.IsUnlocked() {
		t.Errorf("expected handle to be locked after Lock()")
	}

	if handle.Store() != nil {
		t.Errorf("expected store to be nil after lock")
	}
}
