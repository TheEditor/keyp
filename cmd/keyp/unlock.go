package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var unlockTimeout int

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock vault for faster access",
	Long:  "Unlock the vault and keep it open for a period of time, avoiding password re-entry.",
	RunE:  runUnlock,
}

func init() {
	unlockCmd.Flags().IntVar(&unlockTimeout, "timeout", 30, "Minutes before auto-lock (default: 30)")
	rootCmd.AddCommand(unlockCmd)
}

func runUnlock(cmd *cobra.Command, args []string) error {
	vaultPath := getVaultPath()

	// Create handle
	handle := vault.NewHandle(vaultPath)

	// Prompt for password
	password, err := ui.PromptPassword("Enter vault password: ")
	if err != nil {
		return err
	}

	// Unlock
	timeout := time.Duration(unlockTimeout) * time.Minute
	if err := handle.Unlock(password, timeout); err != nil {
		return fmt.Errorf("failed to unlock vault: %w", err)
	}

	// Store handle in global state (in a real app, would use config/state file)
	// For now, just confirm to user
	fmt.Printf("Vault unlocked for %d minutes\n", unlockTimeout)
	fmt.Printf("Time until auto-lock: %v\n", handle.TimeUntilExpire())

	// Keep handle alive for the timeout duration
	// In a real implementation, this would be stored in a state file or config
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-cmd.Context().Done():
			return cmd.Context().Err()
		case <-ticker.C:
			if handle.IsExpired() {
				handle.Lock()
				fmt.Println("Vault auto-locked due to timeout")
				return nil
			}
		}
	}
}
