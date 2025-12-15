package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize       = 32
	IVSize         = 12  // 96 bits for GCM
	KeySize        = 32  // 256 bits
	MinIterations  = 100000
)

// EncryptionResult contains all values needed for decryption
type EncryptionResult struct {
	Ciphertext string `json:"ciphertext"`  // base64
	AuthTag    string `json:"authTag"`     // base64
	IV         string `json:"iv"`          // base64
	Salt       string `json:"salt"`        // base64
}

// DeriveKey derives a 256-bit key from password using PBKDF2-SHA256
func DeriveKey(password string, salt []byte, iterations int) ([]byte, error) {
	if iterations < MinIterations {
		return nil, errors.New("iterations must be at least 100000")
	}
	if len(salt) != SaltSize {
		return nil, errors.New("salt must be 32 bytes")
	}
	key := pbkdf2.Key([]byte(password), salt, iterations, KeySize, sha256.New)
	return key, nil
}

// Encrypt encrypts plaintext using AES-256-GCM with PBKDF2-derived key
func Encrypt(plaintext, password string, iterations int) (*EncryptionResult, error) {
	// Generate random salt
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive key
	key, err := DeriveKey(password, salt, iterations)
	if err != nil {
		return nil, err
	}

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate random IV
	iv := make([]byte, IVSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// Encrypt (GCM appends auth tag to ciphertext)
	sealed := gcm.Seal(nil, iv, []byte(plaintext), nil)

	// Split ciphertext and auth tag
	tagStart := len(sealed) - gcm.Overhead()
	ciphertext := sealed[:tagStart]
	authTag := sealed[tagStart:]

	return &EncryptionResult{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		AuthTag:    base64.StdEncoding.EncodeToString(authTag),
		IV:         base64.StdEncoding.EncodeToString(iv),
		Salt:       base64.StdEncoding.EncodeToString(salt),
	}, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM with PBKDF2-derived key
func Decrypt(result *EncryptionResult, password string, iterations int) (string, error) {
	// Decode base64 values
	ciphertext, err := base64.StdEncoding.DecodeString(result.Ciphertext)
	if err != nil {
		return "", errors.New("invalid ciphertext encoding")
	}
	authTag, err := base64.StdEncoding.DecodeString(result.AuthTag)
	if err != nil {
		return "", errors.New("invalid authTag encoding")
	}
	iv, err := base64.StdEncoding.DecodeString(result.IV)
	if err != nil {
		return "", errors.New("invalid IV encoding")
	}
	salt, err := base64.StdEncoding.DecodeString(result.Salt)
	if err != nil {
		return "", errors.New("invalid salt encoding")
	}

	// Derive key
	key, err := DeriveKey(password, salt, iterations)
	if err != nil {
		return "", err
	}

	// Create cipher
	block, err := aes.NewCipher(key)
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
		return "", errors.New("decryption failed: invalid password or corrupted data")
	}

	return string(plaintext), nil
}
