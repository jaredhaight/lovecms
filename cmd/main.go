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
	"github.com/jaredhaight/lovecms/internal/handlers"
	"github.com/spf13/viper"
)

var debugLogging = flag.Bool("debug", false, "Enable debug logging")
var sitePath = flag.String("site", "", "Path to the site directory")

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

	// get sitepath
	sitePathConfig := v.GetString("SitePath")

	// Use command line flag if provided, otherwise use config
	var finalSitePath string
	if *sitePath != "" {
		finalSitePath = *sitePath
		// Update config with the provided path
		v.Set("SitePath", finalSitePath)
		err = v.SafeWriteConfig()
		if err != nil {
			// If SafeWriteConfig fails, try WriteConfigAs
			configPath := filepath.Join(loveCmsDir, "config.json")
			err = v.WriteConfigAs(configPath)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
		}
	} else {
		finalSitePath = sitePathConfig
	}

	// setup our servers
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))

	// setup handlers
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", handlers.NewHomeHandler(v, logger).Get)
	mux.HandleFunc("GET /post", handlers.NewPostHandler(v, logger).Get)
	mux.HandleFunc("GET /posts/new", handlers.NewPostHandler(v, logger).GetNew)
	mux.HandleFunc("POST /posts/new", handlers.NewPostHandler(v, logger).PostNew)
	mux.HandleFunc("POST /posts/edit", handlers.NewPostHandler(v, logger).PostEdit)

	// start server
	logger.Info(fmt.Sprintf("Starting server on http://localhost:%d", port))

	listen := fmt.Sprintf(":%d", port)
	err = http.ListenAndServe(listen, mux)
	logger.Error("Error starting server", "error", err)
}
