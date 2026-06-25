package file

import (
	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/errorcode"
	"mianshiya-go-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// UploadFileHandler 处理文件上传请求
func (h *Handler) UploadFileHandler(c *gin.Context) {
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

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, "文件不能为空"))
		return
	}

	biz := c.PostForm("biz")

	url, err := h.service.UploadFile(c.Request.Context(), fileHeader, biz, userID)
	if err != nil {
		c.JSON(200, response.ErrorWithMessage(errorcode.ParamsError, err.Error()))
		return
	}

	c.JSON(200, response.Success(url))
}
