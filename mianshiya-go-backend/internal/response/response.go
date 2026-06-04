package response

import "mianshiya-go-backend/internal/errorcode"

type BaseResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func Success(data any) BaseResponse {
	return BaseResponse{
		Code:    0,
		Data:    data,
		Message: "ok",
	}
}

func Error(ec errorcode.ErrorCode) BaseResponse {
	return BaseResponse{
		Code:    ec.Code,
		Data:    nil,
		Message: ec.Message,
	}
}

func ErrorWithMessage(ec errorcode.ErrorCode, message string) BaseResponse {
	return BaseResponse{
		Code:    ec.Code,
		Data:    nil,
		Message: message,
	}
}
