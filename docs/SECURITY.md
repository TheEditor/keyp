# keyp Security Documentation

## Overview

**keyp** implements industry-standard cryptographic practices to protect your secrets. This document explains the security model, threat model, and cryptographic choices.

## Cryptographic Algorithms

### Encryption: AES-256-GCM

- **Algorithm:** Advanced Encryption Standard with 256-bit key size in Galois/Counter Mode
- **Security Level:** 256-bit equivalent strength
- **Authentication:** GCM provides built-in authentication tag to detect tampering
- **Implementation:** Node.js built-in `crypto` module (OpenSSL backend)

**Why AES-256-GCM?**
- NIST-approved algorithm for TOP SECRET information
- Authenticated encryption (detects tampering and corruption)
- Efficient and well-analyzed
- No external cryptographic library dependencies

### Key Derivation: PBKDF2

- **Algorithm:** PBKDF2 with SHA-256 hash function
- **Default Iterations:** 100,000 (configurable per vault, minimum 100,000)
- **Key Size:** 256 bits (32 bytes)
- **Salt:** 256 bits (32 bytes) of cryptographically secure random data per vault

**Why PBKDF2?**
- NIST-recommended KDF
- Resistant to brute-force attacks through salting and iteration counting
- Well-studied and widely deployed
- Deterministic (same password + salt = same key)

## Vault File Format

Vault files are stored as JSON with the following structure:

```json
{
  "version": "1.0.0",
  "crypto": {
    "algorithm": "aes-256-gcm",
    "kdf": "pbkdf2",
    "iterations": 100000,
    "salt": "<base64-encoded-32-byte-salt>"
  },
  "data": "<base64-encoded-encrypted-secrets>",
  "authTag": "<base64-encoded-16-byte-gcm-tag>",
  "iv": "<base64-encoded-12-byte-initialization-vector>",
  "createdAt": "<ISO-8601-timestamp>",
  "updatedAt": "<ISO-8601-timestamp>"
}
```

### Plaintext Data Structure

When decrypted, vault data is a JSON object mapping secret names to values:

```json
{
  "github-token": "ghp_xxx...",
  "database-password": "secret123",
  "api-key": "sk_live_..."
}
```

## Security Properties

### What keyp Protects

✅ **Confidentiality:** Secrets are encrypted with AES-256 at rest
✅ **Integrity:** GCM authentication tag detects tampering
✅ **Authentication:** Password-based key derivation
✅ **Salt Randomness:** Unique random salt per vault prevents pre-computed attacks
✅ **IV Randomness:** Unique random IV per encryption operation prevents patterns

### What keyp Does NOT Protect

❌ **Secret Metadata:** Vault timestamps and creation info are not encrypted
❌ **Secret Names:** Secret keys are stored plaintext (feature, not bug - enables searching)
❌ **In-Memory Secrets:** When vault is unlocked, secrets are in memory (standard limitation)
❌ **Physical Security:** Vault file permissions depend on OS file system permissions
❌ **Quantum Attacks:** Like all current cryptography, vulnerable to future quantum computers

## Threat Model

### Threats Mitigated

**Offline Attack:** If vault file is stolen, brute-force attacks are computationally infeasible due to:
- Strong password derivation (PBKDF2 with 100,000+ iterations)
- 256-bit key size
- Unique salt per vault

**Network Eavesdropping:** keyp doesn't use network. Secrets stay on your machine.

**Tampering:** GCM authentication tag will cause decryption to fail if:
- Vault file is corrupted
- Vault file is modified
- Authentication tag is altered

**Accidental Exposure:** Encrypted vault file is useless without the correct password.

### Threats NOT Mitigated

**Weak Passwords:** keyp cannot protect against weak master passwords. If your password can be guessed in a reasonable time, your vault can be opened. Use a strong, random password (20+ characters recommended).

**Compromised Machine:** If your machine is compromised:
- Keyloggers can capture your password
- Malware can read vault file and memory
- This is inherent to any local-first tool

**Malicious Vault:** keyp doesn't validate vault origin. Only download vault files from trusted sources.

