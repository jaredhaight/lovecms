package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port        int
	CurrentSite string `json:"current_site"`
	Sites       []string
}

func loadConfig(configPath string) (*Config, error) {

	// load file from json
	configFile, err := os.Open(configPath)

	if err != nil {
		log.Fatal(err)
	}

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(configFile)

	config := Config{}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)

	if config.Port == 0 {
		config.Port = 8143
	}
	return &config, err
}

func MustLoadConfig(configPath string) *Config {
	cfg, err := loadConfig(configPath)
	if err != nil {
		panic(err)
	}
	return cfg
}
