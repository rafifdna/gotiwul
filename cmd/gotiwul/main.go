package main

import (
    "flag"
    "log"
    
    "gotiwul/config"
    "gotiwul/internal/server"
)

func main() {
    configPath := flag.String("config", "config/config.yaml", "Path to configuration file")
    flag.Parse()

    // Load configuration
    cfg, err := config.Load(*configPath)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create and start server
    srv := server.New(cfg)
    if err := srv.Start(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}