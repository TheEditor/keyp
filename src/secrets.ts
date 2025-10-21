/**
 * Secret management - CRUD operations for vault secrets
 */

import { VaultData } from './types.js';

/**
 * SecretsManager handles all secret operations (Create, Read, Update, Delete)
 */
export class SecretsManager {
  /**
   * Sets a secret (creates or updates)
   *
   * @param data - Vault data object to modify
   * @param key - Secret name/key
   * @param value - Secret value
   * @returns Result indicating success or failure
   */
  static setSecret(data: VaultData, key: string, value: string): { success: boolean; message?: string; error?: string } {
    try {
      if (!key || key.trim().length === 0) {
        return {
          success: false,
          error: 'Secret name cannot be empty',
        };
      }

      if (!value || value.trim().length === 0) {
        return {
          success: false,
          error: 'Secret value cannot be empty',
        };
      }

      const isUpdate = key in data;
      data[key] = value;

      const action = isUpdate ? 'updated' : 'created';
      return {
        success: true,
        message: `Secret "${key}" ${action}`,
      };
    } catch (error) {
      return {
        success: false,
        error: `Failed to set secret: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * Gets a secret value by key
   *
   * @param data - Vault data object to read from
   * @param key - Secret name/key
   * @returns Secret value or null if not found
   */
  static getSecret(data: VaultData, key: string): string | null {
    return data[key] || null;
  }

  /**
   * Checks if a secret exists
   *
   * @param data - Vault data object
   * @param key - Secret name/key
   * @returns true if secret exists
   */
  static hasSecret(data: VaultData, key: string): boolean {
    return key in data;
  }

  /**
   * Deletes a secret by key
   *
   * @param data - Vault data object to modify
   * @param key - Secret name/key
   * @returns Result indicating success or failure
   */
  static deleteSecret(data: VaultData, key: string): { success: boolean; message?: string; error?: string } {
    try {
      if (!this.hasSecret(data, key)) {
        return {
          success: false,
          error: `Secret "${key}" not found`,
        };
      }

      delete data[key];
      return {
        success: true,
        message: `Secret "${key}" deleted`,
      };
    } catch (error) {
      return {
        success: false,
        error: `Failed to delete secret: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * Lists all secret keys (without revealing values)
   *
   * @param data - Vault data object
   * @returns Array of secret keys, sorted alphabetically
   */
  static listSecrets(data: VaultData): string[] {
    return Object.keys(data).sort();
  }

  /**
   * Gets the count of secrets
   *
   * @param data - Vault data object
   * @returns Number of secrets stored
   */
  static getSecretCount(data: VaultData): number {
    return Object.keys(data).length;
  }

  /**
   * Clears all secrets from vault (dangerous operation)
   * Requires explicit confirmation
   *
   * @param data - Vault data object to modify
   * @param confirmationKey - Special confirmation flag
   * @returns Result indicating success or failure
   */
  static clearAllSecrets(data: VaultData, confirmationKey: string = ''): { success: boolean; message?: string; error?: string } {
    try {
      if (confirmationKey !== 'CONFIRM_DELETE_ALL') {
        return {
          success: false,
          error: 'Confirmation not provided. Pass "CONFIRM_DELETE_ALL" as confirmationKey.',
        };
      }

      const count = Object.keys(data).length;
      for (const key of Object.keys(data)) {
        delete data[key];
      }

      return {
        success: true,
        message: `All ${count} secrets deleted`,
      };
    } catch (error) {
      return {
        success: false,
        error: `Failed to clear secrets: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * Searches for secrets by key pattern (case-insensitive)
   *
   * @param data - Vault data object
   * @param pattern - Search pattern (substring match)
   * @returns Array of matching secret keys
   */
  static searchSecrets(data: VaultData, pattern: string): string[] {
    const lowerPattern = pattern.toLowerCase();
    return Object.keys(data)
      .filter((key) => key.toLowerCase().includes(lowerPattern))
      .sort();
  }
}
