package v1

import (
	"diary-api/internal/protocol/rest/v1/auth"
	"diary-api/internal/protocol/rest/v1/diaries"
	diaryEntries "diary-api/internal/protocol/rest/v1/diary_entries"
	sharingTasks "diary-api/internal/protocol/rest/v1/sharing_tasks"
	"diary-api/internal/protocol/rest/v1/users"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	jwtMw gin.HandlerFunc,
	diaryUc usecase.DiaryUseCase,
	diaryEntriesUc usecase.DiaryEntriesUseCase,
	usersUc usecase.UsersUseCase,
	sharingTasksUc usecase.SharingTasksUseCase,
	authUc usecase.AuthUseCase,
) {
	rg := r.Group("/api/v1")
	registerAuthRoutes(rg, authUc)
	authRg := rg.Group("")
	authRg.Use(jwtMw)
	registerUsersRoutes(authRg, usersUc)
	registerDiariesRoutes(authRg, diaryUc)
	registerDiaryEntriesRoutes(authRg, diaryEntriesUc)
	registerSharingTasksRoutes(authRg, sharingTasksUc)
}

func registerAuthRoutes(rg *gin.RouterGroup, uc usecase.AuthUseCase) {
	authH := auth.New(uc)
	authRoute := rg.Group("/auth")
	authRoute.POST("/login", authH.Login())
	authRoute.POST("/register", authH.Register())
	authRoute.POST("/refresh-token", authH.RefreshToken())
}

func registerUsersRoutes(rg *gin.RouterGroup, uc usecase.UsersUseCase) {
	usersH := users.New(uc)
	usersRoute := rg.Group("/users")
	usersRoute.GET("/me", usersH.GetMe())
	usersRoute.GET("/:id", usersH.GetUserByID())
	usersRoute.GET("/name=:name", usersH.GetUserByName())
}

func registerSharingTasksRoutes(rg *gin.RouterGroup, uc usecase.SharingTasksUseCase) {
	sharingTasksH := sharingTasks.New(uc)
	sharingTasksRoute := rg.Group("/sharing-tasks")
	sharingTasksRoute.POST("", sharingTasksH.Create())
	sharingTasksRoute.GET("", sharingTasksH.GetSharingTasks())
	sharingTasksRoute.POST("/:diaryID/accept", sharingTasksH.AcceptByDiaryID())
}

func registerDiaryEntriesRoutes(rg *gin.RouterGroup, diaryEntriesUc usecase.DiaryEntriesUseCase) {
	diaryEntriesH := diaryEntries.New(diaryEntriesUc)
	diaryEntriesRoute := rg.Group("/diary-entries")
	diaryEntriesRoute.GET("", diaryEntriesH.GetEntriesList())
	diaryEntriesRoute.GET("/:id", diaryEntriesH.GetByID())
	diaryEntriesRoute.POST("", diaryEntriesH.Create())
	diaryEntriesRoute.PATCH("/:id", diaryEntriesH.PatchEntry())
	diaryEntriesRoute.DELETE("/:id", diaryEntriesH.Delete())
}

func registerDiariesRoutes(rg *gin.RouterGroup, diaryUc usecase.DiaryUseCase) {
	diariesH := diaries.New(diaryUc)
	diariesRoute := rg.Group("/diaries")
	diariesRoute.GET("", diariesH.GetMyDiaries())
}
