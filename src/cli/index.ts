#!/usr/bin/env node

/**
 * keyp CLI - Main entry point
 * Local-first secret manager for developers
 */

import { Command } from 'commander';
import chalk from 'chalk';
import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import { initCommand } from './commands/init.js';
import { setCommand } from './commands/set.js';
import { getCommand } from './commands/get.js';
import { listCommand } from './commands/list.js';
import { deleteCommand } from './commands/delete.js';
import { renameCommand } from './commands/rename.js';
import { copyCommand } from './commands/copy.js';
import { exportCommand } from './commands/export.js';
import { importCommand } from './commands/import.js';
import { createSyncCommand } from './commands/sync.js';
import { statsCommand } from './commands/stats.js';
import { configCommand } from './commands/config.js';
import { destroyCommand } from './commands/destroy.js';
import { printBanner } from './utils.js';

// Load package.json in ESM context
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const packageJson = JSON.parse(readFileSync(join(__dirname, '../../package.json'), 'utf-8'));

/**
 * Create and configure the CLI program
 */
function createProgram(): Command {
  const program = new Command();

  program
    .name('keyp')
    .description('🔒 Local-first secret manager for developers')
    .version(packageJson.version)
    .usage('<command> [options]');

  /**
   * keyp init - Initialize vault
   */
  program
    .command('init')
    .description('Initialize a new vault')
    .action(async () => {
      try {
        await initCommand();
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp set - Store a secret
   */
  program
    .command('set <name> [value]')
    .description('Store a secret in the vault')
    .action(async (name: string, value?: string) => {
      try {
        await setCommand(name, value);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp get - Retrieve a secret
   */
  program
    .command('get <name>')
    .description('Retrieve a secret from the vault (copies to clipboard)')
    .option('--stdout', 'Print to stdout instead of clipboard')
    .option('--no-clear', 'Do not auto-clear clipboard')
    .option('--timeout <seconds>', 'Auto-clear timeout in seconds (default 45)')
    .action(async (name: string, options: { stdout?: boolean; noClear?: boolean; timeout?: string }) => {
      try {
        await getCommand(name, options);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp list - List all secrets
   */
  program
    .command('list')
    .description('List all secrets in the vault')
    .option('--search <pattern>', 'Search by pattern')
    .option('--count', 'Show only count')
    .action(async (options: { search?: string; count?: boolean }) => {
      try {
        await listCommand(options);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp delete - Delete a secret (alias: rm)
   */
  program
    .command('delete <name>')
    .alias('rm')
    .description('Delete a secret from the vault')
    .option('-f, --force', 'Skip confirmation')
    .action(async (name: string, options: { force?: boolean }) => {
      try {
        await deleteCommand(name, options);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp rename - Rename a secret
   */
  program
    .command('rename <old-name> <new-name>')
    .description('Rename a secret')
    .action(async (oldName: string, newName: string) => {
      try {
        await renameCommand(oldName, newName);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp copy - Copy a secret
   */
  program
    .command('copy <source> <dest>')
    .description('Copy a secret to a new name')
    .action(async (source: string, dest: string) => {
      try {
        await copyCommand(source, dest);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp export - Export secrets
   */
  program
    .command('export [output-file]')
    .description('Export secrets to file (encrypted by default)')
    .option('--plain', 'Export as plaintext JSON (unencrypted)')
    .option('--stdout', 'Print to stdout instead of file')
    .action(async (outputFile?: string, options?: { plain?: boolean; stdout?: boolean }) => {
      try {
        await exportCommand(outputFile, options);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp import - Import secrets
   */
  program
    .command('import <input-file>')
    .description('Import secrets from file')
    .option('--replace', 'Replace all existing secrets')
    .option('--dry-run', 'Preview without importing')
    .action(async (inputFile: string, options?: { replace?: boolean; dryRun?: boolean }) => {
      try {
        await importCommand(inputFile, options);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp sync - Git synchronization
   */
  program.addCommand(createSyncCommand());

  /**
   * keyp stats - Display vault statistics
   */
  program
    .command('stats')
    .description('Display vault statistics')
    .action(async () => {
      try {
        await statsCommand();
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp config - Manage configuration
   */
  program
    .command('config [action] [key] [value]')
    .description('Manage vault configuration settings')
    .action(async (action?: string, key?: string, value?: string) => {
      try {
        await configCommand(action, key, value);
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * keyp destroy - Permanently delete vault
   */
  program
    .command('destroy')
    .description('Permanently delete vault (cannot be undone)')
    .action(async () => {
      try {
        await destroyCommand();
      } catch (error) {
        console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
        process.exit(1);
      }
    });

  /**
   * Help command with banner
   */
  program.on('--help', () => {
    console.log('');
    console.log(chalk.gray('Examples:'));
    console.log(chalk.gray('  $ keyp init                    Initialize vault'));
    console.log(chalk.gray('  $ keyp set github-token        Store a secret (prompts for value)'));
    console.log(chalk.gray('  $ keyp get github-token        Get secret (copies to clipboard)'));
    console.log(chalk.gray('  $ keyp list                    List all secrets'));
    console.log(chalk.gray('  $ keyp rename old-name new-name Rename a secret'));
    console.log(chalk.gray('  $ keyp copy source dest        Copy a secret'));
    console.log(chalk.gray('  $ keyp delete github-token     Delete a secret'));
    console.log(chalk.gray('  $ keyp export secrets.json     Export secrets'));
    console.log(chalk.gray('  $ keyp import secrets.json     Import secrets'));
    console.log(chalk.gray('  $ keyp sync init <url>         Initialize Git sync'));
    console.log(chalk.gray('  $ keyp sync push               Push to remote'));
    console.log(chalk.gray('  $ keyp sync pull               Pull from remote'));
    console.log(chalk.gray('  $ keyp stats                   Show vault statistics'));
    console.log(chalk.gray('  $ keyp config                  Show configuration'));
    console.log(chalk.gray('  $ keyp config set key value    Set configuration value'));
    console.log(chalk.gray('  $ keyp destroy                 Permanently delete vault'));
    console.log('');
    console.log(chalk.gray('Documentation:'));
    console.log(chalk.gray('  https://github.com/TheEditor/keyp'));
    console.log('');
  });

  return program;
}

/**
 * Main entry point
 */
async function main(): Promise<void> {
  const program = createProgram();

  // Show banner when no command provided
  if (process.argv.length < 3) {
    printBanner();
    program.outputHelp();
    process.exit(0);
  }

  // Parse and execute
  await program.parseAsync(process.argv);
}

// Run main
main().catch((error) => {
  console.error(chalk.red('Fatal error:', error instanceof Error ? error.message : String(error)));
  process.exit(1);
});
