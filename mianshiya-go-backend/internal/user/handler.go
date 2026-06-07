package user

import (
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
)

// Handler 层负责处理 HTTP 请求，调用 Service 层执行业务逻辑，并返回 HTTP 响应
type Handler struct {
	service    *Service
	tokenStore auth.TokenStore
}

// 构造函数
func NewHandler(service *Service, tokenStore auth.TokenStore) *Handler {
	return &Handler{
		service:    service,
		tokenStore: tokenStore,
	}
}

// 注册 Handler
func (h *Handler) RegisterHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.ParamsError))
		return
	}
	// 2. 调用 Service 层注册用户
	userID, err := h.service.Register(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(userID))
}

// 登录 Handler
func (h *Handler) LoginHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.ParamsError))
		return
	}
	// 2. 调用 Service 层登录用户
	user, err := h.service.Login(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	// 3. 生成 token
	token, err := h.tokenStore.CreateToken(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
		return
	}
	c.JSON(200, response.Success(LoginResponse{
		Token: token,
		User:  user,
	}))
}

// 获取当前登录用户信息的 Handler
func (h *Handler) GetLoginUserHandler(c *gin.Context) {
	// 1. 从上下文中获取 userID
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
	// 2. 调用 Service 层获取用户信息
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
		return
	}
	c.JSON(200, response.Success(user))
}

// 退出登录 Handler
func (h *Handler) LogoutHandler(c *gin.Context) {
	// 1. 从上下文中获取 token
	value, exists := c.Get(auth.ContextTokenKey)
	if !exists {
		c.JSON(200, response.Error(errorcode.NotLoginError))
		return
	}
	token, ok := value.(string)
	if !ok {
		c.JSON(200, response.Error(errorcode.SystemError))
		return
	}
	// 2. 从 tokenStore 中删除 token
	err := h.tokenStore.DeleteToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.SystemError, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

// 更新当前登录用户信息的 Handler
func (h *Handler) UpdateMyHandler(c *gin.Context) {
	// 1. 从上下文中获取 userID
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
	// 2. 解析请求参数
	var req UpdateMyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.ParamsError))
		return
	}
	// 3. 调用 Service 层更新用户信息
	err := h.service.UpdateMy(userID, &req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

// 管理员测试接口
func (h *Handler) AdminCheckHandler(c *gin.Context) {
	c.JSON(200, response.Success(true))
}

// 管理员添加用户 Handler
func (h *Handler) AddUserHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.ParamsError))
		return
	}
	// 2. 调用 Service 层添加用户
	userID, err := h.service.AddUser(&req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(userID))
}

// 管理员删除用户 Handler
func (h *Handler) DeleteUserHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Error(errorcode.ParamsError))
		return
	}
	// 2. 调用 Service 层删除用户
	err := h.service.DeleteUser(req.ID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

// 管理员更新用户 Handler
func (h *Handler) UpdateUserHandler(c *gin.Context) {
	// 1. 解析请求参数
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	// 2. 调用 Service 层更新用户信息
	err := h.service.UpdateUser(req.ID, &req)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}
