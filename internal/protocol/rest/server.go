package rest

import (
	"diary-api/internal/auth"
	authUsecase "diary-api/internal/auth/usecase"
	"diary-api/internal/config"
	"diary-api/internal/db"
	diaryRepository "diary-api/internal/diaries/repository"
	diaryUsecase "diary-api/internal/diaries/usecase"
	diaryEntriesRepository "diary-api/internal/diary_entries/repository"
	diaryEntriesUsecase "diary-api/internal/diary_entries/usecase"
	"diary-api/internal/protocol/rest/middleware"
	v1 "diary-api/internal/protocol/rest/v1"
	"diary-api/internal/sharing_tasks"
	sharingTasksRepository "diary-api/internal/sharing_tasks/repository"
	"diary-api/internal/usecase"
	"diary-api/internal/users"
	usersRepository "diary-api/internal/users/repository"
	"github.com/benbjohnson/clock"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func NewServer(cfg *config.Config, l *log.Logger) Server {
	dbConn, err := db.InitDb(cfg.PG)
	if err != nil {
		panic(err)
	}

	myClock := clock.New()
	tokenService := auth.NewTokenService(&cfg.Auth, myClock)
	diaryUc := getDiaryUc(dbConn)
	diaryEntriesUc := getDiaryEntriesUc(dbConn)
	usersUc := getUsersAndAuthUc(dbConn)
	sharingTasksUc := getSharingTasksUc(dbConn)
	authUc := getAuthUc(dbConn, tokenService)

	errorHandlerMw := middleware.ErrorHandler(l)
	s := &server{
		cfg: cfg,
		r:   initRouter(errorHandlerMw),
	}

	s.r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	jwtMw := middleware.JwtMiddleware(tokenService)
	rg := s.r.Group("")
	v1.RegisterRoutes(rg, jwtMw, diaryUc, diaryEntriesUc, usersUc, sharingTasksUc, authUc)
	return s
}

func getSharingTasksUc(conn *sqlx.DB) usecase.SharingTasksUseCase {
	stRepo := sharingTasksRepository.New(conn)
	return sharing_tasks.NewUseCase(stRepo)
}

func initRouter(errorHandler gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Origin"},
		AllowCredentials: false,
	}))
	r.Use(errorHandler)
	return r
}

func getDiaryUc(dbConn *sqlx.DB) usecase.DiaryUseCase {
	diaryRepo := diaryRepository.NewPostgresDiaryRepository(dbConn)
	return diaryUsecase.New(diaryRepo)
}

func getDiaryEntriesUc(dbConn *sqlx.DB) usecase.DiaryEntriesUseCase {
	diaryEntriesRepo := diaryEntriesRepository.New(dbConn)
	return diaryEntriesUsecase.New(diaryEntriesRepo)
}

func getUsersAndAuthUc(conn *sqlx.DB) usecase.UsersUseCase {
	usersRepo := usersRepository.New(conn)
	return users.NewUseCase(usersRepo)
}

func getAuthUc(conn *sqlx.DB, tokenService auth.TokenService) usecase.AuthUseCase {
	usersRepo := usersRepository.New(conn)
	return authUsecase.New(usersRepo, tokenService)
}
