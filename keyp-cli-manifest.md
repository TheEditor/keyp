# keyp CLI Manifest: Phases 2.5-7

## Preface

Complete CLI implementation from storage layer gaps through HTTP server. Phase 2 (SQLite storage with SQLCipher) must be complete before starting.

**For the executing agent:**
1. Read AGENTS.md and TASK-MANIFESTS.md
2. Read BD-Create-Detail-Level-Guideline.md - calibrate description detail for `bd create`
3. Initialize Beads if needed: `bd init`
4. Create issues using `bd create` with descriptions you compose based on the subject lines
5. Replace `<TBD>` with returned IDs
6. Add dependencies per the graphs below

---

## Phase 2.5: Storage Layer Gaps

Phase 2 may be missing components that later phases require. Address these before proceeding.

```
keyp-nlk: epic - Phase 2.5 Storage Layer Gaps
keyp-nlk.8: task - Add context.Context parameter to all store methods [NEW - blocks keyp-nlk.1, keyp-nlk.2]
keyp-nlk.1: task - Add store.Update method for modifying existing secrets [depends: keyp-nlk.8]
keyp-nlk.2: task - Enhance store.Search with SearchOptions (tags filter, limit) [depends: keyp-nlk.8]
keyp-nlk.3: CLOSED - FTS5 virtual table and triggers already exist in store.go
keyp-nlk.4: task - Add missing domain error types to internal/store/errors.go (ErrNotFound exists, add others)
keyp-nlk.5: CLOSED - JSON tags already exist on SecretObject and Field types
keyp-nlk.6: task - Add SecretObject.Redacted() method: returns copy with sensitive field values masked
keyp-nlk.7: task - Add store package unit tests: CRUD operations, search, error cases [depends: keyp-nlk.1, keyp-nlk.2, keyp-nlk.4]
```

**Dependency Graph:**
```
P2.5-epic ─► context-migration ─┬─► store-update
                                └─► store-search-enhance
           ─► error-types
           ─► redacted-method

all-P2.5-tasks ─► store-unit-tests
```

---

## Phase 3: Core CLI Commands

### Utilities

```
keyp-2fx: epic - Phase 3 Core CLI [depends: keyp-nlk]
keyp-2fx.1: task - Create internal/ui package with PromptPassword function (hidden input, non-terminal fallback)
keyp-2fx.2: task - Add PromptConfirmPassword function (prompt twice, verify match) [depends: keyp-2fx.1]
keyp-2fx.3: task - Add PromptVisible function (visible input with trim) [depends: keyp-2fx.1]
keyp-2fx.4: task - Create clipboard functions: CopyToClipboard, CopyWithAutoClear (45s default) using atotto/clipboard
keyp-2fx.5: task - Add ui package tests [depends: keyp-2fx.2, keyp-2fx.3, keyp-2fx.4]
```

### CLI Commands

```
keyp-2fx.6: task - Add cobra initCmd scaffold with --path flag for custom vault location [depends: keyp-2fx]
keyp-2fx.7: task - Implement initCmd.RunE: check exists, prompt password, validate 8+ chars, call vault.Init [depends: keyp-2fx.6, keyp-2fx.2]
keyp-2fx.8: task - Add cobra setCmd scaffold: args <name> [value], --stdin flag [depends: keyp-2fx]
keyp-2fx.9: task - Implement setCmd.RunE: open vault, prompt if no value, create SecretObject with single field [depends: keyp-2fx.8, keyp-2fx.1]
keyp-2fx.10: task - Add cobra getCmd scaffold: args <name>, --stdout flag, --field flag [depends: keyp-2fx]
keyp-2fx.11: task - Implement getCmd.RunE: open vault, fetch secret, copy first field to clipboard or print [depends: keyp-2fx.10, keyp-2fx.4]
keyp-2fx.12: task - Add cobra listCmd scaffold with --tags and --json flags [depends: keyp-2fx]
keyp-2fx.13: task - Implement listCmd.RunE: open vault, list secrets, format output [depends: keyp-2fx.12]
keyp-2fx.14: task - Add cobra deleteCmd scaffold: args <name>, --force flag [depends: keyp-2fx]
keyp-2fx.15: task - Implement deleteCmd.RunE: confirm unless --force, delete secret [depends: keyp-2fx.14, keyp-2fx.3]
keyp-2fx.16: task - Add CLI integration tests for init/set/get/list/delete [depends: keyp-2fx.7, keyp-2fx.9, keyp-2fx.11, keyp-2fx.13, keyp-2fx.15]
```

