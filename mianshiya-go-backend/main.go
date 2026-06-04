package main

import (
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.Success("ok"))
	})
	r.Run("0.0.0.0:8101")
}
