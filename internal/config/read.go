package config

import (
	"encoding/json"
	"os"
)

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
