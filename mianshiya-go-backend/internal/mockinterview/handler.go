package mockinterview

import (
	"errors"
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 模拟面试处理器
type Handler struct {
	service *Service
}

// NewHandler 创建模拟面试处理器实例
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// getLoginUserID 从上下文中获取登录用户 ID
func getLoginUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(auth.ContextUserIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := value.(int64)
	if !ok {
		return 0, false
	}
	return userID, true
}

// AddMockInterviewHandler 处理创建模拟面试请求
func (h *Handler) AddMockInterviewHandler(c *gin.Context) {
	var req AddMockInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	userID, ok := getLoginUserID(c)
	if !ok {
		c.JSON(200, response.Error(errorcode.NotLoginError))
		return
	}

	id, err := h.service.AddMockInterview(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(id))
}

// DeleteMockInterviewHandler 处理删除模拟面试请求
func (h *Handler) DeleteMockInterviewHandler(c *gin.Context) {
	var req DeleteMockInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	userID, ok := getLoginUserID(c)
	if !ok {
		c.JSON(200, response.Error(errorcode.NotLoginError))
		return
	}

	err := h.service.DeleteMockInterview(&req, userID)
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

// GetMockInterviewHandler 处理获取模拟面试详情请求
func (h *Handler) GetMockInterviewHandler(c *gin.Context) {
	var req GetMockInterviewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	mockInterview, err := h.service.GetMockInterviewByID(req.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(mockInterview))
}

// ListMockInterviewHandler 处理管理员分页查询模拟面试请求
func (h *Handler) ListMockInterviewHandler(c *gin.Context) {
	var req ListMockInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	page, err := h.service.ListMockInterviews(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}

// ListMyMockInterviewHandler 处理分页查询我的模拟面试请求
func (h *Handler) ListMyMockInterviewHandler(c *gin.Context) {
	var req ListMockInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	userID, ok := getLoginUserID(c)
	if !ok {
		c.JSON(200, response.Error(errorcode.NotLoginError))
		return
	}

	page, err := h.service.ListMyMockInterviews(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}
