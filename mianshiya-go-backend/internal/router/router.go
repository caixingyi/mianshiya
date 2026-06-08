package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/handler"
	"mianshiya-go-backend/internal/questionbank"
	"mianshiya-go-backend/internal/user"
)

func RegisterRouter(r *gin.Engine, database *gorm.DB, tokenStore auth.TokenStore) {
	api := r.Group("/api")

	api.GET("/health", handler.HealthHandler)
	api.GET("/error-demo", handler.ErrorDemoHandler)

	repo := user.NewRepository(database)
	service := user.NewService(repo)
	userHandler := user.NewHandler(service, tokenStore)
	questionBankHandler := questionbank.NewHandler(questionbank.NewService(questionbank.NewRepository(database)))
	// 公开接口
	api.POST("/user/register", userHandler.RegisterHandler)
	api.POST("/user/login", userHandler.LoginHandler)
	api.GET("/user/get/vo", userHandler.GetUserVOHandler)
	api.GET("/questionBank/get/vo", questionBankHandler.GetQuestionBankVOHandler)
	api.POST("/questionBank/list/page/vo", questionBankHandler.ListQuestionBankVOHandler)
	// 需要认证的接口
	authAPI := api.Group("")
	authAPI.Use(auth.AuthMiddleware(tokenStore))
	authAPI.GET("/user/get/login", userHandler.GetLoginUserHandler)
	authAPI.POST("/user/logout", userHandler.LogoutHandler)
	authAPI.POST("/user/update/my", userHandler.UpdateMyHandler)
	// 管理员接口
	adminAPI := api.Group("")
	adminAPI.Use(auth.AuthMiddleware(tokenStore), user.AdminMiddleware(service))
	adminAPI.GET("/user/admin/check", userHandler.AdminCheckHandler)
	adminAPI.POST("/user/add", userHandler.AddUserHandler)
	adminAPI.POST("/user/delete", userHandler.DeleteUserHandler)
	adminAPI.POST("/user/update", userHandler.UpdateUserHandler)
	adminAPI.POST("/user/list/page", userHandler.ListUserHandler)
	adminAPI.GET("/user/get", userHandler.GetUserHandler)
	adminAPI.POST("/questionBank/add", questionBankHandler.AddQuestionBankHandler)
}
