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
 * Calculates password entropy and returns approximate bits of entropy
 *
 * @param password - Password to analyze
 * @returns Number of entropy bits
 */
function calculateEntropy(password: string): number {
  let charsetSize = 0;

  if (/[a-z]/.test(password)) charsetSize += 26;
  if (/[A-Z]/.test(password)) charsetSize += 26;
  if (/[0-9]/.test(password)) charsetSize += 10;
  if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) charsetSize += 32;

  return Math.log2(Math.pow(charsetSize, password.length));
}

/**
 * Validates password strength and returns detailed feedback with visual indicator
 *
 * @param password - Password to validate
 * @returns Object with isStrong, strength score, and detailed feedback
 */
export function validatePasswordStrength(password: string): {
  isStrong: boolean;
  feedback: string;
  score: number;
  meter: string;
} {
  const issues: string[] = [];
  const strengths: string[] = [];
  let score = 0;

  // Check minimum length
  if (password.length < 8) {
    issues.push('Increase length to at least 8 characters');
  } else if (password.length < 12) {
    issues.push('Increase length to 12+ characters for better security');
  } else {
    strengths.push(`Good length (${password.length} characters)`);
  }

  // Check for uppercase
  if (!/[A-Z]/.test(password)) {
    issues.push('Add uppercase letters (A-Z)');
  } else {
    strengths.push('Uppercase letters included');
  }

  // Check for lowercase
  if (!/[a-z]/.test(password)) {
    issues.push('Add lowercase letters (a-z)');
  } else {
    strengths.push('Lowercase letters included');
  }

  // Check for numbers
  if (!/[0-9]/.test(password)) {
    issues.push('Add numbers (0-9)');
  } else {
    strengths.push('Numbers included');
  }

  // Check for special characters
  if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
    issues.push('Add special characters (!@#$%^&*)');
  } else {
    strengths.push('Special characters included');
  }

  // Calculate entropy-based score
  const entropy = calculateEntropy(password);
  score = Math.min(100, Math.floor((entropy / 50) * 100)); // 50 bits = 100 score

  // Determine if strong (75+ score or no critical issues)
  const isStrong = score >= 75 || issues.length <= 1;

  // Create visual strength meter
  const barLength = 10;
  const filledLength = Math.floor((score / 100) * barLength);
  const meter = chalk.green('█'.repeat(filledLength)) + chalk.gray('░'.repeat(barLength - filledLength));

  // Build feedback string
  let feedbackLines: string[] = [];
  feedbackLines.push(`Strength: ${meter} ${score}%`);

  if (strengths.length > 0) {
    for (const strength of strengths) {
      feedbackLines.push(chalk.green(`  ✓ ${strength}`));
    }
  }

  if (issues.length > 0) {
    for (const issue of issues) {
      feedbackLines.push(chalk.yellow(`  ✗ ${issue}`));
    }
  }

  const feedback = feedbackLines.join('\n');

  return { isStrong, feedback, score, meter };
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
