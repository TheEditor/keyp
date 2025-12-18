package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var version = "2.0.0-dev"

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
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("keyp v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
