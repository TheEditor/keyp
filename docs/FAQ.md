# Frequently Asked Questions

## General Questions

### What is keyp?

keyp is a local-first, open-source secret manager for developers. It lets you securely store API keys, passwords, and tokens using AES-256-GCM encryption with PBKDF2 key derivation.

### How is keyp different from password managers?

| Feature | keyp | Password Manager |
|---------|------|------------------|
| **Storage** | Local only | Cloud + local |
| **Complexity** | Very simple | Feature-rich |
| **Target** | Developers | General users |
| **CLI** | Yes | Usually no |
| **Automation** | Easy with scripts | Difficult |
| **Learning curve** | Minutes | Hours |

keyp is designed for developers who want simple, scriptable secret management.

### Is keyp free?

Yes! keyp is free and open-source under the MIT license.

### Can I use keyp commercially?

Yes. The MIT license allows commercial use, modification, and distribution.

## Security Questions

### How secure is keyp?

keyp uses industry-standard encryption:
- AES-256-GCM for authenticated encryption
- PBKDF2-SHA256 with 100,000+ iterations for key derivation
- Cryptographically secure random number generation

The security is limited primarily by your master password strength.

### Can you recover my password?

No, and that's intentional. Your master password is never stored or transmitted. Only you know it, which means:
- Maximum security (no backdoor)
- No recovery option (if you forget it, data is lost)

Always store your password securely!

### What if my computer gets hacked?

- **Vault file:** Encrypted, useless without master password
- **Clipboard:** Auto-clears after 45 seconds
- **Memory:** Cleared when vault is locked
- **Backup:** Git backup is encrypted; no benefit to attacker

Recommendations:
- Keep your OS and software updated
- Use antivirus/malware protection
- Don't run keyp on untrusted machines
- Lock vault when leaving computer

### Does keyp phone home?

No. keyp is completely offline. It:
- Doesn't make network requests
- Doesn't collect telemetry
- Doesn't require authentication
- Doesn't phone home

Only Git sync connects to the network (to your chosen Git provider).

### Can the maintainers access my secrets?

No. keyp is open-source, runs locally, and never shares data. The maintainers:
- Never receive your secrets
- Cannot decrypt your vault
- Have no access to your system
- Have no backdoor

### Is keyp audited?

keyp hasn't had a professional security audit, but:
- All code is open-source and reviewable
- Cryptography uses Node.js built-in crypto module
- No external crypto dependencies
- Community contributions welcome

## Installation Questions

### What versions of Node.js are supported?

Node.js 18.0.0 or higher.

To check your version:
```bash
node --version
```

To update:
```bash
# Using nvm (recommended)
nvm install 18

# Or visit https://nodejs.org/
```

### How do I install keyp on Windows?

```bash
# Using npm
npm install -g @theeditor/keyp

# Then verify
keyp --version
```

If command not found, restart PowerShell or Command Prompt.

### How do I install keyp on macOS?

```bash
# Using npm
npm install -g @theeditor/keyp

# Or using Homebrew (when available)
brew install keyp

# Then verify
keyp --version
```

### How do I install keyp on Linux?

```bash
# Using npm
npm install -g @theeditor/keyp

# For clipboard support, install xclip
sudo apt-get install xclip  # Ubuntu/Debian
sudo dnf install xclip      # Fedora
sudo pacman -S xclip        # Arch

# Then verify
keyp --version
```

## Usage Questions

### How do I create a vault?

```bash
keyp init
```

You'll be prompted to create a master password. This creates an encrypted vault at `~/.keyp/vault.json`.

### How do I store a secret?

```bash
keyp set my-secret
# or
keyp set my-secret "my-value"
```

### How do I retrieve a secret?

```bash
keyp get my-secret  # Copies to clipboard
keyp get my-secret --stdout  # Prints to terminal
```

### How do I list all secrets?

```bash
keyp list
keyp list --search pattern  # Search
keyp list --count           # Count only
```

### How do I delete a secret?

```bash
keyp delete my-secret       # Prompts for confirmation
keyp delete my-secret -f    # Force (no confirmation)
```

### Can I use special characters in secret names?

Yes! Secret names support:
- Letters: `a-z`, `A-Z`
- Numbers: `0-9`
- Special characters: `-`, `_`, `.`, `:` (recommended)
- Unicode characters: any language

Example:
```bash
keyp set prod-api-key-v2 value
keyp set db.host.prod value
keyp set 生产-api-密钥 value
```

### How do I rename a secret?

```bash
keyp rename old-name new-name
```

### How do I copy a secret?

```bash
keyp copy source-name dest-name
```

### How do I back up my vault?

```bash
# Export encrypted
keyp export backup.keyp

# Export plaintext (less secure)
keyp export backup.json --plain

# Restore
keyp import backup.keyp
```

## Git Sync Questions

### What is Git sync?

Git sync lets you back up your encrypted vault to GitHub, GitLab, or any Git provider. Benefits:
- Multi-machine sync
- Disaster recovery
- Version history
- No cloud vendor lock-in

### Is my vault secure in Git?

