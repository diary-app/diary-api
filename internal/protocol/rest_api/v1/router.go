package v1

import (
	"diary-api/internal/protocol/rest_api/v1/diaries_handler"
	"diary-api/internal/protocol/rest_api/v1/diary_entries_handler"
	"diary-api/internal/protocol/rest_api/v1/sharing_tasks_handler"
	"diary-api/internal/protocol/rest_api/v1/users_handler"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(r *gin.RouterGroup, diaryUc usecase.DiaryUseCase, diaryEntriesUc usecase.DiaryEntriesUseCase,
	usersUc usecase.UsersUseCase, jwtMw gin.HandlerFunc) {
	rg := r.Group("/v1")

	rg.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	usersH := users_handler.New(usersUc)
	diariesH := diaries_handler.New(diaryUc)
	diaryEntriesH := diary_entries_handler.New(diaryEntriesUc)
	sharingTasksH := sharing_tasks_handler.New()

	registerAuthRoutes(rg, usersH)
	rg = rg.Group("")
	rg.Use(jwtMw)
	registerUsersRoutes(rg, usersH)
	registerDiariesRoutes(rg, diariesH)
	registerDiaryEntriesRoutes(rg, diaryEntriesH)
	registerSharingTasksRoutes(rg, sharingTasksH)
}

func registerAuthRoutes(rg *gin.RouterGroup, usersH users_handler.Handler) {
	authGroup := rg.Group("/auth")
	authGroup.POST("/login", usersH.Login())
	authGroup.POST("/register", usersH.Register())
}

func registerUsersRoutes(rg *gin.RouterGroup, usersH users_handler.Handler) {
	users := rg.Group("/users")
	users.GET("/me", usersH.GetMe())
	users.GET("/:id", usersH.GetUser())
}

func registerSharingTasksRoutes(rg *gin.RouterGroup, sharingTasksH sharing_tasks_handler.Handler) {
	sharingTasks := rg.Group("/sharing-tasks")
	sharingTasks.GET("", sharingTasksH.GetAllMine())
	sharingTasks.POST("", sharingTasksH.Create())
	sharingTasks.DELETE("/:id", sharingTasksH.DeleteById())
}

func registerDiaryEntriesRoutes(rg *gin.RouterGroup, diaryEntriesH diary_entries_handler.Handler) {
	diaryEntries := rg.Group("/diary-entries")
	diaryEntries.GET("", diaryEntriesH.GetList())
	diaryEntries.GET("/:id/download", diaryEntriesH.Download())
	diaryEntries.POST("", diaryEntriesH.Create())
	diaryEntries.POST("/:id/upload", diaryEntriesH.Upload())
	diaryEntries.DELETE("/:id", diaryEntriesH.Delete())
	diaryEntries.PATCH("/:id", diaryEntriesH.Patch())
}

func registerDiariesRoutes(rg *gin.RouterGroup, diariesH diaries_handler.Handler) {
	diariesGroup := rg.Group("/diaries")
	diariesGroup.GET("", diariesH.GetMyDiaries())
	diariesGroup.POST("", diariesH.CreateDiary())
}
