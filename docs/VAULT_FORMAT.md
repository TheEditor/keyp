# keyp Vault File Format Specification

## Version

**Format Version:** 1.0.0

This specification describes the structure and encoding of keyp vault files.

## File Location

Default location: `~/.keyp/vault.json`

Can be customized when creating a VaultManager instance.

## Complete Vault File Structure

```json
{
  "version": "1.0.0",
  "crypto": {
    "algorithm": "aes-256-gcm",
    "kdf": "pbkdf2",
    "iterations": 100000,
    "salt": "base64-encoded-salt"
  },
  "data": "base64-encoded-encrypted-secrets",
  "authTag": "base64-encoded-authentication-tag",
  "iv": "base64-encoded-initialization-vector",
  "createdAt": "2025-10-20T23:45:00.000Z",
  "updatedAt": "2025-10-20T23:50:30.000Z"
}
```

## Field Descriptions

### `version` (string, required)

Semantic version of the vault format.

- Current value: `"1.0.0"`
- Used for future migration compatibility
- Must be present for all vaults

**Example:**
```json
"version": "1.0.0"
```

### `crypto` (object, required)

Container for cryptographic configuration and parameters.

#### `crypto.algorithm` (string, required)

The encryption algorithm used for secret data.

- Current value: `"aes-256-gcm"` (only valid value for v1.0.0)
- 256-bit AES in Galois/Counter Mode
- Provides both encryption and authentication

**Example:**
```json
"algorithm": "aes-256-gcm"
```

#### `crypto.kdf` (string, required)

The key derivation function used to convert password to encryption key.

- Current value: `"pbkdf2"` (only valid value for v1.0.0)
- PBKDF2 with SHA-256 hash function
- Resistant to brute-force attacks

**Example:**
```json
"kdf": "pbkdf2"
```

#### `crypto.iterations` (number, required)

Number of PBKDF2 iterations for key derivation.

- Minimum value: `100000` (enforced by implementation)
- Higher values = slower but more secure
- Same value must be used for both encryption and decryption

**Example:**
```json
"iterations": 100000
```

#### `crypto.salt` (string, required)

Base64-encoded salt value used by PBKDF2 for key derivation.

- Encoding: Base64 (RFC 4648)
- Byte length: 32 bytes when decoded (256 bits)
- Unique random value generated during vault initialization
- Stored in plaintext (not secret)
- Prevents rainbow table attacks

**Example:**
```json
"salt": "X7k9mN2+/Zq8vL4pQ6rS3tU1wX9yZ0aB1cD2eF3gH4i5jK6"
```

### `data` (string, required)

Base64-encoded encrypted vault secrets.

- Encoding: Base64 (RFC 4648)
- Encryption: AES-256-GCM
- Original plaintext: JSON object of secrets
- Byte length: varies with number and size of secrets

**Encryption Details:**
1. Vault data is serialized as JSON string
2. Encrypted using AES-256-GCM with:
   - Key: Derived from password using PBKDF2
   - IV: Random 12-byte value
   - Auth tag: 16-byte tag appended by GCM
3. Resulting ciphertext is encoded in Base64

**Example:**
```json
"data": "AbC123dEf456gHi789jKl012mNo345pQr678sTu901vWx234yZ567aBc890dEf123gHi456"
```

### `authTag` (string, required)

Base64-encoded GCM authentication tag for vault data integrity.

- Encoding: Base64 (RFC 4648)
- Byte length: 16 bytes when decoded (128 bits)
- Protects against tampering and corruption
- Decryption fails if vault data or tag is modified
- Generated automatically by GCM mode during encryption

**Example:**
```json
"authTag": "X1yZ2aB3cD4eF5gH6iJ7"
```

### `iv` (string, required)

Base64-encoded initialization vector for AES-256-GCM.

- Encoding: Base64 (RFC 4648)
- Byte length: 12 bytes when decoded (96 bits)
- Random value, unique for each encryption operation
- Stored in plaintext (not secret)
- Prevents patterns when encrypting same data multiple times

**Example:**
```json
"iv": "AbCdEfGhIjKlMnOp"
```

### `createdAt` (string, required)

ISO 8601 timestamp of vault creation.

- Format: RFC 3339 / ISO 8601 format with 'Z' timezone
- Example: `"2025-10-20T23:45:00.000Z"`
- Set once during vault initialization
- Never changes during vault lifetime
- Helps identify vault age and lifecycle

**Example:**
```json
"createdAt": "2025-10-20T23:45:00.000Z"
```

### `updatedAt` (string, required)

ISO 8601 timestamp of last vault modification.

- Format: RFC 3339 / ISO 8601 format with 'Z' timezone
- Updated every time vault secrets are modified
- Indicates last time vault was saved

**Example:**
```json
"updatedAt": "2025-10-20T23:50:30.000Z"
```

## Plaintext Data Structure

When decrypted, `data` field contains the following JSON structure:

```json
{
  "secret-name-1": "secret-value-1",
  "secret-name-2": "secret-value-2",
  "github-token": "ghp_xxx...",
  "database-password": "secret123"
}
```

### Structure Details

- **Format:** JSON object (key-value pairs)
- **Keys:** Secret names (strings)
- **Values:** Secret values (strings)
- **Examples:**
  - API keys: `"github-token": "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
  - Passwords: `"db-password": "super_secret_123!@#"`
  - Tokens: `"jwt-secret": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

