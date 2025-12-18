package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/vault"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Explicitly lock the vault",
	Long:  "Lock the vault immediately, requiring password re-entry for next access.",
	RunE:  runLock,
}

func init() {
	rootCmd.AddCommand(lockCmd)
}

func runLock(cmd *cobra.Command, args []string) error {
	vaultPath := getVaultPath()

	// Create handle
	handle := vault.NewHandle(vaultPath)

	// Lock
	handle.Lock()

	fmt.Println("Vault locked")
	return nil
}
