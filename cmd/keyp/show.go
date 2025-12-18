package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
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

	// Redact sensitive fields if not revealing
	if !showReveal {
		secret = secret.Redacted()
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
