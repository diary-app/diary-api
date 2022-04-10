package main

import (
	"diary-api/internal/config"
	"diary-api/internal/protocol/http"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf(".env file not loaded: %v", err)
	}

	cfg := config.Read()
	server := http.NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Fatalf("Unable to start http server: %v", err)
	}
}
