package main

import (
	"encoding/json"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
)

type Config struct {
	LastSite int
	Sites    []Site
}

type Site struct {
	Id   int
	Path string
}

func (app *application) loadLoveConfig() {

	// load file from json
	configFile, err := os.Open(app.loveConfigPath)

	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

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

	defer tomlFile.Close()
	decoder := toml.NewDecoder(tomlFile)
	err = decoder.Decode(&app.hugoConfig)
	if err != nil {
		log.Fatal(err)
	}
}
