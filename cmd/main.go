package main

import (
	"flag"
	"fmt"
	"github.com/jaredhaight/lovecms/internal/config"
	"github.com/jaredhaight/lovecms/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// Get our paths
	configPath := *flag.String("config", "config.json", "Path to the LoveCMS config file")
	flag.Parse()

	cfg := config.MustLoadConfig(configPath)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", handlers.NewHomeHandler().ServeHTTP)
	log.Println("message", "2", "3")
	log.Printf("Starting server on %d\n", cfg.Port)

	port := fmt.Sprintf(":%d", cfg.Port)

	err := http.ListenAndServe(port, mux)
	log.Fatal("Error starting server", "error", err)
}
