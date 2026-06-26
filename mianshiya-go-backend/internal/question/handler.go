package question

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

// AddQuestionHandler 处理添加题目的请求
func (h *Handler) AddQuestionHandler(c *gin.Context) {
	var req AddQuestionRequest
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

	id, err := h.service.AddQuestion(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(id))
}

// GetQuestionVOHandler 处理获取题目详情的请求
func (h *Handler) GetQuestionVOHandler(c *gin.Context) {
	var req GetQuestionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	question, err := h.service.GetQuestionResponseByID(req.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(question))
}

// ListQuestionVOHandler 处理列出题目的请求
func (h *Handler) ListQuestionVOHandler(c *gin.Context) {
	var req ListQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	page, err := h.service.ListQuestions(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}

// ListMyQuestionVOHandler 处理获取我的题目列表的请求
func (h *Handler) ListMyQuestionVOHandler(c *gin.Context) {
	var req ListQuestionRequest
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

	page, err := h.service.ListMyQuestions(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}

// EditQuestionHandler 处理用户编辑题目的请求
func (h *Handler) EditQuestionHandler(c *gin.Context) {
	var req UpdateQuestionRequest
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

	err := h.service.EditQuestion(&req, userID)
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

// DeleteQuestionHandler 处理删除题目的请求
func (h *Handler) DeleteQuestionHandler(c *gin.Context) {
	var req DeleteQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	err := h.service.DeleteQuestion(req.ID)
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

// UpdateQuestionHandler 处理更新题目的请求
func (h *Handler) UpdateQuestionHandler(c *gin.Context) {
	var req UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	err := h.service.UpdateQuestion(&req)
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

// ListQuestionHandler 处理列出题目的请求
func (h *Handler) ListQuestionHandler(c *gin.Context) {
	var req ListQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	page, err := h.service.ListQuestionPage(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}

// BatchDeleteQuestionsHandler 处理批量删除题目的请求
func (h *Handler) BatchDeleteQuestionsHandler(c *gin.Context) {
	var req BatchDeleteQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	err := h.service.BatchDeleteQuestions(&req)
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

// AIGenerateHandler 处理 AI 生成题目的请求（管理员）
// 对应 Java: POST /question/ai/generate/question
func (h *Handler) AIGenerateHandler(c *gin.Context) {
	var req AIGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	// 从 context 取用户 ID（AdminMiddleware 已确保是管理员）
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

	err := h.service.AIGenerateQuestions(req.QuestionType, req.Number, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}
