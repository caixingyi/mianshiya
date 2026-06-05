package main

import (
	"mianshiya-go-backend/internal/config"
	"mianshiya-go-backend/internal/db"
	"mianshiya-go-backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	database, err := db.InitMySQL(cfg.Database)
	if err != nil {
		panic(err)
	}
	_ = database
	r := gin.Default()

	router.RegisterRouter(r)

	r.Run("0.0.0.0:8101")
}
