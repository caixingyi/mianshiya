package auth

import (
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"
const ContextTokenKey = "token"

func AuthMiddleware(tokenStore TokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		authHeader := c.GetHeader("Authorization")
		// 没有 token 视为未登录
		if authHeader == "" {
			c.JSON(200, response.ErrorWithMessage(errorcode.NotLoginError, "Missing token"))
			c.Abort()
			return
		}
		// token 应该以 "Bearer " 开头
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(200, response.Error(errorcode.NotLoginError))
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		// 空 token 也视为未登录
		if token == "" {
			c.JSON(200, response.Error(errorcode.NotLoginError))
			c.Abort()
			return
		}
		// 验证 token 是否有效
		userID, exists, err := tokenStore.GetUserID(c.Request.Context(), token)
		if err != nil {
			c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
			c.Abort()
			return
		}
		if !exists {
			c.JSON(200, response.ErrorWithMessage(errorcode.NotLoginError, "Invalid token"))
			c.Abort()
			return
		}
		// 将 userID 存储到上下文中
		c.Set(ContextUserIDKey, userID)
		// 将 token 存储到上下文中，方便后续处理
		c.Set(ContextTokenKey, token)
		// 继续处理请求
		c.Next()
	}
}
