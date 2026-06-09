package questionbankquestion

// 该文件定义了与题库问题相关的请求结构体
type BatchAddQuestionsToBankRequest struct {
	QuestionBankID int64   `json:"questionBankId"`
	QuestionIDs    []int64 `json:"questionIdList"`
}

type BatchRemoveQuestionsFromBankRequest struct {
	QuestionBankID int64   `json:"questionBankId"`
	QuestionIDs    []int64 `json:"questionIdList"`
}
