package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/color"
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
	// Clear the global handle and session file
	clearVaultHandle()

	fmt.Println(color.Success("Vault locked"))
	return nil
}
