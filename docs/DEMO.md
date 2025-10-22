# keyp Live Demo Script

A complete demonstration of keyp's features and workflows. You can follow along or use this as a reference.

**Duration:** ~5 minutes

---

## Demo 1: Basic Secret Management (2 minutes)

### Initialize Vault

```bash
$ keyp init

🔒 keyp
────────────────────────────────────

ℹ Creating a new vault...

Enter master password: ●●●●●●●●●●●●
Confirm master password: ●●●●●●●●●●●●

Strength: ██████████ 95%
  ✓ Good length (13 characters)
  ✓ Uppercase letters included
  ✓ Lowercase letters included
  ✓ Numbers included
  ✓ Special characters included

✓ Vault initialized successfully!

ℹ Location: ~/.keyp/vault.json

Next steps:
  1. keyp set <secret-name>   - Store your first secret
  2. keyp list                 - List all secrets
  3. keyp get <secret-name>    - Retrieve a secret
```

**Screenshot placeholder:** [Show vault initialized screen]

### Store Some Secrets

```bash
$ keyp set github-token
Enter master password: ●●●●●●●●
Enter value for "github-token": ●●●●●●●●●●●●●●●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret "github-token" saved
ℹ Total secrets: 1
ℹ Retrieve with: keyp get github-token

$ keyp set api-key sk_live_prod_abc123xyz
Enter master password: ●●●●●●●●
✓ Secret "api-key" saved
ℹ Total secrets: 2

$ keyp set db-password db_secure_pass_123
Enter master password: ●●●●●●●●
✓ Secret "db-password" saved
ℹ Total secrets: 3
```

### List Secrets

```bash
$ keyp list
Enter master password: ●●●●●●●●

  • api-key
  • db-password
  • github-token

3 secrets stored
```

### Retrieve Secret (Clipboard)

```bash
$ keyp get github-token
Enter master password: ●●●●●●●●
✓ Copied to clipboard (clears in 45 seconds)

# Now paste: Ctrl+V or Cmd+V
# Secret is now in your clipboard, automatically cleared in 45 seconds
```

**Screenshot placeholder:** [Show clipboard feedback]

### Search Secrets

```bash
$ keyp list --search "key"
Enter master password: ●●●●●●●●

  • api-key
  • github-token

2 secrets stored
```

---

## Demo 2: Managing Secrets (1.5 minutes)

### Rename a Secret

```bash
$ keyp rename api-key production-api-key
Enter master password: ●●●●●●●●
✓ Secret renamed: api-key → production-api-key
```

### Copy a Secret

```bash
$ keyp copy production-api-key staging-api-key
Enter master password: ●●●●●●●●
✓ Secret copied: production-api-key → staging-api-key
```

### Delete a Secret

```bash
$ keyp delete staging-api-key
Enter master password: ●●●●●●●●
Delete secret "staging-api-key"? (y/N): y
✓ Secret "staging-api-key" deleted
Remaining secrets: 3
```

### Export for Backup

```bash
$ keyp export backup.keyp
Enter master password: ●●●●●●●●
✓ Exported 3 secrets to backup.keyp
ℹ Backup is encrypted and safe to store
```

### Import Secrets

```bash
$ keyp import backup.keyp
Enter master password: ●●●●●●●●
ℹ Import mode: merge (add/update existing)
ℹ Ready to import 3 secrets

Continue? (y/N): y
✓ Successfully imported 3 secrets
```

---

## Demo 3: Vault Statistics & Configuration (1 minute)

### View Vault Statistics

```bash
$ keyp stats
Enter master password: ●●●●●●●●

📊 Vault Statistics
────────────────────────────────────

Secrets
  Total: 3
  Average value length: 24 characters
  Longest name: production-api-key

Storage
  Vault file size: 2.34 KB
  Location: ~/.keyp/vault.json

Dates
  Last modified: 10/21/2025, 2:45:00 PM
  Last synced: Never (Git sync not configured)

Encryption
  Algorithm: AES-256-GCM
  Key derivation: PBKDF2-SHA256
  Iterations: 100,000+
```

**Screenshot placeholder:** [Show stats screen]

### View Configuration

