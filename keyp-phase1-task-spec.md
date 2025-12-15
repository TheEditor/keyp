# keyp Phase 1: Go Scaffold + Crypto Foundation

## Preface

**For the implementing assistant (Haiku)**: This phase establishes the Go project structure and implements the cryptographic foundation. The TypeScript code is archived for reference only—do not port tests or attempt compatibility.

**Before starting:**
1. Read the bd-issue-tracking skill: `view /mnt/skills/user/bd-issue-tracking/SKILL.md`
2. Read AGENTS.md in the repo root
3. Execute Beads commands below in order, capturing returned IDs
4. Work through tasks—check `bd ready --json` to see what's available
5. Close each task IMMEDIATELY after completing it
6. When `bd ready --json` returns empty, phase is complete

---

## Beads Issue Setup

Execute these commands in order. **Capture the returned ID from each command.**

### STEP 1: Create Parent Epic

```bash
bd create "Phase 1: Go scaffold and crypto foundation" -t epic -p 0 -d "Establish Go project structure, archive TypeScript code, implement AES-256-GCM encryption with PBKDF2 key derivation, and create CLI skeleton." --json
```

### STEP 2: Archive TypeScript Code

```bash
bd create "Archive TypeScript code to legacy folder" -t task -p 0 --parent <EPIC_ID> -d "Move existing TypeScript code to legacy/ for reference.

STEPS:
1. Create legacy/ directory in repo root
2. Move these items INTO legacy/:
   - src/
   - bin/
   - completions/
   - docs/
   - package.json
   - tsconfig.json
   - .npmignore
   - PUBLISH.md
   - CHANGELOG.md
3. Keep in repo root (do NOT move):
   - README.md (will be rewritten)
   - LICENSE
   - .gitignore (will be updated)
   - .github/ (keep for now)
   - AGENTS.md
   - .beads/

4. Update .gitignore - append these lines:
   # Go
   keyp
   keyp.exe
   *.test
   coverage.out
   
   # Legacy
   legacy/node_modules/
   legacy/lib/

5. Commit: git commit -m \"chore: archive TypeScript code to legacy/ (bd:<TASK_ID>)\"

ACCEPTANCE: 
- legacy/ contains all TS source
- go.mod can be created in repo root without conflicts
- .gitignore updated for Go" --json
```

### STEP 3: Initialize Go Module

```bash
bd create "Initialize Go module and directory structure" -t task -p 0 --parent <EPIC_ID> -d "Create Go module and establish package structure.

STEPS:
1. In repo root, run:
   go mod init github.com/TheEditor/keyp

2. Create directory structure:
   mkdir -p cmd/keyp
   mkdir -p internal/core
   mkdir -p internal/model
   mkdir -p internal/store
   mkdir -p internal/vault

3. Create placeholder main.go:

// cmd/keyp/main.go
package main

import \"fmt\"

func main() {
    fmt.Println(\"keyp v2.0.0-dev\")
}

4. Verify build:
   go build -o keyp ./cmd/keyp
   ./keyp

5. Commit: git commit -m \"feat: initialize Go module structure (bd:<TASK_ID>)\"

ACCEPTANCE:
- go.mod exists with module github.com/TheEditor/keyp
- Directory structure matches AGENTS.md
- Binary builds and runs" --json
```

### STEP 4: Implement Crypto Module

