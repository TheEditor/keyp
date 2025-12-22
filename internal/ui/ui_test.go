package ui

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/atotto/clipboard"
)

// TestPromptPasswordWithPipe tests PromptPassword with piped input
func TestPromptPasswordWithPipe(t *testing.T) {
	// Create a pipe to simulate stdin
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	// Write password with newline
	go func() {
		w.WriteString("testpass123\n")
		w.Close()
	}()

	// Call PromptPassword (suppresses the printed prompt)
	password, err := PromptPassword("Password: ")
	if err != nil {
		t.Fatalf("PromptPassword failed: %v", err)
	}

	if password != "testpass123" {
		t.Errorf("expected 'testpass123', got '%s'", password)
	}
}

// TestPromptPasswordTrimming tests that PromptPassword trims whitespace
func TestPromptPasswordTrimming(t *testing.T) {
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	go func() {
		w.WriteString("  secret123  \n")
		w.Close()
	}()

	password, err := PromptPassword("Password: ")
	if err != nil {
		t.Fatalf("PromptPassword failed: %v", err)
	}

	if password != "secret123" {
		t.Errorf("expected trimmed password, got '%s'", password)
	}
}

// TestPromptConfirmPasswordMatching tests PromptConfirmPassword with matching passwords
func TestPromptConfirmPasswordMatching(t *testing.T) {
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	go func() {
		w.WriteString("testpass\n")
		w.WriteString("testpass\n")
		w.Close()
	}()

	password, err := PromptConfirmPassword("Password: ", "Confirm: ")
	if err != nil {
		t.Fatalf("PromptConfirmPassword failed: %v", err)
	}

	if password != "testpass" {
		t.Errorf("expected 'testpass', got '%s'", password)
	}
}

// TestPromptConfirmPasswordMismatch tests PromptConfirmPassword with mismatched passwords
func TestPromptConfirmPasswordMismatch(t *testing.T) {
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	go func() {
		w.WriteString("password1\n")
		w.WriteString("password2\n")
		w.Close()
	}()

	password, err := PromptConfirmPassword("Password: ", "Confirm: ")
	if err == nil {
		t.Errorf("expected error for mismatched passwords, got none")
	}

	if password != "" {
		t.Errorf("expected empty password on mismatch, got '%s'", password)
	}
}

// TestPromptVisibleInput tests PromptVisible with piped input
func TestPromptVisibleInput(t *testing.T) {
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	go func() {
		w.WriteString("user input\n")
		w.Close()
	}()

	input, err := PromptVisible("Enter: ")
	if err != nil {
		t.Fatalf("PromptVisible failed: %v", err)
	}

	if input != "user input" {
		t.Errorf("expected 'user input', got '%s'", input)
	}
}

// TestPromptVisibleTrimming tests that PromptVisible trims whitespace
func TestPromptVisibleTrimming(t *testing.T) {
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	go func() {
		w.WriteString("  some input  \n")
		w.Close()
	}()

	input, err := PromptVisible("Enter: ")
	if err != nil {
		t.Fatalf("PromptVisible failed: %v", err)
	}

	if input != "some input" {
		t.Errorf("expected trimmed input, got '%s'", input)
	}
}

// TestCopyToClipboard tests CopyToClipboard sets clipboard content
func TestCopyToClipboard(t *testing.T) {
	text := "secret123"

	err := CopyToClipboard(text)
	if err != nil {
		t.Fatalf("CopyToClipboard failed: %v", err)
	}

	// Verify clipboard contains the text
	content, err := clipboard.ReadAll()
	if err != nil {
		t.Logf("warning: clipboard read failed (may not have clipboard in test environment): %v", err)
		return
	}

	if content != text {
		t.Errorf("expected clipboard to contain '%s', got '%s'", text, content)
	}
}

// TestCopyWithAutoClearClears tests that CopyWithAutoClear clears after duration
func TestCopyWithAutoClearClears(t *testing.T) {
	text := "temporary123"
	clearDuration := 50 * time.Millisecond

	err := CopyWithAutoClear(text, clearDuration)
	if err != nil {
		t.Logf("warning: CopyWithAutoClear failed: %v", err)
		return
	}

	// Verify immediately after copy
	content, err := clipboard.ReadAll()
	if err != nil {
		t.Logf("warning: clipboard read failed (may not have clipboard in test environment): %v", err)
		return
	}

	if content != text {
		t.Errorf("expected clipboard to contain text immediately after copy, got '%s'", content)
	}

	// Wait for auto-clear
	time.Sleep(clearDuration + 10*time.Millisecond)

	// Verify clipboard is cleared
	content, err = clipboard.ReadAll()
	if err != nil {
		t.Logf("warning: clipboard read after clear failed: %v", err)
		return
	}

	if content != "" {
		t.Errorf("expected clipboard to be cleared, but contains '%s'", content)
	}
}

// TestCopyWithAutoClearZeroDuration tests CopyWithAutoClear with zero duration (no clear)
func TestCopyWithAutoClearZeroDuration(t *testing.T) {
	text := "persistent123"

	err := CopyWithAutoClear(text, 0)
	if err != nil {
		t.Logf("warning: CopyWithAutoClear failed: %v", err)
		return
	}

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Verify clipboard still has content (not cleared with zero duration)
	content, err := clipboard.ReadAll()
	if err != nil {
		t.Logf("warning: clipboard read failed: %v", err)
		return
	}

	if content != text {
		t.Errorf("expected clipboard to contain '%s' with zero duration, got '%s'", text, content)
	}
}

// TestPromptVisibleEOF tests PromptVisible with EOF
func TestPromptVisibleEOF(t *testing.T) {
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	go func() {
		w.Close()
	}()

	_, err := PromptVisible("Enter: ")
	if err == nil {
		t.Errorf("expected error on EOF, got none")
	}
}

// TestPromptPasswordEOF tests PromptPassword with EOF
func TestPromptPasswordEOF(t *testing.T) {
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	go func() {
		w.Close()
	}()

	_, err := PromptPassword("Password: ")
	if err == nil {
		t.Errorf("expected error on EOF, got none")
	}
}

// TestCopyWithAutoClearNegativeDuration tests CopyWithAutoClear with negative duration
func TestCopyWithAutoClearNegativeDuration(t *testing.T) {
	text := "negativeduration"

	err := CopyWithAutoClear(text, -1*time.Second)
	if err != nil {
		t.Logf("warning: CopyWithAutoClear failed: %v", err)
		return
	}

	// Should not clear with negative duration (treated as < 0)
	time.Sleep(10 * time.Millisecond)

	content, err := clipboard.ReadAll()
	if err != nil {
		t.Logf("warning: clipboard read failed: %v", err)
		return
	}

	if content != text {
		t.Errorf("expected clipboard to contain text with negative duration, got '%s'", content)
	}
}

// TestPromptConfirmPasswordEmptyMismatch tests PromptConfirmPassword when one is empty
func TestPromptConfirmPasswordEmptyMismatch(t *testing.T) {
	r, w := io.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	go func() {
		w.WriteString("password\n")
		w.WriteString("\n")
		w.Close()
	}()

	password, err := PromptConfirmPassword("Password: ", "Confirm: ")
	if err == nil {
		t.Errorf("expected error for empty confirm password, got none")
	}

	if password != "" {
		t.Errorf("expected empty password on mismatch, got '%s'", password)
	}
}
