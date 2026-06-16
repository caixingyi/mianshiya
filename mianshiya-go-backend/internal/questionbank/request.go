package questionbank

// AddQuestionBankRequest 定义了添加题库的请求结构体
type AddQuestionBankRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Picture     string `json:"picture"`
}

// GetQuestionBankRequest 定义了获取题库详情的请求结构体
type GetQuestionBankRequest struct {
	ID                    int64 `form:"id" binding:"required"`
	NeedQueryQuestionList bool  `form:"needQueryQuestionList"`
	Current               int64 `form:"current"`
	PageSize              int64 `form:"pageSize"`
}

// ListQuestionBankRequest 定义了获取题库列表的请求结构体
type ListQuestionBankRequest struct {
	Current     int64  `json:"current"`
	PageSize    int64  `json:"pageSize"`
	ID          int64  `json:"id"`
	NotID       int64  `json:"notId"`
	SearchText  string `json:"searchText"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	UserID      int64  `json:"userId"`
}

// DeleteQuestionBankRequest 定义了删除题库的请求结构体
type DeleteQuestionBankRequest struct {
	ID int64 `json:"id"`
}

// UpdateQuestionBankRequest 定义了更新题库的请求结构体
type UpdateQuestionBankRequest struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
}
