package main

import (
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/config"
	"mianshiya-go-backend/internal/db"
	"mianshiya-go-backend/internal/router"
	"mianshiya-go-backend/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// 初始化数据库
	database, err := db.InitMySQL(cfg.Database)
	if err != nil {
		panic(err)
	}

	// 初始化 token 存储
	tokenStore := auth.NewMemoryTokenStore()

	// 自动迁移 User 模型
	if err := database.AutoMigrate(&user.User{}); err != nil {
		panic(err)
	}

	// 初始化 Gin 路由
	r := gin.Default()

	router.RegisterRouter(r, database, tokenStore)

	r.Run("0.0.0.0:8101")
}
