package config

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	Port        int      `json:"port"`
	CurrentSite string   `json:"current_site"`
	Sites       []string `json:"sites"`
}

func loadConfig(logger *slog.Logger) (*Config, error) {
	// Create our config variable
	config := Config{}

	// Get our paths. Right now we're not worried about cross platform
	appDataDir := os.Getenv("APPDATA")
	loveCmsDir := filepath.Join(appDataDir, "lovecms")
	configFilePath := filepath.Join(loveCmsDir, "config.json")

	// load file from json
	logger.Debug("Loading config file", "path", configFilePath)
	configFile, err := os.Open(configFilePath)

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(configFile)

	// if we don't have a config file, create it and return the empty config
	if errors.Is(err, os.ErrNotExist) {
		// create our directory if it doesn't exist
		err := os.MkdirAll(loveCmsDir, os.ModePerm)
		if err != nil {
			return nil, err
		}

		// create the config file
		logger.Debug("Creating config file", "path", configFilePath)
		config.Port = 8143
		configFile, err = os.Create(configFilePath)

		// throw error if we can't create the file
		if err != nil {
			return nil, err
		}

		// save file
		logger.Debug("Saving config file", "path", configFilePath)
		encoder := json.NewEncoder(configFile)
		err = encoder.Encode(config)

		// return our config
		return &config, err
	}

	// if we have an error, bomb out
	if err != nil {
		logger.Error("Error loading config file", "path", configFilePath, "error", err)
		return nil, err
	}

	// Parse our config file from disk
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)

	// set a sane default for port
	if config.Port == 0 {
		config.Port = 8143
	}
	return &config, err
}

func MustLoadConfig(logger *slog.Logger) *Config {
	cfg, err := loadConfig(logger)
	if err != nil {
		logger.Error("Error loading config file", "error", err)
		panic(err)
	}
	return cfg
}
