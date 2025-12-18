package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search secrets",
	Long:  "Full-text search across secret names, tags, and notes.",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

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

	// Search secrets
	secrets, err := v.Search(cmd.Context(), query, nil)
	if err != nil {
		return fmt.Errorf("failed to search secrets: %w", err)
	}

	// Display results
	if len(secrets) == 0 {
		fmt.Printf("No secrets found matching '%s'\n", query)
		return nil
	}

	fmt.Printf("Found %d secret(s) matching '%s':\n\n", len(secrets), query)
	fmt.Printf("%-30s %-20s %s\n", "NAME", "TAGS", "UPDATED")
	for _, s := range secrets {
		tags := strings.Join(s.Tags, ", ")
		updated := s.UpdatedAt.Format("2006-01-02 15:04")
		fmt.Printf("%-30s %-20s %s\n", s.Name, tags, updated)
	}

	return nil
}
