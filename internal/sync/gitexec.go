package sync

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitExecSyncer implements Syncer using exec.Command to call git binary
type GitExecSyncer struct {
	vaultPath string
}

// NewGitExecSyncer creates a new git-based syncer for a vault path
func NewGitExecSyncer(vaultPath string) *GitExecSyncer {
	return &GitExecSyncer{vaultPath: vaultPath}
}

// git runs a git command in the vault directory
func (g *GitExecSyncer) git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.vaultPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		return output, fmt.Errorf("git %s failed: %w (output: %s)", args[0], err, strings.TrimSpace(output))
	}

	return strings.TrimSpace(output), nil
}

// Init initializes a git repository in the vault directory and creates .gitignore
func (g *GitExecSyncer) Init() error {
	_, err := g.git("init")
	if err != nil {
		return err
	}

	// Create .gitignore to exclude SQLite database files
	gitignorePath := filepath.Join(g.vaultPath, ".gitignore")
	gitignoreContent := "# Exclude SQLite database files\n*.db\n*.db-journal\n*.db-wal\n*.db-shm\n"
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// AddRemote adds a remote repository URL
func (g *GitExecSyncer) AddRemote(url string) error {
	_, err := g.git("remote", "add", "origin", url)
	return err
}

// RemoveRemote removes a remote by name
func (g *GitExecSyncer) RemoveRemote(name string) error {
	_, err := g.git("remote", "remove", name)
	return err
}

// GetRemoteURL returns the URL for a named remote
func (g *GitExecSyncer) GetRemoteURL(name string) (string, error) {
	url, err := g.git("remote", "get-url", name)
	if err != nil {
		return "", err
	}
	return url, nil
}

// Commit creates a commit with the given message
func (g *GitExecSyncer) Commit(message string) error {
	// First add all files
	if _, err := g.git("add", "."); err != nil {
		return err
	}

	// Check if there are changes to commit
	status, err := g.git("status", "--porcelain")
	if err != nil {
		return err
	}

	if status == "" {
		// No changes to commit
		return nil
	}

	// Commit with message
	_, err = g.git("commit", "-m", message)
	return err
}

// Push pushes commits to the remote repository
func (g *GitExecSyncer) Push() error {
	_, err := g.git("push", "-u", "origin", "main")
	return err
}

// Pull pulls changes from the remote repository
func (g *GitExecSyncer) Pull() error {
	_, err := g.git("pull")
	return err
}

// Status returns the current sync status
func (g *GitExecSyncer) Status() (*SyncStatus, error) {
	status := &SyncStatus{
		Initialized:      false,
		RemoteConfigured: false,
		Clean:            true,
		UnpushedCommits:  0,
		UnpulledCommits:  0,
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return status, fmt.Errorf("git not found: %w", err)
	}

	// Check if initialized
	if _, err := g.git("rev-parse", "--git-dir"); err == nil {
		status.Initialized = true
	}

	if !status.Initialized {
		return status, nil
	}

	// Check remote configuration
	if _, err := g.git("remote", "get-url", "origin"); err == nil {
		status.RemoteConfigured = true
	}

	// Check status (clean vs dirty)
	statusOutput, err := g.git("status", "--porcelain")
	if err == nil && statusOutput == "" {
		status.Clean = true
	} else {
		status.Clean = false
	}

	// Count unpushed commits (commits not in origin/main)
	if status.RemoteConfigured {
		unpushedOutput, err := g.git("log", "origin/main..HEAD", "--oneline")
		if err == nil && unpushedOutput != "" {
			status.UnpushedCommits = len(strings.Split(unpushedOutput, "\n"))
		}

		// Count unpulled commits (commits in origin/main not in HEAD)
		unpulledOutput, err := g.git("log", "HEAD..origin/main", "--oneline")
		if err == nil && unpulledOutput != "" {
			status.UnpulledCommits = len(strings.Split(unpulledOutput, "\n"))
		}
	}

	return status, nil
}

// Verify it implements the Syncer interface
var _ Syncer = (*GitExecSyncer)(nil)