**Dependency Graph:**
```
prompt-password ─► prompt-confirm ─┬─► ui-tests
                 ├─► prompt-visible ─┘
                 │
clipboard ───────┘

P3-epic ─┬─► init-scaffold ─► init-impl
         ├─► set-scaffold ─► set-impl
         ├─► get-scaffold ─► get-impl
         ├─► list-scaffold ─► list-impl
         └─► delete-scaffold ─► delete-impl

prompt-confirm ─► init-impl
prompt-password ─► set-impl
clipboard ─► get-impl
prompt-visible ─► delete-impl

all-cmd-implementations ─► cli-integration-tests
```

---

## Phase 4: Multi-Field Secrets & Search

### Interactive Secret Creation

```
keyp-laz: epic - Phase 4 Advanced CLI [depends: keyp-2fx]
keyp-laz.1: task - Add PromptLoop helper: repeatedly prompt until empty input or max fields [depends: keyp-2fx.3]
keyp-laz.2: task - Add cobra addCmd scaffold: args <name>, interactive multi-field creation [depends: keyp-laz]
keyp-laz.3: task - Implement addCmd.RunE: prompt for fields (label/value/sensitive/type) in loop, create SecretObject [depends: keyp-laz.2, keyp-laz.1]
```

### Secret Display & Editing

```
keyp-laz.4: task - Add cobra showCmd scaffold: args <name>, --reveal flag to unmask sensitive [depends: keyp-laz]
keyp-laz.5: task - Implement showCmd.RunE: display all fields, mask sensitive unless --reveal [depends: keyp-laz.4]
keyp-laz.6: task - Add cobra editCmd scaffold: args <name>, --field flag to target specific field [depends: keyp-laz]
keyp-laz.7: task - Implement editCmd.RunE: load secret, prompt for new values, call store.Update [depends: keyp-laz.6, keyp-laz.5, keyp-nlk.1]
```

### Search & Tags

```
keyp-laz.8: task - Add cobra searchCmd scaffold: args <query>, uses FTS5 [depends: keyp-laz]
keyp-laz.9: task - Implement searchCmd.RunE: call store.Search, format results with match highlights [depends: keyp-laz.8, keyp-nlk.2]
keyp-laz.10: task - Add cobra tagCmd scaffold with subcommands: add, rm, list [depends: keyp-laz]
keyp-laz.11: task - Implement tagCmd subcommands: modify secret.Tags array [depends: keyp-laz.10]
keyp-laz.12: task - Add Phase 4 integration tests [depends: keyp-laz.3, keyp-laz.5, keyp-laz.7, keyp-laz.9, keyp-laz.11]
```

**Dependency Graph:**
```
P4-epic ─┬─► add-scaffold ─► add-impl
         ├─► show-scaffold ─► show-impl ─► edit-impl
         ├─► search-scaffold ─► search-impl
         └─► tag-scaffold ─► tag-impl

prompt-loop ─► add-impl
all-P4-implementations ─► P4-integration-tests
```

---

## Phase 5: Git Sync

```
keyp-mrz: epic - Phase 5 Git Sync [depends: keyp-2fx]
keyp-mrz.1: task - Implement GitExecSyncer using exec.Command (implements Syncer interface from P2)
keyp-mrz.2: task - Add GitExecSyncer.Init: git init in vault dir, create .gitignore for *.db [depends: keyp-mrz.1]
keyp-mrz.3: task - Add GitExecSyncer.AddRemote: git remote add origin <url> [depends: keyp-mrz.1]
keyp-mrz.4: task - Add GitExecSyncer.Commit: git add + git commit with message [depends: keyp-mrz.1]
keyp-mrz.5: task - Add GitExecSyncer.Push and Pull: git push/pull origin main [depends: keyp-mrz.1]
keyp-mrz.6: task - Add GitExecSyncer.Status: parse git status for clean/dirty/ahead/behind [depends: keyp-mrz.1]
keyp-mrz.7: task - Add cobra syncCmd scaffold with subcommands: init, push, pull, status [depends: keyp-mrz]
keyp-mrz.8: task - Wire syncCmd subcommands to GitExecSyncer methods [depends: keyp-mrz.7, keyp-mrz.2, keyp-mrz.3, keyp-mrz.4, keyp-mrz.5, keyp-mrz.6]
keyp-mrz.9: task - Add git sync integration tests (requires git installed) [depends: keyp-mrz.8]
```

**Dependency Graph:**
```
P5-epic ─► git-exec-syncer ─┬─► git-init
                            ├─► git-add-remote
                            ├─► git-commit
                            ├─► git-push-pull
                            └─► git-status

P5-epic ─► sync-scaffold

all-git-exec-methods + sync-scaffold ─► sync-cmd-wiring ─► git-sync-tests
```

