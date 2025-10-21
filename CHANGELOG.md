# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2025-10-21

### Added - Git Sync + Polish

**Core Features:**
- Git synchronization with encrypted backups
  - `keyp sync init` - Initialize Git sync with remote repository
  - `keyp sync push` - Push encrypted vault to remote
  - `keyp sync pull` - Pull vault from remote with conflict detection
  - `keyp sync status` - Display synchronization status
  - `keyp sync config` - Configure sync settings
  - Multi-machine secret synchronization
  - Conflict resolution (keep-local, keep-remote)

**CLI Enhancements:**
- `keyp stats` - Display vault statistics and encryption info
- `keyp config` - Manage configuration settings
  - Clipboard timeout configuration
  - Auto-lock settings
  - Git auto-sync options
- Enhanced password strength validation with visual meter
  - Entropy-based strength scoring (0-100%)
  - Visual progress bar (█████░░░░░)
  - Specific recommendations for improvements

**Developer Experience:**
- Shell completion scripts
  - Bash completion (completions/keyp.bash)
  - Zsh completion (completions/keyp.zsh)
  - Tab completion for commands, secrets, flags, and paths

**Documentation:**
- Comprehensive Git Sync guide (docs/GIT_SYNC.md)
  - Setup instructions for GitHub, GitLab, self-hosted Git
  - SSH and HTTPS authentication guides
  - Multi-machine sync workflows
  - Conflict resolution strategies
  - Security best practices
- Updated CLI reference with all new commands
- Configuration management documentation

### Changed
- Enhanced error messages with contextual help
- Improved password initialization workflow with strength feedback
- Updated README with Week 3 completion status
- Expanded package.json metadata

### Technical
- Added simple-git dependency for Git operations
- Full TypeScript strict mode compliance
- Backward compatible with existing vaults

## [0.1.5] - 2025-10-21

### Added - Optional CLI Enhancements

**New Commands:**
- `keyp rename <old> <new>` - Rename secrets
- `keyp copy <source> <dest>` - Copy secrets to new names
- `keyp export [file]` - Export secrets with encryption or plaintext
- `keyp import <file>` - Import secrets from file
  - Merge mode (default)
  - Replace mode for full vault migration
  - Dry-run preview before importing

**Enhanced Commands:**
- `keyp get` - Added `--timeout` parameter for clipboard clear timeout
  - Default: 45 seconds
  - Customizable per-command
  - Can be disabled with `--no-clear`

**Documentation:**
- Updated CLI.md with 250+ lines for new commands
- Complete export/import workflow documentation
- Real-world examples for all operations

## [0.1.0] - 2025-10-20

### Added - CLI Commands

**Core Commands:**
- `keyp init` - Initialize vault with password prompts and strength validation
- `keyp set <name> [value]` - Store secrets with masked input
- `keyp get <name>` - Retrieve secrets to clipboard (45s auto-clear)
- `keyp list` - List all secrets with search and count options
- `keyp delete <name>` - Delete secrets with confirmation

**Features:**
- Beautiful terminal UI with colors and formatting
- Masked password input for security
- Password strength recommendations
- Clipboard integration (macOS, Linux, Windows)
- Cross-platform support

**Documentation:**
- Comprehensive CLI reference (docs/CLI.md)
- Troubleshooting guide
- Platform-specific instructions

## [0.0.1] - 2025-10-20

### Added - Core Encryption + Vault Management

**Cryptographic Foundation:**
- AES-256-GCM authenticated encryption
- PBKDF2-SHA256 key derivation with 100,000+ iterations
- Secure random salt generation
- Authenticated encryption with auth tags

**Vault Management:**
- Encrypted vault file format (JSON with Base64 encoding)
- Vault initialization with password protection
- Vault locking/unlocking
- Secret CRUD operations
- Vault persistence with encryption

**Features:**
- Complete secret management (Create, Read, Update, Delete)
- Secret searching with pattern matching
- Secret counting and listing
- Secure data clearing on lock
- Full TypeScript with strict type checking

**Documentation:**
- Security guide with threat model and analysis (docs/SECURITY.md)
- Vault format specification (docs/VAULT_FORMAT.md)
- API reference with examples (docs/API.md)
- Complete README with quick start

**Testing:**
- 39 comprehensive tests across crypto and vault modules
- Unit tests for encryption/decryption
- Integration tests for vault operations
- 100% test pass rate

---

## Version Roadmap

### v0.3.0 (Planned)
- Performance optimizations and benchmarks
- Cross-platform CI/CD pipeline
- Expanded testing coverage
- Community feedback integration

### v1.0.0 (Planned)
- Production-ready release
- Complete documentation and examples
- Launch announcement and promotion
- Stable API commitment

