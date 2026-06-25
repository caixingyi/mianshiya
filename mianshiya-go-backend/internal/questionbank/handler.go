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
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	// 2. 从上下文中获取用户 ID
	value, exists := c.Get(auth.ContextUserIDKey)
	if !exists {
		c.JSON(200, response.Error(errorcode.NotFoundError))
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
	questionBank, err := h.service.GetQuestionBankResponseByID(&req)
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
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	// 2. 调用服务层获取题库列表
	page, err := h.service.ListQuestionBanks(&req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(page))
}

// ListMyQuestionBankVOHandler 处理获取我的题库列表的请求
func (h *Handler) ListMyQuestionBankVOHandler(c *gin.Context) {
	var req ListQuestionBankRequest
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

	page, err := h.service.ListMyQuestionBanks(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}

// EditQuestionBankHandler 处理用户编辑题库的请求
func (h *Handler) EditQuestionBankHandler(c *gin.Context) {
	var req UpdateQuestionBankRequest
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

	err := h.service.EditQuestionBank(&req, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}

// DeleteQuestionBankHandler 处理删除题库的请求
func (h *Handler) DeleteQuestionBankHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req DeleteQuestionBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	// 2. 调用服务层删除题库
	err := h.service.DeleteQuestionBank(req.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

// UpdateQuestionBankHandler 处理更新题库的请求
func (h *Handler) UpdateQuestionBankHandler(c *gin.Context) {
	var req UpdateQuestionBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}

	err := h.service.UpdateQuestionBank(&req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}

// ListQuestionBankHandler 处理获取题库列表的请求
func (h *Handler) ListQuestionBankHandler(c *gin.Context) {
	var req ListQuestionBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}

	page, err := h.service.ListQuestionBankPage(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}
