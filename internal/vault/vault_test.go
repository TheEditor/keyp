//go:build cgo

package vault

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/TheEditor/keyp/internal/model"
	"github.com/TheEditor/keyp/internal/store"
)

func TestVaultEncryption(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	password := "testpassword123"

	// Create vault
	v, err := Init(path, password)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer v.Close()

	// Add secret with sensitive field
	secret := model.NewSecretObject("test")
	secret.AddField(model.NewField("password", "supersecret"))
	err = v.Create(context.Background(), secret)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	v.Close()

	// Reopen with correct password
	v2, err := Open(path, password)
	if err != nil {
		t.Fatalf("Open with correct password failed: %v", err)
	}
	defer v2.Close()

	// Retrieve and verify
	retrieved, err := v2.GetByName(context.Background(), "test")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}
	if len(retrieved.Fields) == 0 {
		t.Fatal("No fields in retrieved secret")
	}
	if retrieved.Fields[0].Value != "supersecret" {
		t.Errorf("Expected 'supersecret', got %q", retrieved.Fields[0].Value)
	}
}

func TestVaultWrongPassword(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	password := "correctpassword"

	// Create vault
	v, err := Init(path, password)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	v.Close()

	// Try to open with wrong password
	_, err = Open(path, "wrongpassword")
	if err != store.ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword, got %v", err)
	}
}

func TestVaultDecryptionFailureWithWrongPassword(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	password := "testpassword123"

	// Create vault and add a secret
	v, err := Init(path, password)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer v.Close()

	secret := model.NewSecretObject("test")
	secret.AddField(model.NewField("password", "supersecret"))
	err = v.Create(context.Background(), secret)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	v.Close()

	// Try opening with wrong password - should fail decryption
	v3, err := Open(path, "wrongpassword")
	if err != store.ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword for wrong password, got %v", err)
	}
	if v3 != nil {
		v3.Close()
	}
}

func TestVaultCreateCloseReopenGetByName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	password := "testpassword123"

	// Create vault and add secrets
	v, err := Init(path, password)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	secret1 := model.NewSecretObject("github")
	secret1.AddField(model.NewField("username", "alice"))
	secret1.AddField(model.NewField("token", "ghp_abc123"))
	err = v.Create(context.Background(), secret1)
	if err != nil {
		t.Fatalf("Create secret1 failed: %v", err)
	}

	secret2 := model.NewSecretObject("twitter")
	secret2.AddField(model.NewField("username", "alice"))
	secret2.AddField(model.NewField("password", "secret123"))
	err = v.Create(context.Background(), secret2)
	if err != nil {
		t.Fatalf("Create secret2 failed: %v", err)
	}
	v.Close()

	// Reopen and retrieve
	v2, err := Open(path, password)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer v2.Close()

	// Get first secret
	github, err := v2.GetByName(context.Background(), "github")
	if err != nil {
		t.Fatalf("GetByName github failed: %v", err)
	}
	if github.Fields[0].Label != "username" || github.Fields[0].Value != "alice" {
		t.Errorf("github username mismatch")
	}
	if github.Fields[1].Label != "token" || github.Fields[1].Value != "ghp_abc123" {
		t.Errorf("github token mismatch")
	}

	// Get second secret
	twitter, err := v2.GetByName(context.Background(), "twitter")
	if err != nil {
		t.Fatalf("GetByName twitter failed: %v", err)
	}
	if twitter.Fields[1].Label != "password" || twitter.Fields[1].Value != "secret123" {
		t.Errorf("twitter password mismatch")
	}
}

func TestVaultListDecrypts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	password := "testpassword123"

	// Create vault and add multiple secrets
	v, err := Init(path, password)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	for _, name := range []string{"secret1", "secret2", "secret3"} {
		secret := model.NewSecretObject(name)
		secret.AddField(model.NewField("value", "secret_"+name))
		err = v.Create(context.Background(), secret)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}
	v.Close()

	// Reopen and list
	v2, err := Open(path, password)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer v2.Close()

	list, err := v2.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 secrets, got %d", len(list))
	}

	// Verify all values are decrypted
	for _, secret := range list {
		if len(secret.Fields) == 0 {
			t.Errorf("Secret %s has no fields", secret.Name)
			continue
		}
		expected := "secret_" + secret.Name
		if secret.Fields[0].Value != expected {
			t.Errorf("Secret %s field value mismatch: expected %q, got %q",
				secret.Name, expected, secret.Fields[0].Value)
		}
	}
}

func TestVaultUpdateEncrypts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	password := "testpassword123"

	// Create vault and add secret
	v, err := Init(path, password)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	secret := model.NewSecretObject("myapp")
	secret.AddField(model.NewField("password", "old_password"))
	err = v.Create(context.Background(), secret)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	v.Close()

	// Reopen, update, and verify
	v2, err := Open(path, password)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	retrieved, err := v2.GetByName(context.Background(), "myapp")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}

	// Update the field
	retrieved.Fields[0].Value = "new_password"
	err = v2.Update(context.Background(), retrieved)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	v2.Close()

	// Reopen and verify update persisted
	v3, err := Open(path, password)
	if err != nil {
		t.Fatalf("Open after update failed: %v", err)
	}
	defer v3.Close()

	updated, err := v3.GetByName(context.Background(), "myapp")
	if err != nil {
		t.Fatalf("GetByName after update failed: %v", err)
	}

	if updated.Fields[0].Value != "new_password" {
		t.Errorf("Expected 'new_password', got %q", updated.Fields[0].Value)
	}
}

func TestVaultNonSensitiveFieldsNotEncrypted(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	password := "testpassword123"

	// Create vault and add secret with mixed sensitive/non-sensitive fields
	v, err := Init(path, password)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	secret := model.NewSecretObject("mixed")
	f1 := model.NewField("username", "alice")
	f1.Sensitive = false
	secret.AddField(f1)

	f2 := model.NewField("password", "secret123")
	f2.Sensitive = true
	secret.AddField(f2)

	err = v.Create(context.Background(), secret)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Get back and verify
	retrieved, err := v.GetByName(context.Background(), "mixed")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}

	// Non-sensitive field should be plaintext
	if retrieved.Fields[0].Value != "alice" {
		t.Errorf("Non-sensitive field corrupted: %q", retrieved.Fields[0].Value)
	}

	// Sensitive field should be plaintext after decryption
	if retrieved.Fields[1].Value != "secret123" {
		t.Errorf("Sensitive field decryption failed: %q", retrieved.Fields[1].Value)
	}

	v.Close()
}

func TestVaultFileNotExists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.db")

	_, err := Open(path, "password")
	if err != ErrNotExists {
		t.Errorf("Expected ErrNotExists, got %v", err)
	}
}

func TestVaultAlreadyExists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	// Create first vault
	v1, err := Init(path, "password1")
	if err != nil {
		t.Fatalf("First Init failed: %v", err)
	}
	v1.Close()

	// Try to create again
	_, err = Init(path, "password2")
	if err != ErrAlreadyExists {
		t.Errorf("Expected ErrAlreadyExists, got %v", err)
	}
}

func setupTestVault(t *testing.T, password string) (*Vault, string) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "keyp-vault-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	t.Cleanup(func() {
		os.Remove(tmpfile.Name())
	})

	v, err := Init(tmpfile.Name(), password)
	if err != nil {
		t.Fatal(err)
	}

	return v, tmpfile.Name()
}
