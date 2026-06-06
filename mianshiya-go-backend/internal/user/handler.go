package user

import (
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service    *Service
	tokenStore *auth.MemoryTokenStore
}

func NewHandler(service *Service, tokenStore *auth.MemoryTokenStore) *Handler {
	return &Handler{
		service:    service,
		tokenStore: tokenStore,
	}
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.ParamsError))
		return
	}
	// 2. 调用 Service 层注册用户
	userId, err := h.service.Register(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(userId))
}

func (h *Handler) LoginHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.ParamsError))
		return
	}
	// 2. 调用 Service 层登录用户
	user, err := h.service.Login(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	// 3. 生成 token
	token, err := h.tokenStore.CreateToken(user.ID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
		return
	}
	c.JSON(200, response.Success(LoginResponse{
		Token: token,
		User:  user,
	}))
}

func (h *Handler) GetLoginUserHandler(c *gin.Context) {
	// 1. 从上下文中获取 userID
	value, exists := c.Get("userID")
	if !exists {
		c.JSON(200, response.Error(errorcode.NotLoginError))
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.JSON(200, response.Error(errorcode.SystemError))
		return
	}
	// 2. 调用 Service 层获取用户信息
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
		return
	}
	c.JSON(200, response.Success(user))
}
