package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aerogu/tvchooser"
	"github.com/dusted-go/logging/prettylog"
	"github.com/jaredhaight/lovecms/internal/handlers"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

var debugLogging = flag.Bool("debug", false, "Enable debug logging")

// Get our paths. Right now we're not worried about cross platform
var appDataDir = os.Getenv("APPDATA")
var loveCmsDir = filepath.Join(appDataDir, "lovecms")

func main() {
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

	// make sure our directory exists
	err := os.MkdirAll(loveCmsDir, os.ModePerm)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// setup config stuff
	v := viper.New()
	v.AddConfigPath(loveCmsDir)
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.SetDefault("SitePath", "")
	v.SetDefault("Port", 8143)

	// load config
	err = v.ReadInConfig()
	if errors.Is(err, viper.ConfigFileNotFoundError{}) {
		// if we can't find a config file, create it
		err = viper.WriteConfig()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	var port = v.GetInt("Port")

	// get sitepath
	sitePath := v.GetString("SitePath")

	// prompt for a site directory if we don't have one
	if sitePath == "" {
		sitePath = tvchooser.DirectoryChooser(nil, false)
		v.Set("SitePath", sitePath)
		err = v.WriteConfig()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	// setup our servers
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))

	// setup handlers
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", handlers.NewHomeHandler(v, logger).ServeHTTP)
	mux.HandleFunc("/post/{$}", handlers.NewPostHandler(v, logger).ServeHTTP)

	// start server
	logger.Info(fmt.Sprintf("Starting server on http://localhost:%d", port))

	listen := fmt.Sprintf(":%d", port)
	err = http.ListenAndServe(listen, mux)
	logger.Error("Error starting server", "error", err)
}