**Forced Disclosure:** Law enforcement with a warrant can force you to disclose your password.

## Best Practices

### Protecting Your Vault

1. **Strong Master Password**
   - Use at least 20 random characters
   - Mix uppercase, lowercase, numbers, and symbols
   - Don't reuse passwords from other services
   - Consider using a passphrase of random words

2. **Secure Storage**
   - Store vault file on encrypted disk (FileVault, BitLocker, LUKS)
   - Use standard OS file permissions (~/.keyp has mode 0700)
   - Consider backup encryption

3. **Backup Strategy**
   - Encrypt backups (vault file is already encrypted, but think about backup storage)
   - Store backups securely away from main machine
   - Test that backups can be restored

4. **Machine Security**
   - Keep operating system updated with security patches
   - Use antivirus/malware protection
   - Be cautious with browser extensions and plugins
   - Don't use public WiFi for sensitive operations

### Using keyp Securely

1. **Lock After Use**
   - Explicitly lock vault after operations (clears in-memory data)
   - Don't leave vault unlocked for extended periods

2. **Audit Secret Usage**
   - Know what secrets you have stored
   - Regularly review and remove unused secrets
   - Search for secrets by category

3. **Credential Rotation**
   - Periodically rotate important credentials
   - Remove compromised secrets immediately

4. **Master Password**
   - Never share your master password
   - Don't store master password in files
   - Consider using only in memory (memorize it)

## Cryptographic Analysis

### Key Length Analysis

- **256-bit key for AES-256:** Resistant to all known attacks. Even with quantum computers, 256-bit keys provide ~128-bit security.
- **32-byte (256-bit) salt:** Large enough that collision probability is negligible (2^-128 for random selection of 2^64 salts)
- **12-byte (96-bit) IV:** Standard for GCM mode, not reused within same key

### PBKDF2 Parameter Analysis

- **100,000 iterations (default):**
  - With 1 GHz processor: ~10ms per derivation
  - Reasonable balance between security and usability
  - Can be increased for more security (trades off speed)

- **Future Guidance:**
  - Consider migrating to Argon2 in future versions
  - PBKDF2 has known limitations (hardware optimization attacks)
  - Argon2 is resistant to GPU/ASIC attacks

### Authenticated Encryption Analysis

- **GCM Mode:** NIST and industry standard choice
- **16-byte authentication tag:** 2^-128 collision probability (industry standard)
- **Prevents padding oracle attacks** (unlike CBC mode)

## Implementation Security

### Dependencies

✅ **Zero cryptographic library dependencies**
- Uses only Node.js built-in `crypto` module
- Reduces attack surface
- No dependency vulnerabilities
- Auditable cryptography code

### Code Security

- ✅ TypeScript with strict type checking
- ✅ All functions have documented security properties
- ✅ Comprehensive test coverage
- ✅ Error messages don't leak sensitive information

## Compliance

### Standards Alignment

- **NIST SP 800-132:** PBKDF2 recommendations
- **NIST SP 800-38D:** GCM mode specification
- **FIPS 197:** AES algorithm

### Not Compliant With

- **FIPS 140-3:** Uses non-validated crypto library (Node.js crypto)
- **HIPAA:** Not designed for healthcare use
- **PCI DSS:** Not designed for payment systems
- **SOC 2:** Not formally audited

## Reporting Security Issues

If you discover a security vulnerability:

1. **DO NOT** post it in public issues
2. Email security details to the maintainer
3. Include steps to reproduce
4. Allow time for a fix before public disclosure

## Future Improvements

Potential future enhancements to security:

- [ ] Argon2 KDF (resistant to GPU attacks)
- [ ] Hardware key support (USB security keys)
- [ ] Biometric unlock
- [ ] Master password complexity requirements
- [ ] Formal security audit
- [ ] Encrypted logging/audit trail

## Conclusion

keyp implements strong encryption suitable for protecting developer credentials and API keys. The threat model protects against offline attacks, tampering, and eavesdropping. However, like all local-first tools, security depends on machine integrity and password strength.

**For more information, see [VAULT_FORMAT.md](./VAULT_FORMAT.md) for technical vault structure details.**
