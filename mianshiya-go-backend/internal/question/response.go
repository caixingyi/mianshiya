package question

import "time"

type QuestionResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Answer    string    `json:"answer"`
	UserID    int64     `json:"userId"`
	TagList   []string  `json:"tagList"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
