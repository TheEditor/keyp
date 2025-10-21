# keyp Core API Documentation

Complete API reference for the keyp core library.

## Table of Contents

1. [Crypto Module](#crypto-module)
2. [Vault Manager](#vault-manager)
3. [Secrets Manager](#secrets-manager)
4. [Configuration](#configuration)
5. [Types](#types)

## Crypto Module

Low-level cryptographic functions for encryption/decryption.

### `encrypt(data: string, password: string, keyDerivationIterations?: number): EncryptionResult`

Encrypts plaintext using AES-256-GCM with a password-derived key.

**Parameters:**
- `data` - Plaintext string to encrypt
- `password` - Master password for key derivation
- `keyDerivationIterations` - Optional PBKDF2 iterations (default: 100,000, min: 100,000)

**Returns:**
```typescript
{
  ciphertext: string;      // Base64-encoded ciphertext
  authTag: string;         // Base64-encoded authentication tag
  iv: string;              // Base64-encoded initialization vector
  salt: string;            // Base64-encoded salt
}
```

**Example:**
```typescript
import { encrypt } from '@theeditor/keyp';

const result = encrypt(
  JSON.stringify({ apiKey: "sk-123" }),
  "my-secure-password"
);

console.log(result.ciphertext); // Base64 encrypted data
console.log(result.salt);       // Unique salt for this encryption
```

### `decrypt(encryptionResult: EncryptionResult, password: string, keyDerivationIterations?: number): string`

Decrypts AES-256-GCM encrypted data using a password.

**Parameters:**
- `encryptionResult` - Result object from encrypt()
- `password` - Master password (must match encryption password)
- `keyDerivationIterations` - PBKDF2 iterations (must match encryption)

**Returns:** Decrypted plaintext string

**Throws:** Error if decryption fails (wrong password, tampering, corruption)

**Example:**
```typescript
import { decrypt } from '@theeditor/keyp';

try {
  const plaintext = decrypt(encryptionResult, "my-secure-password");
  const data = JSON.parse(plaintext);
  console.log(data.apiKey); // "sk-123"
} catch (error) {
  console.error("Decryption failed:", error.message);
}
```

### `deriveKey(password: string, salt?: Buffer, iterations?: number): {key: Buffer, salt: Buffer}`

Derives an AES encryption key from a password using PBKDF2.

**Parameters:**
- `password` - Master password
- `salt` - Optional salt (generates random if not provided)
- `iterations` - PBKDF2 iterations (default: 100,000, min: 100,000)

**Returns:**
```typescript
{
  key: Buffer;    // 32-byte (256-bit) encryption key
  salt: Buffer;   // 32-byte salt used for derivation
}
```

**Example:**
```typescript
import { deriveKey } from '@theeditor/keyp';

const { key, salt } = deriveKey("my-password");
console.log(key.length);      // 32 bytes
console.log(salt.length);     // 32 bytes
```

### `verifyPassword(encryptionResult: EncryptionResult, password: string, keyDerivationIterations?: number): boolean`

Verifies if a password can decrypt encryption result without throwing.

**Parameters:**
- `encryptionResult` - Result object from encrypt()
- `password` - Password to test
- `keyDerivationIterations` - PBKDF2 iterations from encryption

**Returns:** true if password is correct, false otherwise

**Example:**
```typescript
import { verifyPassword } from '@theeditor/keyp';

const isCorrect = verifyPassword(encryptionResult, userProvidedPassword);
if (isCorrect) {
  console.log("Password is correct!");
} else {
  console.log("Wrong password");
}
```

## Vault Manager

High-level vault operations.

### `VaultManager` Class

Constructor:
```typescript
new VaultManager(vaultPath?: string, keyDerivationIterations?: number)
```

**Parameters:**
- `vaultPath` - Optional custom path to vault file (default: ~/.keyp/vault.json)
- `keyDerivationIterations` - Optional PBKDF2 iterations (default: 100,000)

### Methods

#### `initializeVault(password: string): {success: boolean, message?: string, error?: string}`

Initializes a new encrypted vault.

**Parameters:**
- `password` - Master password for the vault

**Returns:** Operation result

**Throws:** No exceptions, returns error in result object

**Behavior:**
- Fails if vault already exists at the path
- Creates ~/.keyp directory with mode 0o700
- Initializes vault with empty secrets object
- Automatically unlocks vault in memory

**Example:**
```typescript
import { VaultManager } from '@theeditor/keyp';

const manager = new VaultManager();
const result = manager.initializeVault("my-secure-password");

if (result.success) {
  console.log("Vault created successfully!");
} else {
  console.error(result.error);
}
```

#### `unlockVault(password: string): {success: boolean, message?: string, error?: string}`

Unlocks an existing vault and loads secrets into memory.

**Parameters:**
- `password` - Master password

**Returns:** Operation result

**Behavior:**
- Reads vault file from disk
- Decrypts secrets using password
- Loads data into memory (locked state cleared)
- Fails silently if password is incorrect

**Example:**
```typescript
const manager = new VaultManager();
const result = manager.unlockVault("my-secure-password");

if (result.success) {
  console.log("Vault unlocked!");
  const data = manager.getUnlockedData();
  console.log(data); // { secret1: "value1", ... }
}
```

#### `lockVault(): void`

Clears vault secrets from memory.

**Behavior:**
- Removes all decrypted secrets from memory
- Does NOT affect vault file on disk
- Marks vault as locked
- Safe to call multiple times

**Example:**
```typescript
manager.lockVault();
console.log(manager.isVaultUnlocked()); // false
console.log(manager.getUnlockedData()); // null
```

#### `saveVault(password: string): {success: boolean, message?: string, error?: string}`

Saves in-memory vault data to disk (encrypted).

**Parameters:**
- `password` - Master password (needed to re-encrypt)

**Returns:** Operation result

**Behavior:**
- Fails if vault is not unlocked
- Re-encrypts all secrets
- Updates `updatedAt` timestamp
- Overwrites vault file on disk
- Generates new random IV and salt

**Example:**
```typescript
// Add a secret
const data = manager.getUnlockedData();
if (data) {
  data['new-secret'] = 'secret-value';

  // Save to disk
  const result = manager.saveVault("my-secure-password");
  if (result.success) {
    console.log("Vault saved!");
  }
}
```

#### `getUnlockedData(): VaultData | null`

Gets the in-memory vault data.

**Returns:** Vault data object if unlocked, null if locked

**Example:**
```typescript
const data = manager.getUnlockedData();
if (data) {
  console.log("Secrets:", Object.keys(data));
} else {
  console.log("Vault is locked");
}
```

#### `isVaultUnlocked(): boolean`

Checks if vault is currently unlocked.

**Returns:** true if unlocked and secrets are in memory

**Example:**
```typescript
if (manager.isVaultUnlocked()) {
  console.log("Vault is unlocked");
}
```

#### `vaultFileExists(): boolean`

Checks if vault file exists on disk.

**Returns:** true if vault file exists

**Example:**
```typescript
if (!manager.vaultFileExists()) {
  manager.initializeVault("password");
}
```

#### `getVaultPath(): string`

Gets the path to the vault file.

**Returns:** Full path to vault file

**Example:**
```typescript
console.log(manager.getVaultPath());
// Output: /home/user/.keyp/vault.json
```

## Secrets Manager

Secret CRUD operations.

### `SecretsManager` Class

All methods are static.

#### `setSecret(data: VaultData, key: string, value: string): {success: boolean, message?: string, error?: string}`

Adds or updates a secret.

**Parameters:**
- `data` - Vault data object
- `key` - Secret name
- `value` - Secret value

**Returns:** Operation result

**Validation:**
- Fails if key is empty or whitespace
- Fails if value is empty or whitespace

**Example:**
```typescript
import { SecretsManager } from '@theeditor/keyp';

const result = SecretsManager.setSecret(data, "github-token", "ghp_xxx");
if (result.success) {
  console.log(result.message); // "Secret 'github-token' created"
}
```

#### `getSecret(data: VaultData, key: string): string | null`

Retrieves a secret value.

**Parameters:**
- `data` - Vault data object
- `key` - Secret name

**Returns:** Secret value or null if not found

**Example:**
```typescript
const token = SecretsManager.getSecret(data, "github-token");
if (token) {
  console.log("Found token:", token);
}
```

#### `hasSecret(data: VaultData, key: string): boolean`

Checks if a secret exists.

**Parameters:**
- `data` - Vault data object
- `key` - Secret name

**Returns:** true if secret exists

**Example:**
```typescript
if (SecretsManager.hasSecret(data, "api-key")) {
  console.log("API key exists");
}
```

#### `deleteSecret(data: VaultData, key: string): {success: boolean, message?: string, error?: string}`

Deletes a secret.

**Parameters:**
- `data` - Vault data object
- `key` - Secret name

**Returns:** Operation result

**Example:**
```typescript
const result = SecretsManager.deleteSecret(data, "old-token");
if (result.success) {
  console.log("Secret deleted");
}
```

#### `listSecrets(data: VaultData): string[]`

Lists all secret names.

**Parameters:**
- `data` - Vault data object

**Returns:** Array of secret names, sorted alphabetically

**Example:**
```typescript
const secrets = SecretsManager.listSecrets(data);
console.log("Available secrets:", secrets);
// Output: [ 'api-key', 'github-token', 'password' ]
```

#### `getSecretCount(data: VaultData): number`

Gets the number of secrets.

**Parameters:**
- `data` - Vault data object

**Returns:** Number of secrets

**Example:**
```typescript
const count = SecretsManager.getSecretCount(data);
console.log(`You have ${count} secrets`);
```

#### `searchSecrets(data: VaultData, pattern: string): string[]`

Searches for secrets by name pattern (case-insensitive).

**Parameters:**
- `data` - Vault data object
- `pattern` - Search pattern (substring match)

**Returns:** Array of matching secret names, sorted alphabetically

**Example:**
```typescript
const results = SecretsManager.searchSecrets(data, "github");
console.log(results);
// Output: [ 'github-api-key', 'github-token' ]
```

#### `clearAllSecrets(data: VaultData, confirmationKey?: string): {success: boolean, message?: string, error?: string}`

Clears all secrets (dangerous operation).

**Parameters:**
- `data` - Vault data object
- `confirmationKey` - Must be "CONFIRM_DELETE_ALL" to succeed

**Returns:** Operation result

**Example:**
```typescript
// Dangerous! Requires explicit confirmation
const result = SecretsManager.clearAllSecrets(data, "CONFIRM_DELETE_ALL");
if (result.success) {
  console.log("All secrets deleted");
}
```

## Configuration

Configuration utilities and defaults.

### Functions

#### `getKeypDir(): string`

Gets the keyp configuration directory.

**Returns:** Path to ~/.keyp (creates directory if needed)

**Example:**
```typescript
import { getKeypDir } from '@theeditor/keyp';

const dir = getKeypDir();
console.log(dir); // /home/user/.keyp
```

#### `getVaultPath(customPath?: string): string`

Gets the vault file path.

**Parameters:**
- `customPath` - Optional custom path

**Returns:** Vault file path

**Example:**
```typescript
import { getVaultPath } from '@theeditor/keyp';

const path = getVaultPath();
console.log(path); // /home/user/.keyp/vault.json
```

#### `vaultExists(vaultPath?: string): boolean`

Checks if vault file exists.

**Parameters:**
- `vaultPath` - Optional custom path

**Returns:** true if vault file exists

#### `ensureKeypDirExists(): void`

Ensures ~/.keyp directory exists with proper permissions.

### Constants

```typescript
DEFAULT_CONFIG = {
  KEY_DERIVATION_ITERATIONS: 100000,
  VAULT_FILE_NAME: 'vault.json',
  VAULT_VERSION: '1.0.0',
}
```

## Types

### `EncryptionResult`

```typescript
interface EncryptionResult {
  ciphertext: string;  // Base64-encoded encrypted data
  authTag: string;     // Base64-encoded GCM tag
  iv: string;          // Base64-encoded IV
  salt: string;        // Base64-encoded salt
}
```

### `VaultFile`

```typescript
interface VaultFile {
  version: string;
  crypto: CryptoConfig;
  data: string;
  authTag: string;
  iv: string;
  createdAt: string;
  updatedAt: string;
}
```

### `CryptoConfig`

```typescript
interface CryptoConfig {
  algorithm: 'aes-256-gcm';
  kdf: 'pbkdf2';
  iterations: number;
  salt: string;
}
```

### `VaultData`

```typescript
type VaultData = {
  [key: string]: string;
}
```

### `VaultOperationResult`

```typescript
interface VaultOperationResult {
  success: boolean;
  message?: string;
  error?: string;
}
```

### `VaultConfig`

```typescript
interface VaultConfig {
  vaultPath?: string;
  keyDerivationIterations?: number;
}
```

## Complete Example

```typescript
import {
  VaultManager,
  SecretsManager,
} from '@theeditor/keyp';

// Create and initialize vault
const manager = new VaultManager();
manager.initializeVault("my-secure-password");

// Get vault data
const data = manager.getUnlockedData();
if (!data) throw new Error("Vault not unlocked");

// Add secrets
SecretsManager.setSecret(data, "github-token", "ghp_xxx");
SecretsManager.setSecret(data, "api-key", "sk_live_xxx");

// Save to disk
manager.saveVault("my-secure-password");

// Later... lock vault
manager.lockVault();

// And unlock it again
manager.unlockVault("my-secure-password");

// List all secrets
const secrets = SecretsManager.listSecrets(data);
console.log(secrets); // [ 'api-key', 'github-token' ]

// Search for secrets
const results = SecretsManager.searchSecrets(data, "github");
console.log(results); // [ 'github-token' ]

// Retrieve a secret
const token = SecretsManager.getSecret(data, "github-token");
console.log(token); // "ghp_xxx"
```

## Error Handling

All operations return result objects with `success`, `message`, and `error` fields:

```typescript
const result = manager.initializeVault("password");
if (!result.success) {
  console.error("Operation failed:", result.error);
  // Handle error
} else {
  console.log(result.message);
  // Continue
}
```

Only `decrypt()` throws exceptions (cryptographic failures).

## Performance Considerations

- **PBKDF2 iterations:** Each decryption takes ~100ms per 100,000 iterations
- **Large vaults:** Encryption/decryption scales linearly with data size
- **GCM performance:** Highly optimized in Node.js crypto module

## Security Notes

- Passwords are not stored or logged
- Decryption failures return generic error messages
- Error messages don't reveal password material
- In-memory secrets are cleared on `lockVault()`
- See [SECURITY.md](./SECURITY.md) for threat model

## Version Compatibility

Current API is stable for v0.0.1. Breaking changes will increment major version.

## License

MIT © Dave Fobare

---

**For more information:**
- [Security Documentation](./SECURITY.md)
- [Vault Format Specification](./VAULT_FORMAT.md)
- [GitHub Repository](https://github.com/TheEditor/keyp)
