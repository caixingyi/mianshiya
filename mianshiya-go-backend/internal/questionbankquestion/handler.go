package questionbankquestion

import (
	"errors"
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 定义了题库与题目的关联关系的处理器层
type Handler struct {
	service *Service
}

// NewHandler 创建一个新的 Handler 实例
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// AddQuestionBankQuestionHandler 处理单条添加题库题目关联的请求
func (h *Handler) AddQuestionBankQuestionHandler(c *gin.Context) {
	var req AddQuestionBankQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	userID, ok := getLoginUserID(c)
	if !ok {
		return
	}

	id, err := h.service.AddQuestionBankQuestion(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(id))
}

// DeleteQuestionBankQuestionHandler 处理删除题库题目关联的请求
func (h *Handler) DeleteQuestionBankQuestionHandler(c *gin.Context) {
	var req DeleteQuestionBankQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	err := h.service.DeleteQuestionBankQuestion(req.ID)
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

// UpdateQuestionBankQuestionHandler 处理更新题库题目关联的请求
func (h *Handler) UpdateQuestionBankQuestionHandler(c *gin.Context) {
	var req UpdateQuestionBankQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	err := h.service.UpdateQuestionBankQuestion(&req)
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

// GetQuestionBankQuestionVOHandler 处理获取题库题目关联详情的请求
func (h *Handler) GetQuestionBankQuestionVOHandler(c *gin.Context) {
	var req GetQuestionBankQuestionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	relation, err := h.service.GetQuestionBankQuestionResponseByID(req.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(relation))
}

// ListQuestionBankQuestionHandler 处理分页获取题库题目关联原始列表的请求
func (h *Handler) ListQuestionBankQuestionHandler(c *gin.Context) {
	var req ListQuestionBankQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	page, err := h.service.ListQuestionBankQuestionPage(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}

// ListQuestionBankQuestionVOHandler 处理分页获取题库题目关联 VO 列表的请求
func (h *Handler) ListQuestionBankQuestionVOHandler(c *gin.Context) {
	var req ListQuestionBankQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	page, err := h.service.ListQuestionBankQuestions(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}

// ListMyQuestionBankQuestionVOHandler 处理分页获取当前登录用户创建的题库题目关联 VO 列表的请求
func (h *Handler) ListMyQuestionBankQuestionVOHandler(c *gin.Context) {
	var req ListQuestionBankQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	userID, ok := getLoginUserID(c)
	if !ok {
		return
	}
	req.UserID = userID

	page, err := h.service.ListQuestionBankQuestions(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(page))
}

// RemoveQuestionBankQuestionHandler 处理按题库 ID 和题目 ID 移除关联的请求
func (h *Handler) RemoveQuestionBankQuestionHandler(c *gin.Context) {
	var req RemoveQuestionBankQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	err := h.service.RemoveQuestionBankQuestion(&req)
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

// BatchAddQuestionsToBankHandler 处理批量添加题目到题库的请求
func (h *Handler) BatchAddQuestionsToBankHandler(c *gin.Context) {
	var req BatchAddQuestionsToBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "请求参数错误"))
		return
	}

	userID, ok := getLoginUserID(c)
	if !ok {
		return
	}

	if err := h.service.BatchAddQuestionsToBank(&req, userID); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}

func getLoginUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(auth.ContextUserIDKey)
	if !exists {
		c.JSON(200, response.Error(errorcode.NotLoginError))
		return 0, false
	}
	userID, ok := value.(int64)
	if !ok {
		c.JSON(200, response.Error(errorcode.SystemError))
		return 0, false
	}
	return userID, true
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
