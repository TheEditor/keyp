import simpleGit, { SimpleGit } from 'simple-git';
import { promises as fs } from 'fs';
import path from 'path';
import { getKeypDir, getVaultPath } from './config.js';

export interface GitSyncConfig {
  remoteUrl: string;
  autoPush: boolean;
  autoCommit: boolean;
  branch: string;
  sshKey?: string;
}

export interface SyncStatus {
  synced: boolean;
  lastSync?: Date;
  uncommittedChanges: boolean;
  unpushedCommits: number;
  conflictCount: number;
  localOnly: string[];
  remoteOnly: string[];
}

export interface ConflictResolution {
  strategy: 'keep-local' | 'keep-remote' | 'prompt';
  autoResolve: boolean;
}

export class GitSyncManager {
  private git: SimpleGit;
  private keypDir: string;
  private gitDir: string;
  private config: GitSyncConfig | null = null;

  constructor() {
    this.keypDir = getKeypDir();
    this.gitDir = path.join(this.keypDir, '.git');
    this.git = simpleGit(this.keypDir);
  }

  /**
   * Initialize Git repository in keyp directory
   */
  async initGitRepo(): Promise<{ success: boolean; error?: string }> {
    try {
      // Check if already initialized
      const isRepo = await this.git.checkIsRepo();
      if (isRepo) {
        return { success: true };
      }

      // Initialize new repo
      await this.git.init();

      // Create initial .gitignore
      const gitignorePath = path.join(this.keypDir, '.gitignore');
      await fs.writeFile(gitignorePath, 'config.json\n.DS_Store\n*.swp\n*.swo\n*~\n');

      // Create initial commit
      await this.git.add('.gitignore');
      await this.git.commit('Initial commit: Create vault repository');

      return { success: true };
    } catch (error) {
      const err = error as Error;
      return { success: false, error: err.message };
    }
  }

  /**
   * Configure Git remote
   */
  async configureRemote(
    remoteUrl: string,
    autoPush: boolean = false,
    autoCommit: boolean = false,
    sshKey?: string
  ): Promise<{ success: boolean; error?: string }> {
    try {
      // Initialize if needed
      const initResult = await this.initGitRepo();
      if (!initResult.success) {
        return initResult;
      }

      // Check if remote exists
      const remotes = await this.git.getRemotes();
      const existingRemote = remotes.find((r) => r.name === 'origin');

      if (existingRemote) {
        // Update existing remote
        await this.git.removeRemote('origin');
      }

      // Add new remote
      await this.git.addRemote('origin', remoteUrl);

      // Store configuration
      this.config = {
        remoteUrl,
        autoPush,
        autoCommit,
        branch: 'main',
        sshKey,
      };

      // Save config to file
      await this.saveConfig();

      return { success: true };
    } catch (error) {
      const err = error as Error;
      return { success: false, error: err.message };
    }
  }

  /**
   * Commit vault changes
   */
  async commitVault(message?: string): Promise<{ success: boolean; error?: string }> {
    try {
      const isRepo = await this.git.checkIsRepo();
      if (!isRepo) {
        return { success: false, error: 'Git repository not initialized' };
      }

      // Add vault file
      const vaultPath = getVaultPath();
      const relativeVaultPath = path.relative(this.keypDir, vaultPath);
      await this.git.add(relativeVaultPath);

      // Check if there are changes to commit
      const status = await this.git.status();
      if (!status.staged || status.staged.length === 0) {
        return { success: true }; // No changes
      }

      // Create commit
      const commitMessage = message || `Update vault: ${new Date().toISOString()}`;
      await this.git.commit(commitMessage);

      return { success: true };
    } catch (error) {
      const err = error as Error;
      return { success: false, error: err.message };
    }
  }

  /**
   * Push vault to remote
   */
  async pushToRemote(branch: string = 'main'): Promise<{ success: boolean; error?: string }> {
    try {
      const isRepo = await this.git.checkIsRepo();
      if (!isRepo) {
        return { success: false, error: 'Git repository not initialized' };
      }

      const remotes = await this.git.getRemotes();
      if (!remotes.some((r) => r.name === 'origin')) {
        return { success: false, error: 'Remote "origin" not configured' };
      }

      // Ensure we're on the right branch
      const currentBranch = await this.git.revparse(['--abbrev-ref', 'HEAD']);
      if (currentBranch.trim() !== branch) {
        await this.git.checkout(branch);
      }

      // Push to remote
      await this.git.push('origin', branch);

      // Update last sync time
      await this.updateLastSync();

      return { success: true };
    } catch (error) {
      const err = error as Error;
      return { success: false, error: err.message };
    }
  }

