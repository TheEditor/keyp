/**
 * keyp stats command - Display vault statistics
 */

import chalk from 'chalk';
import { VaultManager } from '../../vault-manager.js';
import { promises as fs } from 'fs';
import path from 'path';
import { getKeypDir, getVaultPath } from '../../config.js';
import { printError, printInfo, unlockVaultWithRetry } from '../utils.js';

/**
 * Display vault statistics
 */
export async function statsCommand(): Promise<void> {
  try {
    const manager = new VaultManager();
    const vaultPath = getVaultPath();

    // Check if vault exists
    if (!manager.vaultFileExists()) {
      printError('No vault found');
      printInfo('Run "keyp init" to create a vault');
      process.exit(1);
    }

    // Unlock vault
    await unlockVaultWithRetry(manager);

    // Get vault data
    const data = manager.getUnlockedData();
    if (!data) {
      printError('Failed to access vault data');
      process.exit(1);
    }

    // Get file stats
    const stats = await fs.stat(vaultPath);
    const secretCount = Object.keys(data).length;
    const vaultSize = stats.size;
    const lastModified = new Date(stats.mtime);

    // Check for Git sync time
    const keypDir = getKeypDir();
    const syncTimePath = path.join(keypDir, '.keyp-sync-time');
    let lastSync: Date | undefined;
    try {
      const syncContent = await fs.readFile(syncTimePath, 'utf-8');
      lastSync = new Date(syncContent);
    } catch {
      lastSync = undefined;
    }

    // Display statistics
    console.log('');
    console.log(chalk.cyan.bold('📊 Vault Statistics'));
    console.log(chalk.gray('─'.repeat(60)));
    console.log('');

    // Secrets
    console.log(`${chalk.bold('Secrets')}`);
    console.log(`  Total: ${chalk.cyan(secretCount.toString())}`);
    if (secretCount > 0) {
      const secretNames = Object.keys(data).sort();
      const avgLength = Math.round(
        secretNames.reduce((sum, name) => sum + (data[name] || '').toString().length, 0) / secretCount
      );
      console.log(`  Average value length: ${chalk.cyan(avgLength.toString())} characters`);
      console.log(`  Longest name: ${chalk.cyan(secretNames.reduce((a, b) => (a.length > b.length ? a : b)))}`);
    }

    console.log('');

    // Storage
    console.log(`${chalk.bold('Storage')}`);
    const vaultSizeKB = (vaultSize / 1024).toFixed(2);
    const vaultSizeMB = (vaultSize / (1024 * 1024)).toFixed(2);
    const sizeDisplay = vaultSize > 1024 * 1024 ? `${vaultSizeMB} MB` : `${vaultSizeKB} KB`;
    console.log(`  Vault file size: ${chalk.cyan(sizeDisplay)}`);
    console.log(`  Location: ${chalk.gray(formatPath(vaultPath))}`);

    console.log('');

    // Dates
    console.log(`${chalk.bold('Dates')}`);
    console.log(`  Last modified: ${chalk.cyan(formatDate(lastModified))}`);
    if (lastSync) {
      console.log(`  Last synced: ${chalk.cyan(formatDate(lastSync))}`);
    } else {
      console.log(`  Last synced: ${chalk.gray('Never (Git sync not configured)')}`);
    }

    console.log('');

    // Encryption info
    console.log(`${chalk.bold('Encryption')}`);
    console.log(`  Algorithm: ${chalk.cyan('AES-256-GCM')}`);
    console.log(`  Key derivation: ${chalk.cyan('PBKDF2-SHA256')}`);
    console.log(`  Iterations: ${chalk.cyan('100,000+')}`);

    console.log('');

    // Lock vault
    manager.lockVault();
  } catch (error) {
    if (error instanceof Error && error.message === 'Password entry cancelled') {
      process.exit(0);
    }

    printError(error instanceof Error ? error.message : 'Unknown error');
    process.exit(1);
  }
}

/**
 * Format file path for display
 */
function formatPath(filePath: string): string {
  const homeDir = process.env.HOME || process.env.USERPROFILE || '~';
  return filePath.replace(homeDir, '~');
}

/**
 * Format date for display
 */
function formatDate(date: Date): string {
  return date.toLocaleString();
}
