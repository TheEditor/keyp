# keyp

> Local-first secret manager for developers
> *"pass for the Node.js generation"*

## What is keyp?

**keyp** is a simple, secure, local-first secret manager designed specifically for developers. Store your API keys, tokens, and credentials with AES-256 encryption, sync across machines via Git, and never worry about leaking secrets again.

## Status

**✅ Week 1 Complete: Core encryption + vault management**

- ✅ AES-256-GCM encryption implementation
- ✅ PBKDF2 key derivation with secure salts
- ✅ Encrypted vault file format
- ✅ Vault initialization and management
- ✅ Secret CRUD operations
- ✅ Comprehensive test suite (39 tests, all passing)
- ✅ Security and format documentation

**✅ Week 2 Complete: CLI Commands**

- ✅ `keyp init` - Initialize vault
- ✅ `keyp set` - Store secrets
- ✅ `keyp get` - Retrieve secrets (clipboard support)
- ✅ `keyp list` - List all secrets
- ✅ `keyp delete` - Delete secrets
- ✅ `keyp rename` / `keyp copy` - Manage secrets
- ✅ `keyp export` / `keyp import` - Backup and migrate
- ✅ Beautiful terminal UI with colors and formatting
- ✅ Masked password input for security
- ✅ All core commands tested and working

**✅ Week 3 Complete: Git sync + polish**

- ✅ Git integration for encrypted backups
- ✅ `keyp sync` command (init, push, pull, status, config)
- ✅ Enhanced password strength validation with visual meter
- ✅ Shell completion scripts (bash and zsh)
- ✅ `keyp stats` - Vault statistics and encryption info
- ✅ `keyp config` - Configuration management
- ✅ Comprehensive Git sync documentation
- ✅ All 39 tests passing with new features

## Features

- 🔒 **Secure** - AES-256-GCM encryption with PBKDF2 key derivation
- 🏠 **Local-first** - No cloud account required, works offline
- 🔄 **Git-based sync** - Encrypted secrets safely committed to Git
- ⚡ **Fast & simple** - Intuitive CLI, zero configuration
- 🔧 **Developer-friendly** - Script integration, clipboard support
- 🆓 **Free & open source** - MIT license

## Quick Start

### Initialize your vault
```bash
$ keyp init
Enter master password: ●●●●●●●●
Confirm master password: ●●●●●●●●
✓ Vault initialized successfully!
```

### Store a secret
```bash
$ keyp set github-token
Enter master password: ●●●●●●●●
Enter value for "github-token": ●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret "github-token" saved
```

### Retrieve a secret (copies to clipboard)
```bash
$ keyp get github-token
Enter master password: ●●●●●●●●
✓ Copied to clipboard (clears in 45 seconds)
```

### List all secrets
```bash
$ keyp list
Enter master password: ●●●●●●●●

  • api-key
  • database-url
  • github-token

3 secrets stored
```

### Delete a secret
```bash
$ keyp delete github-token -f
Enter master password: ●●●●●●●●
Delete secret "github-token"? (y/N): y
✓ Secret "github-token" deleted
Remaining secrets: 2
```

## Why keyp?

Unlike enterprise secret managers (too complex) or traditional Unix password managers (too arcane), **keyp** is designed for the way modern developers work:

- ✅ No GPG complexity
- ✅ No cloud account needed
- ✅ No team features bloat
- ✅ Just simple, secure secret storage

## Installation

```bash
npm install -g @theeditor/keyp
```

## Development

Want to contribute? Watch this repo for updates!

```bash
git clone https://github.com/TheEditor/keyp.git
cd keyp
npm install
```

## Documentation

- 📖 **[CLI Reference](./docs/CLI.md)** - Command-line interface guide
- 🌐 **[Git Sync Guide](./docs/GIT_SYNC.md)** - Multi-machine sync and encrypted backups
- 🔧 **[API Reference](./docs/API.md)** - Library API with examples
- 🔐 **[Security Guide](./docs/SECURITY.md)** - Cryptographic details and threat model
- 📋 **[Vault Format](./docs/VAULT_FORMAT.md)** - Technical vault file specification

## Roadmap

**Week 1: Core encryption + vault management** ✅
- [x] Core encryption implementation (AES-256-GCM)
- [x] PBKDF2 key derivation with 100,000+ iterations
- [x] Vault initialization and management
- [x] Secret CRUD operations
- [x] Comprehensive tests (39 passing, 100%)
- [x] Security and vault format documentation

**Week 2: CLI Commands** ✅
- [x] Beautiful CLI with colors and formatting
- [x] `keyp init` - Initialize vault with password prompts
- [x] `keyp set <name> [value]` - Store secrets
- [x] `keyp get <name>` - Retrieve secrets to clipboard
- [x] `keyp list` - List all secrets with search
- [x] `keyp delete <name>` - Delete secrets (bonus)
- [x] Masked password input for security
- [x] Clipboard auto-clear after 45 seconds

**Week 3: Git sync + polish** ✅
- [x] Git integration for encrypted backups
- [x] `keyp sync init/push/pull/status/config` commands
- [x] Enhanced password strength validation with visual meter
- [x] Shell completion scripts (bash and zsh)
- [x] `keyp stats` - Vault statistics command
- [x] `keyp config` - Configuration management
- [x] Comprehensive Git sync and CLI documentation

**Week 4: v1.0.0 launch** 📅
- [ ] Complete documentation and examples
- [ ] Launch announcement
- [ ] Community feedback and iteration
- [ ] Performance optimization
- [ ] Additional platform support

## Philosophy

**keyp** follows these principles:

1. **Local-first** - Your secrets stay on your machine
2. **Simple** - One command to do one thing, and do it well
3. **Secure** - Industry-standard encryption, no shortcuts
4. **Developer-focused** - Built for developers, by developers

## Inspiration

Inspired by [pass](https://www.passwordstore.org/) but designed for modern Node.js developers who want simplicity without complexity.

## License

MIT © Dave Fobare

---

⭐ **Star this repo to follow development!**

🐛 **Found a bug?** [Open an issue](https://github.com/TheEditor/keyp/issues)

💡 **Have ideas?** [Start a discussion](https://github.com/TheEditor/keyp/discussions)
