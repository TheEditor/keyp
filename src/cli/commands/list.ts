/**
 * keyp list command - List all secrets
 */

import chalk from 'chalk';
import { VaultManager } from '../../vault-manager.js';
import { SecretsManager } from '../../secrets.js';
import {
  unlockVaultWithRetry,
  printError,
  printInfo,
  printSecretList,
} from '../utils.js';

/**
 * List all secrets in the vault
 *
 * @param options - Command options (optional)
 */
export async function listCommand(options?: { search?: string; count?: boolean }): Promise<void> {
  const manager = new VaultManager();

  // Check if vault exists
  if (!manager.vaultFileExists()) {
    printError('Vault not found. Run "keyp init" to create one.');
    process.exit(1);
  }

  try {
    console.log('');

    // Unlock vault with retry
    await unlockVaultWithRetry(manager);

    const data = manager.getUnlockedData();
    if (!data) {
      printError('Failed to load vault data');
      process.exit(1);
    }

    // Get secrets list
    let secrets = SecretsManager.listSecrets(data);

    // Filter by search pattern if provided
    if (options?.search) {
      secrets = SecretsManager.searchSecrets(data, options.search);
      printInfo(`Search results for "${chalk.cyan(options.search)}"`);
      console.log('');
    }

    // Print secrets
    if (options?.count) {
      console.log(`${secrets.length} secret${secrets.length === 1 ? '' : 's'}`);
    } else {
      printSecretList(secrets);
    }

    // Lock vault
    manager.lockVault();
  } catch (error) {
    if (error instanceof Error && error.message === 'Password entry cancelled') {
      printInfo('Operation cancelled');
      manager.lockVault();
      process.exit(0);
    }

    printError(error instanceof Error ? error.message : 'Unknown error');
    manager.lockVault();
    process.exit(1);
  }
}
