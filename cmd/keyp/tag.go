package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/ui"
	"github.com/TheEditor/keyp/internal/vault"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage secret tags",
	Long:  "Add, remove, or list tags on secrets.",
}

var tagAddCmd = &cobra.Command{
	Use:   "add <secret> <tag>",
	Short: "Add a tag to a secret",
	Args:  cobra.ExactArgs(2),
	RunE:  runTagAdd,
}

var tagRmCmd = &cobra.Command{
	Use:   "rm <secret> <tag>",
	Short: "Remove a tag from a secret",
	Args:  cobra.ExactArgs(2),
	RunE:  runTagRm,
}

var tagListCmd = &cobra.Command{
	Use:   "list [secret]",
	Short: "List tags (all or for specific secret)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runTagList,
}

func init() {
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRmCmd)
	tagCmd.AddCommand(tagListCmd)
	rootCmd.AddCommand(tagCmd)
}

func runTagAdd(cmd *cobra.Command, args []string) error {
	secretName := args[0]
	tag := args[1]

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
	secret, err := v.GetByName(cmd.Context(), secretName)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	// Check if tag already exists
	for _, t := range secret.Tags {
		if t == tag {
			fmt.Printf("Tag '%s' already exists on secret '%s'\n", tag, secretName)
			return nil
		}
	}

	// Add tag
	secret.Tags = append(secret.Tags, tag)

	// Update secret
	if err := v.Update(cmd.Context(), secret); err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	fmt.Printf("Tag '%s' added to secret '%s'\n", tag, secretName)
	return nil
}

func runTagRm(cmd *cobra.Command, args []string) error {
	secretName := args[0]
	tag := args[1]

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
	secret, err := v.GetByName(cmd.Context(), secretName)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	// Find and remove tag
	found := false
	newTags := []string{}
	for _, t := range secret.Tags {
		if t != tag {
			newTags = append(newTags, t)
		} else {
			found = true
		}
	}

	if !found {
		fmt.Printf("Tag '%s' not found on secret '%s'\n", tag, secretName)
		return nil
	}

	secret.Tags = newTags

	// Update secret
	if err := v.Update(cmd.Context(), secret); err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	fmt.Printf("Tag '%s' removed from secret '%s'\n", tag, secretName)
	return nil
}

func runTagList(cmd *cobra.Command, args []string) error {
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

	if len(args) == 0 {
		// List all tags across all secrets
		secrets, err := v.List(cmd.Context(), nil)
		if err != nil {
			return fmt.Errorf("failed to list secrets: %w", err)
		}

		tagSet := make(map[string]bool)
		for _, s := range secrets {
			for _, tag := range s.Tags {
				tagSet[tag] = true
			}
		}

		if len(tagSet) == 0 {
			fmt.Println("No tags found")
			return nil
		}

		fmt.Println("All tags:")
		for tag := range tagSet {
			fmt.Printf("  %s\n", tag)
		}
	} else {
		// List tags for specific secret
		secretName := args[0]
		secret, err := v.GetByName(cmd.Context(), secretName)
		if err != nil {
			return fmt.Errorf("failed to get secret: %w", err)
		}

		if len(secret.Tags) == 0 {
			fmt.Printf("Secret '%s' has no tags\n", secretName)
			return nil
		}

		fmt.Printf("Tags for secret '%s':\n", secretName)
		for _, tag := range secret.Tags {
			fmt.Printf("  %s\n", tag)
		}
	}

	return nil
}