```bash
bd create "Implement AES-256-GCM encryption with PBKDF2" -t task -p 0 --parent <EPIC_ID> -d "Implement core cryptographic functions matching security requirements.

CREATE FILE: internal/core/crypto.go

package core

import (
    \"crypto/aes\"
    \"crypto/cipher\"
    \"crypto/rand\"
    \"crypto/sha256\"
    \"encoding/base64\"
    \"errors\"
    \"golang.org/x/crypto/pbkdf2\"
)

const (
    SaltSize       = 32
    IVSize         = 12  // 96 bits for GCM
    KeySize        = 32  // 256 bits
    MinIterations  = 100000
)

// EncryptionResult contains all values needed for decryption
type EncryptionResult struct {
    Ciphertext string \`json:\"ciphertext\"\`  // base64
    AuthTag    string \`json:\"authTag\"\`     // base64
    IV         string \`json:\"iv\"\`          // base64
    Salt       string \`json:\"salt\"\`        // base64
}

// DeriveKey derives a 256-bit key from password using PBKDF2-SHA256
func DeriveKey(password string, salt []byte, iterations int) ([]byte, error) {
    if iterations < MinIterations {
        return nil, errors.New(\"iterations must be at least 100000\")
    }
    if len(salt) != SaltSize {
        return nil, errors.New(\"salt must be 32 bytes\")
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
        return \"\", errors.New(\"invalid ciphertext encoding\")
    }
    authTag, err := base64.StdEncoding.DecodeString(result.AuthTag)
    if err != nil {
        return \"\", errors.New(\"invalid authTag encoding\")
    }
    iv, err := base64.StdEncoding.DecodeString(result.IV)
    if err != nil {
        return \"\", errors.New(\"invalid IV encoding\")
    }
    salt, err := base64.StdEncoding.DecodeString(result.Salt)
    if err != nil {
        return \"\", errors.New(\"invalid salt encoding\")
    }

    // Derive key
    key, err := DeriveKey(password, salt, iterations)
    if err != nil {
        return \"\", err
    }

    // Create cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        return \"\", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return \"\", err
    }

    // Reconstruct sealed message (ciphertext + authTag)
    sealed := append(ciphertext, authTag...)

    // Decrypt
    plaintext, err := gcm.Open(nil, iv, sealed, nil)
    if err != nil {
        return \"\", errors.New(\"decryption failed: invalid password or corrupted data\")
    }

    return string(plaintext), nil
}

THEN run: go mod tidy

ACCEPTANCE:
- File compiles: go build ./internal/core/
- Functions exported: DeriveKey, Encrypt, Decrypt
- Constants defined: SaltSize, IVSize, KeySize, MinIterations" --json
```

### STEP 5: Add Crypto Tests

```bash
bd create "Add comprehensive crypto unit tests" -t task -p 1 --parent <EPIC_ID> -d "Create unit tests for crypto module.

CREATE FILE: internal/core/crypto_test.go

package core

import (
    \"testing\"
)

func TestDeriveKey(t *testing.T) {
    salt := make([]byte, SaltSize)
    for i := range salt {
        salt[i] = byte(i)
    }

    key, err := DeriveKey(\"testpassword\", salt, MinIterations)
    if err != nil {
        t.Fatalf(\"DeriveKey failed: %v\", err)
    }

    if len(key) != KeySize {
        t.Errorf(\"Expected key length %d, got %d\", KeySize, len(key))
    }

    // Same inputs should produce same key
    key2, _ := DeriveKey(\"testpassword\", salt, MinIterations)
    for i := range key {
        if key[i] != key2[i] {
            t.Error(\"Same inputs produced different keys\")
            break
        }
    }

    // Different password should produce different key
    key3, _ := DeriveKey(\"differentpassword\", salt, MinIterations)
    same := true
    for i := range key {
        if key[i] != key3[i] {
            same = false
            break
        }
    }
    if same {
        t.Error(\"Different passwords produced same key\")
    }
}

func TestDeriveKeyValidation(t *testing.T) {
    salt := make([]byte, SaltSize)

    // Too few iterations
    _, err := DeriveKey(\"password\", salt, 1000)
    if err == nil {
        t.Error(\"Expected error for low iterations\")
    }

    // Wrong salt size
    _, err = DeriveKey(\"password\", []byte(\"short\"), MinIterations)
    if err == nil {
        t.Error(\"Expected error for wrong salt size\")
    }
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
    testCases := []struct {
        name      string
        plaintext string
    }{
        {\"simple\", \"hello world\"},
        {\"empty\", \"\"},
        {\"unicode\", \"hello world 🔐\"},
        {\"json\", \`{\"key\": \"value\", \"nested\": {\"a\": 1}}\`},
        {\"long\", string(make([]byte, 10000))},
    }

    password := \"test-master-password\"
    iterations := MinIterations

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := Encrypt(tc.plaintext, password, iterations)
            if err != nil {
                t.Fatalf(\"Encrypt failed: %v\", err)
            }

            decrypted, err := Decrypt(result, password, iterations)
            if err != nil {
                t.Fatalf(\"Decrypt failed: %v\", err)
            }

            if decrypted != tc.plaintext {
                t.Errorf(\"Round trip failed: got %q, want %q\", decrypted, tc.plaintext)
            }
        })
    }
}

func TestDecryptWrongPassword(t *testing.T) {
    result, _ := Encrypt(\"secret data\", \"correct-password\", MinIterations)

    _, err := Decrypt(result, \"wrong-password\", MinIterations)
    if err == nil {
        t.Error(\"Expected error when decrypting with wrong password\")
    }
}

func TestEncryptProducesUniqueOutput(t *testing.T) {
    plaintext := \"same plaintext\"
    password := \"same password\"

    result1, _ := Encrypt(plaintext, password, MinIterations)
    result2, _ := Encrypt(plaintext, password, MinIterations)

    // Salt should differ
    if result1.Salt == result2.Salt {
        t.Error(\"Expected different salts\")
    }

    // IV should differ
    if result1.IV == result2.IV {
        t.Error(\"Expected different IVs\")
    }

    // Ciphertext should differ
    if result1.Ciphertext == result2.Ciphertext {
        t.Error(\"Expected different ciphertexts\")
    }

    // But both should decrypt to same plaintext
    d1, _ := Decrypt(result1, password, MinIterations)
    d2, _ := Decrypt(result2, password, MinIterations)
    if d1 != d2 || d1 != plaintext {
        t.Error(\"Both should decrypt to original plaintext\")
    }
}

RUN: go test ./internal/core/ -v

ACCEPTANCE:
- All tests pass
- Tests cover: key derivation, validation, round-trip, wrong password, uniqueness

Commit: git commit -m \"test(core): add crypto unit tests (bd:<TASK_ID>)\"" --json
```

