# keyp CLI Reference

Complete command-line interface guide for keyp - Local-first secret manager for developers.

## Installation

```bash
npm install -g @theeditor/keyp
```

## Global Options

```bash
keyp --version      # Show version
keyp --help         # Show help
```

---

## Commands

### keyp init

Initialize a new encrypted vault with a master password.

**Usage:**
```bash
keyp init
```

**Interactive prompts:**
- Enter master password (masked)
- Confirm master password (masked)
- Password strength feedback

**Example:**
```bash
$ keyp init
ℹ Creating a new vault...

Enter master password: ●●●●●●●●●●●●●●●●●●
Confirm master password: ●●●●●●●●●●●●●●●●●●

ℹ Password strength: Strong password!

✓ Vault initialized successfully!

ℹ Location: ~/.keyp/vault.json

Next steps:
  1. keyp set <secret-name>   - Store your first secret
  2. keyp list                 - List all secrets
  3. keyp get <secret-name>    - Retrieve a secret
```

**Output:**
- ✓ Success message with vault location
- Next steps guidance
- Password strength feedback

**Errors:**
- "Vault already exists" - if vault is already initialized
- "Password entry cancelled" - if user aborts

---

### keyp set <name> [value]

Store a new secret or update an existing one in the vault.

**Usage:**
```bash
keyp set <name>              # Prompt for value
keyp set <name> <value>      # Provide value as argument
```

**Parameters:**
- `<name>` - Secret name/identifier (required)
- `[value]` - Secret value (optional, prompts if not provided)

**Interactive prompts:**
- Enter master password (masked, with retry on failure)
- Enter value for secret (if not provided as argument, masked)
- Enter master password to save (masked)

**Examples:**
```bash
# Prompt for value
$ keyp set github-token
Enter master password: ●●●●●●●●
Enter value for "github-token": ●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret "github-token" saved
ℹ Total secrets: 1
ℹ Retrieve with: keyp get github-token

# Provide value directly
$ keyp set api-key "sk-1234567890"
Enter master password: ●●●●●●●●
✓ Secret "api-key" saved
ℹ Total secrets: 2

# Update existing secret
$ keyp set github-token "ghp_newtoken123"
Enter master password: ●●●●●●●●
✓ Secret "github-token" saved (updated)
```

**Output:**
- ✓ Success message showing secret was saved/updated
- ℹ Total number of secrets in vault
- ℹ Hint showing how to retrieve the secret

**Errors:**
- "Vault not found" - run `keyp init` first
- "Secret name required" - secret name is mandatory
- "Secret value cannot be empty" - value cannot be empty string
- "Incorrect password" - wrong master password (3 attempts max)
- "Password entry cancelled" - user aborted operation

---

### keyp get <name> [options]

Retrieve a secret from the vault and copy to clipboard.

**Usage:**
```bash
keyp get <name>                # Copy to clipboard (default)
keyp get <name> --stdout       # Print to terminal
keyp get <name> --no-clear     # Don't auto-clear clipboard
```

**Parameters:**
- `<name>` - Secret name to retrieve (required)

**Options:**
- `--stdout` - Print secret to terminal instead of clipboard
  - **Warning:** Secret will be visible on screen
  - Useful for pipes and scripts

- `--no-clear` - Don't auto-clear clipboard after timeout
  - Default behavior: clears after 45 seconds
  - Use this to keep secret in clipboard longer

**Interactive prompts:**
- Enter master password (masked, with retry on failure)

**Examples:**
```bash
# Copy to clipboard (default)
$ keyp get github-token
Enter master password: ●●●●●●●●
✓ Copied to clipboard
ℹ Will clear in 45 seconds

# Print to terminal
$ keyp get github-token --stdout
Enter master password: ●●●●●●●●
⚠ Output to terminal (secret will be visible!)

ghp_abc123xyz789

# Don't auto-clear
$ keyp get api-key --no-clear
Enter master password: ●●●●●●●●
✓ Copied to clipboard
ℹ Will clear in 45 seconds
```

**Output:**
- ✓ Success message when copied to clipboard
- ⚠ Warning when printing to stdout
- ℹ Information about clipboard clearing
- Secret value when using `--stdout`

