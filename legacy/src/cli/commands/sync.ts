import { Command } from 'commander';
import { GitSyncManager } from '../../git-sync.js';
import {
  promptString,
  confirm,
  printSuccess,
  printError,
  printWarning,
  printInfo,
  printHint,
  unlockVaultWithRetry,
  printBanner,
} from '../utils.js';
import { VaultManager } from '../../vault-manager.js';
import chalk from 'chalk';

export function createSyncCommand(): Command {
  const cmd = new Command('sync');
  cmd.description('Synchronize vault with Git remote');

  // keyp sync init <remote-url>
  cmd
    .command('init <remoteUrl>')
    .description('Initialize Git sync with remote repository')
    .option('-a, --auto-push', 'Enable automatic push on vault changes')
    .option('-c, --auto-commit', 'Enable automatic commit on vault changes')
    .action(async (remoteUrl: string, options: { autoPush?: boolean; autoCommit?: boolean }) => {
      try {
        const syncManager = new GitSyncManager();

        printInfo('Initializing Git sync...');

        // Initialize git repo
        const initResult = await syncManager.initGitRepo();
        if (!initResult.success) {
          printError(`Failed to initialize Git repository: ${initResult.error}`);
          process.exit(1);
        }

        printSuccess('✓ Git repository initialized');

        // Configure remote
        const configResult = await syncManager.configureRemote(
          remoteUrl,
          options.autoPush || false,
          options.autoCommit || false
        );

        if (!configResult.success) {
          printError(`Failed to configure remote: ${configResult.error}`);
          printHint('Ensure the remote URL is correct and accessible');
          process.exit(1);
        }

        printSuccess(`✓ Remote configured: ${remoteUrl}`);

        if (options.autoPush) {
          printInfo('Auto-push enabled');
        }
        if (options.autoCommit) {
          printInfo('Auto-commit enabled');
        }

        printSuccess('Git sync initialized successfully!');
      } catch (error) {
        const err = error as Error;
        printError(`Error: ${err.message}`);
        process.exit(1);
      }
    });

  // keyp sync push
  cmd
    .command('push')
    .description('Push encrypted vault to remote')
    .option('-m, --message <message>', 'Custom commit message')
    .action(async (options: { message?: string }) => {
      try {
        const syncManager = new GitSyncManager();

        // Check if sync is configured
        const isConfigured = await syncManager.isConfigured();
        if (!isConfigured) {
          printError('Git sync not configured');
          printHint('Run "keyp sync init <remote-url>" to configure');
          process.exit(1);
        }

        printInfo('Pushing vault to remote...');

        // Commit changes
        const commitResult = await syncManager.commitVault(options.message);
        if (!commitResult.success) {
          printError(`Failed to commit: ${commitResult.error}`);
          process.exit(1);
        }

        // Push to remote
        const pushResult = await syncManager.pushToRemote();
        if (!pushResult.success) {
          printError(`Failed to push: ${pushResult.error}`);
          printHint('Check your internet connection and remote URL');
          process.exit(1);
        }

        printSuccess('✓ Vault pushed to remote');
      } catch (error) {
        const err = error as Error;
        printError(`Error: ${err.message}`);
        process.exit(1);
      }
    });

  // keyp sync pull
  cmd
    .command('pull')
    .description('Pull vault from remote')
    .option('-s, --strategy <strategy>', 'Conflict resolution strategy (keep-local, keep-remote)', 'keep-local')
    .option('--auto-resolve', 'Automatically resolve conflicts')
    .action(async (options: { strategy?: string; autoResolve?: boolean }) => {
      try {
        const syncManager = new GitSyncManager();

        // Check if sync is configured
        const isConfigured = await syncManager.isConfigured();
        if (!isConfigured) {
          printError('Git sync not configured');
          printHint('Run "keyp sync init <remote-url>" to configure');
          process.exit(1);
        }

        // Validate strategy
        const strategy = options.strategy as 'keep-local' | 'keep-remote';
        if (!['keep-local', 'keep-remote'].includes(strategy)) {
          printError('Invalid conflict resolution strategy');
          printHint('Use --strategy keep-local or --strategy keep-remote');
          process.exit(1);
        }

        printInfo('Pulling vault from remote...');

        // Pull from remote
        const pullResult = await syncManager.pullFromRemote(
          'main',
          {
            strategy,
            autoResolve: options.autoResolve || false,
          }
        );

        if (!pullResult.success) {
          printError(`Failed to pull: ${pullResult.error}`);

          if (pullResult.conflicts && pullResult.conflicts.length > 0) {
            printWarning(`Conflicts detected in: ${pullResult.conflicts.join(', ')}`);
            printHint('Run "keyp sync pull --auto-resolve" to automatically resolve');

            // Offer manual resolution
            const shouldAbort = await confirm('Abort pull?');
            if (shouldAbort) {
              await syncManager.abortMerge();
              printInfo('Pull aborted');
            }
          }

          process.exit(1);
        }

        if (pullResult.conflicts && pullResult.conflicts.length > 0) {
          printWarning(`Resolved ${pullResult.conflicts.length} conflict(s) using ${strategy} strategy`);
        }

        printSuccess('✓ Vault pulled from remote');
      } catch (error) {
        const err = error as Error;
        printError(`Error: ${err.message}`);
        process.exit(1);
      }
    });

  // keyp sync status
  cmd
    .command('status')
    .description('Show Git sync status')
    .action(async () => {
      try {
        const syncManager = new GitSyncManager();

        // Check if sync is configured
        const isConfigured = await syncManager.isConfigured();
        if (!isConfigured) {
          printWarning('Git sync not configured');
          printHint('Run "keyp sync init <remote-url>" to configure');
          return;
        }

        console.log(chalk.cyan.bold('\nGit Sync Status'));
        console.log(chalk.gray('─'.repeat(60)));

        const status = await syncManager.getStatus();

        // Status indicator
        const syncIndicator = status.synced ? chalk.green('✓ Synced') : chalk.yellow('⚠ Out of sync');
        console.log(`Status: ${syncIndicator}`);

        if (status.lastSync) {
          const timeAgo = getTimeAgo(status.lastSync);
          console.log(`Last sync: ${timeAgo}`);
        } else {
          console.log('Last sync: Never');
        }

        // Changes
        if (status.uncommittedChanges) {
          console.log(chalk.yellow(`Uncommitted changes: Yes`));
        } else {
          console.log(chalk.green(`Uncommitted changes: No`));
        }

        // Unpushed commits
        if (status.unpushedCommits > 0) {
          console.log(chalk.yellow(`Unpushed commits: ${status.unpushedCommits}`));
        } else {
          console.log(chalk.green(`Unpushed commits: 0`));
        }

        // Conflicts
        if (status.conflictCount > 0) {
          console.log(chalk.red(`Conflicts: ${status.conflictCount}`));
        } else {
          console.log(chalk.green(`Conflicts: 0`));
        }

        // Additional info
        if (status.localOnly.length > 0) {
          console.log(`\nLocal only: ${status.localOnly.join(', ')}`);
        }

        if (!status.synced) {
          printHint('Run "keyp sync push" to push changes or "keyp sync pull" to pull updates');
        }
      } catch (error) {
        const err = error as Error;
        printError(`Error: ${err.message}`);
        process.exit(1);
      }
    });

  // keyp sync config
  cmd
    .command('config')
    .description('Configure Git sync settings')
    .option('--auto-push <enabled>', 'Enable/disable auto-push (true/false)')
    .option('--auto-commit <enabled>', 'Enable/disable auto-commit (true/false)')
    .action(async (options: { autoPush?: string; autoCommit?: string }) => {
      try {
        const syncManager = new GitSyncManager();

        // Load current config
        const config = await syncManager.loadConfig();
        if (!config) {
          printError('Git sync not configured');
          printHint('Run "keyp sync init <remote-url>" to configure');
          process.exit(1);
        }

        // Update config if options provided
        if (options.autoPush !== undefined) {
          config.autoPush = options.autoPush.toLowerCase() === 'true';
          printInfo(`Auto-push: ${config.autoPush ? 'enabled' : 'disabled'}`);
        }

        if (options.autoCommit !== undefined) {
          config.autoCommit = options.autoCommit.toLowerCase() === 'true';
          printInfo(`Auto-commit: ${config.autoCommit ? 'enabled' : 'disabled'}`);
        }

        // Display current config
        console.log(chalk.cyan.bold('\nGit Sync Configuration'));
        console.log(chalk.gray('─'.repeat(60)));
        console.log(`Remote URL: ${config.remoteUrl}`);
        console.log(`Branch: ${config.branch}`);
        console.log(`Auto-push: ${config.autoPush ? chalk.green('enabled') : chalk.gray('disabled')}`);
        console.log(`Auto-commit: ${config.autoCommit ? chalk.green('enabled') : chalk.gray('disabled')}`);
      } catch (error) {
        const err = error as Error;
        printError(`Error: ${err.message}`);
        process.exit(1);
      }
    });

  return cmd;
}

/**
 * Format time difference from now
 */
function getTimeAgo(date: Date): string {
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSecs = Math.floor(diffMs / 1000);
  const diffMins = Math.floor(diffSecs / 60);
  const diffHours = Math.floor(diffMins / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffSecs < 60) return `${diffSecs}s ago`;
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  return `${diffDays}d ago`;
}
