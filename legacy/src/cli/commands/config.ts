/**
 * keyp config command - Manage vault configuration settings
 */

import chalk from 'chalk';
import { promises as fs } from 'fs';
import path from 'path';
import { getKeypDir } from '../../config.js';
import { printSuccess, printError, printInfo, promptString, confirm } from '../utils.js';

export interface KeypConfig {
  clipboardTimeout: number;
  autoLock: number | null;
  gitAutoSync: boolean;
}

const DEFAULT_CONFIG: KeypConfig = {
  clipboardTimeout: 45,
  autoLock: null,
  gitAutoSync: false,
};

/**
 * Get config file path
 */
function getConfigPath(): string {
  return path.join(getKeypDir(), '.keyp-config.json');
}

/**
 * Load configuration from file
 */
async function loadConfig(): Promise<KeypConfig> {
  try {
    const configPath = getConfigPath();
    const content = await fs.readFile(configPath, 'utf-8');
    const loaded = JSON.parse(content);
    // Merge with defaults to handle new config keys
    return { ...DEFAULT_CONFIG, ...loaded };
  } catch {
    return { ...DEFAULT_CONFIG };
  }
}

/**
 * Save configuration to file
 */
async function saveConfig(config: KeypConfig): Promise<void> {
  const configPath = getConfigPath();
  await fs.writeFile(configPath, JSON.stringify(config, null, 2));
}

/**
 * Config command implementation
 */
export async function configCommand(
  action?: string,
  key?: string,
  value?: string
): Promise<void> {
  try {
    const config = await loadConfig();

    if (!action) {
      // Display all config
      displayConfig(config);
      return;
    }

    if (action === 'list') {
      displayConfig(config);
      return;
    }

    if (action === 'set') {
      if (!key || !value) {
        printError('Usage: keyp config set <key> <value>');
        printInfo('Keys: clipboard-timeout, auto-lock, git-auto-sync');
        process.exit(1);
      }

      await setConfigValue(config, key, value);
      return;
    }

    if (action === 'get') {
      if (!key) {
        printError('Usage: keyp config get <key>');
        printInfo('Keys: clipboard-timeout, auto-lock, git-auto-sync');
        process.exit(1);
      }

      getConfigValue(config, key);
      return;
    }

    if (action === 'reset') {
      const shouldReset = await confirm('Reset all settings to defaults?');
      if (shouldReset) {
        await saveConfig({ ...DEFAULT_CONFIG });
        printSuccess('Configuration reset to defaults');
      }
      return;
    }

    printError(`Unknown action: ${action}`);
    printInfo('Usage: keyp config [list|get|set|reset]');
    process.exit(1);
  } catch (error) {
    printError(error instanceof Error ? error.message : 'Unknown error');
    process.exit(1);
  }
}

/**
 * Set a configuration value
 */
async function setConfigValue(config: KeypConfig, key: string, value: string): Promise<void> {
  const normalizedKey = key.replace(/-/g, '');

  switch (normalizedKey.toLowerCase()) {
    case 'clipboardtimeout':
      const timeout = parseInt(value);
      if (isNaN(timeout) || timeout < 0) {
        printError('clipboard-timeout must be a non-negative number');
        process.exit(1);
      }
      config.clipboardTimeout = timeout;
      printSuccess(`Clipboard timeout set to ${timeout} seconds`);
      break;

    case 'autolock':
      if (value.toLowerCase() === 'none' || value.toLowerCase() === 'null') {
        config.autoLock = null;
        printSuccess('Auto-lock disabled');
      } else {
        const seconds = parseInt(value);
        if (isNaN(seconds) || seconds < 30) {
          printError('auto-lock must be "none" or at least 30 seconds');
          process.exit(1);
        }
        config.autoLock = seconds;
        printSuccess(`Auto-lock set to ${seconds} seconds`);
      }
      break;

    case 'gitautosync':
      const enabled = value.toLowerCase() === 'true';
      config.gitAutoSync = enabled;
      printSuccess(`Git auto-sync ${enabled ? 'enabled' : 'disabled'}`);
      break;

    default:
      printError(`Unknown configuration key: ${key}`);
      printInfo('Available keys: clipboard-timeout, auto-lock, git-auto-sync');
      process.exit(1);
  }

  await saveConfig(config);
}

/**
 * Get a configuration value
 */
function getConfigValue(config: KeypConfig, key: string): void {
  const normalizedKey = key.replace(/-/g, '').toLowerCase();

  switch (normalizedKey) {
    case 'clipboardtimeout':
      console.log(`${config.clipboardTimeout}`);
      break;

    case 'autolock':
      console.log(config.autoLock === null ? 'none' : `${config.autoLock}`);
      break;

    case 'gitautosync':
      console.log(config.gitAutoSync ? 'true' : 'false');
      break;

    default:
      printError(`Unknown configuration key: ${key}`);
      process.exit(1);
  }
}

/**
 * Display all configuration settings
 */
function displayConfig(config: KeypConfig): void {
  console.log('');
  console.log(chalk.cyan.bold('⚙️  Keyp Configuration'));
  console.log(chalk.gray('─'.repeat(60)));
  console.log('');

  console.log(`${chalk.bold('Clipboard')}`);
  console.log(`  Timeout: ${chalk.cyan(config.clipboardTimeout.toString())} seconds`);
  console.log(`  (how long before clipboard is auto-cleared)`);

  console.log('');

  console.log(`${chalk.bold('Vault')}`);
  if (config.autoLock === null) {
    console.log(`  Auto-lock: ${chalk.gray('disabled')}`);
  } else {
    console.log(`  Auto-lock: ${chalk.cyan(config.autoLock.toString())} seconds`);
  }
  console.log(`  (lock vault after inactivity)`);

  console.log('');

  console.log(`${chalk.bold('Git Sync')}`);
  console.log(`  Auto-sync: ${config.gitAutoSync ? chalk.green('enabled') : chalk.gray('disabled')}`);
  console.log(`  (automatically push on vault changes)`);

  console.log('');

  console.log(chalk.gray('Configuration stored in: ~/.keyp/.keyp-config.json'));
  console.log('');
}
