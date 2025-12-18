package ui

import (
	"time"

	"github.com/atotto/clipboard"
)

const DefaultClearDuration = 45 * time.Second

// CopyToClipboard copies text to the system clipboard
func CopyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}

// CopyWithAutoClear copies text to clipboard and clears it after the specified duration
func CopyWithAutoClear(text string, clearAfter time.Duration) error {
	if err := clipboard.WriteAll(text); err != nil {
		return err
	}

	if clearAfter > 0 {
		go func() {
			time.Sleep(clearAfter)
			// Only clear if clipboard still contains our text
			current, _ := clipboard.ReadAll()
			if current == text {
				clipboard.WriteAll("")
			}
		}()
	}

	return nil
}
