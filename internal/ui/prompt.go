package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// PromptPassword prompts the user for a password with hidden input
func PromptPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		pass, err := term.ReadPassword(fd)
		fmt.Println() // newline after hidden input
		if err != nil {
			return "", err
		}
		return string(pass), nil
	}

	// Fallback for non-terminal (pipe, file)
	reader := bufio.NewReader(os.Stdin)
	pass, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(pass), nil
}

// PromptConfirmPassword prompts twice and verifies both inputs match
func PromptConfirmPassword(prompt string, confirmPrompt string) (string, error) {
	pass1, err := PromptPassword(prompt)
	if err != nil {
		return "", err
	}

	pass2, err := PromptPassword(confirmPrompt)
	if err != nil {
		return "", err
	}

	if pass1 != pass2 {
		return "", fmt.Errorf("passwords do not match")
	}

	return pass1, nil
}

// PromptVisible prompts for visible input and returns trimmed string
func PromptVisible(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// PromptLoop interactively prompts for multiple fields until user indicates done
// Returns a map of label -> value pairs
func PromptLoop() (map[string]string, error) {
	fields := make(map[string]string)
	for {
		label, err := PromptVisible("Field label (or empty to finish): ")
		if err != nil {
			return nil, err
		}
		if label == "" {
			break
		}

		// Check for duplicates
		if _, exists := fields[label]; exists {
			fmt.Println("Field already exists. Use different label.")
			continue
		}

		value, err := PromptPassword(fmt.Sprintf("Value for '%s': ", label))
		if err != nil {
			return nil, err
		}

		fields[label] = value
	}
	return fields, nil
}
