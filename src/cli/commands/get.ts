/**
 * keyp get command - Retrieve a secret (with clipboard support)
 */

import chalk from 'chalk';
import clipboard from 'clipboardy';
import { VaultManager } from '../../vault-manager.js';
import { SecretsManager } from '../../secrets.js';
import {
  unlockVaultWithRetry,
  printSuccess,
  printError,
  printInfo,
  printWarning,
} from '../utils.js';

/**
 * Retrieve a secret from the vault
 *
 * @param name - Secret name to retrieve
 * @param options - Command options
 */
export async function getCommand(
  name?: string,
  options?: {
    stdout?: boolean;
    noClear?: boolean;
  }
): Promise<void> {
  if (!name) {
    printError('Secret name required. Usage: keyp get <name>');
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

    // Get secret
    const secret = SecretsManager.getSecret(data, name);
    if (!secret) {
      printError(`Secret "${chalk.cyan(name)}" not found`);
      manager.lockVault();
      process.exit(1);
    }

    // Lock vault first
    manager.lockVault();

    // Handle output
    if (options?.stdout) {
      // Print to stdout (warn about visibility)
      printWarning('Output to terminal (secret will be visible!)');
      console.log('');
      console.log(secret);
      console.log('');
    } else {
      // Copy to clipboard
      try {
        await clipboard.write(secret);
        printSuccess('Copied to clipboard');

        const clearTime = 45; // seconds
        printInfo(`Will clear in ${clearTime} seconds`);
        console.log('');

        if (!options?.noClear) {
          // Auto-clear clipboard after specified time
          setTimeout(async () => {
            try {
              const current = await clipboard.read();
              if (current === secret) {
                await clipboard.write('');
                printInfo('Clipboard cleared');
              }
            } catch (err) {
              // Silently fail if clipboard operations fail
            }
          }, clearTime * 1000);
        }
      } catch (err) {
        // Clipboard not available, fall back to stdout
        printWarning('Clipboard not available, showing on screen:');
        console.log('');
        console.log(secret);
        console.log('');
      }
    }
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
