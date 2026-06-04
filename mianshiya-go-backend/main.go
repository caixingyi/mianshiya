package main

import (
	"mianshiya-go-backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.RegisterRouter(r)

	r.Run("0.0.0.0:8101")
}