**Errors:**
- "Vault not found" - run `keyp init` first
- "Secret name required" - secret name is mandatory
- "Secret 'X' not found" - secret doesn't exist
- "Incorrect password" - wrong master password (3 attempts max)
- "Clipboard not available" - falls back to stdout printing

**Security Notes:**
- Default behavior (clipboard copy) is secure - secret doesn't appear on screen
- Clipboard auto-clears after 45 seconds by default
- `--stdout` exposes secret - use only when necessary
- Terminal history may contain `--stdout` calls

---

### keyp list [options]

List all secrets stored in the vault.

**Usage:**
```bash
keyp list                      # List all secrets
keyp list --search <pattern>   # Search for secrets
keyp list --count              # Show only count
```

**Options:**
- `--search <pattern>` - Filter secrets by pattern (substring match, case-insensitive)
  - Shows matching secret names

- `--count` - Show only the total number of secrets
  - Useful for scripts

**Interactive prompts:**
- Enter master password (masked, with retry on failure)

**Examples:**
```bash
# List all secrets
$ keyp list
Enter master password: ●●●●●●●●

  • api-key
  • database-password
  • github-token
  • jwt-secret

4 secrets stored

# Search for secrets
$ keyp list --search github
Enter master password: ●●●●●●●●
ℹ Search results for "github"

  • github-api-key
  • github-token

2 secrets stored

# Count only
$ keyp list --count
Enter master password: ●●●●●●●●
4 secrets
```

**Output:**
- Bullet list of secret names (sorted alphabetically)
- ℹ Search results header (when using --search)
- Total count of secrets
- Works with empty vault ("No secrets yet")

**Errors:**
- "Vault not found" - run `keyp init` first
- "Incorrect password" - wrong master password (3 attempts max)
- "Password entry cancelled" - user aborted

**Notes:**
- Secret names are shown but NOT values (safe operation)
- Results are sorted alphabetically
- Search is case-insensitive (e.g., "github" matches "GitHub-Token")
- Empty vault shows helpful hint: "No secrets yet. Try: keyp set <name>"

---

### keyp delete <name> [options]

Delete a secret from the vault.

**Aliases:** `rm`

**Usage:**
```bash
keyp delete <name>             # Prompt for confirmation
keyp delete <name> -f          # Force delete (skip confirmation)
keyp delete <name> --force     # Force delete (long form)
keyp rm <name>                 # Alias for delete
```

**Parameters:**
- `<name>` - Secret name to delete (required)

**Options:**
- `-f, --force` - Skip confirmation prompt
  - Useful for scripts and automation

**Interactive prompts:**
- Enter master password (masked, with retry on failure)
- Delete confirmation (unless --force flag used)
- Enter master password to save (masked)

**Examples:**
```bash
# Delete with confirmation
$ keyp delete old-token
Enter master password: ●●●●●●●●
Delete secret "old-token"? (y/N): y
Enter master password to save: ●●●●●●●●
✓ Secret "old-token" deleted
ℹ Remaining secrets: 3

# Force delete (no confirmation)
$ keyp delete temporary-secret -f
Enter master password: ●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret "temporary-secret" deleted
ℹ Remaining secrets: 2

# Using alias
$ keyp rm deprecated-api-key -f
Enter master password: ●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret "deprecated-api-key" deleted
ℹ Remaining secrets: 1
```

**Output:**
- Delete confirmation prompt (unless --force)
- ✓ Success message showing deletion
- ℹ Count of remaining secrets

**Errors:**
- "Vault not found" - run `keyp init` first
- "Secret name required" - secret name is mandatory
- "Secret 'X' not found" - secret doesn't exist in vault
- "Incorrect password" - wrong master password (3 attempts max)
- "Password entry cancelled" - user aborted

**Notes:**
- Confirmation defaults to "No" (safe by default)
- Type "y" or "yes" for confirmation
- `-f` flag useful for automation and scripts
- Deletion is permanent - secret cannot be recovered

---

## Common Workflows

### Initialize and Add First Secret

