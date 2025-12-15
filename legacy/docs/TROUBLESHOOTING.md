# Troubleshooting Guide

Solutions for common issues and problems with keyp.

## Installation Issues

### "Command not found: keyp"

**Problem:** After installing keyp, the command is not found in your terminal.

**Solutions:**

1. **Verify npm global bin is in PATH:**
   ```bash
   npm config get prefix
   ```
   This should return a path. Verify this path is in your system PATH.

2. **Reinstall globally:**
   ```bash
   npm uninstall -g @theeditor/keyp
   npm install -g @theeditor/keyp
   ```

3. **On Windows:** Close and reopen your terminal or PowerShell after installation.

4. **On macOS/Linux:** Reload your shell:
   ```bash
   exec $SHELL
   ```

### "EACCES: Permission denied"

**Problem:** Permission error when installing globally on macOS or Linux.

**Solutions:**

1. **Fix npm permissions** (recommended):
   ```bash
   mkdir ~/.npm-global
   npm config set prefix '~/.npm-global'
   export PATH=~/.npm-global/bin:$PATH
   ```
   Then add this to your `.bashrc` or `.zshrc`:
   ```bash
   export PATH=~/.npm-global/bin:$PATH
   ```

2. **Use sudo** (not recommended):
   ```bash
   sudo npm install -g @theeditor/keyp
   ```

## Vault Issues

### "No vault found"

**Problem:** You get an error saying the vault doesn't exist.

**Solution:** Initialize a new vault:
```bash
keyp init
```

You'll be prompted to create a master password. This creates an encrypted vault file at `~/.keyp/vault.json`.

### "Vault already exists"

**Problem:** Error when trying to initialize but vault already exists.

**Solution:** If you want to start fresh:
```bash
# Back up your current vault first
cp ~/.keyp/vault.json ~/.keyp/vault.json.backup

# Remove the old vault
rm ~/.keyp/vault.json

# Initialize new vault
keyp init
```

### "Incorrect password" after multiple retries

**Problem:** You're locked out with "Maximum password attempts exceeded".

**Solutions:**

1. **Remember your password:**
   - The master password is the only way to unlock your vault
   - There's no password recovery (by design)

2. **Start fresh:**
   ```bash
   # Backup your encrypted vault (in case password comes back to you)
   cp ~/.keyp/vault.json ~/.keyp/vault.json.locked

   # Create new vault with password you remember
   rm ~/.keyp/vault.json
   keyp init
   ```

## Password Issues

### Weak password warning

**Problem:** You see a warning about weak password strength.

