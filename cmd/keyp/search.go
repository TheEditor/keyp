package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var searchPorcelain bool

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search secrets",
	Long:  "Full-text search across secret names, tags, and notes.",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().BoolVar(&searchPorcelain, "porcelain", false, "Output tab-separated values (no headers)")
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

	// JSON output
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		return enc.Encode(secrets)
	}

	// Porcelain output (tab-separated, no headers)
	if searchPorcelain {
		for _, s := range secrets {
			tags := strings.Join(s.Tags, ", ")
			updated := s.UpdatedAt.Format("2006-01-02")
			fmt.Printf("%s\t%s\t%s\n", s.Name, tags, updated)
		}
		return nil
	}

	// Display results
	if len(secrets) == 0 {
		fmt.Printf("No secrets match '%s'\n", query)
		return nil
	}

	count := len(secrets)
	secretWord := "secret"
	if count > 1 {
		secretWord = "secrets"
	}
	fmt.Printf("Found %d %s matching '%s':\n\n", count, secretWord, query)
	fmt.Printf("%-30s %-20s %s\n", "NAME", "TAGS", "UPDATED")
	for _, s := range secrets {
		tags := strings.Join(s.Tags, ", ")
		updated := s.UpdatedAt.Format("2006-01-02 15:04")
		fmt.Printf("%-30s %-20s %s\n", s.Name, tags, updated)
	}

	return nil
}