```bash
# 1. Initialize vault
$ keyp init
Enter master password: ●●●●●●●●
Confirm master password: ●●●●●●●●
✓ Vault initialized successfully!

# 2. Add first secret
$ keyp set github-token
Enter master password: ●●●●●●●●
Enter value for "github-token": ●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret "github-token" saved
```

### Store Multiple Secrets

```bash
# Add API key
$ keyp set openai-api-key sk-...
Enter master password: ●●●●●●●●
✓ Secret "openai-api-key" saved

# Add database password
$ keyp set db-password
Enter master password: ●●●●●●●●
Enter value for "db-password": ●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret "db-password" saved
```

### Retrieve and Use Secret in Script

```bash
# Get secret and use in script
$ API_KEY=$(keyp get openai-api-key --stdout)
Enter master password: ●●●●●●●●
✓ Copied to clipboard

# Now use in curl or other commands
$ curl -H "Authorization: Bearer $API_KEY" https://api.openai.com/...
```

### Search and List Secrets

```bash
# List all secrets
$ keyp list
Enter master password: ●●●●●●●●

  • api-key
  • database-url
  • github-token

3 secrets stored

# Search for github-related secrets
$ keyp list --search github
Enter master password: ●●●●●●●●
ℹ Search results for "github"

  • github-api-key
  • github-token

2 secrets stored
```

### Clean Up Old Secrets

```bash
# List secrets to see what to delete
$ keyp list
Enter master password: ●●●●●●●●
  • deprecated-key
  • new-key
  • old-token

# Delete old ones
$ keyp delete deprecated-key -f
Enter master password: ●●●●●●●●
✓ Secret "deprecated-key" deleted

$ keyp delete old-token -f
Enter master password: ●●●●●●●●
✓ Secret "old-token" deleted
```

---

## Password Requirements

### Best Practices

