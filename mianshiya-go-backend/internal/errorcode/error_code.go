package errorcode

type ErrorCode struct {
	Code    int
	Message string
}

var Success = ErrorCode{Code: 0, Message: "ok"}
var ParamsError = ErrorCode{Code: 40000, Message: "请求参数错误"}
var NotLoginError = ErrorCode{Code: 40100, Message: "未登录"}
var NoAuthError = ErrorCode{Code: 40101, Message: "无权限"}
var ForbiddenError = ErrorCode{Code: 40300, Message: "禁止访问"}
var NotFoundError = ErrorCode{Code: 40400, Message: "请求数据不存在"}
var SystemError = ErrorCode{Code: 50000, Message: "系统内部异常"}
var OperationError = ErrorCode{Code: 50001, Message: "操作失败"}
