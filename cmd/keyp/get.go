package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/color"
	"github.com/TheEditor/keyp/internal/store"
	"github.com/TheEditor/keyp/internal/ui"
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

	// Get or unlock vault
	handle, err := getOrUnlockVault(cmd, 0)
	if err != nil {
		return err
	}

	// Get secret
	secret, err := handle.Store().GetByName(cmd.Context(), name)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("secret '%s' not found: %w", name, err)
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

	// JSON output
	if jsonOutput {
		output := map[string]string{"value": value}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		return enc.Encode(output)
	}

	// Output
	if getStdout {
		fmt.Println(value)
	} else {
		if err := ui.CopyWithAutoClear(value, ui.DefaultClearDuration); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		fmt.Println(color.Success("Copied to clipboard (clears in 45s)"))
	}

	return nil
}
