package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/store"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var deleteForce bool

var deleteCmdObj = &cobra.Command{
	Use:     "delete <name>",
	Short:   "Delete a secret",
	Long:    "Remove a secret from the vault. Requires confirmation unless --force is used.",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

func init() {
	deleteCmdObj.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation prompt")
	deleteCmdObj.Flags().BoolVarP(&deleteForce, "f", "f", false, "Skip confirmation prompt (shorthand)")
	rootCmd.AddCommand(deleteCmdObj)
}

func runDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Prompt for vault password
	password, err := ui.PromptPassword("Enter vault password: ")
	if err != nil {
		return err
	}

	// Open vault
	v, err := vault.Open(getVaultPath(), password)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}
	defer v.Close()

	// Verify exists
	_, err = v.GetByName(cmd.Context(), name)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("secret '%s' not found", name)
		}
		return fmt.Errorf("failed to get secret: %w", err)
	}

	// Confirm deletion
	if !deleteForce {
		confirm, err := ui.PromptVisible(fmt.Sprintf("Type '%s' to confirm deletion: ", name))
		if err != nil {
			return err
		}
		if confirm != name {
			return fmt.Errorf("deletion cancelled")
		}
	}

	// Delete
	if err := v.Delete(cmd.Context(), name); err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	fmt.Printf("Secret '%s' deleted\n", name)
	return nil
}
