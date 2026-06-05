package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mianshiya-go-backend/internal/handler"
	"mianshiya-go-backend/internal/user"
)

func RegisterRouter(r *gin.Engine, database *gorm.DB) {
	api := r.Group("/api")

	api.GET("/health", handler.HealthHandler)
	api.GET("/error-demo", handler.ErrorDemoHandler)

	repo := user.NewRepository(database)
	service := user.NewService(repo)
	userHandler := user.NewHandler(service)
	api.POST("/user/register", userHandler.RegisterHandler)
	api.POST("/user/login", userHandler.LoginHandler)
}
