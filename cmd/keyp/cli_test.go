package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/TheEditor/keyp/internal/model"
	"github.com/TheEditor/keyp/internal/store"
	"github.com/TheEditor/keyp/internal/vault"
)

// setupTestVault creates a temporary vault for CLI testing
func setupTestVault(t *testing.T) (string, string) {
	tmpDir, err := os.MkdirTemp("", "keyp-cli-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Initialize vault
	password := "testpassword"
	v, err := vault.Open(tmpDir, password)
	if err != nil {
		t.Fatalf("failed to initialize vault: %v", err)
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

// TestAddWorkflow tests the add command workflow
func TestAddWorkflow(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	t.Run("creates multi-field secret", func(t *testing.T) {
		v, err := vault.Open(tmpDir, password)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}

		secret := model.NewSecretObject("myapp")
		secret.AddField(model.NewField("username", "user@example.com"))
		secret.AddField(model.NewField("password", "secret123"))

		err = v.Create(context.Background(), secret)
		v.Close()
		if err != nil {
			t.Errorf("failed to create secret: %v", err)
		}

		// Verify secret was created
		v, err = vault.Open(tmpDir, password)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer v.Close()

		retrieved, err := v.GetByName(context.Background(), "myapp")
		if err != nil {
			t.Errorf("secret not found: %v", err)
		}

		if len(retrieved.Fields) != 2 {
			t.Errorf("expected 2 fields, got %d", len(retrieved.Fields))
		}
	})

	t.Run("secret with at least one field", func(t *testing.T) {
		v, err := vault.Open(tmpDir, password)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer v.Close()

		secret := model.NewSecretObject("test-secret")
		secret.AddField(model.NewField("value", "data"))

		err = v.Create(context.Background(), secret)
		if err != nil {
			t.Errorf("failed to create secret: %v", err)
		}

		if len(secret.Fields) != 1 {
			t.Errorf("expected at least 1 field")
		}
	})
}

// TestShowWorkflow tests the show command workflow
func TestShowWorkflow(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	// Create a test secret
	v, err := vault.Open(tmpDir, password)
	if err != nil {
		t.Fatalf("failed to open vault: %v", err)
	}

	secret := model.NewSecretObject("test-secret")
	secret.AddField(model.NewField("username", "alice"))
	field := model.NewField("password", "secret123")
	field.Sensitive = true
	secret.AddField(field)

	err = v.Create(context.Background(), secret)
	v.Close()
	if err != nil {
		t.Fatalf("failed to create secret: %v", err)
	}

	t.Run("displays all fields", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		secret, err := v.GetByName(context.Background(), "test-secret")
		if err != nil {
			t.Errorf("failed to get secret: %v", err)
		}

		if len(secret.Fields) != 2 {
			t.Errorf("should have 2 fields")
		}

		// Check fields exist
		hasUsername := false
		hasPassword := false
		for _, f := range secret.Fields {
			if f.Label == "username" {
				hasUsername = true
			}
			if f.Label == "password" {
				hasPassword = true
			}
		}

		if !hasUsername || !hasPassword {
			t.Errorf("secret should have username and password fields")
		}
	})

	t.Run("redacted version masks sensitive", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		secret, _ := v.GetByName(context.Background(), "test-secret")
		redacted := secret.Redacted()

		// Find password field
		for _, f := range redacted.Fields {
			if f.Label == "password" && f.Sensitive {
				if f.Value != model.RedactedValue {
					t.Errorf("sensitive value should be masked")
				}
			}
		}
	})

	t.Run("original version keeps values", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		secret, _ := v.GetByName(context.Background(), "test-secret")

		// Original should have actual values
		for _, f := range secret.Fields {
			if f.Label == "password" {
				if f.Value != "secret123" {
					t.Errorf("original value should be preserved")
				}
			}
		}
	})
}

