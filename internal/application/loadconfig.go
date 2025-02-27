package application

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

func loadConfig(logger slog.Logger) (Config, error) {
	// Create our application variable
	config := Config{}

	// Get our paths. Right now we're not worried about cross platform
	appDataDir := os.Getenv("APPDATA")
	loveCmsDir := filepath.Join(appDataDir, "lovecms")
	configFilePath := filepath.Join(loveCmsDir, "config.json")

	// load file from json
	logger.Debug("Loading application file", "path", configFilePath)
	configFile, err := os.Open(configFilePath)

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(configFile)

	// if we don't have a application file, create it and return the empty application
	if errors.Is(err, os.ErrNotExist) {
		// create our directory if it doesn't exist
		err := os.MkdirAll(loveCmsDir, os.ModePerm)
		if err != nil {
			return config, err
		}

		// create the application file
		logger.Debug("Creating application file", "path", configFilePath)
		config.Port = 8143
		configFile, err = os.Create(configFilePath)

		// throw error if we can't create the file
		if err != nil {
			return config, err
		}

		// save file
		logger.Debug("Saving application file", "path", configFilePath)
		encoder := json.NewEncoder(configFile)
		err = encoder.Encode(config)

		// return our application
		return config, err
	}

	// if we have an error, bomb out
	if err != nil {
		logger.Error("Error loading application file", "path", configFilePath, "error", err)
		return config, err
	}

	// Parse our application file from disk
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)

	// set a sane default for port
	if config.Port == 0 {
		config.Port = 8143
	}
	return config, err
}

func MustLoadConfig(logger slog.Logger) Config {
	cfg, err := loadConfig(logger)
	if err != nil {
		logger.Error("Error loading application file", "error", err)
		panic(err)
	}
	return cfg
}
