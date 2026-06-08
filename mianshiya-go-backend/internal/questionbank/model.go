package questionbank

import "time"

// QuestionBank 定义了题库的结构体
type QuestionBank struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:256;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Picture     string    `gorm:"type:varchar(1024)" json:"picture"`
	UserID      int64     `gorm:"not null;index" json:"userId"`
	EditTime    time.Time `gorm:"autoUpdateTime" json:"editTime"`
	IsDelete    bool      `gorm:"not null;default:false" json:"isDeleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