  /**
   * Pull vault from remote
   */
  async pullFromRemote(
    branch: string = 'main',
    conflictResolution: ConflictResolution = { strategy: 'prompt', autoResolve: false }
  ): Promise<{ success: boolean; error?: string; conflicts?: string[] }> {
    try {
      const isRepo = await this.git.checkIsRepo();
      if (!isRepo) {
        return { success: false, error: 'Git repository not initialized' };
      }

      const remotes = await this.git.getRemotes();
      if (!remotes.some((r) => r.name === 'origin')) {
        return { success: false, error: 'Remote "origin" not configured' };
      }

      // Fetch from remote
      await this.git.fetch('origin', branch);

      // Try to merge
      try {
        await this.git.merge(['origin/' + branch]);
        await this.updateLastSync();
        return { success: true };
      } catch (mergeError) {
        // Handle merge conflicts
        const conflicts = await this.detectConflicts();

        if (conflicts.length > 0) {
          if (conflictResolution.autoResolve) {
            const strategy = conflictResolution.strategy as 'keep-local' | 'keep-remote';
            const resolveResult = await this.resolveConflicts(
              conflicts,
              strategy
            );
            if (!resolveResult.success) {
              return { success: false, error: resolveResult.error, conflicts };
            }
          } else {
            return { success: false, error: 'Merge conflicts detected', conflicts };
          }
        }

        await this.updateLastSync();
        return { success: true, conflicts };
      }
    } catch (error) {
      const err = error as Error;
      return { success: false, error: err.message };
    }
  }

  /**
   * Get current sync status
   */
  async getStatus(): Promise<SyncStatus> {
    try {
      const isRepo = await this.git.checkIsRepo();
      if (!isRepo) {
        return {
          synced: false,
          uncommittedChanges: false,
          unpushedCommits: 0,
          conflictCount: 0,
          localOnly: [],
          remoteOnly: [],
        };
      }

      const status = await this.git.status();
      const conflicts = await this.detectConflicts();

      // Count unpushed commits
      let unpushedCommits = 0;
      try {
        const log = await this.git.log(['origin/main..HEAD']);
        unpushedCommits = log.total;
      } catch {
        unpushedCommits = 0;
      }

      const lastSync = await this.getLastSyncTime();

      return {
        synced: status.behind === 0 && status.ahead === 0 && conflicts.length === 0,
        lastSync,
        uncommittedChanges: status.files.length > 0 || (status.staged && status.staged.length > 0),
        unpushedCommits,
        conflictCount: conflicts.length,
        localOnly: status.created || [],
        remoteOnly: [],
      };
    } catch {
      return {
        synced: false,
        uncommittedChanges: false,
        unpushedCommits: 0,
        conflictCount: 0,
        localOnly: [],
        remoteOnly: [],
      };
    }
  }

  /**
   * Detect merge conflicts
   */
  private async detectConflicts(): Promise<string[]> {
    try {
      const status = await this.git.status();
      return status.conflicted || [];
    } catch {
      return [];
    }
  }

  /**
   * Resolve merge conflicts
   */
  async resolveConflicts(
    conflicts: string[],
    strategy: 'keep-local' | 'keep-remote'
  ): Promise<{ success: boolean; error?: string }> {
    try {
      for (const file of conflicts) {
        if (strategy === 'keep-local') {
          await this.git.checkout([file, '--ours']);
        } else if (strategy === 'keep-remote') {
          await this.git.checkout([file, '--theirs']);
        }
      }

      // Stage resolved files
      await this.git.add(conflicts);

      // Complete merge
      await this.git.commit('Resolve merge conflicts');

      return { success: true };
    } catch (error) {
      const err = error as Error;
      return { success: false, error: err.message };
    }
  }

  /**
   * Load configuration from file
   */
  async loadConfig(): Promise<GitSyncConfig | null> {
    try {
      const configPath = path.join(this.keypDir, '.keyp-git-config.json');
      const content = await fs.readFile(configPath, 'utf-8');
      this.config = JSON.parse(content);
      return this.config;
    } catch {
      return null;
    }
  }

  /**
   * Save configuration to file
   */
  private async saveConfig(): Promise<void> {
    if (!this.config) {
      return;
    }
    const configPath = path.join(this.keypDir, '.keyp-git-config.json');
    await fs.writeFile(configPath, JSON.stringify(this.config, null, 2));
  }

  /**
   * Get last sync time
   */
  private async getLastSyncTime(): Promise<Date | undefined> {
    try {
      const configPath = path.join(this.keypDir, '.keyp-sync-time');
      const content = await fs.readFile(configPath, 'utf-8');
      return new Date(content);
    } catch {
      return undefined;
    }
  }

  /**
   * Update last sync time
   */
  private async updateLastSync(): Promise<void> {
    const syncTimePath = path.join(this.keypDir, '.keyp-sync-time');
    await fs.writeFile(syncTimePath, new Date().toISOString());
  }

  /**
   * Check if Git sync is configured
   */
  async isConfigured(): Promise<boolean> {
    const config = await this.loadConfig();
    return config !== null;
  }

  /**
   * Abort ongoing merge
   */
  async abortMerge(): Promise<{ success: boolean; error?: string }> {
    try {
      await this.git.merge(['--abort']);
      return { success: true };
    } catch (error) {
      const err = error as Error;
      return { success: false, error: err.message };
    }
  }
}
