package questionbank

import (
	"errors"
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 处理器
type Handler struct {
	service *Service
}

// NewHandler 创建一个新的 Handler 实例
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// AddQuestionBankHandler 处理添加题库的请求
func (h *Handler) AddQuestionBankHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req AddQuestionBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}
	// 2. 从上下文中获取用户 ID
	value, exists := c.Get(auth.ContextUserIDKey)
	if !exists {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, "User ID not found"))
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, "Invalid User ID"))
		return
	}
	// 3. 调用服务层添加题库
	id, err := h.service.AddQuestionBank(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(id))
}

// GetQuestionBankVOHandler 处理获取题库详情的请求
func (h *Handler) GetQuestionBankVOHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req GetQuestionBankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}
	// 2. 调用服务层获取题库详情
	questionBank, err := h.service.GetQuestionBankResponseByID(int64(req.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(questionBank))
}

// ListQuestionBankVOHandler 处理获取题库列表的请求
func (h *Handler) ListQuestionBankVOHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req ListQuestionBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}
	// 2. 调用服务层获取题库列表
	page, err := h.service.ListQuestionBanks(&req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.ErrorWithMessage(errorcode.NotFoundError, "No question banks found"))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(page))
}
