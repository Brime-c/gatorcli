package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	db_url *string `json:"db_url"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, configFileName)
	return path, nil
}

func write(cfg Config) error {
	return nil
}
