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
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.ParamsError))
		return
	}
	userId, err := h.service.Register(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(userId))
}
