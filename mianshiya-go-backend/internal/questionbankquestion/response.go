package questionbankquestion

import "time"

// QuestionBankQuestionResponse 是题库题目关联的对外响应结构
type QuestionBankQuestionResponse struct {
	ID             int64     `json:"id"`
	QuestionBankID int64     `json:"questionBankId"`
	QuestionID     int64     `json:"questionId"`
	UserID         int64     `json:"userId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
