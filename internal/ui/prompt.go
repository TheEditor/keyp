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
