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

### Report Issues

- 🐛 [GitHub Issues](https://github.com/TheEditor/keyp/issues)
- 💬 [GitHub Discussions](https://github.com/TheEditor/keyp/discussions)

---

## License

MIT © Dave Fobare

**keyp** is a local-first secret manager for developers, built with ❤️ using Node.js.
