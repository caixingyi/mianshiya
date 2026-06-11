package questionbankquestion

// AddQuestionBankQuestionRequest 定义单条添加题库题目关联的请求参数
type AddQuestionBankQuestionRequest struct {
	QuestionBankID int64 `json:"questionBankId" binding:"required"`
	QuestionID     int64 `json:"questionId" binding:"required"`
}

// DeleteQuestionBankQuestionRequest 定义删除题库题目关联的请求参数
type DeleteQuestionBankQuestionRequest struct {
	ID int64 `json:"id"`
}

// UpdateQuestionBankQuestionRequest 定义更新题库题目关联的请求参数
type UpdateQuestionBankQuestionRequest struct {
	ID             int64 `json:"id"`
	QuestionBankID int64 `json:"questionBankId"`
	QuestionID     int64 `json:"questionId"`
}

// GetQuestionBankQuestionRequest 定义获取题库题目关联详情的请求参数
type GetQuestionBankQuestionRequest struct {
	ID int64 `form:"id" binding:"required"`
}

// ListQuestionBankQuestionRequest 定义分页查询题库题目关联的请求参数
type ListQuestionBankQuestionRequest struct {
	Current        int64 `json:"current"`
	PageSize       int64 `json:"pageSize"`
	ID             int64 `json:"id"`
	NotID          int64 `json:"notId"`
	QuestionBankID int64 `json:"questionBankId"`
	QuestionID     int64 `json:"questionId"`
	UserID         int64 `json:"userId"`
}

// RemoveQuestionBankQuestionRequest 定义按题库 ID 和题目 ID 移除关联的请求参数
type RemoveQuestionBankQuestionRequest struct {
	QuestionBankID int64 `json:"questionBankId"`
	QuestionID     int64 `json:"questionId"`
}

// 该文件定义了与题库问题相关的批量请求结构体
type BatchAddQuestionsToBankRequest struct {
	QuestionBankID int64   `json:"questionBankId"`
	QuestionIDs    []int64 `json:"questionIdList"`
}

type BatchRemoveQuestionsFromBankRequest struct {
	QuestionBankID int64   `json:"questionBankId"`
	QuestionIDs    []int64 `json:"questionIdList"`
}