Yes! Your vault is always encrypted before pushing to Git:
- Only encrypted file is stored
- Your password never leaves your computer
- GitHub/GitLab cannot see your secrets

### How do I set up Git sync?

```bash
# Initialize with GitHub
keyp sync init https://github.com/username/keyp-backup.git

# Or with SSH (recommended)
keyp sync init git@github.com:username/keyp-backup.git

# Push your vault
keyp sync push

# On another machine, pull
keyp sync init git@github.com:username/keyp-backup.git
keyp sync pull
```

### What if I forget my password?

Your vault is encrypted with your password. If you forget:
- You cannot recover it
- No one can help (not even keyp maintainers)
- You must start fresh: `rm ~/.keyp/vault.json && keyp init`

**Prevention:**
- Store your password securely
- Consider a password manager for your master password
- Write it down and lock it in a safe

### Can I use keyp with multiple teams?

Not directly. Each vault is personal. For team secrets:
- Each person creates their own vault
- Each vault backs up to shared Git repository
- Conflicts resolved via `--strategy keep-local` or `--keep-remote`
- Alternative: one team member shares credentials (less secure)

### Is Git sync required?

No! Git sync is optional:
- Works offline without it
- Backups optional
- Multi-machine sync optional

Use just for local development if you prefer.

## Performance Questions

### Is keyp slow?

Typical operations:
- Initialize: ~200ms
- Unlock: ~100-200ms (intentionally slow for security)
- Get secret: ~50ms
- Store secret: ~100ms

These are acceptable for CLI operations.

### Why is unlock slow?

Intentional! Password derivation uses 100,000+ PBKDF2 iterations to prevent brute-force attacks. This is a security feature.

### Can I speed it up?

No, nor should you - it compromises security. If it's too slow, consider:
- Using shell aliases to reduce frequency
- Keeping vault unlocked for batch operations
- Using automation scripts

## Platform Questions

### Does keyp work on Windows?

Yes! Fully supported on Windows 10/11.

### Does keyp work on macOS?

Yes! Fully supported on macOS Intel and Apple Silicon.

### Does keyp work on Linux?

Yes! Fully supported on Ubuntu, Fedora, Arch, Debian, and other distributions.

### Does keyp work on WSL?

Yes! Works great on Windows Subsystem for Linux.

### Can I use keyp on Android/iOS?

Currently only CLI on desktop/server. Mobile support planned for future versions.

## Privacy Questions

### Does keyp collect data?

No. keyp:
- Doesn't track you
- Doesn't collect telemetry
- Doesn't send network requests (except Git sync)
- Doesn't require authentication
- Doesn't report usage

### Is keyp open-source?

Yes! Full source code available at https://github.com/TheEditor/keyp

You can review all code to verify security and privacy claims.

## Contributing Questions

### Can I contribute?

Yes! We welcome contributions:
- Bug reports and feature requests
- Code contributions
- Documentation improvements
- Translations
- Security reviews

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

### How do I report a bug?

Create an issue on GitHub with:
- Clear description
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Node version, keyp version)

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

### How do I request a feature?

Create an issue on GitHub with:
- Clear description of desired feature
- Use case / why you need it
- Possible implementation ideas

### How do I report a security vulnerability?

**Do not create a public issue.** See [SECURITY.md](../SECURITY.md) for responsible disclosure process.

## Comparison Questions

### How does keyp compare to `pass`?

| Feature | keyp | pass |
|---------|------|------|
| **Language** | Node.js | Shell |
| **Encryption** | AES-256-GCM | GPG |
| **Learning curve** | Minutes | Hours |
| **Git integration** | Built-in | Via plugin |
| **GUI** | No | No |
| **Backup** | Easy | Manual |

keyp is simpler and more automated; pass is more mature and flexible.

### How does keyp compare to 1Password?

| Feature | keyp | 1Password |
|---------|------|-----------|
| **Cost** | Free | Paid |
| **Cloud** | No | Yes |
| **Teams** | No | Yes |
| **Mobile** | No | Yes |
| **Advanced features** | No | Yes |
| **CLI** | Yes | Yes |

1Password is enterprise-ready; keyp is for developers.

### How does keyp compare to Bitwarden?

| Feature | keyp | Bitwarden |
|---------|------|-----------|
| **Cost** | Free | Free/Paid |
| **Cloud** | No | Yes (optional self-hosted) |
| **Encryption** | AES-256-GCM | AES-256 CBC |
| **CLI** | Yes | Yes |
| **Mobile** | No | Yes |
| **GUI** | No | Yes |

Bitwarden is more full-featured; keyp is simpler for developers.

## Still Have Questions?

- 📖 **Documentation:** Check out [CLI.md](./CLI.md), [GIT_SYNC.md](./GIT_SYNC.md), [EXAMPLES.md](./EXAMPLES.md)
- 🐛 **Issues:** Browse [GitHub Issues](https://github.com/TheEditor/keyp/issues)
- 💬 **Discussions:** Join [GitHub Discussions](https://github.com/TheEditor/keyp/discussions)
- 📝 **Troubleshooting:** See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)

Happy key managing! 🔒
