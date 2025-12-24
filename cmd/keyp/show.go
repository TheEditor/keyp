package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var showReveal bool

var showCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show secret details",
	Long:  "Display all fields of a secret. Sensitive values are masked unless --reveal is used.",
	Args:  cobra.ExactArgs(1),
	RunE:  runShow,
}

func init() {
	showCmd.Flags().BoolVar(&showReveal, "reveal", false, "Show sensitive values (default: masked)")
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
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

	// Redact sensitive fields if not revealing
	if !showReveal {
		secret = secret.Redacted()
	}

	// JSON output
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(secret)
	}

	// Display secret details
	fmt.Printf("Name: %s\n", secret.Name)
	fmt.Printf("Tags: %v\n", secret.Tags)
	fmt.Printf("Created: %s\n", secret.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("Updated: %s\n", secret.UpdatedAt.Format("2006-01-02 15:04"))
	if secret.Notes != "" {
		fmt.Printf("Notes: %s\n", secret.Notes)
	}
	fmt.Println("\nFields:")
	for _, field := range secret.Fields {
		fmt.Printf("  %s: %s\n", field.Label, field.Value)
	}

	return nil
}
