package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/handler"
	"mianshiya-go-backend/internal/user"
)

func RegisterRouter(r *gin.Engine, database *gorm.DB, tokenStore auth.TokenStore) {
	api := r.Group("/api")

	api.GET("/health", handler.HealthHandler)
	api.GET("/error-demo", handler.ErrorDemoHandler)

	repo := user.NewRepository(database)
	service := user.NewService(repo)
	userHandler := user.NewHandler(service, tokenStore)
	// 公开接口
	api.POST("/user/register", userHandler.RegisterHandler)
	api.POST("/user/login", userHandler.LoginHandler)
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
}
