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
    .action(async (name: string, options: { stdout?: boolean; noClear?: boolean }) => {
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
   * Help command with banner
   */
  program.on('--help', () => {
    console.log('');
    console.log(chalk.gray('Examples:'));
    console.log(chalk.gray('  $ keyp init                    Initialize vault'));
    console.log(chalk.gray('  $ keyp set github-token        Store a secret (prompts for value)'));
    console.log(chalk.gray('  $ keyp set api-key sk-123      Store a secret with value'));
    console.log(chalk.gray('  $ keyp list                    List all secrets'));
    console.log(chalk.gray('  $ keyp get github-token        Get secret (copies to clipboard)'));
    console.log(chalk.gray('  $ keyp delete github-token     Delete a secret'));
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
