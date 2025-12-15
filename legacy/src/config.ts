/**
 * Configuration management for keyp
 * Handles vault paths, directories, and settings
 */

import { existsSync, mkdirSync } from 'fs';
import { join } from 'path';
import { homedir } from 'os';

/**
 * Default vault directory in user's home
 */
const KEYP_DIR = join(homedir(), '.keyp');

/**
 * Default vault file name
 */
const VAULT_FILE = 'vault.json';

/**
 * Gets the path to the .keyp directory, creating it if it doesn't exist
 *
 * @returns Full path to .keyp directory
 */
export function getKeypDir(): string {
  if (!existsSync(KEYP_DIR)) {
    mkdirSync(KEYP_DIR, { recursive: true, mode: 0o700 });
  }
  return KEYP_DIR;
}

/**
 * Gets the path to the vault file with optional override
 *
 * @param customPath - Optional custom vault path
 * @returns Full path to vault file
 */
export function getVaultPath(customPath?: string): string {
  if (customPath) {
    return customPath;
  }
  return join(getKeypDir(), VAULT_FILE);
}

/**
 * Checks if vault exists at the given path
 *
 * @param vaultPath - Path to vault file (uses default if not provided)
 * @returns true if vault file exists
 */
export function vaultExists(vaultPath?: string): boolean {
  const path = getVaultPath(vaultPath);
  return existsSync(path);
}

/**
 * Ensures the .keyp directory exists with proper permissions
 * Creates it with 0o700 (rwx------) for security
 */
export function ensureKeypDirExists(): void {
  const dir = getKeypDir();
  if (!existsSync(dir)) {
    mkdirSync(dir, { recursive: true, mode: 0o700 });
  }
}

/**
 * Default configuration values
 */
export const DEFAULT_CONFIG = {
  /** Default number of PBKDF2 iterations */
  KEY_DERIVATION_ITERATIONS: 100000,
  /** Default vault file name */
  VAULT_FILE_NAME: VAULT_FILE,
  /** Current vault format version */
  VAULT_VERSION: '1.0.0',
} as const;
