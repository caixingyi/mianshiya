package router

import (
	"github.com/gin-gonic/gin"

	"mianshiya-go-backend/internal/handler"
)

func RegisterRouter(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/health", handler.HealthHandler)
	api.GET("/error-demo", handler.ErrorDemoHandler)

}
