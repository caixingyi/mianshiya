package postfavour

import (
	"errors"
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/post"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 帖子收藏处理器
type Handler struct {
	service *Service
	postSvc *post.Service
}

// NewHandler 创建帖子收藏处理器实例
func NewHandler(service *Service, postSvc *post.Service) *Handler {
	return &Handler{service: service, postSvc: postSvc}
}

// DoPostFavourHandler 处理收藏 / 取消收藏请求
func (h *Handler) DoPostFavourHandler(c *gin.Context) {
	var req AddPostFavourRequest
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

	result, err := h.service.DoPostFavour(req.PostID, userID)
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

// ListMyFavourPostHandler 获取当前用户收藏的帖子列表
func (h *Handler) ListMyFavourPostHandler(c *gin.Context) {
	var req post.ListPostsRequest
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

	req.FavourUserID = userID
	result, err := h.postSvc.ListPosts(&req, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
		return
	}

	c.JSON(200, response.Success(result))
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

// ListFavourPostHandler 获取指定用户收藏的帖子列表
func (h *Handler) ListFavourPostHandler(c *gin.Context) {
	var req PostFavourQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "Invalid request parameters"))
		return
	}

	req.ListPostsRequest.FavourUserID = req.UserID
	result, err := h.postSvc.ListPosts(&req.ListPostsRequest, getLoginUserID(c))
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
		return
	}

	c.JSON(200, response.Success(result))
}
