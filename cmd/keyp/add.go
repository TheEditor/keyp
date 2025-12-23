package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/model"
	"github.com/TheEditor/keyp/internal/ui"
)

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new secret with multiple fields",
	Long:  "Create a new secret with interactive prompts for multiple fields.",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Get or unlock vault
	handle, err := getOrUnlockVault(cmd, 0)
	if err != nil {
		return err
	}

	// Create new secret
	secret := model.NewSecretObject(name)

	// Prompt for fields
	fmt.Println("Enter fields (empty label to finish):")
	fields, err := ui.PromptLoop()
	if err != nil {
		return err
	}

	// Add fields to secret
	for label, value := range fields {
		field := model.NewField(label, value)
		field.Sensitive = true
		secret.AddField(field)
	}

	// Ensure at least one field
	if len(secret.Fields) == 0 {
		fmt.Println("Secret must have at least one field")
		return nil
	}

	// Create secret
	if err := handle.Store().Create(cmd.Context(), secret); err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	fmt.Printf("Secret '%s' created with %d field(s)\n", name, len(secret.Fields))
	return nil
}
