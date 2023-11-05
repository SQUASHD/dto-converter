// Package config includes functions for initializing, reading, and writing
// Gopen config files.
package config

import (
	"dto-converter/internal/structs"
	"encoding/json"
	"os"
)

// Literally copy-pasted from https://github.com/wipdev-tech/gopen

// Init checks if the config file exists in configPath. If not, creates an
// empty config file. configDir will also be created if it doesn't exist.
func Init(configDir string, configPath string) error {
	_, err := os.Stat(configPath)
	if err != nil {
		return os.ErrExist
	}

	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(configPath)
	if err != nil {
		return err
	}

	emptyConfig := structs.Config{}
	err = Write(emptyConfig, configPath)
	if err != nil {
		return err
	}

	return nil
}

// Write writes config to configPath (will OVERWRITE if file already exists)
func Write(config structs.Config, configPath string) error {
	jsonFile, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, jsonFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Read reads the configPath file and returns a Config struct
func Read(configPath string) (config structs.Config, err error) {
	f, err := os.ReadFile(configPath)
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &config)
	if err != nil {
		return structs.Config{}, err
	}

	return
}
