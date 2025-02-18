package main

import (
	"encoding/json"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
)

type Config struct {
	CurrentSite string `json:"current_site"`
	Sites       []string
}

func (app *application) loadLoveConfig() {

	// load file from json
	configFile, err := os.Open(app.loveConfigPath)

	if err != nil {
		log.Fatal(err)
	}

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(configFile)

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&app.loveConfig)

	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) loadHugoConfig() {
	tomlFile, err := os.Open(app.hugoConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	defer func(tomlFile *os.File) {
		err := tomlFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tomlFile)

	decoder := toml.NewDecoder(tomlFile)
	err = decoder.Decode(&app.hugoConfig)
	if err != nil {
		log.Fatal(err)
	}
}
