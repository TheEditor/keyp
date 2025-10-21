# keyp v0.2.0 Launch Announcement

## Introducing keyp: Local-First Secret Management for Developers

We're excited to announce **keyp v0.2.0** — a powerful yet simple secret manager designed specifically for developers who want secure, scriptable credential storage without the complexity of enterprise solutions.

### What is keyp?

keyp is the secret manager you've been looking for. It's:

- **Simple** — Initialize a vault, store secrets, retrieve them. Done.
- **Secure** — AES-256-GCM encryption with PBKDF2 key derivation
- **Local-First** — Your secrets stay on your machine (unless you sync them)
- **Developer-Friendly** — Perfect for CLI automation, scripts, and workflows
- **Open Source** — MIT licensed, fully auditable code

Think of it as "pass for the Node.js generation" — familiar simplicity with modern tooling.

### Why keyp?

If you're tired of...

- ❌ Password managers that don't fit developer workflows
- ❌ Cloud vendors asking for access to your credentials
- ❌ Complex enterprise solutions for simple needs
- ❌ Grepping through terminal history for API keys

...then keyp is for you.

### Key Features

#### 🔒 Core Features

```bash
# Initialize vault (one-time setup)
keyp init

# Store secrets
keyp set github-token ghp_abc123xyz
keyp set api-key sk_live_12345

# Retrieve to clipboard
keyp get github-token

# List all secrets
keyp list

# Search
keyp list --search api

# Manage
keyp delete api-key
keyp rename old-name new-name
keyp copy prod-key staging-key
```

#### 🌐 Git Sync

Back up your encrypted vault to GitHub, GitLab, or any Git provider:

```bash
# Initialize sync
keyp sync init https://github.com/username/keyp-backup.git

# Push your vault
keyp sync push

# Pull on another machine
keyp sync pull
```

#### 📊 Statistics & Configuration

```bash
# See vault statistics
keyp stats

# Manage configuration
keyp config
keyp config set clipboard-timeout 120
```

#### 🐚 Shell Completion

Tab completion for bash and zsh — makes keyp faster and more enjoyable to use.

### Security You Can Trust

- **AES-256-GCM:** Military-grade authenticated encryption
- **PBKDF2:** 100,000+ iterations for password derivation
- **No Backdoor:** Open-source code you can review
- **No Cloud:** Your secrets don't leave your machine
- **No Tracking:** No telemetry, no analytics, no phoning home

### Perfect For

✅ **Individual Developers**
- Store all your API keys, tokens, and credentials securely
- Use across multiple machines via Git sync
- Scriptable for automation

✅ **Development Teams**
- Each developer has their own encrypted vault
- Shared backup repository (vault is encrypted)
- Perfect for multi-environment management (dev/staging/prod)

✅ **CI/CD Pipelines**
- Easy integration with GitHub Actions, GitLab CI, etc.
- Perfect for deployment secrets
- Simple enough to understand and audit

✅ **DevOps & Infrastructure**
- Secure credential storage for scripts
- Easy automation and orchestration
- No dependency on external services

### Getting Started

#### Installation

```bash
npm install -g @theeditor/keyp
```

#### Quick Start

```bash
# Initialize your vault
keyp init
# → Set your master password

# Store a secret
keyp set my-api-key secret-value-123

# Retrieve it
keyp get my-api-key
# → Copied to clipboard!

# List all secrets
keyp list
```

#### Next Steps

- Read the [CLI Reference](../docs/CLI.md) for all commands
- Check [Git Sync Guide](../docs/GIT_SYNC.md) to sync across machines
- See [Examples](../docs/EXAMPLES.md) for real-world workflows
- Explore [Troubleshooting](../docs/TROUBLESHOOTING.md) for help

### What's New in v0.2.0

This release adds powerful synchronization and polish:

**🌐 Git Sync**
- Push encrypted vault to any Git provider
- Multi-machine synchronization
- Conflict detection and resolution
- Backup and disaster recovery

**✨ Polish**
- Enhanced password strength validation with visual meter
- Shell completion scripts for bash and zsh
- `keyp stats` for vault statistics
- `keyp config` for settings management
- Improved error messages with helpful hints

### Real-World Examples

**Store Development Credentials**
```bash
keyp set db-host localhost
keyp set db-user developer
keyp set db-password dev_password
```

**Use in Scripts**
```bash
#!/bin/bash
API_KEY=$(keyp get api-key --stdout)
curl -H "Authorization: Bearer $API_KEY" https://api.example.com
```

**Shell Aliases**
```bash
alias getdb='keyp get db-password && echo'
alias getapi='keyp get api-key && echo'
```

**Sync to GitHub**
```bash
# First machine
keyp sync init https://github.com/you/keyp-backup.git
keyp sync push

# Second machine
keyp sync init https://github.com/you/keyp-backup.git
keyp sync pull
```

### Roadmap

We're committed to continuous improvement:

- **v0.3.0** — Performance optimizations and cross-platform CI/CD
- **v1.0.0** — Stable API, comprehensive documentation, community feedback integration

### Open Source & Community

keyp is fully open-source under the MIT license:

- 📖 **[GitHub Repository](https://github.com/TheEditor/keyp)**
- 🐛 **[Report Issues](https://github.com/TheEditor/keyp/issues)**
- 💡 **[Request Features](https://github.com/TheEditor/keyp/discussions)**
- 🤝 **[Contribute](../docs/CONTRIBUTING.md)**

### Why Now?

We realized there was a gap in the tooling landscape:

- **Too Simple:** `pass` is great but not designed for developers
- **Too Complex:** Enterprise password managers are overkill
- **Too Cloud-Dependent:** We want to own our credentials
- **Not Scriptable:** We need to automate without GUI friction

keyp fills that gap perfectly.

### Try It Today

```bash
npm install -g @theeditor/keyp
keyp init
keyp --help
```

### Questions?

- 📖 **Read the docs:** Full documentation in the repo
- 💬 **Start a discussion:** GitHub Discussions
- 🐛 **Report a bug:** GitHub Issues
- 🔐 **Report a vulnerability:** See SECURITY.md

---

## Comparison with Alternatives

### vs. `pass`

| Feature | keyp | pass |
|---------|------|------|
| Language | Node.js/JavaScript | Shell/GPG |
| Learning Curve | Minutes | Hours |
| Git Integration | Built-in | Via plugin |
| Cross-platform | Excellent | Good |
| For Developers | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

### vs. 1Password/Bitwarden

| Feature | keyp | Enterprise |
|---------|------|-----------|
| Cost | Free | Paid |
| Complexity | Simple | Full-featured |
| CLI Quality | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| For Teams | Not designed | Excellent |
| For Solo Devs | ⭐⭐⭐⭐⭐ | Overkill |

### Bottom Line

- **Choose keyp** if you're a developer who values simplicity and control
- **Choose pass** if you prefer traditional Unix tools
- **Choose enterprise solutions** if you manage many users with team features

---

## The keyp Philosophy

keyp is built on these principles:

1. **Keep It Simple** — One command does one thing well
2. **Security First** — Never compromise on encryption or key derivation
3. **Developer Focused** — Built for how we actually work
4. **Local by Default** — Your data is yours to control
5. **Open Source** — Transparency builds trust

---

**Ready to get started?**

```bash
npm install -g @theeditor/keyp
```

**Questions or feedback?** Open an issue on [GitHub](https://github.com/TheEditor/keyp/issues).

**Star the repo** if you find keyp useful! ⭐

---

*keyp: Keep your keys. Keep them safe. Keep them simple.* 🔒
