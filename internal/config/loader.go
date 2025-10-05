package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// FindConfigFile searches for configuration file in standard locations
func FindConfigFile() (string, error) {
	// Search locations in order of precedence
	searchPaths := []string{
		// Current directory
		"media-organizer.json",
		// User config directory
		filepath.Join(os.Getenv("APPDATA"), "media-organizer", "config.json"),
		// Home directory
		filepath.Join(os.Getenv("USERPROFILE"), ".media-organizer.json"),
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no configuration file found in standard locations")
}

// MergeConfigs merges configurations with precedence: CLI > File > Defaults
func MergeConfigs(defaults, file, cli *Config) *Config {
	result := &Config{}

	// Start with defaults
	if defaults != nil {
		*result = *defaults
	}

	// Apply file configuration (overrides defaults)
	if file != nil {
		if file.SourceDir != "" {
			result.SourceDir = file.SourceDir
		}
		if file.TargetDir != "" {
			result.TargetDir = file.TargetDir
		}
		if file.DuplicateDetection != "" {
			result.DuplicateDetection = file.DuplicateDetection
		}
		if file.DuplicateStrategy != "" {
			result.DuplicateStrategy = file.DuplicateStrategy
		}
		if file.Mode != "" {
			result.Mode = file.Mode
		}
		if file.ConfigFile != "" {
			result.ConfigFile = file.ConfigFile
		}
		if file.LogLevel != "" {
			result.LogLevel = file.LogLevel
		}
	}

	// Apply CLI configuration (overrides file and defaults)
	if cli != nil {
		if cli.SourceDir != "" {
			result.SourceDir = cli.SourceDir
		}
		if cli.TargetDir != "" {
			result.TargetDir = cli.TargetDir
		}
		if cli.DuplicateDetection != "" {
			result.DuplicateDetection = cli.DuplicateDetection
		}
		if cli.DuplicateStrategy != "" {
			result.DuplicateStrategy = cli.DuplicateStrategy
		}
		if cli.Mode != "" {
			result.Mode = cli.Mode
		}
		if cli.ConfigFile != "" {
			result.ConfigFile = cli.ConfigFile
		}
		if cli.LogLevel != "" {
			result.LogLevel = cli.LogLevel
		}
	}

	return result
}

// LoadFullConfig loads and merges configuration from all sources
func LoadFullConfig(cliConfig *Config) (*Config, error) {
	// Start with defaults
	defaults := NewDefaultConfig()

	// Determine which config file to use
	var fileConfig *Config
	var configFile string
	var err error

	if cliConfig != nil && cliConfig.ConfigFile != "" {
		// Use explicitly specified config file
		configFile = cliConfig.ConfigFile
	} else {
		// Try to find config file automatically
		configFile, err = FindConfigFile()
		if err != nil {
			// No config file found, that's okay
			fileConfig = nil
		}
	}

	// Load file configuration if available
	if configFile != "" {
		fileConfig, err = LoadConfigFromFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file %s: %w", configFile, err)
		}
	}

	// Merge configurations with precedence: CLI > File > Defaults
	finalConfig := MergeConfigs(defaults, fileConfig, cliConfig)

	// Set the actual config file path that was used
	if configFile != "" {
		finalConfig.ConfigFile = configFile
	}

	return finalConfig, nil
}
