package user

import (
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从上下文中获取 userID
		userID, exists := c.Get(auth.ContextUserIDKey)
		if !exists {
			c.JSON(200, response.Error(errorcode.NotLoginError))
			c.Abort()
			return
		}
		// 2. userID 应该是 int64 类型
		userIDInt64, ok := userID.(int64)
		if !ok {
			c.JSON(200, response.Error(errorcode.SystemError))
			c.Abort()
			return
		}
		// 3. 验证 userID 是否为管理员
		isAdmin, err := s.IsAdmin(userIDInt64)
		if err != nil {
			c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
			c.Abort()
			return
		}
		if !isAdmin {
			c.JSON(200, response.Error(errorcode.NoAuthError))
			c.Abort()
			return
		}
		// 4. 继续处理请求
		c.Next()
	}
}
