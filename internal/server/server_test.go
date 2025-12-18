package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/TheEditor/keyp/internal/vault"
)

// setupTestVault creates a temporary vault for testing
func setupTestVault(t *testing.T) (string, string) {
	tmpDir, err := os.MkdirTemp("", "keyp-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Initialize vault
	password := "testpassword"
	v, err := vault.Open(tmpDir, password)
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

// TestHealthEndpoint tests the /health endpoint
func TestHealthEndpoint(t *testing.T) {
	tmpDir, _ := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	srv := NewServer("localhost:0", tmpDir)
	server := httptest.NewServer(srv.Handler())
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/health", server.URL))
	if err != nil {
		t.Fatalf("failed to get /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

// TestVersionEndpoint tests the /version endpoint
func TestVersionEndpoint(t *testing.T) {
	tmpDir, _ := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	srv := NewServer("localhost:0", tmpDir)
	server := httptest.NewServer(srv.Handler())
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/version", server.URL))
	if err != nil {
		t.Fatalf("failed to get /version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

// TestProtectedEndpointRequiresAuth tests that protected endpoints require authentication
func TestProtectedEndpointRequiresAuth(t *testing.T) {
	tmpDir, _ := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	srv := NewServer("localhost:0", tmpDir)
	server := httptest.NewServer(srv.Handler())
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/secrets", server.URL))
	if err != nil {
		t.Fatalf("failed to get /v1/secrets: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

// TestUnlockAndListSecrets tests unlock followed by listing secrets
func TestUnlockAndListSecrets(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	srv := NewServer("localhost:0", tmpDir)
	server := httptest.NewServer(srv.Handler())
	defer server.Close()

	// Unlock vault
	unlockReq := UnlockRequest{Password: password}
	body, _ := json.Marshal(unlockReq)
	resp, err := http.Post(
		fmt.Sprintf("%s/v1/unlock", server.URL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to unlock: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var unlockResp Response
	if err := json.NewDecoder(resp.Body).Decode(&unlockResp); err != nil {
		t.Fatalf("failed to decode unlock response: %v", err)
	}

	var unlockData UnlockResponse
	if err := json.Unmarshal(unlockResp.Data, &unlockData); err != nil {
		t.Fatalf("failed to unmarshal unlock data: %v", err)
	}

	token := unlockData.Token
	if token == "" {
		t.Fatalf("expected token in response")
	}

	// List secrets with token
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/secrets", server.URL), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to list secrets: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

// TestInvalidPassword tests that invalid password is rejected
func TestInvalidPassword(t *testing.T) {
	tmpDir, _ := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	srv := NewServer("localhost:0", tmpDir)
	server := httptest.NewServer(srv.Handler())
	defer server.Close()

	unlockReq := UnlockRequest{Password: "wrongpassword"}
	body, _ := json.Marshal(unlockReq)
	resp, err := http.Post(
		fmt.Sprintf("%s/v1/unlock", server.URL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to call unlock: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

// TestSessionExpiry tests that sessions expire
func TestSessionExpiry(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	srv := NewServer("localhost:0", tmpDir)
	srv.SetSessionTimeout(100 * time.Millisecond)
	server := httptest.NewServer(srv.Handler())
	defer server.Close()

	unlockReq := UnlockRequest{Password: password}
	body, _ := json.Marshal(unlockReq)
	resp, _ := http.Post(
		fmt.Sprintf("%s/v1/unlock", server.URL),
		"application/json",
		bytes.NewReader(body),
	)

	var unlockResp Response
	json.NewDecoder(resp.Body).Decode(&unlockResp)
	var unlockData UnlockResponse
	json.Unmarshal(unlockResp.Data, &unlockData)
	resp.Body.Close()

	token := unlockData.Token

	time.Sleep(150 * time.Millisecond)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/secrets", server.URL), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, _ = http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for expired session, got %d", resp.StatusCode)
	}
}

// TestLockInvalidatesSession tests that lock invalidates session
func TestLockInvalidatesSession(t *testing.T) {
	tmpDir, password := setupTestVault(t)
	defer cleanupTestVault(t, tmpDir)

	srv := NewServer("localhost:0", tmpDir)
	server := httptest.NewServer(srv.Handler())
	defer server.Close()

	unlockReq := UnlockRequest{Password: password}
	body, _ := json.Marshal(unlockReq)
	resp, _ := http.Post(
		fmt.Sprintf("%s/v1/unlock", server.URL),
		"application/json",
		bytes.NewReader(body),
	)

	var unlockResp Response
	json.NewDecoder(resp.Body).Decode(&unlockResp)
	var unlockData UnlockResponse
	json.Unmarshal(unlockResp.Data, &unlockData)
	resp.Body.Close()

	token := unlockData.Token

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1/lock", server.URL), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, _ = http.DefaultClient.Do(req)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 from lock, got %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/v1/secrets", server.URL), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, _ = http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 after lock, got %d", resp.StatusCode)
	}
}
