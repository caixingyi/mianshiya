package handler

import (
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	c.JSON(200, response.Success("ok"))
}

func ErrorDemoHandler(c *gin.Context) {
	c.JSON(200, response.Error(errorcode.ParamsError))
}
