package config

import (
	"encoding/json"
	"os"
)

// Get config file path
func getConfigFilePath() (string, error) {
	basePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := basePath + "/" + configFileName
	return configPath, nil
}

// Write Config
func writeConfig(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	filePath, _ := getConfigFilePath()
	err = os.WriteFile(filePath, jsonData, 0666)
	if err != nil {
		return err
	}
	return nil
}

// Create Config
func ReadConfig() Config {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}
	}

	return config
}
