package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	DBPath string
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() (*Config, error) {
	dbPath, err := GetDatabasePath()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBPath: dbPath,
	}, nil
}

// GetDatabasePath returns the path to the database file
func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	kotoDir := filepath.Join(homeDir, ".koto")

	// Create .koto directory if it doesn't exist
	if err := os.MkdirAll(kotoDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create koto directory: %w", err)
	}

	dbPath := filepath.Join(kotoDir, "koto.db")
	return dbPath, nil
}
