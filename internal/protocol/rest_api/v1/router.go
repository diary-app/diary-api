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

func RegisterRoutes(r *gin.RouterGroup, diaryUc usecase.DiaryUseCase, diaryEntriesUc usecase.DiaryEntriesUseCase) {
	rg := r.Group("/v1")

	rg.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	registerDiariesRoutes(rg, diaryUc)
	registerDiaryEntriesRoutes(rg, diaryEntriesUc)
	registerUsersRoutes(rg)
	registerSharingTasksRoutes(rg)
}

func registerDiariesRoutes(r *gin.RouterGroup, uc usecase.DiaryUseCase) {
	handler := diaries_handler.New(uc)
	diariesGroup := r.Group("/diaries")
	diariesGroup.GET("", handler.GetMyDiaries())
}

func registerDiaryEntriesRoutes(r *gin.RouterGroup, uc usecase.DiaryEntriesUseCase) {
	handler := diary_entries_handler.New(uc)
	diaryEntries := r.Group("/diary-entries")
	diaryEntries.GET("", handler.GetList())
	diaryEntries.GET("/:id/download", handler.Download())
	diaryEntries.POST("", handler.Create())
	diaryEntries.POST("/:id/upload", handler.Upload())
	diaryEntries.DELETE("/:id", handler.Delete())
	diaryEntries.PATCH("/:id", handler.Patch())
}

func registerUsersRoutes(r *gin.RouterGroup) {
	handler := users_handler.New()
	users := r.Group("/users")
	users.POST("/login", handler.Login())
	users.POST("/register", handler.Register())
	users.GET("/me", handler.GetMe())
	users.GET("/:id", handler.GetUser())
}

func registerSharingTasksRoutes(r *gin.RouterGroup) {
	handler := sharing_tasks_handler.New()
	sharingTasks := r.Group("/sharing-tasks")
	sharingTasks.GET("", handler.GetAllMine())
	sharingTasks.POST("", handler.Create())
	sharingTasks.DELETE("/:id", handler.DeleteById())
}
