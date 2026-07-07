package config

import (
	"encoding/json"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func (c *Config) SetUser(current_user string) error {
	c.CurrentUsername = current_user
	return writeConfig(*c)
}

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
func ReadConfig() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
