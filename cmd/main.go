package main

import (
	"flag"
	"fmt"
	"github.com/dusted-go/logging/prettylog"
	"github.com/jaredhaight/lovecms/internal/application"
	"github.com/jaredhaight/lovecms/internal/handlers"
	"log/slog"
	"net/http"
)

var debugLogging = flag.Bool("debug", false, "Enable debug logging")

func main() {
	// parse flags
	flag.Parse()

	// logging defaults
	logLevel := slog.LevelInfo
	addSource := false

	if *debugLogging {
		logLevel = slog.LevelDebug
		addSource = true
	}

	// setup logging
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: addSource,
	}

	logger := slog.New(prettylog.NewHandler(opts))

	// load config
	cfg := application.MustLoadConfig(logger)

	// setup our servers
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))

	// setup handlers
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", handlers.NewHomeHandler(cfg, logger).ServeHTTP)

	// start server
	logger.Info(fmt.Sprintf("Starting server on http://localhost:%d", cfg.Port))

	port := fmt.Sprintf(":%d", cfg.Port)
	err := http.ListenAndServe(port, mux)
	logger.Error("Error starting server", "error", err)
}
