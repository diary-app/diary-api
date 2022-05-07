package rest_api

import (
	"diary-api/internal/protocol/rest_api/diaries_handler"
	"diary-api/internal/protocol/rest_api/diary_entries_handler"
	"diary-api/internal/protocol/rest_api/sharing_tasks_handler"
	"diary-api/internal/protocol/rest_api/users_handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *server) registerRoutes() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	registerDiariesRoutes(r)
	registerDiaryEntriesRoutes(r)
	registerUsersRoutes(r)
	registerSharingTasksRoutes(r)

	s.router = r
}

func registerDiariesRoutes(r *gin.Engine) {
	handler := diaries_handler.New()
	diariesGroup := r.Group("/diaries")
	diariesGroup.GET("", handler.GetMyDiaries())
	diariesGroup.GET("/:id/entries", handler.GetDiaryEntries())
}

func registerDiaryEntriesRoutes(r *gin.Engine) {
	handler := diary_entries_handler.New()
	diaryEntries := r.Group("/diary-entries")
	diaryEntries.GET("", handler.GetList())
	diaryEntries.GET("/:id/download", handler.Download())
	diaryEntries.POST("", handler.Create())
	diaryEntries.POST("/:id/upload", handler.Upload())
	diaryEntries.DELETE("/:id", handler.Delete())
	diaryEntries.PATCH("/:id", handler.Patch())
}

func registerUsersRoutes(r *gin.Engine) {
	handler := users_handler.New()
	users := r.Group("/users")
	users.POST("/login", handler.Login())
	users.POST("/register", handler.Register())
	users.GET("/:id/sharing-key", handler.GetSharingKey())
}

func registerSharingTasksRoutes(r *gin.Engine) {
	handler := sharing_tasks_handler.New()
	sharingTasks := r.Group("/sharing-tasks")
	sharingTasks.GET("", handler.GetAllMine())
	sharingTasks.POST("", handler.Create())
	sharingTasks.DELETE("/:id", handler.DeleteById())
}
