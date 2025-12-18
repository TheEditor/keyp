//go:build cgo

package store

import (
	"context"
	"os"
	"testing"

	"github.com/TheEditor/keyp/internal/model"
)

func TestCreateAndGetByName(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	secret := model.NewSecretObject("test-secret")
	secret.Notes = "test notes"
	secret.Tags = []string{"tag1", "tag2"}
	secret.AddField(model.NewField("username", "testuser"))
	secret.AddField(model.NewField("password", "testpass"))

	if err := s.Create(context.Background(), secret); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := s.GetByName(context.Background(), "test-secret")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}

	if retrieved.Name != "test-secret" {
		t.Errorf("Name mismatch: got %s", retrieved.Name)
	}
	if retrieved.Notes != "test notes" {
		t.Errorf("Notes mismatch: got %s", retrieved.Notes)
	}
	if len(retrieved.Tags) != 2 {
		t.Errorf("Tags count: got %d, want 2", len(retrieved.Tags))
	}
	if len(retrieved.Fields) != 2 {
		t.Errorf("Fields count: got %d, want 2", len(retrieved.Fields))
	}
}

func TestGetByNameNotFound(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	_, err := s.GetByName(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Create multiple secrets
	for _, name := range []string{"charlie", "alpha", "bravo"} {
		secret := model.NewSecretObject(name)
		s.Create(context.Background(), secret)
	}

	list, err := s.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("List count: got %d, want 3", len(list))
	}

	// Should be alphabetically sorted
	if list[0].Name != "alpha" {
		t.Errorf("Expected alpha first, got %s", list[0].Name)
	}
}

func TestSearch(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	// Create secrets with different content
	s1 := model.NewSecretObject("gmail-account")
	s1.Tags = []string{"email", "google"}
	s.Create(context.Background(), s1)

	s2 := model.NewSecretObject("bank-login")
	s2.Notes = "Chase online banking"
	s.Create(context.Background(), s2)

	s3 := model.NewSecretObject("work-email")
	s3.Tags = []string{"email", "work"}
	s.Create(context.Background(), s3)

	// Search by name
	results, _ := s.Search(context.Background(), "gmail", nil)
	if len(results) != 1 {
		t.Errorf("gmail search: got %d, want 1", len(results))
	}

	// Search by tag
	results, _ = s.Search(context.Background(), "email", nil)
	if len(results) != 2 {
		t.Errorf("email search: got %d, want 2", len(results))
	}

	// Search by notes
	results, _ = s.Search(context.Background(), "Chase", nil)
	if len(results) != 1 {
		t.Errorf("Chase search: got %d, want 1", len(results))
	}
}

func TestDelete(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	secret := model.NewSecretObject("to-delete")
	secret.AddField(model.NewField("data", "value"))
	s.Create(context.Background(), secret)

	if err := s.Delete(context.Background(), "to-delete"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := s.GetByName(context.Background(), "to-delete")
	if err != ErrNotFound {
		t.Error("Secret should be deleted")
	}
}

func TestDeleteNotFound(t *testing.T) {
	s := setupTestStore(t)
	defer s.Close()

	err := s.Delete(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func setupTestStore(t *testing.T) *Store {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "keyp-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	t.Cleanup(func() {
		os.Remove(tmpfile.Name())
	})

	s, err := Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	return s
}
