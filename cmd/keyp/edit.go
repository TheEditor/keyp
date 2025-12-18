package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var editField string

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an existing secret",
	Long:  "Modify fields of an existing secret. Use --field to target a specific field.",
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

func init() {
	editCmd.Flags().StringVar(&editField, "field", "", "Specific field to edit (by label)")
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
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

	// Get secret
	secret, err := v.GetByName(cmd.Context(), name)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	// Determine which field to edit
	if editField == "" {
		// If no field specified, edit all fields interactively
		if len(secret.Fields) == 0 {
			fmt.Println("Secret has no fields")
			return nil
		}

		// Show fields and prompt for each
		for i := range secret.Fields {
			prompt := fmt.Sprintf("Edit '%s' (leave empty to skip): ", secret.Fields[i].Label)
			newValue, err := ui.PromptPassword(prompt)
			if err != nil {
				return err
			}
			if newValue != "" {
				secret.Fields[i].Value = newValue
			}
		}
	} else {
		// Edit specific field
		found := false
		for i, field := range secret.Fields {
			if field.Label == editField {
				newValue, err := ui.PromptPassword(fmt.Sprintf("New value for '%s': ", editField))
				if err != nil {
					return err
				}
				secret.Fields[i].Value = newValue
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("field '%s' not found", editField)
		}
	}

	// Update secret
	if err := v.Update(cmd.Context(), secret); err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	fmt.Printf("Secret '%s' updated\n", name)
	return nil
}
