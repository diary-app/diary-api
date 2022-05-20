package main

import (
	"diary-api/internal/config"
	"diary-api/internal/protocol/rest_api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf(".env file not loaded: %v", err)
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}

	server := rest_api.NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Fatalf("failed to start REST API server: %v", err)
	}
}
