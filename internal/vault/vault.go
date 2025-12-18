package vault

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TheEditor/keyp/internal/core"
	"github.com/TheEditor/keyp/internal/model"
	"github.com/TheEditor/keyp/internal/store"
)

const verificationPlaintext = "keyp-vault-v1"

var (
	ErrLocked        = errors.New("vault is locked")
	ErrAlreadyExists = errors.New("vault already exists")
	ErrNotExists     = errors.New("vault does not exist")
)

// Vault manages the secret store lifecycle
type Vault struct {
	path   string
	store  *store.Store
	key    []byte
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

// Init creates a new vault at the specified path with password protection
func Init(path string, password string) (*Vault, error) {
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

	// Generate random salt
	salt := make([]byte, core.SaltSize)
	if _, err := rand.Read(salt); err != nil {
		s.Close()
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive encryption key
	key, err := core.DeriveKey(password, salt, core.MinIterations)
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Store encryption metadata
	if err := s.SetMeta("salt", base64.StdEncoding.EncodeToString(salt)); err != nil {
		s.Close()
		return nil, err
	}
	if err := s.SetMeta("iterations", fmt.Sprintf("%d", core.MinIterations)); err != nil {
		s.Close()
		return nil, err
	}

	// Create and store verification value (encrypted with the derived key)
	v := &Vault{
		path:   path,
		store:  s,
		key:    key,
		locked: false,
	}
	verifyEncrypted, err := v.encryptValue(verificationPlaintext)
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("failed to create verification value: %w", err)
	}
	if err := s.SetMeta("verify", verifyEncrypted); err != nil {
		s.Close()
		return nil, err
	}

	return v, nil
}

// Open opens an existing vault with password
func Open(path string, password string) (*Vault, error) {
	if !Exists(path) {
		return nil, ErrNotExists
	}

	s, err := store.Open(path)
	if err != nil {
		return nil, err
	}

	// Read encryption parameters
	saltB64, err := s.GetMeta("salt")
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("failed to read vault metadata: %w", err)
	}
	iterStr, err := s.GetMeta("iterations")
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("failed to read vault metadata: %w", err)
	}

	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("corrupted vault metadata: %w", err)
	}
	iterations, err := strconv.Atoi(iterStr)
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("corrupted vault metadata: %w", err)
	}

	// Derive key from password
	key, err := core.DeriveKey(password, salt, iterations)
	if err != nil {
		s.Close()
		return nil, err
	}

	// Create temporary vault to verify password by decrypting verification value
	v := &Vault{
		path:   path,
		store:  s,
		key:    key,
		locked: false,
	}

	// Read and decrypt verification value
	verifyEncrypted, err := s.GetMeta("verify")
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("failed to read verification value: %w", err)
	}

	decrypted, err := v.decryptValue(verifyEncrypted)
	if err != nil || decrypted != verificationPlaintext {
		s.Close()
		return nil, store.ErrInvalidPassword
	}

	return v, nil
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
func (v *Vault) Create(ctx context.Context, secret *model.SecretObject) error {
	if v.IsLocked() {
		return ErrLocked
	}
	// Encrypt sensitive field values before storage
	encrypted := v.encryptSecret(secret)
	return v.store.Create(ctx, encrypted)
}

// GetByName retrieves a secret by name
func (v *Vault) GetByName(ctx context.Context, name string) (*model.SecretObject, error) {
	if v.IsLocked() {
		return nil, ErrLocked
	}
	secret, err := v.store.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	// Decrypt sensitive field values
	return v.decryptSecret(secret)
}

// List returns all secrets
func (v *Vault) List(ctx context.Context, opts *store.SearchOptions) ([]*model.SecretObject, error) {
	if v.IsLocked() {
		return nil, ErrLocked
	}
	secrets, err := v.store.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	// Decrypt all secrets
	var decrypted []*model.SecretObject
	for _, s := range secrets {
		d, err := v.decryptSecret(s)
		if err != nil {
			return nil, err
		}
		decrypted = append(decrypted, d)
	}
	return decrypted, nil
}

