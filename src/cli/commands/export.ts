/**
 * keyp export command - Export secrets to file
 */

import chalk from 'chalk';
import { writeFileSync } from 'fs';
import { join, resolve } from 'path';
import { VaultManager } from '../../vault-manager.js';
import { SecretsManager } from '../../secrets.js';
import {
  unlockVaultWithRetry,
  printSuccess,
  printError,
  printInfo,
  printWarning,
  formatVaultPath,
} from '../utils.js';

/**
 * Export secrets to file
 *
 * @param outputFile - Optional output filename
 * @param options - Export options
 */
export async function exportCommand(
  outputFile?: string,
  options?: {
    plain?: boolean;
    stdout?: boolean;
  }
): Promise<void> {
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

    const secretCount = SecretsManager.getSecretCount(data);
    if (secretCount === 0) {
      printWarning('No secrets to export');
      manager.lockVault();
      process.exit(0);
    }

    // Lock vault after reading data
    manager.lockVault();

    // Determine output format
    let exportData: string;
    let exportType: string;

    if (options?.plain) {
      // Plaintext export
      printWarning('Exporting secrets as PLAINTEXT - this is NOT encrypted!');
      exportData = JSON.stringify(data, null, 2);
      exportType = 'plaintext';
    } else {
      // Encrypted export - copy vault structure
      const vaultContent = JSON.parse(
        require('fs').readFileSync(manager.getVaultPath(), 'utf-8')
      );
      exportData = JSON.stringify(vaultContent, null, 2);
      exportType = 'encrypted';
    }

    // Handle stdout output
    if (options?.stdout) {
      console.log('');
      console.log(exportData);
      console.log('');
      return;
    }

    // Determine output filename
    let finalOutputFile: string;
    if (outputFile) {
      finalOutputFile = resolve(outputFile);
    } else {
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5);
      finalOutputFile = resolve(`keyp-export-${timestamp}.json`);
    }

    // Write file
    writeFileSync(finalOutputFile, exportData, 'utf-8');

    // Success!
    console.log('');
    printSuccess(`Secrets exported (${exportType})`);
    printInfo(`Location: ${formatVaultPath(finalOutputFile)}`);
    printInfo(`Secrets exported: ${secretCount}`);

    if (options?.plain) {
      printWarning('Remember: This file contains PLAINTEXT secrets - keep it safe!');
    }

    console.log('');
  } catch (error) {
    if (error instanceof Error && error.message === 'Password entry cancelled') {
      printInfo('Export cancelled');
      manager.lockVault();
      process.exit(0);
    }

    printError(error instanceof Error ? error.message : 'Unknown error');
    manager.lockVault();
    process.exit(1);
  }
}
