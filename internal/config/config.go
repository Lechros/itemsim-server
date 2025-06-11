package config

import (
	"os"
	"path/filepath"
)

// Config holds all configuration for the application
type Config struct {
	ResourcesPath string
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		ResourcesPath: getEnvOrDefault("RESOURCES_PATH", "resources"),
	}
}

// GetFilePath returns the full path for a resource file
func (c *Config) GetFilePath(fileName string) string {
	return filepath.Join(c.ResourcesPath, fileName)
}

// getEnvOrDefault returns the value of the environment variable or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