---

## Phase 6: Vault Lock/Unlock UX

**SCOPE: IN-PROCESS ONLY.** CLI lock/unlock commands are per-process (no cross-invocation persistence). Each CLI command prompts for password. The HTTP server (Phase 7) manages its own session lifecycle with persistent VaultHandle. Future enhancement: daemon mode for cross-process CLI sessions.

### Vault Handle (shared by CLI and Server)

```
keyp-elp: epic - Phase 6 Vault Lock/Unlock UX [depends: keyp-2fx]
keyp-elp.1: task - Define VaultHandle struct: holds open db, derived key, unlock timestamp
keyp-elp.2: task - Add VaultHandle.Unlock: accept password, derive key, open db, return handle [depends: keyp-elp.1]
keyp-elp.3: task - Add VaultHandle.Lock: zero key memory, close db, invalidate handle [depends: keyp-elp.1]
keyp-elp.4: task - Add VaultHandle.IsExpired(timeout): check if unlock timestamp + timeout < now [depends: keyp-elp.1]
keyp-elp.5: task - Add vault handle unit tests [depends: keyp-elp.2, keyp-elp.3, keyp-elp.4]
```

### CLI Lock/Unlock (in-process only)

```
keyp-elp.6: task - Add cobra unlockCmd: validate password, --timeout flag (primarily for HTTP server prep) [depends: keyp-elp.2]
keyp-elp.7: task - Add cobra lockCmd: call VaultHandle.Lock, clear process-global handle [depends: keyp-elp.3]
keyp-elp.8: task - Refactor CLI commands to use getOrUnlockVault helper (prompts each invocation) [depends: keyp-elp.6]
keyp-elp.9: task - Add CLI auto-lock: check IsExpired before each command, lock if expired [depends: keyp-elp.8, keyp-elp.4]
keyp-elp.10: task - Add CLI lock/unlock integration tests [depends: keyp-elp.7, keyp-elp.6, keyp-elp.9]
```

**Dependency Graph:**
```
P6-epic ─► vault-handle ─┬─► vault-handle-unlock ─► unlock-cmd ─► cli-refactor ─► cli-auto-lock
                         ├─► vault-handle-lock ─► lock-cmd
                         └─► vault-expired ─► cli-auto-lock

vault-handle-tests ─► [vault-unlock, vault-lock, vault-expired]
cli-lock-unlock-tests ─► [lock-cmd, unlock-cmd, cli-auto-lock]
```

---

## Phase 7: HTTP Server Mode

**REQUIREMENT: Go 1.22+** for new http.ServeMux routing patterns (`"GET /path"`, `{name}` path parameters, `r.PathValue()`).

### Foundation

```
keyp-36c: epic - Phase 7 HTTP Server [depends: keyp-2fx, keyp-elp]
keyp-36c.1: task - Define API response envelope: {ok: bool, data?: T, error?: {code, message}} in internal/server/types.go
keyp-36c.2: task - Create internal/server/router.go with stdlib mux (Go 1.22+), middleware chain [depends: keyp-36c.1]
keyp-36c.3: task - Add request logging middleware: method, path, status, duration [depends: keyp-36c.2]
keyp-36c.4: task - Add panic recovery middleware: catch panics, return 500 with error envelope [depends: keyp-36c.2]
```

### Authentication

```
keyp-36c.5: task - Define SessionStore interface and MemorySessionStore (map + mutex + expiry) in internal/server/session.go
keyp-36c.6: task - Add Session struct: token, VaultHandle, created_at, expires_at [depends: keyp-36c.5, keyp-elp.1]
keyp-36c.7: task - Add token generation: 32 bytes from crypto/rand, hex encoded [depends: keyp-36c.6]
keyp-36c.8: task - Implement POST /v1/unlock handler: accept {password}, create Session with VaultHandle, return {token, expires_at} [depends: keyp-36c.7, keyp-36c.2]
keyp-36c.9: task - Implement POST /v1/lock handler: call VaultHandle.Lock, remove session from store [depends: keyp-36c.8]
keyp-36c.10: task - Add auth middleware: extract Bearer token, lookup Session, inject VaultHandle into request context [depends: keyp-36c.8]
keyp-36c.11: task - Implement POST /v1/refresh handler: extend session expiry (not VaultHandle), return new expires_at [depends: keyp-36c.10]
```

### Public Endpoints (no auth)

```
keyp-36c.12: task - Implement GET /health endpoint: return {status: "ok"} [depends: keyp-36c.2]
keyp-36c.13: task - Implement GET /version endpoint: return {version, go_version, build_time} [depends: keyp-36c.2]
```

