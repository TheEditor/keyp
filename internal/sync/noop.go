package sync

// NoopSyncer is a placeholder implementation that does nothing
// Used until actual git sync is implemented in Phase 3
type NoopSyncer struct{}

// NewNoop creates a no-op syncer
func NewNoop() Syncer {
	return &NoopSyncer{}
}

func (n *NoopSyncer) Init() error {
	return nil
}

func (n *NoopSyncer) AddRemote(url string) error {
	return nil
}

func (n *NoopSyncer) Commit(message string) error {
	return nil
}

func (n *NoopSyncer) Push() error {
	return nil
}

func (n *NoopSyncer) Pull() error {
	return nil
}

func (n *NoopSyncer) Status() (*SyncStatus, error) {
	return &SyncStatus{
		Initialized:      false,
		RemoteConfigured: false,
		Clean:            true,
		UnpushedCommits:  0,
		UnpulledCommits:  0,
	}, nil
}
