package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/store"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var (
	listTags []string
	listJSON bool
)

var listCmdObj = &cobra.Command{
	Use:     "list",
	Short:   "List all secrets",
	Long:    "Show all secrets in the vault with optional tag filtering.",
	Aliases: []string{"ls"},
	RunE:    runList,
}

func init() {
	listCmdObj.Flags().StringSliceVar(&listTags, "tags", nil, "Filter by tags (comma-separated)")
	listCmdObj.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(listCmdObj)
}

func runList(cmd *cobra.Command, args []string) error {
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

	// Build SearchOptions for tag filtering
	var opts *store.SearchOptions
	if len(listTags) > 0 {
		opts = &store.SearchOptions{Tags: listTags}
	}

	// List secrets
	secrets, err := v.List(cmd.Context(), opts)
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	// Output
	if listJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(secrets)
	}

	// Table output
	if len(secrets) == 0 {
		fmt.Println("No secrets found")
		return nil
	}

	fmt.Printf("%-30s %-20s %s\n", "NAME", "TAGS", "UPDATED")
	for _, s := range secrets {
		tags := strings.Join(s.Tags, ", ")
		updated := s.UpdatedAt.Format("2006-01-02 15:04")
		fmt.Printf("%-30s %-20s %s\n", s.Name, tags, updated)
	}

	return nil
}