### Protected Endpoints (require auth)

```
keyp-36c.14: task - Define SecretListItem and SecretDetail types for API responses (use SecretObject.Redacted by default) [depends: keyp-36c.1, keyp-nlk.6]
keyp-36c.15: task - Implement GET /v1/secrets endpoint: return [{name, tags, created_at, updated_at}] [depends: keyp-36c.10, keyp-36c.14]
keyp-36c.16: task - Implement GET /v1/secrets/:name endpoint: return full secret with fields [depends: keyp-36c.10, keyp-36c.14]
keyp-36c.17: task - Define CreateSecretRequest and UpdateSecretRequest types [depends: keyp-36c.1]
keyp-36c.18: task - Implement POST /v1/secrets endpoint: create secret from request body [depends: keyp-36c.10, keyp-36c.17]
keyp-36c.19: task - Implement PUT /v1/secrets/:name endpoint: call store.Update [depends: keyp-36c.10, keyp-36c.17, keyp-nlk.1]
keyp-36c.20: task - Implement DELETE /v1/secrets/:name endpoint: delete secret, return 204 [depends: keyp-36c.10]
keyp-36c.21: task - Implement GET /v1/search?q= endpoint: call store.Search, return matching secrets [depends: keyp-36c.10, keyp-36c.14, keyp-nlk.2]
keyp-36c.22: task - Implement POST /v1/secrets/:name/clipboard endpoint: copy field value to server clipboard [depends: keyp-36c.10, keyp-2fx.4]
```

### Server Lifecycle

```
keyp-36c.23: task - Add cobra serveCmd scaffold: --port (default 9999), --bind (default 127.0.0.1), --timeout (default 15m)
keyp-36c.24: task - Implement Server.Start and Server.Shutdown: bind to address, signal handling, graceful drain [depends: keyp-36c.23, keyp-36c.2, keyp-36c.15, keyp-36c.16, keyp-36c.18, keyp-36c.19, keyp-36c.20, keyp-36c.21, keyp-36c.22]
keyp-36c.25: task - Add HTTP server integration tests: auth flow, CRUD operations, error cases [depends: keyp-36c.24]
```

**Dependency Graph:**
```
P7-epic ─► response-envelope ─► router ─┬─► logging-middleware
                                        └─► panic-recovery-middleware

         ─► session-store ─► token-generation ─► unlock-handler ─┬─► lock-handler
                                                                  ├─► auth-middleware ─► refresh-handler
                                                                  │
                                                                  └─► all protected endpoints

router ─► health-endpoint
       ─► version-endpoint

response-envelope ─► secret-types ─► GET/POST/PUT/DELETE secrets
                  ─► request-types ─► POST/PUT secrets

serve-scaffold + all-endpoints ─► server-start ─► server-shutdown ─► http-tests
```

**API versioning:** All protected routes under `/v1/` prefix for future compatibility.

---

## Cross-Phase Dependencies

```
Phase 2 (complete) ─► keyp-nlk ─► keyp-2fx ─┬─► keyp-laz
                                            ├─► keyp-mrz
                                            └─► keyp-elp ─► keyp-36c
```

P2.5 fills storage gaps. Phases 4, 5, 6 can proceed in parallel after Phase 3. Phase 7 requires Phase 6 (VaultHandle).

---

## Summary

| Phase | Epic | Tasks | Focus |
|-------|------|-------|-------|
| 2.5 | keyp-nlk | 6 open (2 closed) | context.Context, Update, Search, errors, Redacted |
| 3 | keyp-2fx | 16 | utilities + init/set/get/list/delete |
| 4 | keyp-laz | 12 | add/show/edit/search/tag |
| 5 | keyp-mrz | 9 | exec-based git operations |
| 6 | keyp-elp | 10 | VaultHandle + CLI commands (in-process only) |
| 7 | keyp-36c | 25 | REST API for GUI (Go 1.22+) |
| **Total** | **6** | **78 open** | |

---

## Completion

- [x] All `<TBD>` replaced with Beads IDs
- [x] All dependencies added via `bd dep add`
- [ ] Phase 2.5: store.Update, store.Search, error types work
- [ ] Phase 3: `keyp init/set/get/list/delete` work
- [ ] Phase 4: `keyp add/show/edit/search/tag` work
- [ ] Phase 5: `keyp sync init/push/pull/status` work
- [ ] Phase 6: `keyp lock/unlock` work, VaultHandle shared by CLI and server
- [ ] Phase 7: `keyp serve` works, all REST endpoints respond
- [ ] `bd ready --json` returns empty
- [ ] `make build && make test` passes
