/**
 * keyp init command - Initialize a new vault
 */

import chalk from 'chalk';
import { VaultManager } from '../../vault-manager.js';
import {
  confirmPassword,
  printSuccess,
  printError,
  printWarning,
  printInfo,
  printHint,
  validatePasswordStrength,
  formatVaultPath,
} from '../utils.js';

/**
 * Initialize a new vault
 */
export async function initCommand(): Promise<void> {
  const manager = new VaultManager();

  // Check if vault already exists
  if (manager.vaultFileExists()) {
    printError('Vault already exists at ' + formatVaultPath(manager.getVaultPath()));
    printHint('To use an existing vault, try: keyp set <name>');
    process.exit(1);
  }

  console.log('');
  printInfo('Creating a new vault...');
  console.log('');

  try {
    // Get master password with confirmation
    const password = await confirmPassword();

    if (password.length === 0) {
      printError('Password cannot be empty');
      process.exit(1);
    }

    // Validate password strength
    const { isStrong, feedback } = validatePasswordStrength(password);

    if (!isStrong) {
      printWarning(`Password is weak: ${feedback}`);
    } else {
      printInfo(`Password strength: ${feedback}`);
    }

    console.log('');

    // Initialize vault
    const result = manager.initializeVault(password);

    if (!result.success) {
      printError(result.error || 'Failed to initialize vault');
      process.exit(1);
    }

    // Success!
    printSuccess('Vault initialized successfully!');
    console.log('');
    printInfo(`Location: ${formatVaultPath(manager.getVaultPath())}`);
    console.log('');

    // Show next steps
    console.log(chalk.bold('Next steps:'));
    printHint('1. keyp set <secret-name>   - Store your first secret');
    printHint('2. keyp list                 - List all secrets');
    printHint('3. keyp get <secret-name>    - Retrieve a secret');
    console.log('');
  } catch (error) {
    if (error instanceof Error && error.message === 'Password entry cancelled') {
      printWarning('Vault initialization cancelled');
      process.exit(0);
    }

    printError(error instanceof Error ? error.message : 'Unknown error');
    process.exit(1);
  }
}
