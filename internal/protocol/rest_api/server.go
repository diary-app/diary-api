package rest_api

import (
	"diary-api/internal/auth"
	"diary-api/internal/config"
	"diary-api/internal/db"
	"diary-api/internal/diaries"
	diaryRepository "diary-api/internal/diaries/repository"
	"diary-api/internal/diary_entries"
	diaryEntriesRepository "diary-api/internal/diary_entries/repository"
	"diary-api/internal/protocol/rest_api/middleware"
	v1 "diary-api/internal/protocol/rest_api/v1"
	"diary-api/internal/sharing_tasks"
	sharingTasksRepository "diary-api/internal/sharing_tasks/repository"
	"diary-api/internal/usecase"
	"diary-api/internal/users"
	usersRepository "diary-api/internal/users/repository"
	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
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
	authService := auth.NewAuthService(&cfg.Auth, myClock)
	diaryUc := getDiaryUc(dbConn)
	diaryEntriesUc := getDiaryEntriesUc(dbConn)
	usersUc := getUsersUc(dbConn, authService)
	sharingTasksUc := getSharingTasksUc(dbConn)

	errorHandlerMw := middleware.ErrorHandler(l)
	s := &server{
		cfg: cfg,
		r:   initRouter(errorHandlerMw),
	}

	jwtMw := middleware.JwtMiddleware(authService)
	v1.RegisterRoutes(s.r.Group(""), jwtMw, diaryUc, diaryEntriesUc, usersUc, sharingTasksUc)
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
	r.Use(errorHandler)
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

func getUsersUc(conn *sqlx.DB, tokensManager auth.TokenService) usecase.UsersUseCase {
	usersRepo := usersRepository.New(conn)
	return users.NewUseCase(tokensManager, usersRepo)
}
