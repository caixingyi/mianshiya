package questionbankquestion

import (
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
)

// Handler 定义了题库与题目的关联关系的处理器层
type Handler struct {
	service *Service
}

// NewHandler 创建一个新的 Handler 实例
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// BatchAddQuestionsToBankHandler 处理批量添加题目到题库的请求
func (h *Handler) BatchAddQuestionsToBankHandler(c *gin.Context) {
	var req BatchAddQuestionsToBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
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

	if err := h.service.BatchAddQuestionsToBank(&req, userID); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}

// BatchRemoveQuestionsFromBankHandler 处理批量从题库移除题目的请求
func (h *Handler) BatchRemoveQuestionsFromBankHandler(c *gin.Context) {
	var req BatchRemoveQuestionsFromBankRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	if err := h.service.BatchRemoveQuestionsFromBank(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}
