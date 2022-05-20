package rest_api

import (
	"diary-api/internal/config"
	"diary-api/internal/db"
	"diary-api/internal/diaries"
	diaryRepository "diary-api/internal/diaries/repository"
	"diary-api/internal/diary_entries"
	diaryEntriesRepository "diary-api/internal/diary_entries/repository"
	v1 "diary-api/internal/protocol/rest_api/v1"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Server interface {
	Run() error
}

type server struct {
	r   *gin.Engine
	cfg *config.Config
}

func (s server) Run() error {
	err := s.r.Run(":" + s.cfg.AppPort)
	if err != nil {
		return err
	}

	return nil
}

func NewServer(cfg *config.Config) Server {
	s := &server{
		cfg: cfg,
		r:   initRouter(),
	}

	dbConn, err := db.InitDb(cfg.PG)
	if err != nil {
		panic(err)
	}

	diaryUc := getDiaryUc(dbConn)
	diaryEntriesUc := getDiaryEntriesUc(dbConn)

	v1.RegisterRoutes(s.r.Group("/api"), diaryUc, diaryEntriesUc)
	return s
}

func initRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	return r
}

func getDiaryUc(dbConn *sqlx.DB) usecase.DiaryUseCase {
	diaryRepo := diaryRepository.NewPostgresDiaryRepository(dbConn)
	return diaries.NewDiaryUseCase(diaryRepo)
}

func getDiaryEntriesUc(dbConn *sqlx.DB) usecase.DiaryEntriesUseCase {
	diaryEntriesRepo := diaryEntriesRepository.New(dbConn)
	return diary_entries.NewUseCase(diaryEntriesRepo)
}
