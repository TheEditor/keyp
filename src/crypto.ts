/**
 * Cryptographic utilities for keyp vault encryption/decryption
 * Uses AES-256-GCM for authenticated encryption and PBKDF2 for key derivation
 */

import { createCipheriv, createDecipheriv, randomBytes, pbkdf2Sync } from 'crypto';

/**
 * Result of encryption operation
 */
export interface EncryptionResult {
  /** Base64-encoded encrypted ciphertext */
  ciphertext: string;
  /** Base64-encoded authentication tag for GCM mode */
  authTag: string;
  /** Base64-encoded initialization vector */
  iv: string;
  /** Base64-encoded salt used for key derivation */
  salt: string;
}

/**
 * Derives a 256-bit encryption key from a password using PBKDF2
 * Uses SHA-256 as the hash function with random salt
 *
 * @param password - Master password for vault
 * @param salt - Optional salt (if not provided, random salt is generated)
 * @param iterations - Number of PBKDF2 iterations (default: 100,000)
 * @returns Object containing derived key and salt used
 */
export function deriveKey(
  password: string,
  salt?: Buffer,
  iterations: number = 100000
): {
  key: Buffer;
  salt: Buffer;
} {
  if (iterations < 100000) {
    throw new Error('PBKDF2 iterations must be at least 100,000 for security');
  }

  const actualSalt = salt || randomBytes(32);
  const key = pbkdf2Sync(password, actualSalt, iterations, 32, 'sha256');

  return { key, salt: actualSalt };
}

/**
 * Encrypts data using AES-256-GCM with PBKDF2-derived key
 * Returns ciphertext, authentication tag, IV, and salt
 *
 * @param data - Plaintext data to encrypt (will be converted to UTF-8)
 * @param password - Master password
 * @param keyDerivationIterations - Number of PBKDF2 iterations (default: 100,000)
 * @returns Encryption result containing ciphertext, authTag, IV, and salt
 */
export function encrypt(
  data: string,
  password: string,
  keyDerivationIterations: number = 100000
): EncryptionResult {
  // Generate random salt and derive key
  const { key, salt } = deriveKey(password, undefined, keyDerivationIterations);

  // Generate random IV (96 bits is recommended for GCM)
  const iv = randomBytes(12);

  // Create cipher
  const cipher = createCipheriv('aes-256-gcm', key, iv);

  // Encrypt data
  let ciphertext = cipher.update(data, 'utf8', 'hex');
  ciphertext += cipher.final('hex');

  // Get authentication tag
  const authTag = cipher.getAuthTag();

  return {
    ciphertext: Buffer.from(ciphertext, 'hex').toString('base64'),
    authTag: authTag.toString('base64'),
    iv: iv.toString('base64'),
    salt: salt.toString('base64'),
  };
}

/**
 * Decrypts AES-256-GCM encrypted data using PBKDF2-derived key
 *
 * @param encryptionResult - Encryption result from encrypt() containing ciphertext, authTag, IV, and salt
 * @param password - Master password (must match the password used for encryption)
 * @param keyDerivationIterations - Number of PBKDF2 iterations (should match encryption)
 * @returns Decrypted plaintext
 * @throws Error if authentication fails or decryption fails
 */
export function decrypt(
  encryptionResult: EncryptionResult,
  password: string,
  keyDerivationIterations: number = 100000
): string {
  try {
    // Reconstruct salt and derive key
    const salt = Buffer.from(encryptionResult.salt, 'base64');
    const { key } = deriveKey(password, salt, keyDerivationIterations);

    // Reconstruct IV and auth tag
    const iv = Buffer.from(encryptionResult.iv, 'base64');
    const authTag = Buffer.from(encryptionResult.authTag, 'base64');

    // Decode ciphertext
    const ciphertext = Buffer.from(encryptionResult.ciphertext, 'base64');

    // Create decipher
    const decipher = createDecipheriv('aes-256-gcm', key, iv);
    decipher.setAuthTag(authTag);

    // Decrypt
    let plaintext = decipher.update(ciphertext, undefined, 'utf8');
    plaintext += decipher.final('utf8');

    return plaintext;
  } catch (error) {
    throw new Error(`Decryption failed: ${error instanceof Error ? error.message : String(error)}`);
  }
}

/**
 * Verifies that a password is correct by attempting to decrypt a known encrypted value
 *
 * @param encryptionResult - Encryption result to verify against
 * @param password - Password to test
 * @param keyDerivationIterations - Number of PBKDF2 iterations used during encryption
 * @returns true if password is correct, false otherwise
 */
export function verifyPassword(
  encryptionResult: EncryptionResult,
  password: string,
  keyDerivationIterations: number = 100000
): boolean {
  try {
    decrypt(encryptionResult, password, keyDerivationIterations);
    return true;
  } catch {
    return false;
  }
}
