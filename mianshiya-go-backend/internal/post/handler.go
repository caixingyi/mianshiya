package post

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

// getLoginUserID 从上下文中获取登录用户 ID，未登录返回 0
func getLoginUserID(c *gin.Context) int64 {
	value, exists := c.Get(auth.ContextUserIDKey)
	if !exists {
		return 0
	}
	userID, ok := value.(int64)
	if !ok {
		return 0
	}
	return userID
}

// AddPostHandler 处理添加帖子的请求
func (h *Handler) AddPostHandler(c *gin.Context) {
	var req AddPostRequest
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

	id, err := h.service.AddPost(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(id))
}

// DeletePostHandler 处理删除帖子的请求
func (h *Handler) DeletePostHandler(c *gin.Context) {
	var req DeletePostRequest
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

	err := h.service.DeletePost(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}

// ListPostVOHandler 处理列出帖子的请求
func (h *Handler) ListPostVOHandler(c *gin.Context) {
	var req ListPostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	resp, err := h.service.ListPosts(&req, getLoginUserID(c))
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(resp))
}

// GetPostVOHandler 处理获取帖子详情的请求
func (h *Handler) GetPostVOHandler(c *gin.Context) {
	var req GetPostRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	post, err := h.service.GetPostByID(req.ID, getLoginUserID(c))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, response.Error(errorcode.NotFoundError))
		return
	}
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(post))
}

// UpdatePostHandler 处理更新帖子的请求(管理员)
func (h *Handler) UpdatePostHandler(c *gin.Context) {
	var req UpdatePostRequest
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

	err := h.service.UpdatePost(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}

// EditPostHandler 处理编辑帖子的请求(用户)
func (h *Handler) EditPostHandler(c *gin.Context) {
	var req EditPostRequest
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

	err := h.service.EditPost(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(true))
}

// ListMyPostsVOHandler 处理列出我的帖子的请求
func (h *Handler) ListMyPostsVOHandler(c *gin.Context) {
	var req ListPostsRequest
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

	req.UserID = userID

	resp, err := h.service.ListMyPostsVO(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(resp))
}

// ListPostPageHandler（管理员列表）
func (h *Handler) ListPostPageHandler(c *gin.Context) {
	var req ListPostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	resp, err := h.service.ListPostPage(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(resp))
}

func (h *Handler) SearchPostsHandler(c *gin.Context) {
	var req SearchPostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	resp, err := h.service.SearchPosts(req.Keyword, req.Current, req.PageSize)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(resp))
}
