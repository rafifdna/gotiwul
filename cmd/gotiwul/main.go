package main

import (
    "context"
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"
    
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

    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Start server in goroutine
    errChan := make(chan error, 1)
    go func() {
        errChan <- srv.StartWithAutoTLS()
    }()

    // Wait for interrupt signal or server error
    select {
    case err := <-errChan:
        log.Fatalf("Server error: %v", err)
    case sig := <-sigChan:
        log.Printf("Received signal: %v", sig)
        // Graceful shutdown
        if err := srv.Shutdown(context.Background()); err != nil {
            log.Printf("Shutdown error: %v", err)
        }
    }
}