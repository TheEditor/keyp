package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config holds the keyp configuration
type Config struct {
	SessionTimeout time.Duration
}

// Load loads the configuration from ~/.keyp/config.yaml
// Environment variable KEYP_SESSION_TIMEOUT overrides config file
func Load() (*Config, error) {
	cfg := &Config{
		SessionTimeout: 15 * time.Minute, // Default
	}

	// Check environment variable first
	if envTimeout := os.Getenv("KEYP_SESSION_TIMEOUT"); envTimeout != "" {
		duration, err := time.ParseDuration(envTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid KEYP_SESSION_TIMEOUT: %w", err)
		}
		cfg.SessionTimeout = duration
		return cfg, nil
	}

	// Load from config file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil // Use defaults
	}

	configPath := filepath.Join(homeDir, ".keyp", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Use defaults if file doesn't exist
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Simple YAML-like parsing for session_timeout
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "session_timeout:") {
			// Extract the value after the colon
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				timeoutStr := strings.TrimSpace(parts[1])
				duration, err := time.ParseDuration(timeoutStr)
				if err != nil {
					return nil, fmt.Errorf("invalid session_timeout in config: %w", err)
				}
				cfg.SessionTimeout = duration
				break
			}
		}
	}

	return cfg, nil
}
