package config

import (
	"encoding/json"
	"os"
)

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	return write(c)
}

func write(c *Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(c); err != nil {
		return err
	}

	return nil
}
