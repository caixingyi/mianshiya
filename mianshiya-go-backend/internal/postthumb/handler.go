package postthumb

import (
	"errors"
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 帖子点赞处理器
type Handler struct {
	service *Service
}

// NewHandler 创建帖子点赞处理器实例
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// DoPostThumbHandler 处理点赞 / 取消点赞请求
func (h *Handler) DoPostThumbHandler(c *gin.Context) {
	var req AddPostThumbRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	value, exists := c.Get(auth.ContextUserIDKey)
	if !exists {
		c.JSON(200, response.Error(errorcode.NotLoginError))
		return
	}

	userID, ok := value.(int64)
	if !ok {
		c.JSON(200, response.Error(errorcode.SystemError))
		return
	}

	result, err := h.service.DoPostThumb(req.PostID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(result))
}
