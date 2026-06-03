package main

import (
	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func main() {
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, BaseResponse{
			Code:    0,
			Data:    "ok",
			Message: "ok",
		})
	})
	r.Run("0.0.0.0:8101")
}
