package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/store"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var (
	getStdout bool
	getField  string
)

var getCmdObj = &cobra.Command{
	Use:   "get <name>",
	Short: "Get a secret value",
	Long:  "Retrieve a secret and copy to clipboard (or print with --stdout).",
	Args:  cobra.ExactArgs(1),
	RunE:  runGet,
}

func init() {
	getCmdObj.Flags().BoolVar(&getStdout, "stdout", false, "Print to stdout instead of clipboard")
	getCmdObj.Flags().StringVar(&getField, "field", "", "Specific field to retrieve (default: first field)")
	rootCmd.AddCommand(getCmdObj)
}

func runGet(cmd *cobra.Command, args []string) error {
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
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("secret '%s' not found", name)
		}
		return fmt.Errorf("failed to get secret: %w", err)
	}

	// Find field
	var value string
	if getField != "" {
		for _, f := range secret.Fields {
			if f.Label == getField {
				value = f.Value
				break
			}
		}
		if value == "" {
			return fmt.Errorf("field '%s' not found", getField)
		}
	} else if len(secret.Fields) > 0 {
		value = secret.Fields[0].Value
	} else {
		return fmt.Errorf("secret has no fields")
	}

	// Output
	if getStdout {
		fmt.Println(value)
	} else {
		if err := ui.CopyWithAutoClear(value, ui.DefaultClearDuration); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		fmt.Printf("Copied to clipboard (clears in 45s)\n")
	}

	return nil
}
