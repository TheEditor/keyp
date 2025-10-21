# keyp Git Sync Guide

Git sync allows you to securely back up your encrypted vault to Git repositories like GitHub, GitLab, or self-hosted Git servers. Your vault remains encrypted at rest and in transit.

## Quick Start

### 1. Initialize Git Sync

```bash
# Initialize git sync with your remote repository
keyp sync init https://github.com/username/keyp-backup.git

# Or with SSH (recommended for security)
keyp sync init git@github.com:username/keyp-backup.git
```

### 2. Push Your Vault

```bash
# Commit and push vault to remote
keyp sync push

# Or with a custom commit message
keyp sync push -m "Update secrets: added new API key"
```

### 3. Check Sync Status

```bash
# View synchronization status
keyp sync status
```

### 4. Pull Updates

```bash
# Pull latest vault from remote
keyp sync pull
```

## Features

### Encrypted Backups

- Your vault file is encrypted with AES-256-GCM before syncing
- Passwords are never stored in Git
- All encryption keys are derived from your master password
- Git history contains only encrypted data

### Multi-Machine Synchronization

Keep your secrets in sync across multiple machines:

```bash
# Machine A: Push new secret
keyp set github-token abc123xyz
keyp sync push

# Machine B: Pull the new secret
keyp sync pull
keyp get github-token  # Returns: abc123xyz
```

### Conflict Detection and Resolution

When pulling from remote with local changes:

```bash
# Automatic resolution - keep local changes
keyp sync pull --strategy keep-local

# Automatic resolution - accept remote changes
keyp sync pull --strategy keep-remote

# Manual resolution - review conflicts first
keyp sync pull
# Then manually merge if needed
```

## Setup Instructions

### GitHub

1. **Create Private Repository**
   ```bash
   # Visit https://github.com/new
   # Create repository (can be empty)
   # Choose "Private" for security
   # Do NOT add README, .gitignore, or license
   ```

2. **Initialize Sync**
   ```bash
   keyp sync init https://github.com/username/keyp-backup.git
   ```

3. **Configure Authentication**

   **Option A: HTTPS with Personal Access Token (Recommended for CI/CD)**
   ```bash
   # Generate PAT at: https://github.com/settings/tokens/new
   # Scopes needed: repo (full control of private repositories)

   # Git will prompt for credentials on first push
   keyp sync push

   # Save credentials (one-time setup):
   git config credential.helper store  # Linux/Mac
   # or
   git config credential.helper osxkeychain  # macOS
   git config credential.helper wincred  # Windows
   ```

   **Option B: SSH (Recommended for local development)**
   ```bash
   # Generate SSH key if you don't have one:
   ssh-keygen -t ed25519 -C "keyp-backup"

   # Add public key to GitHub:
   # https://github.com/settings/keys

   # Use SSH URL:
   keyp sync init git@github.com:username/keyp-backup.git
   ```

### GitLab

```bash
# SSH method (recommended)
keyp sync init git@gitlab.com:username/keyp-backup.git

# Or HTTPS method
keyp sync init https://gitlab.com/username/keyp-backup.git
```

### Self-Hosted Git

```bash
# With SSH
keyp sync init git@git.example.com:username/keyp-backup.git

# Or HTTPS
keyp sync init https://git.example.com/username/keyp-backup.git
```

## Command Reference

### `keyp sync init <remote-url>`

Initialize Git sync with a remote repository.

**Options:**
- `-a, --auto-push` - Enable automatic push on vault changes
- `-c, --auto-commit` - Enable automatic commit on vault changes

**Example:**
```bash
keyp sync init https://github.com/user/backup.git
keyp sync init git@github.com:user/backup.git --auto-push
```

### `keyp sync push [options]`

Push encrypted vault to remote repository.

**Options:**
- `-m, --message <message>` - Custom commit message

**Example:**
```bash
keyp sync push
keyp sync push -m "Add production database credentials"
```

### `keyp sync pull [options]`

Pull vault from remote repository.

**Options:**
- `-s, --strategy <strategy>` - Conflict resolution strategy (`keep-local` or `keep-remote`)
- `--auto-resolve` - Automatically resolve conflicts

**Example:**
```bash
keyp sync pull
keyp sync pull --strategy keep-local
keyp sync pull --auto-resolve --strategy keep-remote
```

### `keyp sync status`

Display current synchronization status.

**Output:**
```
Status: ✓ Synced
Last sync: 2h ago
Uncommitted changes: No
Unpushed commits: 0
Conflicts: 0
```

### `keyp sync config`

View or configure sync settings.

**Options:**
- `--auto-push <enabled>` - Enable/disable auto-push (true/false)
- `--auto-commit <enabled>` - Enable/disable auto-commit (true/false)

**Example:**
```bash
keyp sync config                           # Show current settings
keyp sync config --auto-push true          # Enable auto-push
keyp sync config --auto-commit false       # Disable auto-commit
```

## Workflows

### Daily Backup

```bash
# After creating/updating secrets
keyp set new-api-key secret123
keyp sync push -m "Added API key for service X"

# Check sync status
keyp sync status
```

### Multi-Machine Sync

**Setup on Machine A (initial):**
```bash
keyp sync init https://github.com/user/backup.git
keyp sync push -m "Initial vault backup"
```

**Setup on Machine B:**
```bash
# Clone the synced vault
keyp sync init https://github.com/user/backup.git
keyp sync pull

# Now both machines have the same secrets
keyp list
```

**Ongoing sync:**
```bash
# Before leaving machine A
keyp sync push

# After returning to machine B
keyp sync pull

# Create new secret and push
keyp set mobile-api-key xyz789
keyp sync push
```

