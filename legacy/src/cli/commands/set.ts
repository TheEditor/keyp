/**
 * keyp set command - Store a secret
 */

import chalk from 'chalk';
import { VaultManager } from '../../vault-manager.js';
import { SecretsManager } from '../../secrets.js';
import {
  promptPassword,
  promptString,
  unlockVaultWithRetry,
  printSuccess,
  printError,
  printInfo,
  printHint,
} from '../utils.js';

/**
 * Store a secret in the vault
 *
 * @param name - Secret name/key
 * @param value - Optional secret value (prompts if not provided)
 */
export async function setCommand(name?: string, value?: string): Promise<void> {
  if (!name) {
    printError('Secret name required. Usage: keyp set <name> [value]');
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

    // Get secret value
    let secretValue = value;
    if (!secretValue) {
      secretValue = await promptString(`Enter value for "${chalk.cyan(name)}"`, '');
    }

    if (!secretValue) {
      printError('Secret value cannot be empty');
      process.exit(1);
    }

    // Set secret
    const result = SecretsManager.setSecret(data, name, secretValue);
    if (!result.success) {
      printError(result.error || 'Failed to set secret');
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
    printSuccess(`Secret "${chalk.cyan(name)}" saved`);
    printInfo(`Total secrets: ${SecretsManager.getSecretCount(data)}`);
    console.log('');
    printHint(`Retrieve with: keyp get ${name}`);
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
