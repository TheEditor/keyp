# Security Policy

## Reporting Security Vulnerabilities

If you discover a security vulnerability in keyp, please report it responsibly. **Do not** create a public GitHub issue for security vulnerabilities.

### Reporting Process

1. **Email:** Send a detailed report to security@example.com (or maintainer's email)
   - Include: vulnerability description, affected versions, and proof of concept
   - Do not include actual exploit code that could be used maliciously
   - Allow 24-48 hours for initial response

2. **Alternative:** GitHub Security Advisory
   - Go to Settings → Security → Report a vulnerability
   - Private advisories let us coordinate disclosure

### What Happens Next

- We acknowledge receipt within 24 hours
- We'll investigate and determine severity
- We'll develop and test a fix
- We'll prepare a security release
- We'll credit the researcher (unless you prefer anonymity)
- We'll publish a CVE (if applicable)

## Supported Versions

| Version | Supported | Status |
|---------|-----------|---------|
| 0.2.x   | ✅ Yes    | Current |
| 0.1.x   | ✅ Yes    | Maintenance |
| 0.0.x   | ❌ No     | Deprecated |

Security fixes are provided for current and previous minor versions.

## Security Best Practices

### Using keyp Securely

1. **Strong Master Password**
   - Use 12+ characters
   - Mix uppercase, lowercase, numbers, special characters
   - Avoid dictionary words
   - Never use the same password elsewhere

2. **Protect Your Vault File**
   ```bash
   # Vault file permissions (should be 600)
   chmod 600 ~/.keyp/vault.json

   # Directory permissions (should be 700)
   chmod 700 ~/.keyp
   ```

3. **Git Sync Security**
   - Use private repositories for backups
   - Use SSH keys for authentication (recommended over HTTPS)
   - Enable 2FA on GitHub/GitLab account
   - Rotate SSH keys periodically

4. **System Security**
   - Ensure your computer has updated OS and security patches
   - Use antivirus/malware protection
   - Don't install keyp on untrusted machines
   - Beware of keyloggers in shared environments

5. **Terminal Practices**
   - Don't share terminal sessions with untrusted users
   - Clear command history after using keyp with sensitive operations
   - Use `--no-clear` carefully (clipboard may be visible)

### For CI/CD Pipelines

1. **Store Vault Securely**
   - Use your CI/CD provider's secret management
   - Never commit vault.json to version control
   - Store vault password as a secret (not in config)

2. **Limit Secret Access**
   - Only run sensitive commands in protected environments
   - Restrict CI/CD job access
   - Audit who can access secrets

3. **Network Security**
   - Use HTTPS/SSH for all Git operations
   - Ensure CI/CD runner has secure network access
   - Consider VPN for internal services

## Known Limitations

### Design Choices (Not Bugs)

1. **No Password Recovery**
   - By design - ensures only you can access your data
   - No backdoor means no recovery even for us
   - Write down your password somewhere secure

2. **No Shared Passwords**
   - Each vault has unique encryption key derived from master password
   - Not designed for team shared credentials
   - Use separate machines or vaults for different users

3. **Master Password Required**
   - Every operation requires master password
   - No persistent unlock across processes
   - This is a security feature, not a limitation

4. **Local Storage Only**
   - Vault stored unencrypted in memory while unlocked
   - Could theoretically be accessed via memory dump
   - Mitigation: lock vault when not in use

## Cryptographic Details

For complete cryptographic analysis, see [docs/SECURITY.md](./docs/SECURITY.md)

**Summary:**
- **Encryption:** AES-256-GCM (authenticated encryption)
- **Key Derivation:** PBKDF2-SHA256 with 100,000+ iterations
- **Authentication:** GCM auth tags prevent tampering
- **Randomness:** Cryptographically secure random number generation

## Version Disclosure

We practice responsible disclosure:
- Security fixes are released promptly
- Fixes are applied to all supported versions
- Security releases are clearly labeled
- No details of vulnerabilities before patches are available

## Security Audits

This is a small, open-source project. While we follow security best practices:
- Code is publicly available for review
- Security analysis welcome (see CONTRIBUTING.md)
- Not professionally audited (contributions welcome!)

## Compliance

keyp is designed for:
- ✅ Personal secret management
- ✅ Developer credential storage
- ✅ Local development environments

keyp is NOT designed for:
- ❌ HIPAA compliance
- ❌ SOC 2 requirements
- ❌ Enterprise multi-user scenarios
- ❌ Regulatory compliance (GDPR, etc)

If you need enterprise secret management, consider:
- HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault
- 1Password
- Bitwarden

## Questions?

- **Technical questions:** Create an issue on GitHub
- **Security concerns:** Email security privately (don't create issues)
- **Best practices:** See docs/SECURITY.md and docs/GIT_SYNC.md

---

**Last Updated:** October 2025

Thank you for helping keep keyp secure! 🔒