**Recommendation:** While keyp allows weak passwords, we recommend:
- At least 12 characters
- Mix of uppercase, lowercase, numbers
- Special characters (!@#$%^&*)
- Avoid dictionary words

**Example strong password:**
```
MyP@ssw0rd2024!
```

### Password contains special characters

**Problem:** Password with special characters isn't working.

**Solution:** Make sure you're escaping special characters correctly when entering:
- In terminal: No escaping needed, just type it
- In scripts: May need to escape with backslash or use quotes

```bash
# Good
keyp init
Enter master password: MyP@ss!word  # Just type it

# Script
password="MyP@ss\!word"  # Escape ! in double quotes
```

## Clipboard Issues

### Clipboard not working

**Problem:** Secret doesn't copy to clipboard, or you see "clipboard not available".

**Solutions:**

1. **Linux - install clipboard tools:**
   ```bash
   # Ubuntu/Debian
   sudo apt-get install xclip

   # Fedora/RHEL
   sudo dnf install xclip

   # Arch
   sudo pacman -S xclip

   # Alpine
   apk add xclip
   ```

2. **Use stdout instead:**
   ```bash
   keyp get my-secret --stdout
   ```

3. **Check X11 or Wayland:**
   ```bash
   echo $DISPLAY      # X11
   echo $WAYLAND_DISPLAY  # Wayland
   ```

4. **macOS - ensure Terminal has clipboard access:**
   - System Preferences → Security & Privacy → Accessibility
   - Add Terminal to allowed apps

### Clipboard clears too quickly

**Problem:** Secret clears from clipboard before you can paste.

**Solution:** Adjust the timeout:
```bash
# Change to 120 seconds
keyp get my-secret --timeout 120

# Disable auto-clear
keyp get my-secret --no-clear
```

Or set a global default:
```bash
keyp config set clipboard-timeout 120
```

## Git Sync Issues

### "Git repository not initialized"

**Problem:** Error when trying to sync before initializing git.

**Solution:** Initialize Git sync first:
```bash
keyp sync init https://github.com/username/keyp-backup.git
```

### "Permission denied" during Git sync

**Problem:** Authentication error when pushing/pulling from remote.

**Solutions:**

1. **SSH key setup:**
   ```bash
   # Generate SSH key (if you don't have one)
   ssh-keygen -t ed25519 -C "your-email@example.com"

   # Add public key to GitHub/GitLab
   cat ~/.ssh/id_ed25519.pub  # Copy this to your Git provider

   # Test SSH connection
   ssh -T git@github.com
   ```

2. **Use HTTPS with Personal Access Token:**
   ```bash
   # GitHub: https://github.com/settings/tokens
   # GitLab: https://gitlab.com/-/profile/personal_access_tokens

   # Use SSH URL or HTTPS with token
   keyp sync init https://github.com/username/keyp-backup.git
   # Git will prompt for username and token
   ```

3. **Store Git credentials:**
   ```bash
   # Linux/macOS
   git config --global credential.helper store

   # macOS (recommended)
   git config --global credential.helper osxkeychain

   # Windows
   git config --global credential.helper wincred
   ```

### "Failed to push: remote not found"

**Problem:** Error when trying to push to a non-existent repository.

**Solution:** Ensure your remote repository exists and is accessible:
```bash
# On GitHub/GitLab: create a new repository first
# Then initialize sync with the correct URL

keyp sync init https://github.com/username/existing-repo.git
keyp sync push
```

### Merge conflicts

**Problem:** Conflicts when pulling from remote with local changes.

**Solutions:**

1. **Keep your local changes:**
   ```bash
   keyp sync pull --strategy keep-local --auto-resolve
   ```

2. **Accept remote changes:**
   ```bash
   keyp sync pull --strategy keep-remote --auto-resolve
   ```

3. **Manual resolution:**
   ```bash
   # See what conflicts exist
   keyp sync pull

   # Manually merge if needed, then push
   keyp sync push
   ```

## Performance Issues

### "keyp" is slow to start

**Problem:** Command takes longer than expected to start.

**Causes:**
- First run (Node.js initialization)
- Large vault file
- Slow disk

**Solutions:**
1. Subsequent runs should be faster (Node.js cache)
2. Consider splitting large vaults into multiple files
3. Ensure you have sufficient disk space

### Password unlock is slow

**Problem:** Unlocking vault takes a long time.

**Note:** This is intentional - password derivation uses 100,000+ PBKDF2 iterations for security. This should take 0.1-0.5 seconds on modern hardware.

If it takes significantly longer:
- Check CPU usage
- Try on a faster machine
- Ensure no other processes are consuming resources

## File Permission Issues

### Permission denied (~/.keyp)

**Problem:** Error accessing files in ~/.keyp directory.

**Solution:** Fix directory permissions:
```bash
chmod 700 ~/.keyp
chmod 600 ~/.keyp/vault.json
```

### Vault file inaccessible

**Problem:** Error reading or writing vault file.

**Solution:** Check file permissions and ownership:
```bash
# View permissions
ls -la ~/.keyp/vault.json

# Fix permissions (if needed)
chmod 600 ~/.keyp/vault.json

# Fix ownership (if needed)
chown $(whoami) ~/.keyp/vault.json
```

## Platform-Specific Issues

### macOS

**Clipboard not working:**
- Grant Terminal clipboard access in System Preferences
- System Preferences → Security & Privacy → Accessibility

**Vault location:**
- Vault stored at: `~/.keyp/vault.json`
- Config stored at: `~/.keyp/.keyp-config.json`

### Windows

**Command not recognized:**
- Close and reopen PowerShell or Command Prompt
- Verify npm is in PATH: `npm --version`

**Backslash in paths:**
- Use forward slashes in commands
- Or escape backslashes: `C:\\Users\\username\\.keyp\\vault.json`

**Clipboard:**
- Should work automatically with wl-copy or native Windows clipboard
- If not working, use `--stdout` flag

### Linux

**Clipboard tools not installed:**
- Install xclip: `sudo apt-get install xclip`
- Or use Wayland alternatives

**Permission denied:**
- Use `sudo` carefully - vault may end up owned by root
- Better to fix npm permissions as described above

## Secret-Related Issues

### Can't find a secret

**Problem:** Secret exists but you can't retrieve it.

**Solution:** List all secrets to verify:
```bash
keyp list
```

If your secret appears:
- Check exact spelling (names are case-sensitive)
- Try searching: `keyp list --search pattern`

If not listed:
- Secret may be in a different vault
- Try: `keyp list --count` to see total

### Secret contains garbage/corrupted data

**Problem:** Secret value is corrupted or contains unexpected data.

**Cause:** Vault file may be corrupted.

**Solution:**
1. Restore from backup if available
2. Check vault file integrity
3. Start fresh if necessary

## Error Messages

### "Vault file corrupted"

**Cause:** Vault encryption failed or file was modified externally.

**Solution:**
1. Restore from backup: `cp ~/.keyp/vault.json.backup ~/.keyp/vault.json`
2. Or start fresh: `rm ~/.keyp/vault.json && keyp init`

### "Decryption failed: auth tag mismatch"

**Cause:** Password is wrong or file is corrupted.

**Solution:**
1. Try the correct password
2. If forgotten, restore from backup
3. Or start with new vault

### "Operation not permitted"

**Cause:** Insufficient permissions or file locked by another process.

**Solution:**
1. Check file permissions: `ls -la ~/.keyp/vault.json`
2. Ensure no other keyp instances are running
3. Try closing and reopening terminal

## Getting More Help

If you can't find a solution:

1. **Check related documentation:**
   - [CLI Reference](./CLI.md)
   - [Git Sync Guide](./GIT_SYNC.md)
   - [Security Guide](./SECURITY.md)

2. **Search existing issues:**
   - [GitHub Issues](https://github.com/TheEditor/keyp/issues)

3. **Create a new issue:**
   - [Report a bug](https://github.com/TheEditor/keyp/issues/new?template=bug_report.md)
   - [Request help](https://github.com/TheEditor/keyp/discussions)

4. **Include:**
   - keyp version: `keyp --version`
   - Node version: `node --version`
   - OS and version: `uname -a` or `systeminfo`
   - Steps to reproduce
   - Error messages

---

**Remember:** keyp is designed to be secure first. Some limitations (like no password recovery) are intentional security features, not bugs.
