package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/cli"
	"github.com/TheEditor/keyp/internal/store"
)

var version = "2.0.0-dev"

// Global flags
var jsonOutput bool

// getVaultPath returns the vault path from flag or default
func getVaultPath() string {
	// Check init command path flag
	if initCmdPath != "" {
		return initCmdPath
	}
	// Default to ~/.keyp/vault.db
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".keyp", "vault.db")
}

var rootCmd = &cobra.Command{
	Use:   "keyp",
	Short: "Local-first secret manager",
	Long:  `keyp is a local-first secret manager for developers and families.
Securely store structured secrets with full-text search.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip auto-lock check for these commands
		skipAutoLock := map[string]bool{
			"init":    true,
			"version": true,
			"help":    true,
		}

		if skipAutoLock[cmd.Name()] {
			return nil
		}

		// Check if vault handle exists and is expired
		handle := getVaultHandle()
		if handle != nil && handle.IsUnlocked() && handle.IsExpired() {
			handle.Lock()
		}

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("keyp v%s\n", version)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode := getExitCode(err)
		os.Exit(exitCode)
	}
}

// getExitCode determines the exit code based on error type
func getExitCode(err error) int {
	if err == nil {
		return cli.ExitSuccess
	}

	// Check for specific error types
	if errors.Is(err, store.ErrNotFound) {
		return cli.ExitNotFound
	}

	if errors.Is(err, store.ErrInvalidPassword) {
		return cli.ExitAuthFailed
	}

	if errors.Is(err, store.ErrVaultClosed) || errors.Is(err, store.ErrDatabaseLocked) {
		return cli.ExitVaultLocked
	}

	// Check error message for specific conditions
	errMsg := err.Error()
	if strings.Contains(errMsg, "not found") {
		return cli.ExitNotFound
	}

	if strings.Contains(errMsg, "vault is locked") || strings.Contains(errMsg, "database is locked") {
		return cli.ExitVaultLocked
	}

	// Default to generic error
	return cli.ExitError
}
