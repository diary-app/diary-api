package main

import (
	"diary-api/internal/config"
	"diary-api/internal/protocol/rest"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.New()
	if err := godotenv.Load(".env"); err != nil {
		logger.Printf(".env file not loaded: %v", err)
	}

	cfg, err := config.Read()
	if err != nil {
		logger.Fatalf("cannot start server: %v", err)
	}
	if cfg.LogDebug {
		logger.SetLevel(log.DebugLevel)
	}

	server := rest.NewServer(cfg, logger)
	if err := server.Run(); err != nil {
		logger.Fatalf("failed to start REST API server: %v", err)
	}
}
