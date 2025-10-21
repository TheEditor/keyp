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
- ✅ Secret CRUD operations (in-memory)
- ✅ Comprehensive test suite (39 tests, all passing)
- ✅ Security and format documentation

**🚧 Week 2: CLI Commands** (coming soon)
- `keyp init` - Initialize vault
- `keyp set` - Store secrets
- `keyp get` - Retrieve secrets
- `keyp list` - List all secrets

## Features

- 🔒 **Secure** - AES-256-GCM encryption with PBKDF2 key derivation
- 🏠 **Local-first** - No cloud account required, works offline
- 🔄 **Git-based sync** - Encrypted secrets safely committed to Git
- ⚡ **Fast & simple** - Intuitive CLI, zero configuration
- 🔧 **Developer-friendly** - Script integration, clipboard support
- 🆓 **Free & open source** - MIT license

## Planned Commands

```bash
# Initialize vault
keyp init

# Store a secret
keyp set github-token
# Enter value: ●●●●●●●●

# Retrieve a secret (copies to clipboard)
keyp get github-token

# List all secrets
keyp list

# Sync across machines
keyp sync
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

- 📖 **[API Reference](./docs/API.md)** - Complete library API with examples
- 🔐 **[Security Guide](./docs/SECURITY.md)** - Cryptographic details and threat model
- 📋 **[Vault Format](./docs/VAULT_FORMAT.md)** - Technical specification of vault file structure

## Roadmap

**Week 1: Core encryption + vault management** ✅
- [x] Core encryption implementation (AES-256-GCM)
- [x] Vault initialization and management
- [x] Secret CRUD operations
- [x] Comprehensive tests (39 passing)
- [x] Security documentation

**Week 2: CLI Commands** 🚧
- [ ] Beautiful CLI with colors and prompts
- [ ] `keyp init` command
- [ ] `keyp set` / `keyp get` commands
- [ ] `keyp list` command

**Week 3: Git sync + polish** 📅
- [ ] Git integration for encrypted backups
- [ ] Clipboard support with auto-clear
- [ ] Secret categories and search
- [ ] Error handling refinement

**Week 4: v1.0.0 launch** 📅
- [ ] Full documentation and examples
- [ ] Launch announcement
- [ ] Community feedback

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
