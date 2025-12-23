package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
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

	// Get or unlock vault
	handle, err := getOrUnlockVault(cmd, 0)
	if err != nil {
		return err
	}

	// Search secrets
	secrets, err := handle.Store().Search(cmd.Context(), query, nil)
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
