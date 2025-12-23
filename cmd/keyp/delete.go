package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/store"
	"github.com/TheEditor/keyp/internal/ui"
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

	// Get or unlock vault
	handle, err := getOrUnlockVault(cmd, 0)
	if err != nil {
		return err
	}

	// Verify exists
	_, err = handle.Store().GetByName(cmd.Context(), name)
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
	if err := handle.Store().Delete(cmd.Context(), name); err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	fmt.Printf("Secret '%s' deleted\n", name)
	return nil
}
