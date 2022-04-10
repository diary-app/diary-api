package http

import (
	"diary-api/internal/config"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Run() error
}

type server struct {
	router *gin.Engine
	cfg    *config.Config
}

func (s server) Run() error {
	err := s.router.Run(":" + s.cfg.Port)
	if err != nil {
		return err
	}

	return nil
}

func NewServer(cfg *config.Config) Server {
	router := SetupRouter()
	return &server{
		router: router,
		cfg:    cfg,
	}
}
