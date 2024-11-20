package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot read user home dir: %v", err)
	}

	return filepath.Join(homeDir, configFileName), nil
}