// TestEditWorkflow tests the edit command workflow
func TestEditWorkflow(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	// Create test secret
	v, err := vault.Open(tmpDir, password)
	if err != nil {
		t.Fatalf("failed to open vault: %v", err)
	}

	secret := model.NewSecretObject("test-edit")
	secret.AddField(model.NewField("username", "alice"))
	secret.AddField(model.NewField("password", "secret123"))

	err = v.Create(context.Background(), secret)
	v.Close()
	if err != nil {
		t.Fatalf("failed to create secret: %v", err)
	}

	t.Run("update specific field", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		secret, _ := v.GetByName(context.Background(), "test-edit")

		// Update password field
		for i, f := range secret.Fields {
			if f.Label == "password" {
				secret.Fields[i].Value = "newpassword"
			}
		}

		err := v.Update(context.Background(), secret)
		if err != nil {
			t.Errorf("failed to update secret: %v", err)
		}
	})

	t.Run("verify field update", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		secret, _ := v.GetByName(context.Background(), "test-edit")

		// Check password was updated
		for _, f := range secret.Fields {
			if f.Label == "password" && f.Value != "newpassword" {
				t.Errorf("password field should be updated to 'newpassword'")
			}
		}
	})
}

// TestSearchWorkflow tests the search command workflow
func TestSearchWorkflow(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	// Create test secrets with tags
	v, err := vault.Open(tmpDir, password)
	if err != nil {
		t.Fatalf("failed to open vault: %v", err)
	}

	secret1 := model.NewSecretObject("github-account")
	secret1.AddField(model.NewField("username", "alice"))
	secret1.Tags = []string{"devops", "important"}
	v.Create(context.Background(), secret1)

	secret2 := model.NewSecretObject("gitlab-account")
	secret2.AddField(model.NewField("username", "bob"))
	secret2.Tags = []string{"devops"}
	v.Create(context.Background(), secret2)

	v.Close()

	t.Run("searches by name", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		results, err := v.Search(context.Background(), "github", nil)
		if err != nil {
			t.Errorf("search failed: %v", err)
		}

		if len(results) == 0 {
			t.Errorf("should find secret by name")
		}

		// Verify github-account is in results
		found := false
		for _, s := range results {
			if s.Name == "github-account" {
				found = true
			}
		}

		if !found {
			t.Errorf("search should find github-account")
		}
	})

	t.Run("searches by tag", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		results, err := v.Search(context.Background(), "important", nil)
		if err != nil {
			t.Errorf("search failed: %v", err)
		}

		if len(results) == 0 {
			t.Errorf("should find secret by tag")
		}

		// Should find github-account (tagged with "important")
		found := false
		for _, s := range results {
			if s.Name == "github-account" {
				found = true
			}
		}

		if !found {
			t.Errorf("search should find secret tagged with 'important'")
		}
	})

	t.Run("no results message", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		results, err := v.Search(context.Background(), "nonexistent", nil)
		if err != nil {
			t.Errorf("search should not error on no results: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("should return empty results for nonexistent query")
		}
	})
}

// TestTagWorkflow tests tag management
func TestTagWorkflow(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	// Create test secret
	v, err := vault.Open(tmpDir, password)
	if err != nil {
		t.Fatalf("failed to open vault: %v", err)
	}

	secret := model.NewSecretObject("test-tags")
	secret.AddField(model.NewField("value", "data"))
	v.Create(context.Background(), secret)
	v.Close()

	t.Run("add tag", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		secret, _ := v.GetByName(context.Background(), "test-tags")

		// Add tag
		secret.Tags = append(secret.Tags, "production")
		v.Update(context.Background(), secret)

		// Verify
		secret, _ = v.GetByName(context.Background(), "test-tags")

		found := false
		for _, tag := range secret.Tags {
			if tag == "production" {
				found = true
			}
		}

		if !found {
			t.Errorf("tag should be added")
		}
	})

	t.Run("remove tag", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		secret, _ := v.GetByName(context.Background(), "test-tags")

		// Remove tag
		newTags := []string{}
		for _, tag := range secret.Tags {
			if tag != "production" {
				newTags = append(newTags, tag)
			}
		}
		secret.Tags = newTags
		v.Update(context.Background(), secret)

		// Verify
		secret, _ = v.GetByName(context.Background(), "test-tags")

		for _, tag := range secret.Tags {
			if tag == "production" {
				t.Errorf("tag should be removed")
			}
		}
	})

	t.Run("list tags", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)
		defer v.Close()

		secret, _ := v.GetByName(context.Background(), "test-tags")
		secret.Tags = []string{"api", "security"}
		v.Update(context.Background(), secret)

		// List secrets and collect tags
		secrets, _ := v.List(context.Background(), nil)
		tagSet := make(map[string]bool)
		for _, s := range secrets {
			for _, t := range s.Tags {
				tagSet[t] = true
			}
		}

		if !tagSet["api"] || !tagSet["security"] {
			t.Errorf("tags should be in tag set")
		}
	})
}