// Search performs full-text search
func (v *Vault) Search(ctx context.Context, query string, opts *store.SearchOptions) ([]*model.SecretObject, error) {
	if v.IsLocked() {
		return nil, ErrLocked
	}
	secrets, err := v.store.Search(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	// Decrypt all secrets
	var decrypted []*model.SecretObject
	for _, s := range secrets {
		d, err := v.decryptSecret(s)
		if err != nil {
			return nil, err
		}
		decrypted = append(decrypted, d)
	}
	return decrypted, nil
}

// Update updates an existing secret
func (v *Vault) Update(ctx context.Context, secret *model.SecretObject) error {
	if v.IsLocked() {
		return ErrLocked
	}
	// Encrypt sensitive field values before storage
	encrypted := v.encryptSecret(secret)
	return v.store.Update(ctx, encrypted)
}

// Delete removes a secret
func (v *Vault) Delete(ctx context.Context, name string) error {
	if v.IsLocked() {
		return ErrLocked
	}
	return v.store.Delete(ctx, name)
}

// Path returns the vault file path
func (v *Vault) Path() string {
	return v.path
}

// encryptSecret encrypts sensitive field values in a secret
func (v *Vault) encryptSecret(secret *model.SecretObject) *model.SecretObject {
	copy := *secret
	copy.Fields = make([]model.Field, len(secret.Fields))
	for i, f := range secret.Fields {
		copy.Fields[i] = f
		if f.Sensitive {
			// Encrypt using the derived vault key
			encrypted, err := v.encryptValue(f.Value)
			if err != nil {
				// If encryption fails, keep original (shouldn't happen in normal operation)
				continue
			}
			copy.Fields[i].Value = encrypted
		}
	}
	return &copy
}

// decryptSecret decrypts sensitive field values in a secret
func (v *Vault) decryptSecret(secret *model.SecretObject) (*model.SecretObject, error) {
	copy := *secret
	copy.Fields = make([]model.Field, len(secret.Fields))
	for i, f := range secret.Fields {
		copy.Fields[i] = f
		if f.Sensitive {
			// Decrypt using the derived vault key
			decrypted, err := v.decryptValue(f.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt field %q: %w", f.Label, err)
			}
			copy.Fields[i].Value = decrypted
		}
	}
	return &copy, nil
}

// encryptValue encrypts a single value using the vault's derived key
func (v *Vault) encryptValue(plaintext string) (string, error) {
	// Create a cipher block from the derived key
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate random IV
	iv := make([]byte, core.IVSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// Encrypt (GCM appends auth tag to ciphertext)
	sealed := gcm.Seal(nil, iv, []byte(plaintext), nil)

	// Split ciphertext and auth tag
	tagStart := len(sealed) - gcm.Overhead()
	ciphertext := sealed[:tagStart]
	authTag := sealed[tagStart:]

	// Return as JSON-like string: "iv:ciphertext:authTag" (all base64)
	result := fmt.Sprintf("%s:%s:%s",
		base64.StdEncoding.EncodeToString(iv),
		base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(authTag),
	)
	return result, nil
}

// decryptValue decrypts a single value using the vault's derived key
func (v *Vault) decryptValue(encrypted string) (string, error) {
	// Parse the format: "iv:ciphertext:authTag"
	parts := strings.Split(encrypted, ":")
	if len(parts) != 3 {
		return "", errors.New("invalid encrypted value format")
	}

	iv, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.New("invalid IV encoding")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("invalid ciphertext encoding")
	}
	authTag, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return "", errors.New("invalid authTag encoding")
	}

	// Create a cipher block from the derived key
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Reconstruct sealed message (ciphertext + authTag)
	sealed := append(ciphertext, authTag...)

	// Decrypt
	plaintext, err := gcm.Open(nil, iv, sealed, nil)
	if err != nil {
		return "", store.ErrInvalidPassword
	}

	return string(plaintext), nil
}
