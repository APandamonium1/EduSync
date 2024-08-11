package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config struct to hold the configuration
type Config struct {
	GoogleClientID     string `json:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret"`
	SessionSecret      string `json:"session_secret"`
	AuthKey            string `json:"auth_key"`
	EncryptKey         string `json:"encrypt_key"`
}

// LoadConfig reads the configuration from a file
func LoadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file: %w", err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("cannot decode config JSON: %w", err)
	}

	return &config, nil
}
