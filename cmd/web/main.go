package main

import (
	"flag"
	"github.com/jaredhaight/lovecms/internal/hugo"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger         *slog.Logger
	loveConfigPath string
	loveConfig     Config
	hugoConfigPath string
	hugoConfig     hugo.Config
}

func main() {
	// create our app state
	app := &application{}

	// Setup logging
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	app.logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))

	// Get our paths
	app.loveConfigPath = *flag.String("config", "config.json", "Path to the LoveCMS config file")
	app.hugoConfigPath = *flag.String("hugo-config", "config.toml", "Path to the Hugo config file")
	flag.Parse()

	// load configs
	app.loadLoveConfig()
	app.loadHugoConfig()

	// print contents
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.listPosts)

	app.logger.Info("Starting server on 8143")
	err := http.ListenAndServe(":8143", mux)
	app.logger.Error("Error starting server", "error", err)
}