// TestVaultHandleIntegration tests VaultHandle unlock/lock operations
func TestVaultHandleIntegration(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	t.Run("unlock with correct password", func(t *testing.T) {
		handle := vault.NewHandle(tmpDir)

		// Initially locked
		if handle.IsUnlocked() {
			t.Errorf("expected vault handle to be locked initially")
		}

		// Unlock with correct password
		err := handle.Unlock(password, 30*time.Minute)
		if err != nil {
			t.Fatalf("failed to unlock with correct password: %v", err)
		}

		// Should be unlocked
		if !handle.IsUnlocked() {
			t.Errorf("expected vault handle to be unlocked after unlock")
		}

		// Store should be accessible
		if handle.Store() == nil {
			t.Errorf("expected store to be non-nil after unlock")
		}
	})

	t.Run("unlock fails with wrong password", func(t *testing.T) {
		handle := vault.NewHandle(tmpDir)

		// Try to unlock with wrong password
		err := handle.Unlock("wrongpassword", 30*time.Minute)
		if err == nil {
			t.Errorf("expected unlock to fail with wrong password")
		}

		// Should remain locked
		if handle.IsUnlocked() {
			t.Errorf("expected vault handle to remain locked after failed unlock")
		}
	})

	t.Run("lock unlocked vault", func(t *testing.T) {
		handle := vault.NewHandle(tmpDir)

		// Unlock
		if err := handle.Unlock(password, 30*time.Minute); err != nil {
			t.Fatalf("failed to unlock: %v", err)
		}

		// Lock
		handle.Lock()

		// Should be locked
		if handle.IsUnlocked() {
			t.Errorf("expected vault handle to be locked after lock")
		}

		// Store should be nil
		if handle.Store() != nil {
			t.Errorf("expected store to be nil after lock")
		}
	})

	t.Run("lock when already locked is safe", func(t *testing.T) {
		handle := vault.NewHandle(tmpDir)

		// Lock when already locked - should not panic
		handle.Lock()

		if handle.IsUnlocked() {
			t.Errorf("expected vault handle to be locked")
		}
	})

	t.Run("session persistence within timeout", func(t *testing.T) {
		handle := vault.NewHandle(tmpDir)

		// Unlock
		if err := handle.Unlock(password, 100*time.Millisecond); err != nil {
			t.Fatalf("failed to unlock: %v", err)
		}

		// Should not be expired immediately
		if handle.IsExpired() {
			t.Errorf("expected handle to not be expired immediately after unlock")
		}

		// Wait partially
		time.Sleep(50 * time.Millisecond)

		// Still should not be expired
		if handle.IsExpired() {
			t.Errorf("expected handle to not be expired within timeout window")
		}
	})

	t.Run("auto-lock after timeout", func(t *testing.T) {
		handle := vault.NewHandle(tmpDir)

		// Unlock with short timeout
		if err := handle.Unlock(password, 50*time.Millisecond); err != nil {
			t.Fatalf("failed to unlock: %v", err)
		}

		// Should be unlocked
		if !handle.IsUnlocked() {
			t.Errorf("expected handle to be unlocked")
		}

		// Wait for timeout
		time.Sleep(100 * time.Millisecond)

		// Should be expired
		if !handle.IsExpired() {
			t.Errorf("expected handle to be expired after timeout")
		}
	})

	t.Run("multiple unlocks reset timer", func(t *testing.T) {
		handle := vault.NewHandle(tmpDir)

		// First unlock
		if err := handle.Unlock(password, 100*time.Millisecond); err != nil {
			t.Fatalf("first unlock failed: %v", err)
		}

		firstTime := handle.UnlockedTime()

		// Wait 40ms
		time.Sleep(40 * time.Millisecond)

		// Unlock again (should reset timer)
		if err := handle.Unlock(password, 100*time.Millisecond); err != nil {
			t.Fatalf("second unlock failed: %v", err)
		}

		secondTime := handle.UnlockedTime()

		// Second unlock should have newer timestamp
		if !secondTime.After(firstTime) {
			t.Errorf("expected second unlock to reset timer with newer timestamp")
		}

		// Wait 60ms more (total 100ms from first, but 60ms from second)
		time.Sleep(60 * time.Millisecond)

		// Should still not be expired (only 100ms from second unlock)
		if handle.IsExpired() {
			t.Errorf("expected handle to not be expired - timer should have been reset")
		}

		// Wait more to ensure actual expiry
		time.Sleep(50 * time.Millisecond)
		if !handle.IsExpired() {
			t.Errorf("expected handle to be expired after full timeout from second unlock")
		}
	})

	t.Run("time until expire decreases", func(t *testing.T) {
		handle := vault.NewHandle(tmpDir)

		// Unlock
		if err := handle.Unlock(password, 200*time.Millisecond); err != nil {
			t.Fatalf("failed to unlock: %v", err)
		}

		// Get initial remaining time
		initial := handle.TimeUntilExpire()
		if initial <= 0 {
			t.Errorf("expected initial time remaining to be positive")
		}

		// Wait
		time.Sleep(50 * time.Millisecond)

		// Get new remaining time
		later := handle.TimeUntilExpire()
		if later >= initial {
			t.Errorf("expected time remaining to decrease")
		}
	})

	t.Run("timeout parameter is respected", func(t *testing.T) {
		shortHandle := vault.NewHandle(tmpDir)
		longHandle := vault.NewHandle(tmpDir)

		// Unlock both
		if err := shortHandle.Unlock(password, 50*time.Millisecond); err != nil {
			t.Fatalf("short handle unlock failed: %v", err)
		}
		if err := longHandle.Unlock(password, 200*time.Millisecond); err != nil {
			t.Fatalf("long handle unlock failed: %v", err)
		}

		// After 100ms
		time.Sleep(100 * time.Millisecond)

		// Short should be expired
		if !shortHandle.IsExpired() {
			t.Errorf("expected short timeout handle to be expired")
		}

		// Long should not be expired
		if longHandle.IsExpired() {
			t.Errorf("expected long timeout handle to not be expired")
		}
	})
}

