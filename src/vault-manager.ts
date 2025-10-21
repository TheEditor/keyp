/**
 * Vault management - handles initialization, unlocking, and persistence
 */

import { readFileSync, writeFileSync } from 'fs';
import { encrypt, decrypt, EncryptionResult } from './crypto';
import { VaultFile, VaultData, CryptoConfig } from './types';
import { getVaultPath, vaultExists, DEFAULT_CONFIG, ensureKeypDirExists } from './config';

/**
 * VaultManager handles all vault lifecycle operations
 * Manages encryption, persistence, and in-memory state
 */
export class VaultManager {
  private vaultPath: string;
  private keyDerivationIterations: number;
  private unlockedData: VaultData | null = null;
  private isUnlocked: boolean = false;

  /**
   * Creates a new VaultManager instance
   *
   * @param vaultPath - Optional custom path to vault file
   * @param keyDerivationIterations - Optional number of PBKDF2 iterations
   */
  constructor(vaultPath?: string, keyDerivationIterations?: number) {
    this.vaultPath = getVaultPath(vaultPath);
    this.keyDerivationIterations = keyDerivationIterations || DEFAULT_CONFIG.KEY_DERIVATION_ITERATIONS;
  }

  /**
   * Initializes a new vault with a master password
   * Fails if vault already exists
   *
   * @param password - Master password for vault
   * @returns Result indicating success or failure
   */
  initializeVault(password: string): { success: boolean; message?: string; error?: string } {
    try {
      if (vaultExists(this.vaultPath)) {
        return {
          success: false,
          error: 'Vault already exists. Use unlock() instead.',
        };
      }

      ensureKeypDirExists();

      // Start with empty vault data
      const emptyData: VaultData = {};

      // Encrypt empty data to get encryption result
      const encryptionResult = encrypt(JSON.stringify(emptyData), password, this.keyDerivationIterations);

      // Build crypto config
      const cryptoConfig: CryptoConfig = {
        algorithm: 'aes-256-gcm',
        kdf: 'pbkdf2',
        iterations: this.keyDerivationIterations,
        salt: encryptionResult.salt,
      };

      // Build vault file
      const now = new Date().toISOString();
      const vaultFile: VaultFile = {
        version: DEFAULT_CONFIG.VAULT_VERSION,
        crypto: cryptoConfig,
        data: encryptionResult.ciphertext,
        authTag: encryptionResult.authTag,
        iv: encryptionResult.iv,
        createdAt: now,
        updatedAt: now,
      };

      // Write vault to disk
      writeFileSync(this.vaultPath, JSON.stringify(vaultFile, null, 2));

      // Mark vault as unlocked with empty data
      this.unlockedData = emptyData;
      this.isUnlocked = true;

      return { success: true, message: 'Vault initialized successfully' };
    } catch (error) {
      return {
        success: false,
        error: `Failed to initialize vault: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * Unlocks the vault using the master password
   * Decrypts and loads vault data into memory
   *
   * @param password - Master password for vault
   * @returns Result indicating success or failure
   */
  unlockVault(password: string): { success: boolean; message?: string; error?: string } {
    try {
      if (!vaultExists(this.vaultPath)) {
        return {
          success: false,
          error: 'Vault does not exist. Use initializeVault() first.',
        };
      }

      if (this.isUnlocked) {
        return { success: true, message: 'Vault is already unlocked' };
      }

      // Read vault file
      const vaultContent = readFileSync(this.vaultPath, 'utf-8');
      const vaultFile: VaultFile = JSON.parse(vaultContent);

      // Prepare encryption result for decryption
      const encryptionResult: EncryptionResult = {
        ciphertext: vaultFile.data,
        authTag: vaultFile.authTag,
        iv: vaultFile.iv,
        salt: vaultFile.crypto.salt,
      };

      // Decrypt data
      const decryptedData = decrypt(encryptionResult, password, this.keyDerivationIterations);

      // Parse decrypted JSON
      this.unlockedData = JSON.parse(decryptedData);
      this.isUnlocked = true;

      return { success: true, message: 'Vault unlocked successfully' };
    } catch (error) {
      return {
        success: false,
        error: `Failed to unlock vault: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * Locks the vault, clearing in-memory data
   * Does not affect the encrypted file on disk
   */
  lockVault(): void {
    this.unlockedData = null;
    this.isUnlocked = false;
  }

  /**
   * Saves the current in-memory vault data to disk (encrypted)
   *
   * @param password - Master password (needed to re-encrypt data)
   * @returns Result indicating success or failure
   */
  saveVault(password: string): { success: boolean; message?: string; error?: string } {
    try {
      if (!this.isUnlocked || this.unlockedData === null) {
        return {
          success: false,
          error: 'Vault is not unlocked. Cannot save.',
        };
      }

      // Encrypt current data
      const encryptionResult = encrypt(JSON.stringify(this.unlockedData), password, this.keyDerivationIterations);

      // Read existing vault to preserve metadata
      const vaultContent = readFileSync(this.vaultPath, 'utf-8');
      const vaultFile: VaultFile = JSON.parse(vaultContent);

      // Update encrypted data and timestamp
      vaultFile.data = encryptionResult.ciphertext;
      vaultFile.authTag = encryptionResult.authTag;
      vaultFile.iv = encryptionResult.iv;
      vaultFile.crypto.salt = encryptionResult.salt;
      vaultFile.updatedAt = new Date().toISOString();

      // Write updated vault to disk
      writeFileSync(this.vaultPath, JSON.stringify(vaultFile, null, 2));

      return { success: true, message: 'Vault saved successfully' };
    } catch (error) {
      return {
        success: false,
        error: `Failed to save vault: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * Gets the unlocked vault data
   *
   * @returns Vault data if unlocked, null otherwise
   */
  getUnlockedData(): VaultData | null {
    return this.unlockedData;
  }

  /**
   * Checks if vault is currently unlocked
   *
   * @returns true if unlocked and data is in memory
   */
  isVaultUnlocked(): boolean {
    return this.isUnlocked;
  }

  /**
   * Checks if vault file exists on disk
   *
   * @returns true if vault file exists
   */
  vaultFileExists(): boolean {
    return vaultExists(this.vaultPath);
  }

  /**
   * Gets the vault file path
   *
   * @returns Path to vault file
   */
  getVaultPath(): string {
    return this.vaultPath;
  }
}
