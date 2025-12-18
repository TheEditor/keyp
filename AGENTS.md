# Agent Instructions for keyp

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
bd create "Subtask" --parent <epic-id> --json  # Hierarchical subtask (gets ID like epic-id.1)
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

### Issue Types

- `bug` - Something broken
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature with subtasks
- `chore` - Maintenance (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: `bd ready` shows unblocked issues
2. **Claim your task**: `bd update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `bd create "Found bug" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `bd close <id> --reason "Done"`
6. **Commit together**: Always commit the `.beads/issues.jsonl` file together with the code changes so issue state stays in sync with code state

### Auto-Sync

bd automatically syncs with git:
- Exports to `.beads/issues.jsonl` after changes (5s debounce)
- Imports from JSONL when newer (e.g., after `git pull`)
- No manual export/import needed!

### Managing AI-Generated Planning Documents

AI assistants often create planning and design documents during development:
- PLAN.md, IMPLEMENTATION.md, ARCHITECTURE.md
- DESIGN.md, CODEBASE_SUMMARY.md, INTEGRATION_PLAN.md
- TESTING_GUIDE.md, TECHNICAL_DESIGN.md, and similar files

**Best Practice: Use a dedicated directory for these ephemeral files**

**Recommended approach:**
- Create a `history/` directory in the project root
- Store ALL AI-generated planning/design docs in `history/`
- Keep the repository root clean and focused on permanent project files
- Only access `history/` when explicitly asked to review past planning

### CLI Help

Run `bd <command> --help` to see all available flags for any command.
For example: `bd create --help` shows `--parent`, `--deps`, `--assignee`, etc.

### Important Rules

- ✅ Use bd for ALL task tracking
- ✅ Always use `--json` flag for programmatic use
- ✅ Link discovered work with `discovered-from` dependencies
- ✅ Check `bd ready` before asking "what should I work on?"
- ✅ Store AI planning docs in `history/` directory
- ✅ Run `bd <cmd> --help` to discover available flags
- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT use external issue trackers
- ❌ Do NOT duplicate tracking systems
- ❌ Do NOT clutter repo root with planning documents

---

### Understanding Beads ID Assignment

**CRITICAL CONCEPT**: Beads assigns issue IDs automatically. You CANNOT specify them.

#### When Creating Issues

```bash
# You run this command:
bd create "Implement SQLCipher storage layer" -t task -p 0 -d "Description..." --json

# Beads returns JSON like this:
{"id":"keyp-008","title":"Implement SQLCipher storage layer",...}

# The ID "keyp-008" was ASSIGNED by the system
# You must capture it and use it in subsequent commands
```

#### Using Captured IDs

When task specs show commands like:
```bash
bd create "Task" -t task -p 0 --parent keyp-001 -d "..." --json
```

The `keyp-001` is a **placeholder**. Replace it with the **actual ID** returned from the parent epic creation.

---

### Daily Workflow

#### Starting Work

1. **Check ready work**:
   ```bash
   bd ready --json
   ```

2. **Pick an issue**: Choose based on priority (P0 = highest, P3 = lowest)

3. **Update status**:
   ```bash
   bd update <issue-id> --status in_progress
   ```

4. **Do the work**: Implement according to the task spec or issue description

#### During Work

- **Discovered new work?** File an issue immediately:
  ```bash
  bd create "Found bug in crypto" -t bug -p 1 --json
  bd dep add <new-issue-id> <current-issue-id> --type discovered-from
  ```

- **Need to check dependencies?**
  ```bash
  bd dep tree <issue-id>
  ```

#### Completing Work

1. **Verify acceptance criteria** in the issue

2. **Run quality checks**:
   ```bash
   go build ./...
   go test ./...
   ```

3. **Commit your work**:
   ```bash
   git add .
   git commit -m "feat: implement feature (bd:<issue-id>)"
   git push
   ```

4. **Close the issue**:
   ```bash
   bd close <issue-id> --reason "Implemented and tested"
   ```

---

### Session Ending Protocol

**CRITICAL**: Before ending ANY session, complete this checklist:

#### 1. Issue Tracker Hygiene

- [ ] File issues for any discovered bugs, TODOs, or follow-up work
- [ ] Close all completed issues with `bd close <issue-id>`
- [ ] Update status for any in-progress work
- [ ] Run `bd ready` to confirm state

#### 2. Quality Gates

- [ ] Code compiles: `go build ./...`
- [ ] Tests pass: `go test ./...`
- [ ] If broken: File P0 issue immediately

#### 3. Sync Issue Tracker

```bash
bd sync
git add .beads/
git commit -m "chore: sync issue tracker"
git push
```

#### 4. Next Session Context

Provide a brief summary:
- What was completed
- What issues remain
- What's ready next
- Any blockers

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
│   │   ├── errors.go         # ErrNotFound, ErrAlreadyExists
│   │   ├── schema.sql        # Table definitions
│   │   └── store_test.go
│   ├── model/                # Domain types
│   │   ├── secret.go         # SecretObject, Field
│   │   └── secret_test.go
│   ├── vault/                # Vault lifecycle
│   │   ├── vault.go          # Init, Open, Close
│   │   ├── handle.go         # VaultHandle (shared by CLI and server)
│   │   └── vault_test.go
│   ├── ui/                   # Terminal UI utilities
│   │   ├── prompt.go         # PromptPassword, PromptConfirm, PromptVisible
│   │   ├── clipboard.go      # CopyToClipboard, CopyWithAutoClear
│   │   └── ui_test.go
│   ├── sync/                 # Git sync operations
│   │   ├── syncer.go         # Syncer interface
│   │   ├── gitexec.go        # exec.Command implementation
│   │   └── sync_test.go
│   └── server/               # HTTP server (Phase 7)
│       ├── types.go          # Response envelope, request/response types
│       ├── router.go         # Routes and middleware chain
│       ├── session.go        # SessionStore, token management
│       ├── handlers.go       # Endpoint implementations
│       └── server_test.go
├── pkg/
│   └── keyp/                 # Public Go API (for embedding)
├── .beads/                   # Issue tracker database
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── AGENTS.md                 # This file
├── TASK-MANIFESTS.md         # Explains manifest format
└── keyp-cli-manifest.md      # Current work: Phases 2.5-7
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

### Core (Phase 3)

```bash
keyp init [--path]            # Create vault
keyp set <n> [value] [--stdin] # Quick single-field secret
keyp get <n> [--stdout] [--field] # Copy to clipboard
keyp list [--tags] [--json]   # List all secrets
keyp delete <n> [--force]     # Remove secret
```

### Advanced (Phase 4)

```bash
keyp add <n>                  # Interactive multi-field creation
keyp show <n> [--reveal]      # Display all fields (mask sensitive)
keyp edit <n> [--field]       # Modify fields
keyp search <query>           # Full-text search
keyp tag add|rm <n> <tag>     # Manage tags
```

### Git Sync (Phase 5)

```bash
keyp sync init                # Initialize git in vault dir
keyp sync push                # Push to remote
keyp sync pull                # Pull from remote
keyp sync status              # Show sync state
```

### Lock/Unlock (Phase 6)

```bash
keyp unlock [--timeout]       # Unlock vault, keep handle in process
keyp lock                     # Explicitly lock vault
```

### HTTP Server (Phase 7)

```bash
keyp serve [--port] [--bind] [--timeout]  # Start REST API
```

Default: `--port 9999 --bind 127.0.0.1 --timeout 15m`

### Compatibility Note

`keyp set github-token abc123` creates a SecretObject with a single field:
- Name: "github-token"
- Fields: [{ Label: "value", Value: "abc123", Sensitive: true }]

This preserves the simple mental model while enabling structured storage.


---

## HTTP API (Phase 7)

The HTTP server enables a future GUI application. All protected routes require Bearer token auth.

### Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /health | No | Health check |
| GET | /version | No | Version info |
| POST | /v1/unlock | No | Unlock vault, get session token |
| POST | /v1/lock | Yes | Lock vault, invalidate session |
| POST | /v1/refresh | Yes | Extend session expiry |
| GET | /v1/secrets | Yes | List all secrets |
| GET | /v1/secrets/:name | Yes | Get secret with fields |
| POST | /v1/secrets | Yes | Create secret |
| PUT | /v1/secrets/:name | Yes | Update secret |
| DELETE | /v1/secrets/:name | Yes | Delete secret |
| GET | /v1/search?q= | Yes | Full-text search |
| POST | /v1/secrets/:name/clipboard | Yes | Copy to server clipboard |

### Response Envelope

```json
{"ok": true, "data": {...}}
{"ok": false, "error": {"code": "NOT_FOUND", "message": "Secret not found"}}
```
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

## Task Manifests

Work is organized via task manifests (see TASK-MANIFESTS.md). The current manifest is:
- `keyp-cli-manifest.md` - Phases 2.5-7 (storage gaps through HTTP server)

### Using Manifests

1. Read the manifest to understand scope and dependencies
2. Create issues in Beads using `bd create` with descriptions you compose
3. Replace `<TBD>` placeholders with returned IDs
4. Add dependencies using `bd dep add`
5. Work through issues: `bd ready` → `bd show <id>` → implement → `bd close <id>`
6. When `bd ready --json` returns empty, phase is complete

### Common Pitfalls

❌ **DON'T**: Use example IDs from manifest (keyp-001, etc.)
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
    github.com/google/uuid v1.5.0           // Secret/field IDs
    golang.org/x/crypto v0.18.0             // PBKDF2
    golang.org/x/term v0.16.0               // Password input
    github.com/atotto/clipboard v0.1.4      // Clipboard
)
```

**Notes:**
- SQLCipher support comes via build tags or SQLCipher-enabled driver
- go-git is optional (Phase 5 uses exec.Command by default)
- CGO_ENABLED=1 required for SQLite/SQLCipher

---

## Crypto Specification

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
2. **Read keyp-cli-manifest.md** - Understand current work scope
3. **Capture actual IDs** - Never use placeholder IDs from manifests
4. **Close issues when done** - Verify with `bd ready`
5. **CGO_ENABLED=1** - Required for SQLCipher builds
6. **Include issue ID in commits** - `(bd:<issue-id>)` format
7. **VaultHandle is shared** - CLI and HTTP server both use it (Phase 6)
---

**Last Updated**: December 2024
**Project Status**: v2 Go migration - Phases 2.5-7 defined in keyp-cli-manifest.md
