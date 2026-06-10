package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds the structure of our JSON configuration file,
// tracking the database URL connection string and the currently active username.
type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// configFileName defines the name of the hidden configuration file stored in the user's home directory.
const configFileName = ".gatorconfig.json"

// getConfigFilePath resolves the absolute file path to the config file (e.g., /home/username/.gatorconfig.json).
func getConfigFilePath() (string, error) {
	// 1. Retrieve the current user's system home directory in a cross-platform friendly way
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// 2. Join the home directory path with our configuration file name
	path := filepath.Join(home, configFileName)
	return path, nil
}

// Read reads, parses, and returns the application configuration from the local configuration file.
func Read() (Config, error) {
	// 1. Resolve the path to the configuration file
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	// 2. Read the raw byte data from the file
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	// 3. Parse (unmarshal) the JSON raw bytes into a Config struct
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// write marshals the Config struct into JSON format and saves it to the configuration file.
func write(cfg Config) error {
	// 1. Resolve the path to the configuration file
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// 2. Convert (marshal) the Config struct into indented/raw JSON bytes
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	// 3. Write the JSON data to the file with standard read/write permissions (0644)
	return os.WriteFile(path, jsonData, 0644)
}

// SetUser updates the active user's name in the configuration memory
// and instantly persists that update to the local config file on disk.
func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	return write(*c)
}
