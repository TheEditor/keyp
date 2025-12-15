package sync

// Syncer defines the interface for vault synchronization
// This abstraction allows swapping implementations:
// - go-git library (feature-rich, many dependencies)
// - exec.Command to git binary (fewer deps, requires git installed)
type Syncer interface {
	// Init initializes git repo in vault directory
	Init() error

	// AddRemote configures the remote repository URL
	AddRemote(url string) error

	// Commit creates a commit with the given message
	Commit(message string) error

	// Push pushes commits to remote
	Push() error

	// Pull pulls changes from remote
	Pull() error

	// Status returns current sync status
	Status() (*SyncStatus, error)
}

// SyncStatus represents the current state of sync
type SyncStatus struct {
	Initialized      bool
	RemoteConfigured bool
	Clean            bool
	UnpushedCommits  int
	UnpulledCommits  int
}
