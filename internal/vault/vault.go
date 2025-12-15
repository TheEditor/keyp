package vault

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/TheEditor/keyp/internal/model"
	"github.com/TheEditor/keyp/internal/store"
)

var (
	ErrLocked        = errors.New("vault is locked")
	ErrAlreadyExists = errors.New("vault already exists")
	ErrNotExists     = errors.New("vault does not exist")
)

// Vault manages the secret store lifecycle
type Vault struct {
	path   string
	store  *store.Store
	locked bool
}

// DefaultPath returns the default vault path (~/.keyp/vault.db)
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".keyp", "vault.db")
}

// Exists checks if a vault file exists at path
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Init creates a new vault at the specified path
func Init(path string) (*Vault, error) {
	if Exists(path) {
		return nil, ErrAlreadyExists
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	// Create and initialize store (creates schema)
	s, err := store.Open(path)
	if err != nil {
		return nil, err
	}

	return &Vault{
		path:   path,
		store:  s,
		locked: false,
	}, nil
}

// Open opens an existing vault
func Open(path string) (*Vault, error) {
	if !Exists(path) {
		return nil, ErrNotExists
	}

	s, err := store.Open(path)
	if err != nil {
		return nil, err
	}

	return &Vault{
		path:   path,
		store:  s,
		locked: false,
	}, nil
}

// Close closes the vault
func (v *Vault) Close() error {
	if v.store != nil {
		err := v.store.Close()
		v.store = nil
		v.locked = true
		return err
	}
	return nil
}

// Lock closes the store and marks vault as locked
func (v *Vault) Lock() error {
	return v.Close()
}

// IsLocked returns true if vault is locked
func (v *Vault) IsLocked() bool {
	return v.locked || v.store == nil
}

// Create adds a new secret to the vault
func (v *Vault) Create(secret *model.SecretObject) error {
	if v.IsLocked() {
		return ErrLocked
	}
	return v.store.Create(secret)
}

// GetByName retrieves a secret by name
func (v *Vault) GetByName(name string) (*model.SecretObject, error) {
	if v.IsLocked() {
		return nil, ErrLocked
	}
	return v.store.GetByName(name)
}

// List returns all secrets
func (v *Vault) List() ([]*model.SecretObject, error) {
	if v.IsLocked() {
		return nil, ErrLocked
	}
	return v.store.List()
}

// Search performs full-text search
func (v *Vault) Search(query string) ([]*model.SecretObject, error) {
	if v.IsLocked() {
		return nil, ErrLocked
	}
	return v.store.Search(query)
}

// Delete removes a secret
func (v *Vault) Delete(name string) error {
	if v.IsLocked() {
		return ErrLocked
	}
	return v.store.Delete(name)
}

// Path returns the vault file path
func (v *Vault) Path() string {
	return v.path
}