// TestCLIInitWorkflow tests the init command workflow
func TestCLIInitWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "keyp-cli-init-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer cleanupTestVault(t, tmpDir)

	// Vault path doesn't exist yet
	vaultPath := tmpDir + "/vault.db"

	// Create and initialize vault
	v, err := vault.Init(vaultPath, "testpassword")
	if err != nil {
		t.Fatalf("failed to initialize vault: %v", err)
	}
	defer v.Close()

	// Verify vault was created
	if _, err := os.Stat(vaultPath); err != nil {
		t.Errorf("vault file not created: %v", err)
	}

	// Verify we can open it again with same password
	v2, err := vault.Open(vaultPath, "testpassword")
	if err != nil {
		t.Fatalf("failed to open existing vault: %v", err)
	}
	v2.Close()
}

// TestCLISetGetWorkflow tests set and get commands
func TestCLISetGetWorkflow(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	ctx := context.Background()

	t.Run("set creates secret and get retrieves it", func(t *testing.T) {
		v, err := vault.Open(tmpDir, password)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}

		// Create a secret
		secret := model.NewSecretObject("myapp")
		secret.AddField(model.NewField("password", "secret123"))

		err = v.Create(ctx, secret)
		if err != nil {
			t.Fatalf("failed to create secret: %v", err)
		}
		v.Close()

		// Retrieve the secret
		v2, err := vault.Open(tmpDir, password)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}

		retrieved, err := v2.GetByName(ctx, "myapp")
		if err != nil {
			t.Fatalf("failed to retrieve secret: %v", err)
		}
		v2.Close()

		if retrieved.Name != "myapp" {
			t.Errorf("expected name 'myapp', got '%s'", retrieved.Name)
		}

		if len(retrieved.Fields) != 1 {
			t.Errorf("expected 1 field, got %d", len(retrieved.Fields))
		}

		if retrieved.Fields[0].Value != "secret123" {
			t.Errorf("expected password 'secret123', got '%s'", retrieved.Fields[0].Value)
		}
	})

	t.Run("set with existing name updates secret", func(t *testing.T) {
		v, err := vault.Open(tmpDir, password)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}

		// Create initial secret
		secret := model.NewSecretObject("updatetest")
		secret.AddField(model.NewField("value", "initial"))

		err = v.Create(ctx, secret)
		if err != nil {
			t.Fatalf("failed to create secret: %v", err)
		}

		// Get and update it
		existing, _ := v.GetByName(ctx, "updatetest")
		existing.Fields[0].Value = "updated"
		err = v.Update(ctx, existing)
		if err != nil {
			t.Fatalf("failed to update secret: %v", err)
		}
		v.Close()

		// Verify update persisted
		v2, _ := vault.Open(tmpDir, password)
		retrieved, _ := v2.GetByName(ctx, "updatetest")
		v2.Close()

		if retrieved.Fields[0].Value != "updated" {
			t.Errorf("expected updated value, got '%s'", retrieved.Fields[0].Value)
		}
	})
}