### Special Characters

- Keys and values can contain any UTF-8 characters
- No escaping required (JSON handles escaping)
- Supports emojis and international characters

**Example with special characters:**
```json
{
  "api-key-🔒": "secret",
  "db-password": "密码123",
  "token": "token\\with\\backslashes"
}
```

## File Format Details

### JSON Structure Rules

- Valid JSON 5 (formatted with 2-space indentation)
- All required fields must be present
- No additional fields are allowed (for forwards compatibility)
- Formatted with newlines for readability

### Base64 Encoding

All binary data (salt, iv, data, authTag) is Base64-encoded per RFC 4648:

- Standard Base64 alphabet: A-Z, a-z, 0-9, +, /
- Padding with = characters as needed
- No line breaks in encoded output

**Example encoding:**
```
Binary:  0xFF 0x3E 0x7C 0xA1 0x2B 0x5D
Base64:  /z58oStd
```

### File Encoding

- Character encoding: UTF-8
- Line endings: System-dependent (\n on Unix, \r\n on Windows)
- Newline after closing brace: Yes

## Encryption Process

```
User Master Password
         ↓
    [PBKDF2 KDF]
    (100,000 iterations, SHA-256, random salt)
         ↓
    256-bit Key
         ↓
    [AES-256-GCM Encrypt]
    (random 12-byte IV)
         ↓
    Ciphertext (16-byte auth tag appended)
         ↓
    [Base64 Encode]
         ↓
    Vault File Data Field
```

## Decryption Process

```
Vault File
    ↓
Read: salt, iterations, iv, data, authTag
    ↓
User Master Password + salt
    ↓
[PBKDF2 KDF] (same iterations, SHA-256)
    ↓
256-bit Key
    ↓
[AES-256-GCM Decrypt]
(verify authTag, check for tampering)
    ↓
Original JSON Object
    ↓
Parse Secrets
```

## Security Properties

### Data at Rest
- ✅ Encrypted with AES-256-GCM
- ✅ Password protected with PBKDF2
- ✅ Unique random salt per vault
- ✅ Integrity protected with GCM auth tag

### Metadata Security
- ℹ️ Timestamps and algorithm info are **not** encrypted
- ℹ️ This allows checking vault age without decryption
- ℹ️ Not considered sensitive (no credentials revealed)

### Vault File Permissions
- Linux/macOS: `0700` (rwx------)
- Windows: Inherited from parent directory (recommend NTFS encryption)

## Example Vault File

```json
{
  "version": "1.0.0",
  "crypto": {
    "algorithm": "aes-256-gcm",
    "kdf": "pbkdf2",
    "iterations": 100000,
    "salt": "42Tq0NL/1X0EqM0YCd+8rw7S8bH+FJyKJ/1pL9qWvN0="
  },
  "data": "KlGf7H2xP8NqQvW3YhF5zJ+M8X0kL9pR5sK2dN+vYq/TrS4wT9u0X7c1Y8fG2cJ3dQ+sV5wO9x2a8gH+j3kQ6b",
  "authTag": "vL8j0oK3mN6pQ9rS2tV5",
  "iv": "AbCdEfGhIjKlMnOp",
  "createdAt": "2025-10-20T23:45:00.000Z",
  "updatedAt": "2025-10-20T23:50:30.000Z"
}
```

## Migration Strategy

### Version Compatibility

**v1.0.0 → v2.0.0 Example:**
- New version would introduce new `algorithm` or `kdf` values
- Implementation would support both v1.0.0 and v2.0.0
- User could opt into migration with password

**Migration Process:**
1. Read vault with current version handler
2. Decrypt using old settings
3. Re-encrypt using new settings
4. Update version field
5. Save updated vault file

## Implementation Notes

### When Creating a Vault

1. Generate random 32-byte salt
2. Derive key using PBKDF2 with default 100,000 iterations
3. Generate random 12-byte IV
4. Create empty data object: `{}`
5. Encrypt empty data with AES-256-GCM
6. Generate vault file with current timestamps
7. Write to disk

### When Saving Vault

1. Serialize vault data to JSON string
2. Generate new random 12-byte IV
3. Encrypt JSON string with AES-256-GCM
4. Update `updatedAt` timestamp
5. Write to disk (overwrite previous file)

### When Loading Vault

1. Read vault file as JSON
2. Validate all required fields present
3. Check version compatibility
4. Reconstruct encryption result from file
5. Decrypt with user-provided password
6. Parse decrypted JSON as vault data
7. Load into memory (locked by default)

## Constants

```typescript
const VAULT_VERSION = "1.0.0";
const CIPHER_ALGORITHM = "aes-256-gcm";
const KDF_ALGORITHM = "pbkdf2";
const KDF_HASH = "sha256";
const KDF_ITERATIONS_MIN = 100000;
const KDF_ITERATIONS_DEFAULT = 100000;
const SALT_LENGTH = 32; // bytes
const KEY_LENGTH = 32; // bytes (256-bits)
const IV_LENGTH = 12; // bytes (96-bits for GCM)
const AUTH_TAG_LENGTH = 16; // bytes (128-bits)
```

## Related Documentation

- See [SECURITY.md](./SECURITY.md) for cryptographic analysis
- See API documentation for programmatic vault operations
