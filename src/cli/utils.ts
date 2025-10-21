/**
 * Shared CLI utilities for password prompts, error handling, and formatting
 */

import prompts from 'prompts';
import chalk from 'chalk';
import { VaultManager } from '../vault-manager.js';

/**
 * Prompts user for a password with masking
 *
 * @param message - Prompt message to display
 * @returns The entered password
 */
export async function promptPassword(message: string = 'Enter master password'): Promise<string> {
  const response = await prompts({
    type: 'password',
    name: 'password',
    message: message,
  });

  if (response.password === undefined) {
    throw new Error('Password entry cancelled');
  }

  return response.password;
}

/**
 * Prompts user to confirm password by entering it twice
 *
 * @returns The password if both entries match
 */
export async function confirmPassword(): Promise<string> {
  const password1 = await prompts({
    type: 'password',
    name: 'password',
    message: 'Enter master password',
  });

  if (password1.password === undefined) {
    throw new Error('Password entry cancelled');
  }

  const password2 = await prompts({
    type: 'password',
    name: 'password',
    message: 'Confirm master password',
  });

  if (password2.password === undefined) {
    throw new Error('Password entry cancelled');
  }

  if (password1.password !== password2.password) {
    throw new Error('Passwords do not match');
  }

  return password1.password;
}

/**
 * Prompts user for a string value
 *
 * @param message - Prompt message
 * @param initial - Optional initial value
 * @returns The entered string
 */
export async function promptString(message: string, initial?: string): Promise<string> {
  const response = await prompts({
    type: 'text',
    name: 'value',
    message: message,
    initial: initial,
  });

  if (response.value === undefined) {
    throw new Error('Input cancelled');
  }

  return response.value;
}

/**
 * Prompts user for confirmation (yes/no)
 *
 * @param message - Confirmation message
 * @returns true if user confirms, false otherwise
 */
export async function confirm(message: string): Promise<boolean> {
  const response = await prompts({
    type: 'confirm',
    name: 'value',
    message: message,
    initial: false,
  });

  return response.value === true;
}

/**
 * Unlocks vault with password, retrying on failure
 *
 * @param manager - VaultManager instance
 * @param maxAttempts - Maximum number of password attempts (default: 3)
 * @throws Error after max attempts reached
 */
export async function unlockVaultWithRetry(manager: VaultManager, maxAttempts: number = 3): Promise<void> {
  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      const password = await promptPassword('Enter master password');
      const result = manager.unlockVault(password);

      if (result.success) {
        return;
      }

      const remaining = maxAttempts - attempt;
      if (remaining > 0) {
        printError(`Incorrect password (${remaining} attempt${remaining === 1 ? '' : 's'} remaining)`);
      } else {
        throw new Error('Maximum password attempts exceeded');
      }
    } catch (error) {
      if (error instanceof Error && error.message === 'Password entry cancelled') {
        throw error;
      }

      const remaining = maxAttempts - attempt;
      if (remaining > 0) {
        printError(`Incorrect password (${remaining} attempt${remaining === 1 ? '' : 's'} remaining)`);
      } else {
        throw new Error('Maximum password attempts exceeded');
      }
    }
  }
}

/**
 * Prints a success message in green with checkmark
 *
 * @param message - Message to display
 */
export function printSuccess(message: string): void {
  console.log(chalk.green(`✓ ${message}`));
}

/**
 * Prints an error message in red with X
 *
 * @param message - Message to display
 */
export function printError(message: string): void {
  console.error(chalk.red(`✗ ${message}`));
}

/**
 * Prints a warning message in yellow with exclamation
 *
 * @param message - Message to display
 */
export function printWarning(message: string): void {
  console.log(chalk.yellow(`⚠ ${message}`));
}

/**
 * Prints an info message in cyan with info icon
 *
 * @param message - Message to display
 */
export function printInfo(message: string): void {
  console.log(chalk.cyan(`ℹ ${message}`));
}

/**
 * Prints a hint message in gray
 *
 * @param message - Message to display
 */
export function printHint(message: string): void {
  console.log(chalk.gray(`  ${message}`));
}

/**
 * Prints a list of secret names with colors
 *
 * @param secrets - Array of secret names
 */
export function printSecretList(secrets: string[]): void {
  if (secrets.length === 0) {
    printInfo('No secrets yet. Try: keyp set <name>');
    return;
  }

  console.log('');
  for (const secret of secrets) {
    console.log(chalk.cyan(`  • ${secret}`));
  }
  console.log('');
  console.log(chalk.gray(`${secrets.length} secret${secrets.length === 1 ? '' : 's'} stored`));
}

/**
 * Validates password strength and returns feedback
 *
 * @param password - Password to validate
 * @returns Object with isStrong and feedback message
 */
export function validatePasswordStrength(password: string): { isStrong: boolean; feedback: string } {
  const issues: string[] = [];

  if (password.length < 8) {
    issues.push('at least 8 characters');
  }

  if (password.length < 12) {
    issues.push('consider 12+ characters for better security');
  }

  if (!/[A-Z]/.test(password)) {
    issues.push('mix in uppercase letters');
  }

  if (!/[a-z]/.test(password)) {
    issues.push('mix in lowercase letters');
  }

  if (!/[0-9]/.test(password)) {
    issues.push('add some numbers');
  }

  if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
    issues.push('add special characters');
  }

  const isStrong = issues.length <= 1; // Allow one weakness
  const feedback = issues.length === 0 ? 'Strong password!' : `Consider: ${issues.join(', ')}`;

  return { isStrong, feedback };
}

/**
 * Formats vault path for display
 *
 * @param vaultPath - Full vault path
 * @returns Formatted path with home directory shortened
 */
export function formatVaultPath(vaultPath: string): string {
  const homeDir = process.env.HOME || process.env.USERPROFILE || '~';
  return vaultPath.replace(homeDir, '~');
}

/**
 * Clears the terminal screen
 */
export function clearScreen(): void {
  console.clear();
}

/**
 * Prints a divider line
 */
export function printDivider(): void {
  console.log(chalk.gray('─'.repeat(60)));
}

/**
 * Prints the keyp header/banner
 */
export function printBanner(): void {
  console.log(chalk.cyan.bold('🔒 keyp'));
  console.log(chalk.gray('Local-first secret manager for developers'));
  printDivider();
  console.log('');
}

/**
 * Prints vault info (path, secret count)
 *
 * @param manager - VaultManager instance
 */
export function printVaultInfo(manager: VaultManager): void {
  const vaultPath = formatVaultPath(manager.getVaultPath());
  const data = manager.getUnlockedData();
  const count = data ? Object.keys(data).length : 0;

  printInfo(`Vault: ${vaultPath}`);
  printInfo(`Secrets: ${count}`);
}
