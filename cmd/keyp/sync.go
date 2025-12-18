package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/sync"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync vault with git",
	Long:  "Initialize, push, pull, or check status of vault git repository.",
}

var syncInitCmd = &cobra.Command{
	Use:   "init [remote-url]",
	Short: "Initialize git sync",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSyncInit,
}

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push vault to remote",
	RunE:  runSyncPush,
}

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull vault from remote",
	RunE:  runSyncPull,
}

var syncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show sync status",
	RunE:  runSyncStatus,
}

func init() {
	syncCmd.AddCommand(syncInitCmd)
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncPullCmd)
	syncCmd.AddCommand(syncStatusCmd)
	rootCmd.AddCommand(syncCmd)
}

func runSyncInit(cmd *cobra.Command, args []string) error {
	vaultPath := getVaultPath()

	// Create syncer
	syncer := sync.NewGitExecSyncer(vaultPath)

	// Initialize git repo
	if err := syncer.Init(); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	fmt.Println("Git repository initialized")

	// If remote URL provided, add remote
	if len(args) > 0 {
		remoteURL := args[0]
		if err := syncer.AddRemote(remoteURL); err != nil {
			return fmt.Errorf("failed to add remote: %w", err)
		}
		fmt.Printf("Remote 'origin' added: %s\n", remoteURL)
	} else {
		fmt.Println("No remote URL provided. Use 'keyp sync init <url>' to add one later.")
	}

	return nil
}

func runSyncPush(cmd *cobra.Command, args []string) error {
	vaultPath := getVaultPath()

	// Create syncer
	syncer := sync.NewGitExecSyncer(vaultPath)

	// Check status first
	status, err := syncer.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if !status.Initialized {
		return fmt.Errorf("git not initialized. Run 'keyp sync init' first")
	}

	if !status.RemoteConfigured {
		return fmt.Errorf("remote not configured. Run 'keyp sync init <url>' first")
	}

	// Try commit first (in case there are uncommitted changes)
	if !status.Clean {
		fmt.Println("Uncommitted changes detected. Creating commit...")
		if err := syncer.Commit("Auto-commit from keyp"); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}
	}

	// Push
	if err := syncer.Push(); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Println("Pushed to remote successfully")
	return nil
}

func runSyncPull(cmd *cobra.Command, args []string) error {
	vaultPath := getVaultPath()

	// Create syncer
	syncer := sync.NewGitExecSyncer(vaultPath)

	// Check status first
	status, err := syncer.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if !status.Initialized {
		return fmt.Errorf("git not initialized. Run 'keyp sync init' first")
	}

	if !status.RemoteConfigured {
		return fmt.Errorf("remote not configured. Run 'keyp sync init <url>' first")
	}

	// Pull
	if err := syncer.Pull(); err != nil {
		return fmt.Errorf("failed to pull: %w", err)
	}

	fmt.Println("Pulled from remote successfully")
	return nil
}

func runSyncStatus(cmd *cobra.Command, args []string) error {
	vaultPath := getVaultPath()

	// Create syncer
	syncer := sync.NewGitExecSyncer(vaultPath)

	// Get status
	status, err := syncer.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	// Display status
	fmt.Println("Vault Sync Status:")
	fmt.Printf("  Initialized: %v\n", status.Initialized)
	fmt.Printf("  Remote Configured: %v\n", status.RemoteConfigured)
	fmt.Printf("  Clean: %v\n", status.Clean)
	fmt.Printf("  Unpushed Commits: %d\n", status.UnpushedCommits)
	fmt.Printf("  Unpulled Commits: %d\n", status.UnpulledCommits)

	return nil
}
