package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/model"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
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
	if err := v.Create(cmd.Context(), secret); err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	fmt.Printf("Secret '%s' created with %d field(s)\n", name, len(secret.Fields))
	return nil
}