**Strong master password should:**
- ✅ Be at least 20 characters long
- ✅ Mix uppercase, lowercase, numbers, and symbols
- ✅ Be unique (not reused from other services)
- ✅ Be memorized (don't write it down)
- ✅ Avoid dictionary words and patterns

**Examples of good passwords:**
- `Tr0pical!Flamingo#2025$Sunset`
- `P@ssw0rd!Is*VeryLong&Complex`
- `MyDog+Loves3Clouds?Rainbow!`

**Examples of weak passwords:**
- ❌ `password123` (too common)
- ❌ `qwerty` (keyboard pattern)
- ❌ `123456` (sequential numbers)
- ❌ `abc123` (too short)

### keyp Will Warn About

If your password is weak, keyp will show:
```
⚠ Password is weak: Consider: mix in uppercase letters, add some numbers
```

This is just guidance - you can still use the password.

---

## Error Handling

### Common Errors and Solutions

**"Vault not found"**
- Run `keyp init` to initialize a vault first

**"Incorrect password"**
- Double-check your master password
- Vault won't unlock with wrong password
- You have 3 attempts before the operation fails

**"Secret 'X' not found"**
- Check the exact secret name with `keyp list`
- Secret names are case-sensitive
- Search with `keyp list --search <pattern>` to find it

**"Clipboard not available"**
- Try using `--stdout` flag: `keyp get secret --stdout`
- Or install clipboard utility for your OS

---

## Enhanced Commands

### keyp rename <old-name> <new-name>

Rename an existing secret to a new name while preserving its value.

**Usage:**
```bash
keyp rename <old-name> <new-name>
```

**Example:**
```bash
$ keyp rename github-token github-api-token
Enter master password: ●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret renamed: "github-token" → "github-api-token"
ℹ Total secrets: 5
```

**Use cases:**
- Correct typos in secret names
- Reorganize secret naming conventions
- Update names for clarity

**Errors:**
- "Secret 'X' not found" - old name doesn't exist
- "Secret 'X' already exists" - new name is already taken
- "New name must be different" - old and new names are the same

---

### keyp copy <source> <dest>

Copy a secret to a new name (duplicate with new identifier).

**Usage:**
```bash
keyp copy <source> <dest>
```

**Example:**
```bash
$ keyp copy github-token gitlab-token
Enter master password: ●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret copied: "github-token" → "gitlab-token"
ℹ Total secrets: 6
```

**Use cases:**
- Create similar secrets for different services
- Backup important secrets under new names
- Test secret values before permanently replacing them

**Notes:**
- Source secret value is preserved exactly
- Destination must not already exist
- Both secrets are independent after copying

---

### keyp export [output-file]

Export secrets to a file for backup or migration.

**Usage:**
```bash
keyp export                           # Auto-generate filename
keyp export secrets-backup.json       # Specify output file
keyp export secrets.json --plain      # Export as plaintext (unencrypted)
keyp export --stdout                  # Print to stdout
```

**Options:**
- `--plain` - Export as plaintext JSON (WARNING: unencrypted)
  - Use for migration or testing only
  - Keep plaintext exports secure!
  - Not suitable for long-term storage

- `--stdout` - Print to stdout instead of file
  - Useful for pipes and scripts
  - Can be redirected to file: `keyp export --stdout > backup.json`

**Examples:**
```bash
# Encrypted export to file
$ keyp export secrets-backup.json
Enter master password: ●●●●●●●●
✓ Secrets exported (encrypted)
ℹ Location: ~/secrets-backup.json
ℹ Secrets exported: 5

# Plaintext export (with warning)
$ keyp export --plain secrets.json
Enter master password: ●●●●●●●●
⚠ Exporting secrets as PLAINTEXT - this is NOT encrypted!
✓ Secrets exported (plaintext)
ℹ Location: ~/secrets.json
⚠ Remember: This file contains PLAINTEXT secrets - keep it safe!

# Export to stdout for piping
$ keyp export --stdout | gzip > encrypted-backup.gz
Enter master password: ●●●●●●●●
```

**Output:**
- Default: `keyp-export-{timestamp}.json`
- Encrypted exports can be safely stored/backed up
- Plaintext exports should be encrypted or deleted after use

**Use cases:**
- Create backups of vault
- Migrate to different machine
- Share with team (encrypted exports only)
- Version control (encrypted exports only)

---

### keyp import <input-file>

Import secrets from a file (merge or replace mode).

**Usage:**
```bash
keyp import secrets.json              # Import and merge
keyp import secrets.json --dry-run    # Preview changes
keyp import secrets.json --replace    # Replace all existing secrets
```

**Options:**
- `--dry-run` - Show what would be imported without making changes
  - Preview import results
  - Verify file format
  - Check for conflicts

- `--replace` - Delete all existing secrets and import
  - DANGEROUS - use with caution!
  - Requires confirmation before proceeding
  - Useful for full vault migration

**Default behavior:** Merge
- Add new secrets
- Update existing secrets with same names
- Preserve secrets not in import file

**Examples:**
```bash
# Import and merge
$ keyp import backup.json
Enter master password: ●●●●●●●●
ℹ Import summary:
ℹ   New secrets: 3
ℹ   Updated secrets: 2
Enter master password to save: ●●●●●●●●
✓ Secrets imported successfully
ℹ New secrets: 3
ℹ Updated secrets: 2
ℹ Total secrets now: 10

# Dry run preview
$ keyp import backup.json --dry-run
Enter master password: ●●●●●●●●
ℹ Import summary:
ℹ   New secrets: 3
ℹ   Updated secrets: 2
ℹ Dry run mode - no changes made

# Replace all secrets
$ keyp import new-vault.json --replace
Enter master password: ●●●●●●●●
⚠ REPLACE mode: This will delete all existing secrets!
Delete all existing secrets and import? (y/N): y
✓ Secrets imported successfully
ℹ New secrets: 5
ℹ Updated secrets: 0
ℹ Total secrets now: 5
```

**Supported formats:**
- Plaintext JSON: `{"key": "value", ...}`
- Encrypted vault exports (future support)

**Use cases:**
- Restore from backup
- Migrate from another tool
- Merge vaults from team members
- One-time import of bulk secrets

**Safety features:**
- Dry-run mode for preview
- Confirmation for replace mode
- Shows summary of changes
- Preserves existing secrets (by default)

---

### keyp get <name> [options] - Timeout Control

Extended `keyp get` command with clipboard timeout configuration.

**New Option:**
```bash
keyp get <name> --timeout <seconds>   # Custom clear timeout
```

**Example:**
```bash
# Use default 45-second timeout
$ keyp get api-key
Enter master password: ●●●●●●●●
✓ Copied to clipboard
ℹ Will clear in 45 seconds

# Custom timeout (60 seconds)
$ keyp get api-key --timeout 60
Enter master password: ●●●●●●●●
✓ Copied to clipboard
ℹ Will clear in 60 seconds

# Quick timeout (10 seconds)
$ keyp get secret --timeout 10
Enter master password: ●●●●●●●●
✓ Copied to clipboard
ℹ Will clear in 10 seconds

# Disable auto-clear entirely
$ keyp get api-key --no-clear
Enter master password: ●●●●●●●●
✓ Copied to clipboard
ℹ Will clear in 45 seconds
```

**Use cases:**
- Longer timeout for complex copy-paste operations
- Shorter timeout for high-security environments
- Disable clearing for scripts that need clipboard persistence

**Notes:**
- Default: 45 seconds
- Minimum: 1 second (practical minimum ~5s)
- Works with `--no-clear` to disable entirely
- Invalid values fall back to default

---

## Tips & Tricks

### Bash Aliases

```bash
# Add to ~/.bashrc or ~/.zshrc
alias ks='keyp set'
alias kg='keyp get'
alias kl='keyp list'
alias kd='keyp delete'
alias ki='keyp init'
```

Then use:
```bash
$ ks github-token        # equals: keyp set github-token
$ kg github-token        # equals: keyp get github-token
```

### Bash Function for Clipboard with Timeout

```bash
# Get secret and show countdown
kg-countdown() {
  keyp get "$1" --no-clear
  for i in {45..1}; do
    echo -ne "\rClearing in ${i}s...  "
    sleep 1
  done
  echo -ne "\rClipboard cleared              \n"
}

# Use: kg-countdown my-secret
```

### Scripting with keyp

```bash
#!/bin/bash
# Deploy script that uses keyp secrets

API_KEY=$(keyp get production-api-key --stdout)
DB_PASS=$(keyp get db-production --stdout)

# Use secrets
export API_KEY
export DB_PASSWORD="$DB_PASS"

# Run deployment
npm run deploy
```

---

## Security Notes

### Clipboard Clearing

- Default: Clipboard clears after 45 seconds
- Only affects the secret value, not other clipboard content
- Won't clear if another app changes clipboard before timeout

### Password Storage

- Master password is NEVER stored or logged
- Password is only used to derive encryption key
- Wrong password won't unlock vault (GCM authentication check)

### Vault File

- Located at `~/.keyp/vault.json` by default
- Contains encrypted secrets (AES-256-GCM)
- Safe to commit to Git or store on cloud (encrypted)
- Only readable if someone has your master password

### Terminal History

- Commands are logged in shell history
- Secret values are masked (shown as ●●●●●●●)
- Command names (`keyp set`, `keyp get`) are visible
- Consider: `history -c` to clear history if needed

---

## Platform Support

### Operating Systems

- ✅ macOS (tested on 10.15+)
- ✅ Linux (most distributions)
- ✅ Windows (via WSL or native Node.js)

### Node.js Versions

- ✅ Node.js 14.0.0 or higher
- ✅ Tested on: 14.x, 16.x, 18.x, 20.x

### Clipboard Support

- ✅ macOS: `pbcopy` / `pbpaste`
- ✅ Linux: `xclip`, `xsel`, or wayland clipboard
- ✅ Windows: Native clipboard
- ⚠ Over SSH: Use `--stdout` flag

---

## Troubleshooting

### Vault Won't Initialize

```bash
# Check if vault already exists
$ ls -la ~/.keyp/vault.json

# If it exists but you want to start over:
$ rm ~/.keyp/vault.json
$ keyp init
```

### Can't Unlock Vault

```bash
# Try again - check password carefully
# Remember: password is case-sensitive

$ keyp list
Enter master password: ●●●●●●●●
✗ Incorrect password (2 attempts remaining)
```

### Clipboard Not Working on Linux

```bash
# Install clipboard tool
$ sudo apt install xclip        # Ubuntu/Debian
$ sudo pacman -S xclip          # Arch
$ brew install xclip            # macOS

# Or use --stdout flag
$ keyp get secret --stdout
```

### Permission Denied (~/.keyp)

```bash
# Fix permissions on keyp directory
$ chmod 700 ~/.keyp
$ chmod 600 ~/.keyp/vault.json
```

---

### keyp sync <subcommand>

Synchronize your vault with Git remote repositories for encrypted backups.

**Subcommands:**

#### keyp sync init <remote-url>

Initialize Git synchronization with a remote repository.

**Usage:**
```bash
keyp sync init https://github.com/username/backup.git
keyp sync init git@github.com:username/backup.git --auto-push
```

**Options:**
- `-a, --auto-push` - Enable automatic push on vault changes
- `-c, --auto-commit` - Enable automatic commit on vault changes

**Example:**
```bash
$ keyp sync init https://github.com/user/keyp-backup.git
ℹ Initializing Git sync...
✓ Git repository initialized
✓ Remote configured: https://github.com/user/keyp-backup.git
✓ Git sync initialized successfully!
```

#### keyp sync push

Push encrypted vault to remote repository.

**Usage:**
```bash
keyp sync push
keyp sync push -m "Updated production API key"
```

**Options:**
- `-m, --message <message>` - Custom commit message

#### keyp sync pull

Pull vault from remote repository with conflict detection.

**Usage:**
```bash
keyp sync pull
keyp sync pull --strategy keep-local --auto-resolve
```

**Options:**
- `-s, --strategy <strategy>` - Conflict resolution (`keep-local` or `keep-remote`)
- `--auto-resolve` - Automatically resolve conflicts

#### keyp sync status

Display current Git sync status.

**Usage:**
```bash
$ keyp sync status
Status: ✓ Synced
Last sync: 2h ago
Uncommitted changes: No
Unpushed commits: 0
Conflicts: 0
```

#### keyp sync config

View or configure Git sync settings.

**Usage:**
```bash
keyp sync config                      # Show current settings
keyp sync config --auto-push true     # Enable auto-push
keyp sync config --auto-commit false  # Disable auto-commit
```

**For detailed Git sync guide:** See [Git Sync Documentation](./GIT_SYNC.md)

---

### keyp stats

Display vault statistics and encryption information.

**Usage:**
```bash
keyp stats
```

**Example:**
```bash
$ keyp stats
Enter master password: ●●●●●●●●

📊 Vault Statistics
────────────────────────────────────

Secrets
  Total: 42
  Average value length: 68 characters
  Longest name: database_connection_string

Storage
  Vault file size: 12.34 KB
  Location: ~/.keyp/vault.json

Dates
  Last modified: 10/20/2025, 3:45:00 PM
  Last synced: 2h ago

Encryption
  Algorithm: AES-256-GCM
  Key derivation: PBKDF2-SHA256
  Iterations: 100,000+
```

**Information displayed:**
- Total number of secrets
- Average secret value length
- Longest secret name
- Vault file size
- Last modified date
- Last sync date (if Git sync configured)
- Encryption algorithm and parameters

---

### keyp config [action] [key] [value]

Manage keyp configuration settings.

**Usage:**
```bash
keyp config                           # Show all settings
keyp config list                      # Show all settings
keyp config get <key>                 # Get specific setting
keyp config set <key> <value>         # Set a configuration
keyp config reset                     # Reset to defaults
```

**Configuration Keys:**

- `clipboard-timeout` - How long before clipboard is auto-cleared (seconds, default: 45)
- `auto-lock` - Auto-lock vault after inactivity (seconds, or "none" to disable)
- `git-auto-sync` - Automatically push on vault changes (true/false)

**Examples:**
```bash
# View current configuration
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

# Change clipboard timeout
$ keyp config set clipboard-timeout 60
✓ Clipboard timeout set to 60 seconds

# Get specific value
$ keyp config get clipboard-timeout
60

# Enable auto-sync
$ keyp config set git-auto-sync true
✓ Git auto-sync enabled

# Reset to defaults
$ keyp config reset
✓ Configuration reset to defaults
```

**Default Values:**
```json
{
  "clipboardTimeout": 45,
  "autoLock": null,
  "gitAutoSync": false
}
```

---

### keyp destroy

Permanently delete the entire vault and all associated configuration files. **This action cannot be undone.**

⚠️ **WARNING:** This command will:
- Delete the vault file (vault.json)
- Delete all stored secrets
- Delete configuration files
- Delete Git sync configuration (if configured)
- Delete sync time tracking

This is irreversible and cannot be recovered.

**Usage:**
```bash
keyp destroy
```

**Interactive prompts:**
1. Displays severe warning message about permanent deletion
2. Requires explicit confirmation (must type "destroy")
3. Requires master password verification
4. Final confirmation before deletion

**Example:**
```bash
$ keyp destroy

⚠️  DANGER ⚠️
────────────────────────────────────
Permanent Vault Deletion

You are about to PERMANENTLY DELETE your entire vault.
This action CANNOT be undone.

Vault location: ~/.keyp/vault.json
Secrets: 42
Last modified: 10/20/2025, 3:45:00 PM

This will delete:
  • All secrets in the vault
  • Vault configuration files
  • Git sync configuration
  • All sync history

Are you absolutely sure? Type 'destroy' to confirm: destroy
Enter master password to verify: ●●●●●●●●

✓ Vault destroyed successfully
ℹ Deleted: ~/.keyp/vault.json
ℹ Deleted: ~/.keyp/.keyp-config.json
ℹ Deleted: ~/.keyp/.keyp-git-config.json
ℹ Deleted: ~/.keyp/.keyp-sync-time
ℹ The ~/.keyp directory may still exist (you can remove it manually)
```

**Safety Features:**
- Requires typing "destroy" to confirm (prevents accidental deletion)
- Requires master password verification (ensures authorization)
- Clear warnings about irreversible nature
- Displays vault size and secret count before deletion
- Shows deletion summary after completion

**Errors:**
- "Vault not found" - if no vault is initialized
- "Confirmation cancelled" - if user doesn't type "destroy"
- "Invalid master password" - if password verification fails

**Use Cases:**
- Complete vault reset (initialize fresh vault with new password)
- Secure removal before uninstalling keyp
- Migration to different secret management tool
- Compromised vault that needs complete removal

---

## Shell Completion

Enable tab completion for faster command entry.

### Bash

```bash
# Add to ~/.bashrc or ~/.bash_profile
source /path/to/keyp/completions/keyp.bash

# Or manually enable
complete -o bashdefault -o default -o nospace -F _keyp_completion keyp
```

### Zsh

```bash
# Add to ~/.zshrc
fpath=(/path/to/keyp/completions $fpath)
autoload -U compinit && compinit

# Or copy completion file to zsh directory
cp /path/to/keyp/completions/keyp.zsh ~/.zsh/completions/_keyp
```

### Features

- Command name completion
- Secret name completion (for get, delete, rename, copy)
- Flag and option completion
- File path completion (for export/import)

**Example:**
```bash
keyp get git<TAB>     # Completes to: keyp get github-token
keyp list --se<TAB>   # Completes to: keyp list --search
keyp sync pull -s<TAB> # Completes to: keyp sync pull --strategy
```

---

## Getting Help

### Available Commands

```bash
$ keyp --help          # Show all commands
$ keyp init --help     # Help for specific command
```

### Documentation

- 📖 [Full API Reference](./API.md)
- 🔐 [Security Guide](./SECURITY.md)
- 📋 [Vault Format](./VAULT_FORMAT.md)
- 🌐 [Git Sync Guide](./GIT_SYNC.md)

### Report Issues

- 🐛 [GitHub Issues](https://github.com/TheEditor/keyp/issues)
- 💬 [GitHub Discussions](https://github.com/TheEditor/keyp/discussions)

---

## License

MIT © Dave Fobare

**keyp** is a local-first secret manager for developers, built with ❤️ using Node.js.
