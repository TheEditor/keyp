# keyp

> Local-first secret manager for developers and families  
> *Structured secrets. Encrypted storage. No cloud required.*

## What is keyp?

**keyp** stores your passwords, API keys, PINs, and sensitive notes in an encrypted local database. Unlike simple password managers, keyp handles *structured secrets* — store your AT&T account with its login, account PIN, support PIN, and billing email all in one place.

Built for developers who want something simpler than enterprise vaults but more powerful than plaintext files. Designed to eventually serve families who need one secure place for household credentials.

## Features

- 🔐 **SQLCipher encryption** — Industry-standard AES-256, whole-database encryption
- 📦 **Structured secrets** — Multiple fields per secret (passwords, PINs, notes, URLs)
- 🏷️ **Tag-based organization** — Flexible categorization without rigid folders
- 🔍 **Full-text search** — Find secrets by name, tags, notes, or field labels
- 🔄 **Git sync** — Backup your encrypted vault to any Git remote
- 🖥️ **HTTP API** — Built-in server mode for GUI integration
- ⏱️ **Auto-lock** — Configurable session timeout
- 📋 **Clipboard support** — Copy secrets without displaying them

## Quick Start

### Initialize your vault

```bash
keyp init
# Enter master password: ••••••••
# Confirm password: ••••••••
# ✓ Vault created at ~/.keyp/vault.db
```

### Store a simple secret

```bash
keyp set github-token
# Enter master password: ••••••••
# Enter value: ••••••••
# ✓ Secret 'github-token' saved
```

### Store a structured secret

```bash
keyp add "AT&T"
# Enter master password: ••••••••
# Enter fields (empty label to finish):
#   Label: Account PIN
#   Value: ••••••••
#   Label: Support PIN  
#   Value: ••••••••
#   Label: Email
#   Value: billing@example.com
#   Label: 
# ✓ Secret 'AT&T' created with 3 field(s)
```

### Retrieve a secret

```bash
keyp get github-token
# Enter master password: ••••••••
# ✓ Copied to clipboard (clears in 45s)

keyp get github-token --stdout
# Enter master password: ••••••••
# ghp_xxxxxxxxxxxxxxxxxxxx
```

### View a structured secret

```bash
keyp show "AT&T"
# Enter master password: ••••••••
# 
# AT&T
# ────────────────────
# Account PIN: ********
# Support PIN: ********
# Email: billing@example.com
# 
# Tags: (none)
# Notes: (none)

keyp show "AT&T" --reveal
# Shows actual values instead of ********
```

### Search your vault

```bash
keyp search "PIN"
# Enter master password: ••••••••
# 
# Found 2 secrets:
#   AT&T (Account PIN, Support PIN)
#   Verizon (Account PIN)

keyp list
# Enter master password: ••••••••
# 
# Secrets (3):
#   AT&T [telecom, family]
#   github-token
#   Verizon [telecom]
```

### Organize with tags

```bash
keyp tag add "AT&T" telecom family
keyp tag add "Verizon" telecom

keyp list --tag telecom
# AT&T
# Verizon
```

## Installation

### From Source (requires Go 1.21+ and CGO)

```bash
git clone https://github.com/TheEditor/keyp.git
cd keyp
make build
sudo mv keyp /usr/local/bin/
```

### Build Requirements

keyp uses SQLCipher for encryption, which requires CGO:

