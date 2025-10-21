/**
 * keyp import command - Import secrets from file
 */

import chalk from 'chalk';
import { readFileSync, existsSync } from 'fs';
import { resolve } from 'path';
import { VaultManager } from '../../vault-manager.js';
import { SecretsManager } from '../../secrets.js';
import {
  promptPassword,
  unlockVaultWithRetry,
  confirm,
  printSuccess,
  printError,
  printInfo,
  printWarning,
} from '../utils.js';

/**
 * Import secrets from file
 *
 * @param inputFile - Input filename (required)
 * @param options - Import options
 */
export async function importCommand(
  inputFile?: string,
  options?: {
    replace?: boolean;
    dryRun?: boolean;
  }
): Promise<void> {
  if (!inputFile) {
    printError('Input file required. Usage: keyp import <file> [options]');
    process.exit(1);
  }

  const filePath = resolve(inputFile);

  // Check if file exists
  if (!existsSync(filePath)) {
    printError(`File not found: ${filePath}`);
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

    // Read import file
    let importData: Record<string, string>;
    try {
      const fileContent = readFileSync(filePath, 'utf-8');
      const parsed = JSON.parse(fileContent);

      // Detect format (plaintext vs encrypted vault)
      if (parsed.version && parsed.crypto && parsed.data) {
        // This is an encrypted vault export - extract the plaintext
        printError('Encrypted vault imports not yet supported. Please export as plaintext.');
        process.exit(1);
      } else if (typeof parsed === 'object' && !Array.isArray(parsed)) {
        // This is plaintext JSON
        importData = parsed;
      } else {
        printError('Invalid import file format. Expected JSON object.');
        process.exit(1);
      }
    } catch (err) {
      printError('Failed to parse import file. Must be valid JSON.');
      process.exit(1);
    }

    // Validate import data (all values should be strings)
    for (const [key, value] of Object.entries(importData)) {
      if (typeof value !== 'string') {
        printError(`Invalid import data: "${key}" value must be a string`);
        process.exit(1);
      }
    }

    const importCount = Object.keys(importData).length;
    if (importCount === 0) {
      printWarning('No secrets in import file');
      process.exit(0);
    }

    // Unlock vault with retry
    await unlockVaultWithRetry(manager);

    const data = manager.getUnlockedData();
    if (!data) {
      printError('Failed to load vault data');
      process.exit(1);
    }

    // Analyze what would be imported
    const existingSecrets = SecretsManager.listSecrets(data);
    let newCount = 0;
    let updateCount = 0;
    const updates: string[] = [];

    for (const key of Object.keys(importData)) {
      if (existingSecrets.includes(key)) {
        updateCount++;
        updates.push(key);
      } else {
        newCount++;
      }
    }

    // Show import summary
    console.log('');
    printInfo('Import summary:');
    printInfo(`  New secrets: ${newCount}`);
    printInfo(`  Updated secrets: ${updateCount}`);
    if (updates.length > 0 && updates.length <= 10) {
      printInfo(`  Secrets to update: ${updates.join(', ')}`);
    }

    // Dry run mode
    if (options?.dryRun) {
      printInfo('Dry run mode - no changes made');
      manager.lockVault();
      process.exit(0);
    }

    // Replace mode confirmation
    if (options?.replace) {
      printWarning('REPLACE mode: This will delete all existing secrets!');
      const confirmed = await confirm('Delete all existing secrets and import?');
      if (!confirmed) {
        printInfo('Import cancelled');
        manager.lockVault();
        process.exit(0);
      }

      // Clear all existing secrets
      const existingKeys = SecretsManager.listSecrets(data);
      for (const key of existingKeys) {
        SecretsManager.deleteSecret(data, key);
      }
    }

    // Import secrets
    for (const [key, value] of Object.entries(importData)) {
      SecretsManager.setSecret(data, key, value);
    }

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
    printSuccess('Secrets imported successfully');
    printInfo(`New secrets: ${newCount}`);
    printInfo(`Updated secrets: ${updateCount}`);
    printInfo(`Total secrets now: ${SecretsManager.getSecretCount(data)}`);
    console.log('');
  } catch (error) {
    if (error instanceof Error && error.message === 'Password entry cancelled') {
      printInfo('Import cancelled');
      manager.lockVault();
      process.exit(0);
    }

    printError(error instanceof Error ? error.message : 'Unknown error');
    manager.lockVault();
    process.exit(1);
  }
}