// TestCLIListWorkflow tests the list command
func TestCLIListWorkflow(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	ctx := context.Background()

	t.Run("list shows all secrets", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)

		// Create multiple secrets
		for i := 1; i <= 3; i++ {
			name := "secret" + string(rune('0'+i))
			secret := model.NewSecretObject(name)
			secret.AddField(model.NewField("data", "value"+string(rune('0'+i))))
			v.Create(ctx, secret)
		}

		secrets, err := v.List(ctx, nil)
		if err != nil {
			t.Fatalf("failed to list secrets: %v", err)
		}
		v.Close()

		if len(secrets) < 3 {
			t.Errorf("expected at least 3 secrets, got %d", len(secrets))
		}
	})

	t.Run("list with tags filters correctly", func(t *testing.T) {
		v, _ := vault.Open(tmpDir, password)

		// Create secret with tags
		secret := model.NewSecretObject("tagged-secret")
		secret.AddField(model.NewField("data", "value"))
		secret.Tags = []string{"api", "security"}
		v.Create(ctx, secret)

		// List with tag filter
		opts := &store.SearchOptions{Tags: []string{"api"}}
		secrets, err := v.List(ctx, opts)
		if err != nil {
			t.Fatalf("failed to list with tags: %v", err)
		}
		v.Close()

		found := false
		for _, s := range secrets {
			if s.Name == "tagged-secret" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find tagged-secret in filtered list")
		}
	})
}

// TestCLIDeleteWorkflow tests the delete command
func TestCLIDeleteWorkflow(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	ctx := context.Background()

	v, _ := vault.Open(tmpDir, password)

	// Create a secret
	secret := model.NewSecretObject("to-delete")
	secret.AddField(model.NewField("data", "value"))
	v.Create(ctx, secret)

	// Verify it exists
	_, err := v.GetByName(ctx, "to-delete")
	if err != nil {
		t.Fatalf("secret should exist before delete: %v", err)
	}

	// Delete the secret
	err = v.Delete(ctx, "to-delete")
	if err != nil {
		t.Fatalf("failed to delete secret: %v", err)
	}
	v.Close()

	// Verify it's gone
	v2, _ := vault.Open(tmpDir, password)
	_, err = v2.GetByName(ctx, "to-delete")
	v2.Close()

	if err == nil {
		t.Errorf("secret should not exist after delete")
	}
}
