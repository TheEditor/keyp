package sync

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a temporary test repository
func setupTestRepo(t *testing.T) (string, *GitExecSyncer) {
	tmpDir, err := os.MkdirTemp("", "keyp-sync-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	syncer := NewGitExecSyncer(tmpDir)
	return tmpDir, syncer
}

// cleanupTestRepo removes the test repository
func cleanupTestRepo(t *testing.T, tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Logf("warning: failed to cleanup temp dir: %v", err)
	}
}

// TestInit tests that Init() creates .git directory and .gitignore
func TestInit(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Initialize repo
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check .git directory exists
	gitDir := filepath.Join(tmpDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf("expected .git directory to exist")
	}

	// Check .gitignore exists
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Errorf("expected .gitignore to exist")
	}

	// Check .gitignore content
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	contentStr := string(content)
	expectedPatterns := []string{"*.db", "*.db-journal", "*.db-wal", "*.db-shm"}
	for _, pattern := range expectedPatterns {
		if !contains(contentStr, pattern) {
			t.Errorf("expected .gitignore to contain %q", pattern)
		}
	}
}

// TestAddRemote tests that AddRemote() adds a remote repository
func TestAddRemote(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Initialize repo first
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Add remote
	testURL := "https://github.com/test/repo.git"
	if err := syncer.AddRemote(testURL); err != nil {
		t.Fatalf("AddRemote() failed: %v", err)
	}

	// Verify remote URL
	url, err := syncer.GetRemoteURL("origin")
	if err != nil {
		t.Fatalf("GetRemoteURL() failed: %v", err)
	}

	if url != testURL {
		t.Errorf("expected remote URL %q, got %q", testURL, url)
	}
}

// TestRemoveRemote tests that RemoveRemote() removes a remote
func TestRemoveRemote(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Initialize repo and add remote
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	testURL := "https://github.com/test/repo.git"
	if err := syncer.AddRemote(testURL); err != nil {
		t.Fatalf("AddRemote() failed: %v", err)
	}

	// Remove remote
	if err := syncer.RemoveRemote("origin"); err != nil {
		t.Fatalf("RemoveRemote() failed: %v", err)
	}

	// Verify remote is gone (GetRemoteURL should fail or return empty)
	url, _ := syncer.GetRemoteURL("origin")
	if url != "" {
		t.Errorf("expected remote URL to be removed, got %q", url)
	}
}

// TestCommit tests that Commit() stages and commits changes
func TestCommit(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Initialize repo
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Commit changes
	if err := syncer.Commit("Test commit"); err != nil {
		t.Fatalf("Commit() failed: %v", err)
	}

	// Check status is clean (all changes committed)
	status, err := syncer.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if !status.Clean {
		t.Errorf("expected status to be clean after commit")
	}
}

// TestCommitNoChanges tests that Commit() is no-op when nothing to commit
func TestCommitNoChanges(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Initialize repo
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Commit with nothing staged
	if err := syncer.Commit("Empty commit"); err != nil {
		t.Fatalf("Commit() should not fail when nothing to commit: %v", err)
	}
}

// TestStatus tests that Status() returns correct sync status
func TestStatus(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Check status before init
	status, err := syncer.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if status.Initialized {
		t.Errorf("expected Initialized to be false before Init()")
	}

	// Initialize repo
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check status after init
	status, err = syncer.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if !status.Initialized {
		t.Errorf("expected Initialized to be true after Init()")
	}

	if status.RemoteConfigured {
		t.Errorf("expected RemoteConfigured to be false without remote")
	}
}

// TestStatusWithRemote tests that Status() detects remote configuration
func TestStatusWithRemote(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Initialize repo
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Add remote
	testURL := "https://github.com/test/repo.git"
	if err := syncer.AddRemote(testURL); err != nil {
		t.Fatalf("AddRemote() failed: %v", err)
	}

	// Check status with remote
	status, err := syncer.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if !status.RemoteConfigured {
		t.Errorf("expected RemoteConfigured to be true after AddRemote()")
	}
}

// TestStatusDirtyWorking tests that Status() detects uncommitted changes
func TestStatusDirtyWorking(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Initialize repo
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create a test file (untracked)
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Check status is dirty
	status, err := syncer.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if status.Clean {
		t.Errorf("expected Clean to be false with untracked changes")
	}
}

// TestStatusClean tests that Status() shows clean working directory after commit
func TestStatusClean(t *testing.T) {
	tmpDir, syncer := setupTestRepo(t)
	defer cleanupTestRepo(t, tmpDir)

	// Initialize repo
	if err := syncer.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Create and commit a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if err := syncer.Commit("Add test file"); err != nil {
		t.Fatalf("Commit() failed: %v", err)
	}

	// Check status is clean
	status, err := syncer.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if !status.Clean {
		t.Errorf("expected Clean to be true after commit")
	}

	if !status.Initialized {
		t.Errorf("expected Initialized to be true after init and commit")
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != "" && (s == substr || len(substr) > 0 && findSubstring(s, substr))
}

// findSubstring searches for substring in string
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
