package core

import (
	"testing"
)

func TestDeriveKey(t *testing.T) {
	salt := make([]byte, SaltSize)
	for i := range salt {
		salt[i] = byte(i)
	}

	key, err := DeriveKey("testpassword", salt, MinIterations)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	if len(key) != KeySize {
		t.Errorf("Expected key length %d, got %d", KeySize, len(key))
	}

	// Same inputs should produce same key
	key2, _ := DeriveKey("testpassword", salt, MinIterations)
	for i := range key {
		if key[i] != key2[i] {
			t.Error("Same inputs produced different keys")
			break
		}
	}

	// Different password should produce different key
	key3, _ := DeriveKey("differentpassword", salt, MinIterations)
	same := true
	for i := range key {
		if key[i] != key3[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("Different passwords produced same key")
	}
}

func TestDeriveKeyValidation(t *testing.T) {
	salt := make([]byte, SaltSize)

	// Too few iterations
	_, err := DeriveKey("password", salt, 1000)
	if err == nil {
		t.Error("Expected error for low iterations")
	}

	// Wrong salt size
	_, err = DeriveKey("password", []byte("short"), MinIterations)
	if err == nil {
		t.Error("Expected error for wrong salt size")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
	}{
		{"simple", "hello world"},
		{"empty", ""},
		{"unicode", "hello world 🔐"},
		{"json", `{"key": "value", "nested": {"a": 1}}`},
		{"long", string(make([]byte, 10000))},
	}

	password := "test-master-password"
	iterations := MinIterations

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Encrypt(tc.plaintext, password, iterations)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			decrypted, err := Decrypt(result, password, iterations)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if decrypted != tc.plaintext {
				t.Errorf("Round trip failed: got %q, want %q", decrypted, tc.plaintext)
			}
		})
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	result, _ := Encrypt("secret data", "correct-password", MinIterations)

	_, err := Decrypt(result, "wrong-password", MinIterations)
	if err == nil {
		t.Error("Expected error when decrypting with wrong password")
	}
}

func TestEncryptProducesUniqueOutput(t *testing.T) {
	plaintext := "same plaintext"
	password := "same password"

	result1, _ := Encrypt(plaintext, password, MinIterations)
	result2, _ := Encrypt(plaintext, password, MinIterations)

	// Salt should differ
	if result1.Salt == result2.Salt {
		t.Error("Expected different salts")
	}

	// IV should differ
	if result1.IV == result2.IV {
		t.Error("Expected different IVs")
	}

	// Ciphertext should differ
	if result1.Ciphertext == result2.Ciphertext {
		t.Error("Expected different ciphertexts")
	}

	// But both should decrypt to same plaintext
	d1, _ := Decrypt(result1, password, MinIterations)
	d2, _ := Decrypt(result2, password, MinIterations)
	if d1 != d2 || d1 != plaintext {
		t.Error("Both should decrypt to original plaintext")
	}
}