| Platform | Requirements |
|----------|-------------|
| **macOS** | `brew install sqlcipher` |
| **Ubuntu/Debian** | `sudo apt install libsqlcipher-dev` |
| **Windows** | See [Windows Build Guide](#windows-build) |

Build with:
```bash
CGO_ENABLED=1 go build -o keyp ./cmd/keyp
```

### Pre-built Binaries

Coming soon. See [Releases](https://github.com/TheEditor/keyp/releases).

## Command Reference

### Core Commands

| Command | Description |
|---------|-------------|
| `keyp init` | Create a new vault |
| `keyp set <name> [value]` | Store a simple key-value secret |
| `keyp get <name>` | Copy secret to clipboard |
| `keyp list` | List all secrets |
| `keyp delete <name>` | Remove a secret |

### Structured Secrets

| Command | Description |
|---------|-------------|
| `keyp add <name>` | Create secret with multiple fields (interactive) |
| `keyp show <name>` | Display all fields of a secret |
| `keyp edit <name>` | Modify an existing secret |
| `keyp search <query>` | Full-text search across all secrets |

### Organization

| Command | Description |
|---------|-------------|
| `keyp tag add <name> <tags...>` | Add tags to a secret |
| `keyp tag rm <name> <tags...>` | Remove tags from a secret |
| `keyp list --tag <tag>` | Filter by tag |

### Session Management

| Command | Description |
|---------|-------------|
| `keyp unlock` | Unlock vault for session |
| `keyp lock` | Lock vault immediately |

### Git Sync

| Command | Description |
|---------|-------------|
| `keyp sync init [remote-url]` | Initialize git in vault directory |
| `keyp sync push` | Push encrypted vault to remote |
| `keyp sync pull` | Pull vault from remote |
| `keyp sync status` | Show sync status |

### HTTP Server

| Command | Description |
|---------|-------------|
| `keyp serve` | Start REST API server |
| `keyp serve --port 9999` | Custom port |
| `keyp serve --timeout 30m` | Session timeout |

## Common Flags

| Flag | Description |
|------|-------------|
| `--stdout` | Print secret to terminal instead of clipboard |
| `--reveal` | Show actual values (with `show` command) |
| `--tag <tag>` | Filter by tag (with `list` command) |
| `--json` | Output as JSON (for scripting) |
| `--field <label>` | Get specific field (with `get` command) |

## Git Sync

Sync your encrypted vault across machines using Git:

```bash
# Initialize (first machine)
keyp sync init https://github.com/you/keyp-vault.git
keyp sync push

# Clone (second machine)
keyp sync init https://github.com/you/keyp-vault.git
keyp sync pull
```

The vault file is encrypted — safe for public repos if you trust your master password. But private repos are recommended.

## HTTP API

Start the server:

```bash
keyp serve --port 8080 --timeout 15m
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/v1/unlock` | Unlock vault, get session token |
| `POST` | `/v1/lock` | Lock vault |
| `GET` | `/v1/secrets` | List all secrets |
| `POST` | `/v1/secrets` | Create secret |
| `GET` | `/v1/secrets/:name` | Get secret by name |
| `PUT` | `/v1/secrets/:name` | Update secret |
| `DELETE` | `/v1/secrets/:name` | Delete secret |
| `GET` | `/v1/search?q=<query>` | Search secrets |
| `GET` | `/health` | Health check |

All protected endpoints require `Authorization: Bearer <token>` header.

## Security

### Encryption

- **Algorithm**: AES-256-GCM (via SQLCipher)
- **Key derivation**: PBKDF2-SHA256, 100,000 iterations
- **Storage**: Entire database encrypted at rest
- **In memory**: Decrypted only while vault is unlocked

### Threat Model

keyp protects against:
- ✅ Stolen laptop (encrypted at rest)
- ✅ Backup exposure (encrypted in git)
- ✅ Shoulder surfing (clipboard, not display)

keyp does NOT protect against:
- ❌ Keyloggers capturing your master password
- ❌ Memory forensics while unlocked
- ❌ Weak master passwords

### Best Practices

1. Use a strong, unique master password (16+ characters)
2. Lock your vault when not in use (`keyp lock`)
3. Use private Git repos for sync
4. Don't use `--stdout` in shared terminals

## Data Model

keyp stores **structured secrets**, not just key-value pairs:

```
Secret: "AT&T"
├── Fields:
│   ├── Account PIN: 1234 (sensitive)
│   ├── Support PIN: 5678 (sensitive)
│   └── Email: billing@example.com
├── Tags: [telecom, family]
└── Notes: "Account holder: John Doe"
```

Each field can be marked sensitive (masked in output) or visible.

## Configuration

Vault location: `~/.keyp/vault.db`

Override with `--path` flag:
```bash
keyp --path /custom/path/vault.db list
```

## Migrating from v1 (TypeScript)

The Go version (v2) uses a different storage format. Migration:

```bash
# Export from v1
keyp list --json > secrets-backup.json

# Manually recreate in v2
keyp init
keyp set secret-name value  # for each secret
```

A migration tool may be provided in a future release.

## Development

```bash
# Clone
git clone https://github.com/TheEditor/keyp.git
cd keyp

# Build
make build

# Test
make test

# Run
./keyp --help
```

### Project Structure

```
keyp/
├── cmd/keyp/          # CLI commands
├── internal/
│   ├── core/          # Crypto operations
│   ├── model/         # SecretObject, Field types
│   ├── store/         # SQLCipher database layer
│   ├── vault/         # Vault handle abstraction
│   ├── server/        # HTTP API
│   ├── sync/          # Git sync
│   └── ui/            # Terminal prompts, clipboard
├── legacy/            # TypeScript v1 (archived)
└── .beads/            # Issue tracking database
```

## Roadmap

- [x] Core CLI (init, set, get, list, delete)
- [x] Structured secrets (add, show, edit)
- [x] Full-text search
- [x] Tag-based organization
- [x] Session management (lock/unlock)
- [x] Git sync
- [x] HTTP server mode
- [ ] Pre-built binaries for all platforms
- [ ] Shell completions (bash, zsh, fish)
- [ ] Import/export (JSON, CSV)
- [ ] GUI application (separate project)

## FAQ

**Q: Why not just use 1Password/Bitwarden?**  
A: Those are great for teams and cross-device sync. keyp is for people who want local-first, no-account, no-subscription secret storage with developer-friendly CLI.

**Q: Is it safe to push my vault to GitHub?**  
A: The vault is encrypted with your master password. If your password is strong, it's safe. But private repos are still recommended.

**Q: Can multiple family members use one vault?**  
A: Yes, that's the intended use case. Share the master password with trusted family members. Each machine can sync via Git.

**Q: What if I forget my master password?**  
A: There's no recovery. The encryption is real. Write it down somewhere safe.

## License

MIT © Dave Fobare

---

**GitHub**: [TheEditor/keyp](https://github.com/TheEditor/keyp)  
**npm (v1)**: [@theeditor/keyp](https://www.npmjs.com/package/@theeditor/keyp)
