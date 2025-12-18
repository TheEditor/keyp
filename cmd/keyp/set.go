package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/model"
	"github.com/TheEditor/keyp/internal/store"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var setStdin bool

var setCmdObj = &cobra.Command{
	Use:   "set <name> [value]",
	Short: "Set a secret value",
	Long:  "Create or update a secret. Value can be provided as argument or via stdin.",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runSet,
}

func init() {
	setCmdObj.Flags().BoolVar(&setStdin, "stdin", false, "Read value from stdin")
	rootCmd.AddCommand(setCmdObj)
}

func runSet(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Get value
	var value string
	if setStdin {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		value = strings.TrimSpace(string(bytes))
	} else if len(args) > 1 {
		value = args[1]
	} else {
		// Prompt for value
		var err error
		value, err = ui.PromptPassword("Enter value: ")
		if err != nil {
			return err
		}
	}

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

	// Create secret with single field
	secret := model.NewSecretObject(name)
	field := model.NewField("value", value)
	field.Sensitive = true
	secret.AddField(field)

	// Try create, if exists then update
	if err := v.Create(cmd.Context(), secret); err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			// Update existing
			existing, err := v.GetByName(cmd.Context(), name)
			if err != nil {
				return fmt.Errorf("failed to get existing secret: %w", err)
			}
			existing.Fields = secret.Fields
			if err := v.Update(cmd.Context(), existing); err != nil {
				return fmt.Errorf("failed to update secret: %w", err)
			}
		} else {
			return fmt.Errorf("failed to create secret: %w", err)
		}
	}

	fmt.Printf("Secret '%s' saved\n", name)
	return nil
}
