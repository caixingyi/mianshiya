package response

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
