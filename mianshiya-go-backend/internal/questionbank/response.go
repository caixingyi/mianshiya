package questionbank

import "time"

// QuestionBankResponse 定义了题库响应的结构体
type QuestionBankResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Picture     string    `json:"picture"`
	UserID      int64     `json:"userId"`
	EditTime    time.Time `json:"editTime"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
