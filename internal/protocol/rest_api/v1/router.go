package v1

import (
	"diary-api/internal/protocol/rest_api/v1/auth_handler"
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

	registerAuthRoutes(rg, usersUc)
	authRg := rg.Group("")
	authRg.Use(jwtMw)
	registerUsersRoutes(authRg, usersUc)
	registerDiariesRoutes(authRg, diaryUc)
	registerDiaryEntriesRoutes(authRg, diaryEntriesUc)
	registerSharingTasksRoutes(authRg)
}

func registerAuthRoutes(rg *gin.RouterGroup, uc usecase.UsersUseCase) {
	authH := auth_handler.New(uc)
	authGroup := rg.Group("/auth")
	authGroup.POST("/login", authH.Login())
	authGroup.POST("/register", authH.Register())
}

func registerUsersRoutes(rg *gin.RouterGroup, uc usecase.UsersUseCase) {
	usersH := users_handler.New(uc)
	users := rg.Group("/users")
	users.GET("/me", usersH.GetMe())
	users.GET("/:id", usersH.GetUserById())
	users.GET("/name=:name", usersH.GetUserByName())
}

func registerSharingTasksRoutes(rg *gin.RouterGroup) {
	sharingTasksH := sharing_tasks_handler.New()
	sharingTasks := rg.Group("/sharing-tasks")
	sharingTasks.GET("", sharingTasksH.GetAllMine())
	sharingTasks.POST("", sharingTasksH.Create())
	sharingTasks.DELETE("/:id", sharingTasksH.DeleteById())
}

func registerDiaryEntriesRoutes(rg *gin.RouterGroup, diaryEntriesUc usecase.DiaryEntriesUseCase) {
	diaryEntriesH := diary_entries_handler.New(diaryEntriesUc)
	diaryEntries := rg.Group("/diary-entries")
	diaryEntries.GET("", diaryEntriesH.GetList())
	diaryEntries.GET("/:id/download", diaryEntriesH.Download())
	diaryEntries.POST("", diaryEntriesH.Create())
	diaryEntries.POST("/:id/upload", diaryEntriesH.Upload())
	diaryEntries.DELETE("/:id", diaryEntriesH.Delete())
	diaryEntries.PATCH("/:id", diaryEntriesH.Patch())
}

func registerDiariesRoutes(rg *gin.RouterGroup, diaryUc usecase.DiaryUseCase) {
	diariesH := diaries_handler.New(diaryUc)
	diariesGroup := rg.Group("/diaries")
	diariesGroup.GET("", diariesH.GetMyDiaries())
	diariesGroup.POST("", diariesH.CreateDiary())
}