### STEP 6: Create CLI Skeleton with Cobra

```bash
bd create "Create CLI skeleton with cobra" -t task -p 1 --parent <EPIC_ID> -d "Set up cobra CLI framework with version and help commands.

STEP 1: Add cobra dependency
go get github.com/spf13/cobra@v1.8.0
go mod tidy

STEP 2: Replace cmd/keyp/main.go with:

package main

import (
    \"fmt\"
    \"os\"

    \"github.com/spf13/cobra\"
)

var version = \"2.0.0-dev\"

var rootCmd = &cobra.Command{
    Use:   \"keyp\",
    Short: \"Local-first secret manager\",
    Long:  \`keyp is a local-first secret manager for developers and families.
Securely store structured secrets with full-text search.\`,
}

var versionCmd = &cobra.Command{
    Use:   \"version\",
    Short: \"Print version information\",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf(\"keyp v%s\\n\", version)
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

STEP 3: Build and test
go build -o keyp ./cmd/keyp
./keyp --help
./keyp version

ACCEPTANCE:
- keyp --help shows usage with description
- keyp version prints \"keyp v2.0.0-dev\"
- Unknown commands show error

Commit: git commit -m \"feat(cli): add cobra CLI skeleton (bd:<TASK_ID>)\"" --json
```

### STEP 7: Create Makefile

```bash
bd create "Create Makefile for common operations" -t task -p 2 --parent <EPIC_ID> -d "Add Makefile for build, test, and clean operations.

CREATE FILE: Makefile

.PHONY: build test clean

# Build binary
build:
	go build -o keyp ./cmd/keyp

# Run all tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# Clean build artifacts
clean:
	rm -f keyp keyp.exe coverage.out

# Install locally
install: build
	cp keyp /usr/local/bin/keyp

VERIFY:
make build
make test
make clean

ACCEPTANCE:
- make build produces keyp binary
- make test runs all tests
- make clean removes artifacts

Commit: git commit -m \"chore: add Makefile (bd:<TASK_ID>)\"" --json
```

### STEP 8: Add Dependencies

After capturing all task IDs from steps 2-7, add dependencies:

```bash
bd dep add <TASK_3_ID> <TASK_2_ID> --type blocks
bd dep add <TASK_4_ID> <TASK_3_ID> --type blocks
bd dep add <TASK_5_ID> <TASK_4_ID> --type blocks
bd dep add <TASK_6_ID> <TASK_3_ID> --type blocks
bd dep add <TASK_7_ID> <TASK_6_ID> --type blocks
```

Dependency graph:
```
Archive TS (Task 2)
    └── Init Go (Task 3)
            ├── Crypto (Task 4)
            │       └── Crypto Tests (Task 5)
            └── CLI Skeleton (Task 6)
                    └── Makefile (Task 7)
```

### STEP 9: Verify Setup

```bash
bd ready --json
bd dep tree <TASK_7_ID>
```

---

## Working Through Tasks

After creating all issues:

1. Run `bd ready --json` to see available tasks (should show Archive TS first)
2. Complete the task per its description
3. Run quality checks: `go build ./...` and `go test ./...`
4. Close immediately: `bd close <TASK_ID> --reason "..."`
5. Repeat until `bd ready --json` returns empty

When all tasks closed, close the epic:

```bash
bd close <EPIC_ID> --reason "Phase 1 complete: Go scaffold and crypto foundation"
```

---

## Completion

Work is complete when:
- [ ] `bd ready --json` returns no issues
- [ ] `make build` succeeds
- [ ] `make test` passes
- [ ] `./keyp version` prints version
- [ ] legacy/ contains archived TypeScript

Final commit:
```bash
git add .
git commit -m "feat: complete Phase 1 Go migration (bd:<EPIC_ID>)"
git push
```