```bash
$ keyp config

⚙️  Keyp Configuration
────────────────────────────────────

Clipboard
  Timeout: 45 seconds
  (how long before clipboard is auto-cleared)

Vault
  Auto-lock: disabled
  (lock vault after inactivity)

Git Sync
  Auto-sync: disabled
  (automatically push on vault changes)

Configuration stored in: ~/.keyp/.keyp-config.json

# Change configuration
$ keyp config set clipboard-timeout 120
✓ Clipboard timeout set to 120 seconds
```

---

## Demo 4: Git Synchronization (Optional - 2 minutes)

### Initialize Git Sync

```bash
$ keyp sync init https://github.com/username/keyp-backup.git
ℹ Initializing Git sync...
✓ Git repository initialized
✓ Remote configured: https://github.com/username/keyp-backup.git
✓ Git sync initialized successfully!
```

### Push Vault to GitHub

```bash
$ keyp sync push
ℹ Pushing vault to remote...
ℹ Committing vault changes...
✓ Vault pushed to remote
ℹ Backup is encrypted and secure

# You can now view the backup on GitHub
# (only encrypted file is visible)
```

### Check Sync Status

```bash
$ keyp sync status

Status: ✓ Synced
Last sync: just now
Uncommitted changes: No
Unpushed commits: 0
Conflicts: 0
```

### Pull on Another Machine

```bash
# On second machine
$ keyp sync init https://github.com/username/keyp-backup.git
ℹ Initializing Git sync...
✓ Git repository initialized
✓ Remote configured: https://github.com/username/keyp-backup.git
✓ Git sync initialized successfully!

$ keyp sync pull
ℹ Pulling vault from remote...
✓ Vault pulled from remote

# Now second machine has all the same secrets!
$ keyp list
Enter master password: ●●●●●●●●

  • db-password
  • github-token
  • production-api-key

3 secrets stored
```

---

## Demo 5: Using in Scripts (Optional - 1 minute)

### Get Secret to Variable

```bash
#!/bin/bash
API_KEY=$(keyp get production-api-key --stdout)
echo "Using API Key: ${API_KEY:0:10}..."

# Use in curl
curl -H "Authorization: Bearer $API_KEY" \
  https://api.example.com/v1/status

# Output: {"status": "ok", "version": "1.2.0"}
```

### Shell Alias

```bash
# Add to ~/.bashrc or ~/.zshrc
alias get-api='keyp get production-api-key --stdout'

# Then use it
$ get-api
sk_live_prod_abc123xyz
```

---

## Key Points to Highlight

### Security ✓
- Master password never stored or transmitted
- AES-256-GCM encryption (military-grade)
- PBKDF2 key derivation (100,000+ iterations)
- Clipboard auto-clears for safety
- No telemetry or tracking

### Simplicity ✓
- One command does one thing well
- Intuitive command structure
- Clear, helpful error messages
- Fast and responsive

### Developer-Friendly ✓
- Perfect for scripts and automation
- Shell completion for faster typing
- Works across all platforms (Windows, macOS, Linux)
- Open-source and auditable

### Multi-Machine ✓
- Git sync for encrypted backups
- Easy to sync across machines
- Version control for secrets
- No cloud vendor lock-in

---

## Common Questions During Demo

**Q: Is my password stored?**
A: No! Your password is never stored. It's used only to derive your encryption key.

**Q: Can you access my secrets?**
A: No, and I can't recover your password either. That's a feature, not a limitation.

**Q: How is this different from 1Password?**
A: keyp is simpler and local-first. 1Password is more full-featured. Choose based on your needs.

**Q: Can I use keyp at work?**
A: Absolutely! Many developers do. It's perfect for development credentials.

**Q: Is it production-ready?**
A: Yes! keyp uses industry-standard encryption and has been thoroughly tested.

---

## Wrapping Up

**Key Takeaways:**

1. keyp makes secret management simple and secure
2. Perfect for developers who want control and simplicity
3. Git sync enables multi-machine workflows
4. Open-source and fully auditable
5. Free and MIT licensed

**Next Steps:**

```bash
# Get started
npm install -g @theeditor/keyp

# Initialize
keyp init

# Learn more
keyp --help
keyp list --help

# Read the docs
# https://github.com/TheEditor/keyp/docs/
```

**Questions?**

- 📖 [Full Documentation](../README.md#documentation)
- 🐛 [GitHub Issues](https://github.com/TheEditor/keyp/issues)
- 💬 [GitHub Discussions](https://github.com/TheEditor/keyp/discussions)

---

**Thank you for following along!** 🔒

*keyp: Keep your keys. Keep them safe. Keep them simple.*
