package main

import (
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/config"
	"mianshiya-go-backend/internal/db"
	"mianshiya-go-backend/internal/post"
	"mianshiya-go-backend/internal/question"
	"mianshiya-go-backend/internal/questionbank"
	"mianshiya-go-backend/internal/questionbankquestion"
	"mianshiya-go-backend/internal/router"
	"mianshiya-go-backend/internal/user"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// 初始化mysql数据库
	database, err := db.InitMySQL(cfg.Database)
	if err != nil {
		panic(err)
	}

	// 初始化redis数据库
	rdb, err := db.InitRedis(cfg.Redis)
	if err != nil {
		panic(err)
	}

	// 初始化 token 存储
	tokenStore := auth.NewRedisTokenStore(rdb, 7*24*time.Hour)

	// 自动迁移 User 和 QuestionBank 模型
	if err := database.AutoMigrate(&user.User{}, &questionbank.QuestionBank{}, &question.Question{}, &questionbankquestion.QuestionBankQuestion{}, &post.Post{}); err != nil {
		panic(err)
	}

	// 初始化 Gin 路由
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	router.RegisterRouter(r, database, tokenStore)

	r.Run("0.0.0.0:8101")
}
