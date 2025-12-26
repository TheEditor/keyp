package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/color"
	"github.com/TheEditor/keyp/internal/ui"
)

var editField string
var editNotes string

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an existing secret",
	Long:  "Modify fields of an existing secret. Use --field to target a specific field.",
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

func init() {
	editCmd.Flags().StringVar(&editField, "field", "", "Specific field to edit (by label)")
	editCmd.Flags().StringVar(&editNotes, "notes", "", "Update notes for the secret")
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Get or unlock vault
	handle, err := getOrUnlockVault(cmd, 0)
	if err != nil {
		return err
	}

	// Get secret
	secret, err := handle.Store().GetByName(cmd.Context(), name)
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

	// Update notes if provided
	if editNotes != "" {
		secret.Notes = editNotes
	}

	// Update secret
	if err := handle.Store().Update(cmd.Context(), secret); err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	msg := fmt.Sprintf("Secret '%s' updated", name)
	fmt.Println(color.Success(msg))
	return nil
}