### Disaster Recovery

```bash
# Fresh machine setup
keyp init  # Create new vault with same password

# Pull backed up secrets from remote
keyp sync pull

# All secrets restored!
keyp list
```

## Conflict Resolution Strategies

### Keep Local
Keeps your local vault changes and ignores remote changes.

```bash
keyp sync pull --strategy keep-local --auto-resolve
```

**Best for:** You want to preserve local edits

### Keep Remote
Accepts remote vault and discards local changes.

```bash
keyp sync pull --strategy keep-remote --auto-resolve
```

**Best for:** Syncing a newly setup machine

### Manual Resolution
Review conflicts before resolving.

```bash
# This will fail and show conflicts
keyp sync pull

# Manually edit and resolve, then retry
keyp sync pull --auto-resolve --strategy keep-local
```

## Security Considerations

### What's Encrypted?

✅ **Encrypted:**
- All secret names and values
- Vault file contents
- Backup files in Git

❌ **Not Encrypted:**
- Git repository metadata (commits, history)
- File timestamps and paths
- Sync configuration (remote URL, branch name)

### Best Practices

1. **Use Private Repositories**
   - Always keep backup repositories private
   - Restrict access to collaborators only

2. **Use Strong Passwords**
   ```bash
   keyp init
   # Use complex master password (12+ characters, mixed case, numbers, symbols)
   ```

3. **Protect SSH Keys**
   - Use SSH key passphrase for added security
   - Store SSH keys in secure location
   - Never commit SSH keys to Git

4. **Regular Verification**
   ```bash
   # Regularly verify backups are working
   keyp sync status
   keyp sync push  # Test connectivity
   ```

5. **Monitor Access**
   - Review GitHub/GitLab access logs
   - Use branch protection rules (if available)
   - Set up notifications for suspicious activity

### Password Security

The master password is the only encryption key needed. It's:
- Never transmitted to remote
- Never stored in Git
- Used to derive vault encryption key via PBKDF2-SHA256
- Required for every unlock and sync operation

### Key Rotation

To change your encryption password:

```bash
# Current workflow (no direct password change)
1. Export secrets (encrypted)
   keyp export backup.keyp

2. Reinitialize vault
   keyp init  # Creates new vault with new password

3. Import secrets
   keyp import backup.keyp

4. Sync to remote
   keyp sync push -m "Rotated master password"
```

## Troubleshooting

### Connection Errors

**Problem:** `Failed to push: Permission denied`

**Solutions:**
- Verify SSH key is added to GitHub/GitLab
- Check SSH config: `ssh -T git@github.com`
- Verify remote URL: `keyp sync config`

### Authentication Errors

**Problem:** `Failed to push: Authentication failed`

**Solutions:**
- Verify credentials are correct
- For SSH: Generate new key pair
- For HTTPS: Check Personal Access Token hasn't expired
- Test connectivity: `git push origin main` (if in keyp directory)

### Merge Conflicts

**Problem:** Conflicting changes between machines

**Solutions:**
```bash
# Option 1: Keep local changes
keyp sync pull --strategy keep-local --auto-resolve

# Option 2: Accept remote version
keyp sync pull --strategy keep-remote --auto-resolve

# Option 3: Manual merge (for advanced users)
cd ~/.keyp
git status  # See conflicts
# Edit and resolve conflicts manually
git add vault.json.encrypted
git commit -m "Resolved merge conflicts"
```

### Git Not Initialized

**Problem:** `Git repository not initialized`

**Solutions:**
```bash
keyp sync init <remote-url>  # Initialize first
```

## Advanced Usage

### Custom Commit Messages

```bash
keyp sync push -m "Added AWS credentials for production"
```

### Viewing Git Log

```bash
cd ~/.keyp
git log --oneline -10  # View last 10 commits
git show HEAD  # View latest commit details
```

### Manual Git Operations

```bash
# If needed, you can use git directly on ~/.keyp directory
cd ~/.keyp
git status
git log
git branch -a
```

### SSH Key Management

```bash
# Generate SSH key with custom path
ssh-keygen -t ed25519 -C "keyp" -f ~/.ssh/keyp_ed25519

# Add to ssh-agent
ssh-add ~/.ssh/keyp_ed25519

# Test connection
ssh -T git@github.com
```

## FAQ

**Q: Is my password sent to GitHub?**
A: No. Your password never leaves your machine and is never transmitted.

**Q: Can GitHub employees access my secrets?**
A: No. The vault file is encrypted. GitHub only stores encrypted data.

**Q: Can I use the same backup repo on multiple machines?**
A: Yes. Each machine can push/pull from the same repo, keeping vaults in sync.

**Q: What if I lose my master password?**
A: You cannot recover it. However, you can restore from an encrypted backup and reinitialize with a new password.

**Q: How often should I sync?**
A: After adding/updating important secrets. Or configure auto-sync in `keyp config`.

**Q: Can I use multiple remote repositories?**
A: Currently, keyp supports one remote (`origin`). You can manually add more Git remotes if needed.

**Q: Is there versioning/history?**
A: Yes. Git stores complete history of all pushes. Use `cd ~/.keyp && git log` to view.

**Q: How large can my vault be?**
A: Theoretically unlimited. GitHub allows repos up to 100GB. Performance may vary with very large vaults.

## Next Steps

- [API Reference](./API.md) - Library usage guide
- [Security Details](./SECURITY.md) - Cryptographic analysis
- [Vault Format](./VAULT_FORMAT.md) - Technical specification
- [CLI Reference](./CLI.md) - All CLI commands
