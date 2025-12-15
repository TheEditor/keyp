# Agent Instructions for keyp

## Issue Tracking with bd (beads)

**IMPORTANT**: This project uses **bd (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why bd?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Auto-syncs to JSONL for version control
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**Check for ready work:**
```bash
bd ready --json
```

**Create new issues:**
```bash
bd create "Issue title" -t bug|feature|task -p 0-4 --json
bd create "Issue title" -p 1 --deps discovered-from:bd-123 --json
bd create "Subtask" --parent <epic-id> --json  # Hierarchical subtask
```

**Claim and update:**
```bash
bd update bd-42 --status in_progress --json
bd update bd-42 --priority 1 --json
```

**Complete work:**
```bash
bd close bd-42 --reason "Completed" --json
```

### Workflow for Agents

1. **Check ready work**: `bd ready --json` shows unblocked issues
2. **Claim your task**: `bd update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `bd create "Found bug" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `bd close <id> --reason "Done"`
6. **Commit together**: Always commit the `.beads/issues.jsonl` file together with code changes

### Important Rules

- ✅ Use bd for ALL task tracking
- ✅ Always use `--json` flag for programmatic use
- ✅ Link discovered work with `discovered-from` dependencies
- ✅ Check `bd ready` before asking "what should I work on?"
- ✅ Run `bd <cmd> --help` to discover available flags
- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT use external issue trackers
- ❌ Do NOT duplicate tracking systems

---

## Project Overview

**keyp** is a local-first secret manager transitioning from a developer-focused CLI to a foundation for a family-friendly GUI application.

**Primary Goal**: Migrate from TypeScript/Node.js to Go, replacing the flat key-value data model with structured secrets backed by SQLCipher-encrypted SQLite.

**Tagline**: *"pass for the Node.js generation"* → evolving toward *"secrets for the whole family"*

### Current State

- **v1.x (TypeScript)**: Published to npm as `@theeditor/keyp`. Functional CLI with AES-256-GCM encryption, PBKDF2 key derivation, and Git sync. Flat `key: value` data model.
- **v2.x (Go)**: In development. Structured secrets, SQLCipher storage, full-text search, HTTP server mode for future GUI.

### Key Decisions (Already Made)

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Language | Go | Single-binary distribution, easy cross-compilation |
| Storage | SQLCipher | Whole-DB encryption, battle-tested, full FTS on encrypted data |
| Organization | Tags only | Flexible, flat, search-friendly; folders can come later |
| Data model | Structured secrets | Multiple fields per secret (the AT&T problem) |
| Release | Immediate v2 | No meaningful user base to migrate |
| CLI compatibility | Keep `keyp set` | Simple cases remain simple |

---

## Project Structure (Go v2)

```
keyp/
├── cmd/
│   └── keyp/
│       └── main.go           # CLI entry point (cobra)
├── internal/
│   ├── core/                 # Encryption, key derivation
│   │   ├── crypto.go         # AES-256-GCM, PBKDF2
│   │   └── crypto_test.go
│   ├── store/                # SQLCipher database operations
│   │   ├── store.go          # CRUD, FTS queries
│   │   ├── schema.sql        # Table definitions
│   │   └── store_test.go
│   ├── model/                # Domain types
│   │   ├── secret.go         # SecretObject, Field
│   │   └── secret_test.go
│   ├── vault/                # Vault lifecycle (open, lock, save)
│   │   ├── vault.go
│   │   └── vault_test.go
│   └── sync/                 # Git sync operations
│       ├── git.go
│       └── git_test.go
├── pkg/
│   └── keyp/                 # Public Go API (for embedding)
├── legacy/                   # Archived TypeScript code (reference only)
├── npm/                      # npm wrapper package
│   ├── package.json
│   ├── bin/keyp.js           # Binary downloader/executor
│   └── install.js            # Postinstall hook
├── .beads/                   # Issue tracker database
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── AGENTS.md                 # This file
```

---

## Data Model

### SecretObject

```go
type SecretObject struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`       // "AT&T" - required, searchable
    Tags      []string  `json:"tags"`       // ["telecom", "family"]
    Fields    []Field   `json:"fields"`     // Ordered, user-defined
    Notes     string    `json:"notes"`      // Free-form, searchable
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Field struct {
    Label     string `json:"label"`     // "Account PIN"
    Value     string `json:"value"`     // Encrypted at rest
    Sensitive bool   `json:"sensitive"` // Mask in UI?
    Type      string `json:"type"`      // "password", "pin", "text", "url"
}
```

### SQLite Schema (with SQLCipher)

```sql
CREATE TABLE secrets (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    tags TEXT DEFAULT '[]',    -- JSON array
    notes TEXT,
    created_at TEXT,
    updated_at TEXT
);

CREATE TABLE fields (
    id TEXT PRIMARY KEY,
    secret_id TEXT NOT NULL REFERENCES secrets(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    value TEXT NOT NULL,       -- Encrypted by SQLCipher
    sensitive INTEGER DEFAULT 1,
    type TEXT DEFAULT 'text',
    sort_order INTEGER,
    UNIQUE(secret_id, label)
);

-- Full-text search
CREATE VIRTUAL TABLE secrets_fts USING fts5(
    name, tags, notes,
    content='secrets',
    content_rowid='rowid'
);
```

---

## CLI Commands

### Preserved from v1

```bash
keyp init                     # Create vault
keyp set <name> [value]       # Quick single-field secret
keyp get <name>               # Copy to clipboard
keyp list                     # List all secrets
keyp delete <name>            # Remove secret
keyp sync push|pull           # Git sync
```

### New in v2

```bash
keyp add <name>               # Interactive multi-field creation
keyp show <name>              # Display all fields
keyp edit <name>              # Modify fields
keyp search <query>           # Full-text search
keyp tag add|rm <name> <tag>  # Manage tags
keyp serve                    # HTTP API (future GUI)
```

### Compatibility Note

`keyp set github-token abc123` creates a SecretObject with a single field:
- Name: "github-token"
- Fields: [{ Label: "value", Value: "abc123", Sensitive: true }]

This preserves the simple mental model while enabling structured storage.

---

## Build & Test

```bash
# Build
go build -o keyp ./cmd/keyp

# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Cross-compile (example)
GOOS=darwin GOARCH=arm64 go build -o keyp-darwin-arm64 ./cmd/keyp
GOOS=linux GOARCH=amd64 go build -o keyp-linux-amd64 ./cmd/keyp
GOOS=windows GOARCH=amd64 go build -o keyp-windows-amd64.exe ./cmd/keyp
```

### CGO Requirement

SQLCipher requires CGO. Ensure `CGO_ENABLED=1` for builds:

```bash
CGO_ENABLED=1 go build ./cmd/keyp
```

Cross-compilation requires appropriate C toolchains or Docker-based builds.

---

## Git Commit Conventions

Follow Conventional Commits with Beads reference:

```
<type>(<scope>): <description> (bd:<issue-id>)
```

**Types**: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `perf`

**Examples**:
```bash
git commit -m "feat(core): port AES-256-GCM encryption from TypeScript (bd:keyp-003)"
git commit -m "feat(store): implement SQLCipher storage layer (bd:keyp-004)"
git commit -m "test(core): add crypto round-trip tests (bd:keyp-003)"
```

---

## Task Specification Workflow

### Reading Task Specs

Task specifications are in the project root:
- `keyp-phase1-task-spec.md` - Go scaffold + crypto port
- `keyp-phase2-task-spec.md` - SQLite + data model
- etc.

### Executing Task Specs

1. Read "Beads Issue Setup" section
2. Execute commands **in order**, capturing returned IDs
3. Run `bd ready --json` to verify setup
4. Work through tasks, closing each when complete
5. When `bd ready --json` returns empty, phase is done

### Common Pitfalls

❌ **DON'T**: Use example IDs from spec (keyp-001, etc.)
✅ **DO**: Use actual IDs returned by bd create

❌ **DON'T**: Forget to close issues when done
✅ **DO**: Run `bd close <issue-id>` immediately after completing

❌ **DON'T**: Skip the bd-issue-tracking skill
✅ **DO**: Read it first, every time

---

## Key Dependencies

```go
// go.mod
module github.com/TheEditor/keyp

go 1.21

require (
    github.com/spf13/cobra v1.8.0           // CLI framework
    github.com/mattn/go-sqlite3 v1.14.19    // SQLite driver (CGO)
    golang.org/x/crypto v0.18.0             // PBKDF2
    golang.org/x/term v0.16.0               // Password input
    github.com/atotto/clipboard v0.1.4      // Clipboard
    github.com/go-git/go-git/v5 v5.11.0     // Git operations
)
```

Note: SQLCipher support comes via build tags or a SQLCipher-enabled SQLite driver.

---

## Reference: TypeScript Crypto (Port Source)

The Go crypto implementation must produce **bit-identical** output to the TypeScript version for migration compatibility. Reference `legacy/src/crypto.ts`:

- Algorithm: AES-256-GCM
- Key derivation: PBKDF2-SHA256, 100,000 iterations minimum
- IV: 12 bytes (96 bits), random per encryption
- Salt: 32 bytes, random per vault
- Auth tag: 16 bytes (128 bits)

---

## Getting Help

### Beads Commands Reference

```bash
bd help                    # General help
bd ready                   # Show ready work
bd show <issue-id>         # Show issue details
bd list --status open      # List open issues
bd dep tree <issue-id>     # Show dependency tree
```

### When Stuck

1. Read the task spec - implementation details are there
2. Check issue description - `bd show <issue-id>`
3. Check dependencies - `bd dep tree <issue-id>`
4. Review the bd-issue-tracking skill
5. Ask the human when genuinely blocked

---

## Critical Reminders

1. **Read bd-issue-tracking skill FIRST** - Every session
2. **Capture actual IDs** - Never use placeholder IDs from specs
3. **Close issues when done** - Verify with `bd ready`
4. **CGO_ENABLED=1** - Required for SQLCipher builds
5. **Port crypto exactly** - Must match TypeScript output
6. **Include issue ID in commits** - `(bd:<issue-id>)` format

---

**Last Updated**: December 2024
**Project Status**: v2 Go migration - Phase 1 (scaffold + crypto) pending
