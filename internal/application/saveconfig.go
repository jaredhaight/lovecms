package application

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

func saveConfig(logger slog.Logger, config Config) (bool, error) {
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

	// if we don't have a application file, create it
	if errors.Is(err, os.ErrNotExist) {
		// create our directory if it doesn't exist
		err := os.MkdirAll(loveCmsDir, os.ModePerm)
		if err != nil {
			return false, err
		}

		// create the application file
		logger.Debug("Creating application file", "path", configFilePath)
		configFile, err = os.Create(configFilePath)

		// throw error if we can't create the file
		if err != nil {
			return false, err
		}
	}

	// set sane default for port
	if config.Port == 0 {
		config.Port = 8143
	}

	// save file
	logger.Debug("Saving application file", "path", configFilePath)
	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(config)

	if err != nil {
		return false, err
	}

	// return our application
	return true, err
}
