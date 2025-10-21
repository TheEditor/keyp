/**
 * keyp rename command - Rename a secret
 */

import chalk from 'chalk';
import { VaultManager } from '../../vault-manager.js';
import { SecretsManager } from '../../secrets.js';
import {
  promptPassword,
  unlockVaultWithRetry,
  printSuccess,
  printError,
  printInfo,
  printHint,
} from '../utils.js';

/**
 * Rename a secret in the vault
 *
 * @param oldName - Current secret name
 * @param newName - New secret name
 */
export async function renameCommand(oldName?: string, newName?: string): Promise<void> {
  if (!oldName || !newName) {
    printError('Both names required. Usage: keyp rename <old-name> <new-name>');
    process.exit(1);
  }

  if (oldName === newName) {
    printError('New name must be different from old name');
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

    // Check if old secret exists
    if (!SecretsManager.hasSecret(data, oldName)) {
      printError(`Secret "${chalk.cyan(oldName)}" not found`);
      manager.lockVault();
      process.exit(1);
    }

    // Check if new name already exists
    if (SecretsManager.hasSecret(data, newName)) {
      printError(`Secret "${chalk.cyan(newName)}" already exists`);
      manager.lockVault();
      process.exit(1);
    }

    // Get the secret value
    const secretValue = SecretsManager.getSecret(data, oldName);
    if (!secretValue) {
      printError('Failed to retrieve secret value');
      manager.lockVault();
      process.exit(1);
    }

    // Delete old, add new
    SecretsManager.deleteSecret(data, oldName);
    SecretsManager.setSecret(data, newName, secretValue);

    // Save vault
    const password = await promptPassword('Enter master password to save');
    const saveResult = manager.saveVault(password);
    if (!saveResult.success) {
      printError(saveResult.error || 'Failed to save vault');
      process.exit(1);
    }

    // Lock vault
    manager.lockVault();

    // Success!
    console.log('');
    printSuccess(`Secret renamed: "${chalk.cyan(oldName)}" → "${chalk.cyan(newName)}"`);
    printInfo(`Total secrets: ${SecretsManager.getSecretCount(data)}`);
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
