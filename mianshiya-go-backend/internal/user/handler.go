package user

import (
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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
	c.JSON(200, response.Success(user))
}
