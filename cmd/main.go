package main

import (
    "context"
    "log"
    "auction-system/internal/app"
    "auction-system/internal/config"
)

func main() {
    ctx := context.Background()

    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    application, err := app.NewApp(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
    }

    log.Println("Starting application...")
    if err := application.Run(ctx); err != nil {
        log.Fatalf("Error running app: %v", err)
    }
}
