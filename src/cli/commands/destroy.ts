/**
 * keyp destroy command - Permanently delete entire vault
 * WARNING: This action cannot be undone!
 */

import chalk from 'chalk';
import { promises as fs } from 'fs';
import path from 'path';
import prompts from 'prompts';
import { getKeypDir, getVaultPath } from '../../config.js';
import { printError, printWarning, promptPassword } from '../utils.js';

/**
 * Destroy the entire vault permanently
 */
export async function destroyCommand(): Promise<void> {
  try {
    const keypDir = getKeypDir();
    const vaultPath = getVaultPath();

    // Display severe warning
    console.log('');
    printWarning('╔════════════════════════════════════════════════════╗');
    printWarning('║                   DANGER ZONE                       ║');
    printWarning('╚════════════════════════════════════════════════════╝');
    console.log('');
    printWarning('This action will PERMANENTLY DELETE your entire vault!');
    printWarning('All secrets will be lost FOREVER.');
    printWarning('This action CANNOT be undone or recovered.');
    console.log('');

    // Require explicit confirmation
    let confirmed = false;

    while (!confirmed) {
      const response = await prompts({
        type: 'text',
        name: 'confirm',
        message: 'Type "destroy" to confirm deletion',
      });

      if (response.confirm === 'destroy') {
        confirmed = true;
      } else if (response.confirm === undefined) {
        // User cancelled (Ctrl+C or similar)
        console.log('');
        printWarning('Vault destruction cancelled');
        console.log('');
        return;
      } else {
        // User entered something else, show error and loop
        printWarning('Type "destroy" exactly to confirm');
        console.log('');
      }
    }

    // Require password verification
    const password = await promptPassword('Enter master password to verify');

    // Check if vault exists
    const vaultExists = await checkFileExists(vaultPath);
    if (!vaultExists) {
      printError('Vault file not found');
      process.exit(1);
    }

    // Delete vault file
    console.log('');
    console.log(chalk.gray('Destroying vault...'));

    try {
      await fs.unlink(vaultPath);
      console.log(chalk.green(`✓ Deleted vault: ${formatPath(vaultPath)}`));
    } catch (err) {
      const error = err as Error;
      printError(`Failed to delete vault: ${error.message}`);
      process.exit(1);
    }

    // Delete config files if they exist
    const configPath = path.join(keypDir, '.keyp-config.json');
    const gitConfigPath = path.join(keypDir, '.keyp-git-config.json');
    const syncTimePath = path.join(keypDir, '.keyp-sync-time');

    try {
      if (await checkFileExists(configPath)) {
        await fs.unlink(configPath);
        console.log(chalk.green(`✓ Deleted config: ${formatPath(configPath)}`));
      }
    } catch {
      // Config file might not exist, that's okay
    }

    try {
      if (await checkFileExists(gitConfigPath)) {
        await fs.unlink(gitConfigPath);
        console.log(chalk.green(`✓ Deleted git config: ${formatPath(gitConfigPath)}`));
      }
    } catch {
      // Git config might not exist, that's okay
    }

    try {
      if (await checkFileExists(syncTimePath)) {
        await fs.unlink(syncTimePath);
        console.log(chalk.green(`✓ Deleted sync time: ${formatPath(syncTimePath)}`));
      }
    } catch {
      // Sync time might not exist, that's okay
    }

    console.log('');
    printWarning('╔════════════════════════════════════════════════════╗');
    printWarning('║          Vault permanently destroyed                ║');
    printWarning('╚════════════════════════════════════════════════════╝');
    console.log('');
    printWarning('All secrets have been deleted and cannot be recovered.');
    console.log('');
  } catch (error) {
    if (error instanceof Error && error.message === 'Password entry cancelled') {
      printWarning('Destruction cancelled');
      return;
    }

    printError(error instanceof Error ? error.message : 'Unknown error');
    process.exit(1);
  }
}

/**
 * Check if file exists without throwing
 */
async function checkFileExists(filePath: string): Promise<boolean> {
  try {
    await fs.stat(filePath);
    return true;
  } catch {
    return false;
  }
}

/**
 * Format file path for display
 */
function formatPath(filePath: string): string {
  const homeDir = process.env.HOME || process.env.USERPROFILE || '~';
  return filePath.replace(homeDir, '~');
}
