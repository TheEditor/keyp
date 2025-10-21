/**
 * keyp copy command - Copy a secret to a new name
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
} from '../utils.js';

/**
 * Copy a secret to a new name
 *
 * @param sourceName - Name of secret to copy
 * @param destName - Name for the copy
 */
export async function copyCommand(sourceName?: string, destName?: string): Promise<void> {
  if (!sourceName || !destName) {
    printError('Both names required. Usage: keyp copy <source-name> <dest-name>');
    process.exit(1);
  }

  if (sourceName === destName) {
    printError('Destination must be different from source');
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

    // Check if source secret exists
    if (!SecretsManager.hasSecret(data, sourceName)) {
      printError(`Secret "${chalk.cyan(sourceName)}" not found`);
      manager.lockVault();
      process.exit(1);
    }

    // Check if destination name already exists
    if (SecretsManager.hasSecret(data, destName)) {
      printError(`Secret "${chalk.cyan(destName)}" already exists`);
      manager.lockVault();
      process.exit(1);
    }

    // Get the secret value
    const secretValue = SecretsManager.getSecret(data, sourceName);
    if (!secretValue) {
      printError('Failed to retrieve secret value');
      manager.lockVault();
      process.exit(1);
    }

    // Create copy
    SecretsManager.setSecret(data, destName, secretValue);

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
    printSuccess(`Secret copied: "${chalk.cyan(sourceName)}" → "${chalk.cyan(destName)}"`);
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
