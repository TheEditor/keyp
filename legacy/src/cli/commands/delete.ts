/**
 * keyp delete command - Delete a secret from vault
 */

import chalk from 'chalk';
import { VaultManager } from '../../vault-manager.js';
import { SecretsManager } from '../../secrets.js';
import {
  promptPassword,
  confirm,
  unlockVaultWithRetry,
  printSuccess,
  printError,
  printInfo,
} from '../utils.js';

/**
 * Delete a secret from the vault
 *
 * @param name - Secret name to delete
 * @param options - Command options
 */
export async function deleteCommand(
  name?: string,
  options?: {
    force?: boolean;
  }
): Promise<void> {
  if (!name) {
    printError('Secret name required. Usage: keyp delete <name>');
    process.exit(1);
  }

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

    // Check if secret exists
    if (!SecretsManager.hasSecret(data, name)) {
      printError(`Secret "${chalk.cyan(name)}" not found`);
      manager.lockVault();
      process.exit(1);
    }

    // Confirm deletion (unless --force)
    if (!options?.force) {
      const confirmed = await confirm(`Delete secret "${chalk.cyan(name)}"?`);
      if (!confirmed) {
        printInfo('Deletion cancelled');
        manager.lockVault();
        process.exit(0);
      }
    }

    // Delete secret
    const result = SecretsManager.deleteSecret(data, name);
    if (!result.success) {
      printError(result.error || 'Failed to delete secret');
      manager.lockVault();
      process.exit(1);
    }

    // Save vault
    const password = await promptPassword('Enter master password to save');
    const saveResult = await manager.saveVault(password);
    if (!saveResult.success) {
      printError(saveResult.error || 'Failed to save vault');
      process.exit(1);
    }

    // Lock vault
    manager.lockVault();

    // Success!
    console.log('');
    printSuccess(`Secret "${chalk.cyan(name)}" deleted`);
    printInfo(`Remaining secrets: ${SecretsManager.getSecretCount(data)}`);
    console.log('');
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
