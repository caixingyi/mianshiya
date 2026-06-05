package main

import (
	"mianshiya-go-backend/internal/config"
	"mianshiya-go-backend/internal/db"
	"mianshiya-go-backend/internal/router"
	"mianshiya-go-backend/internal/user"

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
	if err := database.AutoMigrate(&user.User{}); err != nil {
		panic(err)
	}
	r := gin.Default()

	router.RegisterRouter(r, database)

	r.Run("0.0.0.0:8101")
}
