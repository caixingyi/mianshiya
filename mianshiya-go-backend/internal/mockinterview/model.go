package mockinterview

import "time"

const (
	StatusToStart    = 0
	StatusInProgress = 1
	StatusEnded      = 2
)

// MockInterview 表示一次模拟面试记录
type MockInterview struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkExperience string    `gorm:"size:256" json:"workExperience"`
	JobPosition    string    `gorm:"size:256" json:"jobPosition"`
	Difficulty     string    `gorm:"size:50" json:"difficulty"`
	Messages       string    `gorm:"type:mediumtext" json:"messages"`
	Status         int       `gorm:"not null;default:0" json:"status"`
	UserID         int64     `gorm:"not null;index" json:"userId"`
	IsDelete       int       `gorm:"not null;default:0" json:"isDelete"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
