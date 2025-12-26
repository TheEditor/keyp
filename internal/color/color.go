package color

import (
	"io"
	"os"
)

// Color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
)

// isTTY checks if output is a terminal
func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	// Check if file descriptor is a terminal
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// Success wraps text in green color if output is a TTY
func Success(text string) string {
	if !isTTY(os.Stdout) {
		return text
	}
	return colorGreen + text + colorReset
}

// Error wraps text in red color if output is a TTY
func Error(text string) string {
	if !isTTY(os.Stderr) {
		return text
	}
	return colorRed + text + colorReset
}

// Warning wraps text in yellow color if output is a TTY
func Warning(text string) string {
	if !isTTY(os.Stdout) {
		return text
	}
	return colorYellow + text + colorReset
}

// Header wraps text in cyan color if output is a TTY
func Header(text string) string {
	if !isTTY(os.Stdout) {
		return text
	}
	return colorCyan + text + colorReset
}
