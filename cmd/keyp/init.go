package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var initCmdPath string

var initCmdObj = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new keyp vault",
	Long:  "Create a new encrypted vault for storing secrets.",
	RunE:  runInit,
}

func init() {
	initCmdObj.Flags().StringVar(&initCmdPath, "path", "", "Path to vault directory (default: ~/.keyp)")
	rootCmd.AddCommand(initCmdObj)
}

func runInit(cmd *cobra.Command, args []string) error {
	path := initCmdPath
	if path == "" {
		path = vault.DefaultPath()
	}

	// Check if exists
	if vault.Exists(path) {
		return fmt.Errorf("vault already exists at %s", path)
	}

	// Prompt for password with confirmation
	password, err := ui.PromptConfirmPassword(
		"Enter vault password: ",
		"Confirm password: ",
	)
	if err != nil {
		return err
	}

	// Validate length
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Create vault with encryption
	v, err := vault.Init(path, password)
	if err != nil {
		return fmt.Errorf("failed to initialize vault: %w", err)
	}
	v.Close()

	// Auto-unlock vault after successful init since user just proved they know the password
	handle := vault.NewHandle(path)
	if err := handle.Unlock(password, 0); err != nil {
		return fmt.Errorf("failed to unlock vault after init: %w", err)
	}
	setVaultHandle(handle)

	// Save session to avoid prompting for password on subsequent commands
	if derivedKey := handle.GetDerivedKey(); derivedKey != nil {
		sessionMgr.Save(derivedKey)
	}

	fmt.Printf("Vault initialized at %s\n", path)
	return nil
}
