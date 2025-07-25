package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/jaredhaight/lovecms/internal/application"
	"github.com/spf13/viper"
)

var debugLogging = flag.Bool("debug", false, "Enable debug logging")

// Get our paths - cross platform
func getConfigDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "lovecms")
	}
	// For macOS/Linux, use home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(homeDir, ".lovecms")
}

var loveCmsDir = getConfigDir()

func main() {
	flag.Parse()

	// logging defaults
	logLevel := log.InfoLevel
	addSource := false

	if *debugLogging {
		logLevel = log.DebugLevel
		addSource = true
	}

	// setup logging
	opts := log.Options{
		Level:        logLevel,
		ReportCaller: addSource,
	}

	logger := slog.New(log.NewWithOptions(os.Stderr, opts))

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
		err = v.SafeWriteConfig()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	var port = v.GetInt("Port")

	// setup our servers
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))

	// setup handlers
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", application.NewCmsHandler(v, logger).GetHome)
	mux.HandleFunc("GET /editor", application.NewCmsHandler(v, logger).GetEditor)
	mux.HandleFunc("POST /editor", application.NewCmsHandler(v, logger).PostEditor)

	// start server
	logger.Info(fmt.Sprintf("Starting server on http://localhost:%d", port))

	listen := fmt.Sprintf(":%d", port)
	err = http.ListenAndServe(listen, mux)
	logger.Error("Error starting server", "error", err)
}
