/**
 * Type definitions for keyp vault structure and operations
 */

/**
 * Metadata about the cryptographic configuration used for the vault
 */
export interface CryptoConfig {
  /** Encryption algorithm used (currently always "aes-256-gcm") */
  algorithm: 'aes-256-gcm';
  /** Key derivation function used (currently always "pbkdf2") */
  kdf: 'pbkdf2';
  /** Number of PBKDF2 iterations for key derivation */
  iterations: number;
  /** Base64-encoded salt used for key derivation */
  salt: string;
}

/**
 * Complete encrypted vault file structure
 * Stored as JSON and persisted to disk
 */
export interface VaultFile {
  /** Vault format version for future migration compatibility */
  version: string;
  /** Cryptographic configuration and parameters */
  crypto: CryptoConfig;
  /** Base64-encoded encrypted secrets data */
  data: string;
  /** Base64-encoded GCM authentication tag for data integrity verification */
  authTag: string;
  /** Base64-encoded initialization vector used for encryption */
  iv: string;
  /** ISO 8601 timestamp of when vault was created */
  createdAt: string;
  /** ISO 8601 timestamp of when vault was last modified */
  updatedAt: string;
}

/**
 * In-memory representation of vault secrets while unlocked
 * Maps secret names to their string values
 */
export type VaultData = {
  [key: string]: string;
};

/**
 * Result of vault operations
 */
export interface VaultOperationResult {
  success: boolean;
  message?: string;
  error?: string;
}

/**
 * Configuration for vault manager
 */
export interface VaultConfig {
  /** Path to vault file (default: ~/.keyp/vault.json) */
  vaultPath?: string;
  /** Number of PBKDF2 iterations (default: 100,000) */
  keyDerivationIterations?: number;
}
