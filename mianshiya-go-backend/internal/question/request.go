package question

// AddQuestionRequest 添加题目的请求参数
type AddQuestionRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
	Answer  string   `json:"answer"`
}

// GetQuestionRequest 获取题目详情的请求参数
type GetQuestionRequest struct {
	ID int64 `form:"id" binding:"required"`
}

// ListQuestionRequest 获取题目列表的请求参数
type ListQuestionRequest struct {
	Current        int64    `json:"current"`
	PageSize       int64    `json:"pageSize"`
	ID             int64    `json:"id"`
	NotID          int64    `json:"notId"`
	SearchText     string   `json:"searchText"`
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	Tags           []string `json:"tags"`
	Answer         string   `json:"answer"`
	QuestionBankID int64    `json:"questionBankId"`
	UserID         int64    `json:"userId"`
}

// DeleteQuestionRequest 删除题目的请求参数
type DeleteQuestionRequest struct {
	ID int64 `json:"id"`
}

// UpdateQuestionRequest 更新题目的请求参数
type UpdateQuestionRequest struct {
	ID      int64    `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Answer  string   `json:"answer"`
}

// 批量删除题目的请求参数
type BatchDeleteQuestionRequest struct {
	QuestionIDList []int64 `json:"questionIdList"`
}
