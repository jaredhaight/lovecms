package main

import (
	"flag"
	"github.com/jaredhaight/lovecms/internal/hugo"
	"log"
)

type application struct {
	loveConfigPath string
	loveConfig     Config
	hugoConfigPath string
	hugoConfig     hugo.Config
}

func main() {
	app := &application{}
	app.loveConfigPath = *flag.String("config", "config.json", "Path to the LoveCMS config file")
	app.hugoConfigPath = *flag.String("hugo-config", "config.toml", "Path to the Hugo config file")
	flag.Parse()

	// load configs
	app.loadLoveConfig()
	app.loadHugoConfig()

	// print contents
	log.Println(app.loveConfig)
	log.Printf("Site Name: %s\n", app.hugoConfig.Title)
	log.Printf("Theme: %s\n", app.hugoConfig.Theme)
	log.Printf("Taxonomies:")
	for k, t := range app.hugoConfig.Taxonomies {
		log.Printf("\t%s - %s\n", k, t)
	}
}
