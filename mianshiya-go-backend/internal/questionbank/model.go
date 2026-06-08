package questionbank

import "time"

// QuestionBank 定义了题库的结构体
type QuestionBank struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Picture     string    `json:"picture"`
	UserID      int64     `json:"userId"`
	EditTime    time.Time `gorm:"autoUpdateTime" json:"editTime"`
	IsDelete    bool      `json:"isDeleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
